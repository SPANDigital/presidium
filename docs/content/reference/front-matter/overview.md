---
title: Overview
weight: 1
---

frontmatter serves two purposes:

- Indicates that the file should be included as an article in the build process
- Allows you to set properties for your article

For example:

```md
---
title: Presidium
slug: overview
url: front-matter/overview
weight: 1
authors: github-username
roles: Developer
status: Draft
---
```

Presidium includes the following frontmatter variables:

| **Variable** | **Description** |
|--------------|-----------------|
| title | A string representing the article or subsection title or the title of the subsection. |
| slug | A string representing the slug for deep linking the article |
| url | A string representing the URL for the article |
| weight | A number to provide ordering of articles, a higher number means the article will appear later in the section |
| [author](#authors) | A string, generally the Github username of the main author of the article. |
| github | A string, representing the Github username of the main author of the article. |
| [roles](#user-roles) | A list of strings representing the appropriate roles for an article. |
| [status](#statuses) | A string that indicates the status of the article (draft, complete, in progress, etc.) |
