# Vanity IPNS generator

This tool generates Ed25519 private keys for which the base36 IPNS ends with a given suffix.
The first match found is stored in the current working directory.
You can then move the generate file to `~/.ipfs/keystore/` to use it with Kubo.

## Installation

```sh
go install github.com/codesensus-org/vanity-ipns@latest
```

## Usage

```sh
vanity-ipns {suffix}
```

Example:
```sh
vanity-ipns hello
mv 12D3KooW* ~/.ipfs/keystore/hello
ipfs name publish --key=hello /ipfs/bafkreidfdrlkeq4m4xnxuyx6iae76fdm4wgl5d4xzsb77ixhyqwumhz244
```
