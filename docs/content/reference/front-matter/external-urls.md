---
title: External URLs
weight: "6"
---

Sections can redirect to an external URL either in the same tab or a new one.


This can be done by setting the `external-url` in the `site config` like so:

```
sections:
  - title: External Link
    collapsed: true
    external-url:
      href: "http://www.anothersite.com"
      new-tab: true
    collection: testing
    
  - title: Another External Link
    collapsed: true
    external-url:
      href: "http://www.anothersite.com"
      new-tab: false
    collection: testing
```

The `new-tab` attribute is optional and defaults to true