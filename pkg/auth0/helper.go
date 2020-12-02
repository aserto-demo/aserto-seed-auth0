package auth0

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/aserto-demo/aserto-seed-auth0/pkg/csv"
	"github.com/sirupsen/logrus"
	"gopkg.in/auth0.v4"
	"gopkg.in/auth0.v5/management"
)

const (
	objectTypeUser  = "user"
	counterInterval = 10
)

// Helper -.
type Helper struct {
	inputfile    string
	domain       string
	clientID     string
	clientSecret string
	emailDomain  string
	setPassword  string
	mgnt         *management.Management
	identityMap  map[string]string
	spew         bool
	exec         bool
	counter      Counter
}

// NewHelper -.
func NewHelper(filename string) *Helper {
	seeder := Helper{
		inputfile:    filename,
		domain:       os.Getenv("AUTH0_DOMAIN"),
		clientID:     os.Getenv("AUTH0_CLIENT_ID"),
		clientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		emailDomain:  os.Getenv("EMAIL_DOMAIN"),
		setPassword:  os.Getenv("SET_PASSWORD"),
		identityMap:  make(map[string]string),
		exec:         true,
		counter:      Counter{},
	}
	return &seeder
}

// Init -.
func (h *Helper) Init() error {
	if err := h.loadIdentityMap(); err != nil {
		return fmt.Errorf("%w", err)
	}

	mgnt, err := management.New(
		h.domain,
		management.WithClientCredentials(
			h.clientID,
			h.clientSecret,
		),
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	h.mgnt = mgnt

	return nil
}

// Spew -.
func (h *Helper) Spew(f bool) {
	h.spew = f
}

// Dryrun -.
func (h *Helper) Dryrun(f bool) {
	h.exec = !f
}

// Seed -.
func (h *Helper) Seed() error {
	cr := csv.NewCsvReader()
	if err := cr.Open(h.inputfile); err != nil {
		return err
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	corporationRole := strings.Split(h.emailDomain, ".")[0]

	for {
		if err := cr.Read(); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("%w", err)
		}

		// if not a user object skip to next
		if s := cr.Get("ObjectClass"); s != objectTypeUser {
			h.counter.AddSkipped()
			continue
		}

		id := cr.GetToLower("ObjectGUID")

		email := formatEmail(cr.GetToLower("mail"), h.emailDomain)

		u := management.User{
			ID:            auth0.String(id),
			Connection:    auth0.String("Username-Password-Authentication"),
			Email:         auth0.String(email),
			EmailVerified: auth0.Bool(true),
			GivenName:     auth0.String(cr.Get("GivenName")),
			FamilyName:    auth0.String(cr.Get("sn")),
			Nickname:      auth0.String(cr.Get("Name")),
			Password:      auth0.String(h.setPassword),
			UserMetadata:  make(map[string]interface{}),
			AppMetadata:   make(map[string]interface{}),
			Picture:       auth0.String(picURL(cr.Get("Name"))),
		}

		u.UserMetadata["phone"] = formatPhone(cr.Get("OfficePhone"))
		u.UserMetadata["title"] = cr.Get("Title")
		u.UserMetadata["department"] = cr.Get("Department")
		u.UserMetadata["manager"] = h.identityMap[cr.GetToLower("Manager")]
		u.UserMetadata["username"] = cr.GetToLower("SamAccountName")

		departmentRole := strings.ReplaceAll(cr.GetToLower("Department"), " ", "-")
		u.AppMetadata["roles"] = [...]string{"user", corporationRole, departmentRole}

		if val := cr.GetToLower("DistinguishedName"); val != "" {
			u.UserMetadata["dn"] = val
		}

		if h.spew {
			_ = enc.Encode(&u)
		}

		if h.exec {
			if h.userExists(id) {
				u.ID = nil
				u.Password = nil
				logrus.Infof("update user_id: %s\n", "auth0|"+id)
				if err := h.mgnt.User.Update("auth0|"+id, &u); err != nil {
					h.counter.AddError()
					logrus.Error(err)
					continue
				}
			} else {
				logrus.Infof("create user_id: %s\n", "auth0|"+id)
				if err := h.mgnt.User.Create(&u); err != nil {
					h.counter.AddError()
					logrus.Error(err)
					continue
				}
			}
		}

		h.counter.AddRow()
		h.counter.Print(counterInterval)
	}

	h.counter.Print(0) // final count

	return nil
}

// Reset -.
func (h *Helper) Reset() error {
	cr := csv.NewCsvReader()
	if err := cr.Open(h.inputfile); err != nil {
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
			h.counter.AddSkipped()
			continue
		}

		id := cr.GetToLower("ObjectGUID")

		if h.exec {
			if h.userExists(id) {
				if err := h.mgnt.User.Delete("auth0|" + id); err != nil {
					h.counter.AddError()
					logrus.Error(err)
				}
			}
		}

		h.counter.AddRow()
		h.counter.Print(counterInterval)
	}

	h.counter.Print(0) // final count

	return nil
}

func (h *Helper) userExists(id string) bool {
	if _, err := h.mgnt.User.Read("auth0|" + id); err != nil {
		return false
	}
	return true
}

func (h *Helper) loadIdentityMap() error {
	cr := csv.NewCsvReader()
	if err := cr.Open(h.inputfile); err != nil {
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
		h.identityMap[id] = id

		if val := cr.GetToLower("mail"); val != "" {
			h.identityMap[val] = id
		}

		if val := cr.GetToLower("DistinguishedName"); val != "" {
			h.identityMap[val] = id
		}

		if val := cr.GetToLower("UserPrincipalName"); val != "" {
			h.identityMap[val] = id
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

// Counter - accumulator for row, skipped and error counts.
type Counter struct {
	rowCounter  int32
	skipCounter int32
	errCounter  int32
}

// AddRow - increase row counter.
func (c *Counter) AddRow() {
	atomic.AddInt32(&c.rowCounter, 1)
}

// AddSkipped - increase skipped row counter.
func (c *Counter) AddSkipped() {
	atomic.AddInt32(&c.skipCounter, 1)
}

// AddError - increase error counter.
func (c *Counter) AddError() {
	atomic.AddInt32(&c.errCounter, 1)
}

// Print - print counter at interval % m.
func (c *Counter) Print(m int32) {
	// i := m
	linefeed := ""

	if m == 0 {
		linefeed = "\n"
		m = 1
	}

	if d := c.rowCounter % m; d == 0 {
		fmt.Fprintf(os.Stdout, "\033[2K\rrow count: %d skip count %d error count: %d%s",
			c.rowCounter,
			c.skipCounter,
			c.errCounter,
			linefeed,
		)
	}
}
