---
title: Overview
slug: overview
url: front-matter/overview
weight: 1
---

Front Matter serves two purposes:

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

Presidium includes the following front matter variables:

| **Variable** | **Description** |
|--------------|-----------------|
| title | A string representing the article or subsection title or the title of the subsection. [TODO] // Add a link to the directory content structure article |
| slug | A string representing the slug for deep linking the article |
| url | A string representing the URL for the article |
| weight | A number to provide ordering of articles, a higher number means the article will appear later in the section |
| author | A string, generally the Github username of the main author of the article. |
| roles | A list of strings representing the appropriate roles for an article. |
| status | A string that indicates the status of the article (draft, complete, in progress, etc.) |

[TODO] // Add appropriate link for authors and roles
