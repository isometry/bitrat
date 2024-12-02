# bitrat

Lightning-fast, multi-algorithm file checksums.

## Overview

`bitrat` is a command-line tool to quickly calculate checksums for nested file hierarchies, such that subsequent changes to those files can be easily detected, with the goal of identifying files corrupted by bitrot or other unexpected change. It supports multiple hashing algorithms and provides a fast and efficient way to ensure the integrity of your files.

## Features

- Cross-platform compatibility (macOS, Linux, FreeBSD, Windows).
- Supports a wide range of hashing algorithms including BLAKE3, BLAKE2, SHA3, SHA2, and more (`bitrat list-algorithms`).
- Supports HMAC for added security (`bitrat --hmac my-secret`).
- Exploits multi-core systems with fast, parallel processing pipelines.

## Installation

### Go

```shell
go install github.com/isometry/bitrat@latest
```

### Homebrew

```shell
brew install isometry/tap/bitrat
```

## Usage

```shell
bitrat
```

## Future

In future, the tool should be extended to support:

- [ ] support direct verification of checksums in a file (`--check`).
- [ ] support storage of file checksums in extended attributes to streamline the detection of changes.
- [ ] support client/server mode to support locally-stateless centralisation of integrity checks.
