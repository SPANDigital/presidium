---
title: "Authors"
---

Each article may include an author in the `front matter`:
```
---
author: {author name}
---
```

Links to authors are optional and may be enabled by setting a base `authors-url` in the site config:

```
external:
  authors-url: https://github.com/orgs/SPANDigital/people/
```

To hide or show authors on your generated site, 
simply enable or disable the component in the `site config`:

```
show:
  author: true|false
```