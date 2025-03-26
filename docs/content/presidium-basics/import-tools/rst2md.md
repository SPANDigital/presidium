---
title: Restructed Text to Markdown
weight: 40
---

[RST2MD](https://github.com/SPANDigital/presidium-rst-to-markdown) is a tool that allows you to convert RST files into Presidium Markdown. It’s a custom wrapper around Pandoc, that slightly modifies the output to be more compatible with Presidium. (Link to RST2MD binary still needs to be added)

### Installation


Installation methods are listed in order of preference:

#### Via homebrew (recommended)

This requires [homebrew](https://brew.sh/) to be installed.

```shell
brew tap SPANDigital/homebrew-tap
brew update
brew install rst2md
```

#### Via go install (for go developers)

This requires as least [Go 1.22.x](https://go.dev/doc/install) to be installed for your operating system.

```shell
go install github.com/spandigital/presidium-rst-to-markdown/cmd/rst2md
```

### Usage

The tool accepts an input (RST source files) and output parameter (location where generated Markdown files should be saved).

```
./rst2md -v -input /path/to/rst/source/ -output /path/to/output/dir`
```