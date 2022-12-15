
## SUMMARY

- function "create new user identity"
- implemented in [init.go](../../proto/id/init.go)

## PSEUDOCODE

1. provision user's public and private repos
	- ask the user to create two git repositories, one called "public" and one "private"
		- the public repo: everyone can read, only the user can write
		- the private repo: only the user can read or write
	- let USER_PUBLIC be the url for authenticated writable access to the user's public repo
		- this would be a git ssh URL, e.g. `git@github.com:petar/gov4git.public.git`
	- let USER_PRIVATE be the url for authenticated writable access to the user's private repo

2. generate signing keys
     - generate a pair of ED25519 signing keys: USER_ED25519_PUBLIC and USER_ED25519_PRIVATE
     - generate a new UUID for the user: USER_UUID

3. prepare a JSON encoding of the user's public and private identities

     - prepare the content of the user's public credentials file:

          ```json
          {
               "id": STRING_OF(USER_UUID),
               "public_key_ed25519": BASE64_STD_ENCODING(USER_ED25519_PUBLIC),
          }
          ```

          for example

          ```json
          {
               "id":"c5e0b300-82b5-41c3-9bce-24f8c159363e",
               "public_key_ed25519":"yTr4w8DwUXaEZFXpTBRxfaG6F62JD7Ol1j034uRDY+4="
          }
          ```

     - prepare the content of the user's private credentials file:

          ```json
          {
               "private_key_ed25519": BASE64_STD_ENCODING(USER_ED25519_PRIVATE),
               "public_credentials": {
                    "id": STRING_OF(USER_UUID),
                    "public_key_ed25519": BASE64_STD_ENCODING(USER_ED25519_PUBLIC),
               }
          }
          ```

4. write user's private credentials to user's private repo
     - clone branch `main` of USER_PRIVATE (note that the repo may be empty)
     - if file `id/private_credentials.json` already exists, abort with an error "identity already exists"
     - write the user's private credentials to file `id/private_credentials.json`
     - commit changes with message "gov4git: Initialized private credentials."
     - push to origin

5. write user's public credentials to user's public repo
     - clone branch `main` of USER_PUBLIC (note that the repo may be empty)
     - write the user's public credentials to file `id/public_credentials.json`
     - commit changes with message "gov4git: Initialized public credentials."
     - push to origin
