# aserto-seed-auth0

A simple tool for populating users in an Auth0 domain. 

##Configuration
The tool requires a .env file contain configuration information for Auth0 as well as some additional settings. 

A template .env file is available in the root of the source repository.

.env file 

    AUTH0_DOMAIN="mydomain.us.auth0.com"
    AUTH0_CLIENT_ID="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    AUTH0_CLIENT_SECRET="yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
	EMAIL_DOMAIN="acmecorp.com"
	SET_PASSWORD="V@erySec#re123!"


##Tool options

Command line options:

	Usage of ./bin/aserto-seed-auth0:
      --input string   inputfile
      --reset          reset
      --seed           seed
      --spew           spew
      --version        version
	

#Seed

The seed data file used, resides in the [aserto-demo/contoso-ad-sample](https://github.com/aserto-demo/contoso-ad-sample) repository and provides a dataset of 272 users. 

	aserto-seed-auth0 --input ./data/ADUsers.csv --seed
	

Seed with output:

	aserto-seed-auth0 --input ./data/ADUsers.csv --seed --spew
	
Spew output:

	{
	  "user_id": "b7de08a6-8417-491b-be62-85945a538f46",
	  "connection": "Username-Password-Authentication",
	  "email": "danj@acmecorp.com",
	  "given_name": "Dan",
	  "family_name": "Jump",
	  "nickname": "Dan Jump",
	  "password": "***********",
	  "user_metadata": {
	    "department": "Executive",
	    "dn": "cn=dan jump",
	    "manager": "",
	    "phone": "+1-425-555-0179",
	    "title": "CEO",
	    "username": "danj"
	  },
	  "email_verified": true,
	  "app_metadata": {
	    "roles": [
	      "user",
	      "acmecorp",
	      "executive"
	    ]
	  },
	  "picture": "https://github.com/aserto-demo/contoso-ad-sample/raw/main/UserImages/Dan%20Jump.jpg"
	}...

##Reset

Reset, will remove the all users from the seed input file from the user domain, based on their ID. 

**NOTE:** Any user not in the seed input file will not be removed or changed.

	aserto-seed-auth0 --input ./data/ADUsers.csv --reset
	
