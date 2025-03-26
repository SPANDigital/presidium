---
title: JSON Schema
weight: 50
---

presidium-json-schema is a CLI tool for importing your JSON Schema spec into Presidium documentation.

### Install

```
brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git brew install presidium-json-schema
```

### Usage

`presidium-json-schema convert [path] [flags]`

Flags:
* `-d` `--destination <string>`: The output directory (default `.`)
* `-e` `--extension <string>`: The schema extension (default `*.schema.json`)
* `-o` `--ordered`: Preserve the schema order (defaults to alphabetical)
* `-w` `--walk`: Walk through sub-directories

To convert a file, simply run:

`presidium-json-schema convert <PATH_TO_SCHEMA_DIR> -d <THE_DESTINATION_DIR>`

