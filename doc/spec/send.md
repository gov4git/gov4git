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

     - sending a message entails writing a file with the message contents into a specific "sent" directory within the sender's repo

     - the path to the "sent" directory is uniquely determined by the receiver's identity and the topic; it is given by

          `/id/mail/sent/{ HashOf(receiver_id) }/{ HashOf(topic) }`
     
     where the function `HashOf` (implemented [here](https://github.com/gov4git/lib4git/blob/main/form/bytes.go#L38)) returns the lowercase Base32 Standard encoding with no padding of the SHA256 hash of the argument.

     Let's call the resulting path `SENT_DIR`.

     - read the next available message sequence number from the file `SENT_DIR/next.json`; if present, this file contains a single positive integer, otherwise the next available sequence number is 0. Let's call the string representation of this number `NEXT_SEQNO`.
     
     - write the JSON encoding of the signed message (computed in the previous step) to the file `SENT_DIR/NEXT_SEQNO`

     - write the number `NEXT_SEQNO+1` in JSON format to the file `SENT_DIR/next.json`
