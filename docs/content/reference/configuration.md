---
title: Configuration
weight: "1"
---

This is a sample Presidium configuration file:

```yml
languageCode: "en-us"
title: "{{ .Title }}"
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
  main:
    - identifier: overview
      name: Overview
      url: /
      weight: 1
    - identifier: key-concepts
      name: Key Concepts
      url: /key-concepts/
      weight: 2
    - identifier: prerequisites
      name: Prerequisites
      url: /prerequisites/
      weight: 3
    - identifier: getting-started
      name: Getting Started
      url: /getting-started/
      weight: 4
    - identifier: best-practices
      name: Best Practices
      url: /best-practices/
      weight: 5
    - identifier: reference
      name: Reference
      url: /reference/
      weight: 6
    - identifier: glossary
      name: Glossary
      url: /glossary/
      weight: 7
    - identifier: recipes
      name: Recipes
      url: /recipes/
      weight: 8
    - identifier: tools
      name: Tools
      url: /tools/
      weight: 9
    - identifier: updates
      name: Updates
      url: /updates/
      weight: 10
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
    {{ if .Brand -}}
    - path: {{ .Brand }}
      mounts:
        - source: assets
          target: assets
        - source: static
          target: static
    {{- end }}
    - path: {{ .Theme }}
enableInlineShortcodes: true
frontmatter:
  lastmod:
    - lastmod
    - :fileModTime
    - :default
```
