# aserto-seed-auth0

A simple tool for populating users in an Auth0 domain. 

##Configuration
The tool requires a .env file contain configuration information for Auth0 as well as some additional settings. 

A [.env.template](https://raw.githubusercontent.com/aserto-demo/aserto-seed-auth0/main/.env.template) .env file is available in the root of the source repository. Copy the file to .env and adjust the values accordingly.

.env file 

	AUTH0_DOMAIN="mydomain.us.auth0.com"
	AUTH0_CLIENT_ID="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	AUTH0_CLIENT_SECRET="yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"
	EMAIL_DOMAIN="contoso.com"
	SET_PASSWORD="V@erySec#ret321!"


##Tool options

Command line options:

	> aserto-seed-auth0 --help

	NAME:
	aserto-seed-auth0 - seed Auth0 user data

	USAGE:
	main [global options] command [command options] [arguments...]

	COMMANDS:
	seed     seed
	reset    reset
	version  display verion information

	GLOBAL OPTIONS:
	--help, -h  show help (default: false)

#Seed

The seed data file used, resides in the [aserto-demo/contoso-ad-sample](https://github.com/aserto-demo/contoso-ad-sample) repository and provides a dataset of 272 users. 

	aserto-seed-auth0 seed --input ./data/ADUsers.csv
	

Seed with output:

	aserto-seed-auth0 seed --input ./data/ADUsers.csv --spew --dryrun
	
Spew output:

	{
	  "user_id": "b7de08a6-8417-491b-be62-85945a538f46",
	  "connection": "Username-Password-Authentication",
	  "email": "danj@contoso.com",
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

	aserto-seed-auth0 reset --input ./data/ADUsers.csv
	
