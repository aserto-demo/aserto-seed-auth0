package auth0

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/aserto-demo/aserto-seed-auth0/pkg/config"
	"github.com/aserto-demo/aserto-seed-auth0/pkg/counter"
	"github.com/aserto-demo/aserto-seed-auth0/pkg/csv"
	"github.com/sirupsen/logrus"
	"gopkg.in/auth0.v4"
	"gopkg.in/auth0.v5/management"
)

const (
	objectTypeUser  = "user"
	finalCount      = 0
	counterInterval = 10
)

// Manager -.
type Manager struct {
	config      *config.Config
	inputfile   string
	mgnt        *management.Management
	identityMap map[string]string
	spew        bool
	exec        bool
	counter     counter.Counter
}

// NewManager -.
func NewManager(cfg *config.Config, filename string) *Manager {
	seeder := Manager{
		config:      cfg,
		inputfile:   filename,
		identityMap: make(map[string]string),
		exec:        true,
		counter:     counter.Counter{},
	}
	return &seeder
}

// Init -.
func (m *Manager) Init() error {
	if err := m.loadIdentityMap(); err != nil {
		return fmt.Errorf("%w", err)
	}

	mgnt, err := management.New(
		m.config.Auth0.Domain,
		management.WithClientCredentials(
			m.config.Auth0.ClientID,
			m.config.Auth0.ClientSecret,
		),
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	m.mgnt = mgnt

	return nil
}

// Spew -.
func (m *Manager) Spew(f bool) {
	m.spew = f
}

// Dryrun -.
func (m *Manager) Dryrun(f bool) {
	m.exec = !f
}

// Seed -.
func (m *Manager) Seed() error {
	cr := csv.NewCsvReader()
	if err := cr.Open(m.inputfile); err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	for {
		if err := cr.Read(); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%w", err)
		}

		// if not a user object skip to next
		if s := cr.Get("ObjectClass"); s != objectTypeUser {
			m.counter.IncrSkipped()
			continue
		}

		u := m.makeUser(cr)

		if m.spew {
			_ = enc.Encode(u)
		}

		if m.userExists(*u.ID) {
			if err := m.updateUser(*u.ID, u); err != nil {
				continue
			}
		} else {
			if err := m.createUser(*u.ID, u); err != nil {
				continue
			}
		}

		m.counter.IncrRows()
		m.counter.Print(counterInterval)
	}

	m.counter.Print(finalCount)

	return nil
}

// Reset -.
func (m *Manager) Reset() error {
	cr := csv.NewCsvReader()
	if err := cr.Open(m.inputfile); err != nil {
		return err
	}

	for {
		if err := cr.Read(); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%w", err)
		}

		// if not a user object skip to next
		if cr.Get("ObjectClass") != objectTypeUser {
			m.counter.IncrSkipped()
			continue
		}

		id := cr.GetToLower("ObjectGUID")

		if err := m.deleteUser(id); err != nil {
			continue
		}

		m.counter.IncrRows()
		m.counter.Print(counterInterval)
	}

	m.counter.Print(finalCount)

	return nil
}

func (m *Manager) makeUser(cr *csv.Reader) *management.User {
	id := cr.GetToLower("ObjectGUID")

	email := formatEmail(cr.GetToLower("mail"), m.config.EmailDomain)

	u := management.User{
		ID:            auth0.String(id),
		Connection:    auth0.String("Username-Password-Authentication"),
		Email:         auth0.String(email),
		EmailVerified: auth0.Bool(true),
		GivenName:     auth0.String(cr.Get("GivenName")),
		FamilyName:    auth0.String(cr.Get("sn")),
		Nickname:      auth0.String(cr.Get("Name")),
		Password:      auth0.String(m.config.SetPassword),
		UserMetadata:  make(map[string]interface{}),
		AppMetadata:   make(map[string]interface{}),
		Picture:       auth0.String(picURL(cr.Get("Name"))),
	}

	u.UserMetadata["phone"] = formatPhone(cr.Get("OfficePhone"))
	u.UserMetadata["title"] = cr.Get("Title")
	u.UserMetadata["department"] = cr.Get("Department")
	u.UserMetadata["manager"] = m.identityMap[cr.GetToLower("Manager")]
	u.UserMetadata["username"] = cr.GetToLower("SamAccountName")

	corporationRole := strings.Split(m.config.EmailDomain, ".")[0]
	departmentRole := strings.ReplaceAll(cr.GetToLower("Department"), " ", "-")
	u.AppMetadata["roles"] = [...]string{"user", corporationRole, departmentRole}

	if val := cr.GetToLower("DistinguishedName"); val != "" {
		u.UserMetadata["dn"] = val
	}

	return &u
}

func (m *Manager) userExists(id string) bool {
	if _, err := m.mgnt.User.Read("auth0|" + id); err != nil {
		return false
	}
	return true
}

func (m *Manager) createUser(id string, u *management.User) error {
	logrus.Debugf("createUser: %s\n", "auth0|"+id)
	if !m.exec {
		return nil
	}

	if err := m.mgnt.User.Create(u); err != nil {
		m.counter.IncrError()
		logrus.Error(err)
		return err
	}

	return nil
}

func (m *Manager) updateUser(id string, u *management.User) error {
	logrus.Debugf("updateUser: %s\n", "auth0|"+id)
	if !m.exec {
		return nil
	}

	// reset fields which cannot be changed
	u.ID = nil
	u.Password = nil

	if err := m.mgnt.User.Update("auth0|"+id, u); err != nil {
		m.counter.IncrError()
		logrus.Error(err)
		return err
	}

	return nil
}

func (m *Manager) deleteUser(id string) error {
	if !m.exec {
		return nil
	}

	if m.userExists(id) {
		if err := m.mgnt.User.Delete("auth0|" + id); err != nil {
			m.counter.IncrError()
			logrus.Error(err)
			return err
		}
	}

	return nil
}

func (m *Manager) loadIdentityMap() error {
	cr := csv.NewCsvReader()
	if err := cr.Open(m.inputfile); err != nil {
		return err
	}

	for {
		if err := cr.Read(); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%w", err)
		}

		// if not a user object skip to next
		if cr.Get("ObjectClass") != objectTypeUser {
			continue
		}

		id := cr.GetToLower("ObjectGUID")
		m.identityMap[id] = id

		if val := cr.GetToLower("mail"); val != "" {
			m.identityMap[val] = id
		}

		if val := cr.GetToLower("DistinguishedName"); val != "" {
			m.identityMap[val] = id
		}

		if val := cr.GetToLower("UserPrincipalName"); val != "" {
			m.identityMap[val] = id
		}
	}

	return nil
}

var (
	regexPhone1 *regexp.Regexp = regexp.MustCompile(`\)`)
	regexPhone2 *regexp.Regexp = regexp.MustCompile(`[\(\s]`)
)

func formatPhone(s string) string {
	phone := "+1-" + s
	phone = regexPhone1.ReplaceAllLiteralString(phone, "-")
	phone = regexPhone2.ReplaceAllLiteralString(phone, "")
	return phone
}

func formatEmail(s, emailDomain string) string {
	return strings.Replace(s, "@contoso.com", "@"+emailDomain, 1)
}

func picURL(s string) string {
	u, err := url.Parse(
		fmt.Sprintf("https://github.com/aserto-demo/contoso-ad-sample/raw/main/UserImages/%s.jpg",
			s,
		))
	if err != nil {
		return ""
	}

	return u.String()
}
