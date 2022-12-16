
## SUMMARY

- function "create new user identity"
- implemented in [init.go](../../proto/id/init.go)

## PSEUDOCODE

1. provision user's public and private repos
	- ask the user to create two git repositories, one called "public" and one "private"
		- the public repo: everyone can read, only the user can write
		- the private repo: only the user can read or write
     - the input to the identity initialization procedure is:
     	- the git url for authenticated writable access to the user's public repo.
          this would be a git ssh URL, e.g. `git@github.com:petar/gov4git.public.git`
     	- the git url for authenticated writable access to the user's private repo.

1. generate signing keys
     - generate a pair (public and private) of ED25519 signing keys
     - generate a new UUID for the user

2. prepare a JSON encoding of the user's public and private identities

     - prepare the user's public credentials JSON structure; it looks like this

          ```json
          {
               "id": "c5e0b300-82b5-41c3-9bce-24f8c159363e",
               "public_key_ed25519": "yTr4w8 ... 34uRDY+4="
          }
          ```

          `public_key_ed25519` is the Base64 Standard encoding of the ED25519 public key.

     - prepare the content of the user's private credentials file:

          ```json
          {
               "private_key_ed25519": "Tfi5EN ... kPs6XWPTfi5ENj7g==",
               "public_credentials": {
                    "id": "c5e0b300-82b5-41c3-9bce-24f8c159363e",
                    "public_key_ed25519": "yTr4w8 ... 34uRDY+4="
               }
          }
          ```

          `private_key_ed25519` is the Base64 Standard encoding of the ED25519 private key.

          `public_credentials` is a copy of the public credentials structure computed on the previous step.

3. write user's private credentials to user's private repo
     - clone branch `main` of the user's private repo (note that the repo may be empty)
     - if file `id/private_credentials.json` already exists, abort with an error "identity already exists"
     - write the user's private credentials to file `id/private_credentials.json`
     - commit changes with message "gov4git: Initialized private credentials."
     - push to origin

4. write user's public credentials to user's public repo
     - clone branch `main` of the user's public repo (note that the repo may be empty)
     - write the user's public credentials to file `id/public_credentials.json`
     - commit changes with message "gov4git: Initialized public credentials."
     - push to origin
