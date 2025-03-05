---
title: Statuses
weight: 4
---

Each article can be assigned a status to track its lifecycle:

- draft
- review
- published
- retired

Statuses can be set in the site’s `frontmatter`:

```
---
status: draft|review|published|retired
---
```

To show or hide statuses on your generated site, use the following setting in the `site config`:

```
params:
  show:
    status: true|false
```
