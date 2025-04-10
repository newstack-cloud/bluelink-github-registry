# Signing Keys

### Generating signing keys environment variable

You can use the tool in `tools/signing-keys` to generate/insert the environment variable for the signing keys.
This tool takes one or more existing signing public key files as arguments and generates a JSON object of keys that is formatted into a string suitable for use in an environment variable.
The input files are expected to contain the output of the `gpg --armor --export "{name} <{email}>"` command for each GPG key.

To generate the environment variable, run the following command:

```bash
go run tools/signing-keys/main.go <key_file_1> <key_file_2> ...
```

To insert into your existing .env file, you can use the `-insert` flag:

```bash
go run tools/signing-keys/main.go -insert=.env <key_file_1> <key_file_2> ...
```
