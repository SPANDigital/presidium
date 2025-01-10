---
title: JSON Schema
weight: 4
status: GOOD
---

[presidium-json-schema](https://github.com/SPANDigital/presidium-json-schema) is a CLI tool for importing your [JSON Schema](http://json-schema.org/) spec into
[Presidium](http://presidium.spandigital.net) documentation.

## Install

brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git
brew install presidium-json-schema

## Usage

```text
Usage:
  presidium-json-schema convert [path] [flags]

Flags:
  -d, --destination string   the output directory (default ".")
  -e, --extension string     the schema extension (default "*.schema.json")
  -o, --ordered              preserve the schema order (defaults to alphabetical)
  -w, --walk                 walk through sub-directories
```

To convert a file you simply:

```shell
presidium-json-schema convert <PATH_TO_SCHEMA_DIR> -d <THE_DESTINATION_DIR>
```
