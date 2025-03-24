---
title: Restructed Text to Markdown
weight: 40
---

[RST2MD](https://github.com/SPANDigital/presidium-rst-to-markdown) is a tool that allows you to convert RST files into Presidium Markdown. It’s a custom wrapper around Pandoc, that slightly modifies the output to be more compatible with Presidium. (Link to RST2MD binary still needs to be added)

### Installation
Installation is dependent on how the tool will be hosted.

### Usage

The tool accepts an input (RST source files) and output parameter (location where generated MD files should be saved).

```
./rst2md -v -input /path/to/rst/source/ -output /path/to/output/dir`
```