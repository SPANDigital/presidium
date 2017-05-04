---
title: "Article Statuses"
---

Each article may be assigned a status to track its lifecycle:

- draft
- review
- published
- retired

Statuses can be set in the site's `front matter`:
```
---
status: draft|review|published|retired
---
```

To enable or disable statuses from showing on your generated site, 
simply enable or disable the following setting in the site `config`:

```
show-status: true|false
```