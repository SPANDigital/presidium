---
title: Overview
weight: 1
status: OUTDATED
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
authors: someone@example.com
status: Draft
draft: true
---
```
  
Presidium uses the following standard Hugo front matter variables:

| **Variable** | **Description** |
|--------------|-----------------|
| title | A string representing the article or subsection title or the title of the subsection. |
| slug | A string representing the slug for deep linking the article |
| url | A string representing the URL for the article |
| weight | A number to provide ordering of articles, a higher number means the article will appear later in the section |
| [author](#authors) | A string, must be a valid email address, normally of the main author of the article. |

Presidium includes the following custom front matter variables:

| **Variable** | **Description** |
|--------------|-----------------|
| [status](#statuses) | A string that indicates the status of the article (draft, complete, in progress, etc.) |
