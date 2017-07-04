---
title: "User Roles"
---

Various roles can be added to a site to allow readers to filter by a target audience. 
Articles and menu items are filtered based on the selected role. 

This is an optional feature that can be enabled by defining user roles in the `site config`

```
roles:
  label: "Show documentation for role"
  all: "All Roles"
  options:
    - "Business Analyst"
    - "Developer"
    - "Tester"
    - "..."
```


If a role is not specified, an article defaults to `roles.all`.

Articles may have one or more roles defined in article `front matter`:

```
---
roles: [Developer, Business Analyst]
---
```

```
---
roles: [Developer]
---
```

```
---
roles: Business Analyst
---
```