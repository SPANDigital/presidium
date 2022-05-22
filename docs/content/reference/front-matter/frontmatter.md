---
title: Overview
weight: "1"
---

Front Matter serves two purposes:
* Indicates that the file should be included as an article in the build process
* Allows you to set properties for your article

For example:


```md
---
title: Presidium
authors: github-username
roles: Developer
status: Draft
---
```


Presidium includes the following front matter variables:

| variable                                       | description |
|------------------------------------------------|-------------|
|title | A string representing the article or subsection title or [the title of the subsection]({{% baseurl %}}/recipes/content-structure/#directory-structure). |
|[author]({{% baseurl %}}/reference/front-matter/#authors)   | A string, generally the Github username of the main author of the article. |
|[roles]({{% baseurl %}}/reference/front-matter/#user-roles) | A list of strings representing the appropriate roles for an article. |
|[status]({{% baseurl %}}/reference/front-matter/#statuses)  | A string that indicates the status of the article (draft, complete, in progress, etc.) |
