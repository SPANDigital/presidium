---
title: Statuses
weight: 4
status: GOOD
---

Each article can be assigned a status to track its lifecycle:

- draft
- review
- published
- retired

Statuses can be set in the site’s `front matter`:

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

Please note that while this `status` front matter field appears similar to the Hugo `draft` front matter field, it performs a different function.
The Hugo `draft: true` frontmatter field switches an article to only `render` if a special draft build flag is used, otherwise an article with this front matter field will be hidden from the `rendered` site by default.
Then if the `show status:true` param is set in the config file then all `status` fields that are populated in every article's front matter, will always be visible on any `rendered` article.
This performs the function of indicating the maturity status of the article to the reader for any rendered article.
