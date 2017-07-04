---
title: Overview
---

Front Matter serves two purposes, it denotes that the file should be included as an article in the build process, and, it alows you to set properties on your article. For example:

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

Front matter variables available for use in Presidium are listed below:

| variable                                       | description |
|------------------------------------------------|-------------|
|title | A string representing the article title, or [the title of the subsection]({{ site.baseurl }}/recipes/content-structure/#directory-structure). |
|[author]({{ site.baseurl }}/reference/front-matter/#authors)   | A string, generally the Github username of the main author of the article. |
|[roles]({{ site.baseurl }}/reference/front-matter/#user-roles) | A list of strings representing the roles best associated with an article. |
|[status]({{ site.baseurl }}/reference/front-matter/#statuses)  | A string indicating whether the article is complete, in progress or needing review etc. |