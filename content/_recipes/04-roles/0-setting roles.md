---
title: "Using Roles"
---

Various roles can be added to a site to allow readers to 
filter by a target audience. This is an optional feature that can be 
enabled by setting roles in your site config:

```
roles:
  label: "Show documentation for role"
  all: "All Roles"
  options:
    - "Business Analyst"
    - "Developer"
    - "Tester"
```

Articles and menu items are shown depending on the selected role. 
If no role is provided, an article will default to `all`.



