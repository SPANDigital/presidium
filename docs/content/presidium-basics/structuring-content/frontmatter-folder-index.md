---
title: Frontmatter
weight: 40
---
Frontmatter is a required element in Presidium articles that 
* Indicates that the file should be included as an article during the build process
* Sets various properties for the article

Here is an example:
```
---
title: Presidium Overview
slug: overview
url: front-matter/overview
weight: 10
author: Kim
roles: Developer
status: Draft
---
```
The only required elements are the enclosing lines of three hyphens and `title`.

In the above example,
* title: Must be enclosed in double quotes if it contains special characters such as colon or parentheses.
* slug: Overwrites the default, which is the last segment of the URL.
* url: Overwrites the whole default URL.
* weight: Sets the order of this article relative to others in its section.
* author, roles, status: Self-explanatory; you can define these as you wish.

For a full list of possible frontmatter fields, see https://gohugo.io/content-management/front-matter/.