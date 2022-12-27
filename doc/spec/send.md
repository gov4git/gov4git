## SUMMARY

- function "send a JSON structure"
- implemented in [mail/send.go](../../proto/mail/send.go)

## CONTEXT

The send operation implements sending a JSON structure from one identity to another on a given topic.

## PSEUDOCODE

The input to a send operation is:
- a sender, identified by their public and private repos
- a receiver, identified by their public repo
- a topic, any string
- a message, any JSON structure

1. retrieve sender's private credentials

     - clone the sender's private repo
     - read the sender's private credentials (the [procedure for initializing credentials](init-user-id.md) describes the location and format of the private credentials)
     - the sender's private credentials include the sender's ED25519 public and private signing keys

2. sign the message (this is implemented [here](https://github.com/gov4git/gov4git/blob/main/proto/id/crypto.go#L48))

     - encode the message in JSON format; let's call the resulting encoding the "plaintext"
     - sign the plaintext using the sender's private ED25519 key; let's call the resulting bytes the "signature"
     - prepare a JSON structure holding the message and the signature; the structure looks like this

     ```json
     {
          "plaintext": "ewogICAiYmFsb ... CBdCn0=",
          "signature": "6FrwoTLb ... vbC9j4JbSBA==",
          "ed25519_public_key": "yTr4w8DwUXaEZFXpTBRxfaG6F62JD7Ol1j034uRDY+4="
     }
     ```

     `plaintext` is a Base64 Standard encoding of the plaintext.

     `signature` is a Base64 Standard encoding of the signature.

     `ed25519_public_key` is a Base64 Standard encoding of the sender's public ED25519 key.

3. drop the message in the sender's public repo (this is [implemented here](https://github.com/gov4git/gov4git/blob/main/proto/mail/send.go#L16))

     - XXX