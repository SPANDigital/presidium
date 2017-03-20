---
title: "With an Author"
category: Authors
author: pacomendes
---

Each article may have an author in the `front matter`:
```
---
author: {author name}
---
```

Links to authors are optional. To enable, set a base url in your site config:

```
external:
  authors-url: https://github.com/orgs/SPANDigital/people/
```

To hide or show authors on your generated site, 
simply enable or disable the component in your site config:

```
show:
  author: true|false
```