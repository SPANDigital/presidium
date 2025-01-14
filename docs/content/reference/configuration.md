---
title: Configuration
weight: 1
status: UPTDATED
---

This is a sample Presidium configuration file:

```yaml
languageCode: en-us
title: "Sample Title"
copyright: Template Footer
pluralizelisttitles: false
markup:
  goldmark:
    renderer:
      Unsafe: true
    parser:
      attribute:
        block: true
  highlight:
    style: github
menu:
  Main:
  - identifier: organization-overview
    name: Organization Overview
    url: /
    weight: 1
  - identifier: solution-overview
    name: Solution Overview
    url: /solution-overview/
    weight: 2
  - identifier: technology-stack
    name: Technology Stack
    url: /technology-stack/
    weight: 3
  - identifier: tool-chain
    name: Tool Chain
    url: /tool-chain/
    weight: 4
  - identifier: dev-environment-setup
    name: Dev Environment Setup
    url: /dev-environment-setup/
    weight: 5
  - identifier: getting-started
    name: Getting Started
    url: /getting-started/
    weight: 6
  - identifier: reference
    name: Reference
    url: /reference/
    weight: 7
  - identifier: updates
    name: Updates
    url: /updates/
    weight: 8
outputFormats:
  MenuIndex:
    baseName: menu
    mediaType: application/json
  SearchMap:
    baseName: searchmap
    mediaType: application/json
outputs:
  home:
  - HTML
  - RSS
  - MenuIndex
  - SearchMap
module:
  imports:
  - path: github.com/spandigital/presidium-styling-base
  - path: github.com/spandigital/presidium-layouts-base
enableInlineShortcodes: true
```
