---
title: Overview
---

Front Matter serves two purposes:
* Indicates that the file should be included as an article in the build process
* Allows you to set properties for your article

For example:

{% raw %}
```md
---
title: Presidium
authors: github-username
roles: Developer
status: Draft
---
```
{% endraw %}

Presidium includes the following front matter variables:

| variable                                       | description |
|------------------------------------------------|-------------|
|title | A string representing the article or subsection title or [the title of the subsection]({{ site.baseurl }}/recipes/content-structure/#directory-structure). |
|[author]({{ site.baseurl }}/reference/front-matter/#authors)   | A string, generally the Github username of the main author of the article. |
|[roles]({{ site.baseurl }}/reference/front-matter/#user-roles) | A list of strings representing the appropriate roles for an article. |
|[status]({{ site.baseurl }}/reference/front-matter/#statuses)  | A string that indicates the status of the article (draft, complete, in progress, etc.) |
