---
title: External URLs
slug: external-urls
url: front-matter/external-urls
weight: 6
---

Sections can redirect to an external URL either in the same tab or a new one.

This can be done by setting the `externalUrl` in the `site config` like so:

```yaml
menu:
  main:
    - identifier: external-link
      name: External Link
      url: /
      weight: 1
      params:
        externalUrl:
          href: "https://www.anothersite.com"
          newTab: true
    - identifier: another-external-link
      name: Another External Link
      url: /
      weight: 2
      params:
        externalUrl:
          href: "https://www.anothersite.com"
          newTab: false 
```

The `newTab` attribute is optional and defaults to true
