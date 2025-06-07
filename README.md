# Vanity IPNS generator

This tool generates Ed25519 private keys for which the base36 IPNS ends with a given suffix.
The first match found is stored in the current working directory.
You can then import them to use on Kubo with `ipfs key import`.

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
ipfs key import hello k51qzi5uqu5d*hello
ipfs name publish --key=hello /ipfs/bafkreidfdrlkeq4m4xnxuyx6iae76fdm4wgl5d4xzsb77ixhyqwumhz244
```
