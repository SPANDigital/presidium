---
title: Article Statuses
weight: 30
author: helgard.meyer@spandigital.com
status: published
---

Each article can be assigned a status to track its lifecycle:

- draft
- review
- published
- retired

This displays one of the following article status indicators next to the author's name:
<html>
<div class="article-status">  
<ul>  
<li><span title="Article Status" class="label label-success status-draft">Draft</span></li> We’re still working on this article, so feedback is premature.
<li><span title="Article Status" class="label label-success status-review">Review</span></li> This article is ready for review. Please add your comments and feedback.    
<li><span title="Article Status" class="label label-success status-published">Published</span></li> This article has been published and should be correct, but please provide feedback, if you have any.
<li><span title="Article Status" class="label label-success status-retired">Deprecated</span></li> This article has been deprecated—it’s either no longer relevant or it has been replaced by more up-to-date information.
</ul>
</div>
</html>
  

To show or hide statuses on your generated site, use the following setting in the `site config`:

```
params:
  show:
    status: true|false
```


Note that while this `status` frontmatter field appears similar to the Hugo `draft` frontmatter field, its behavior is different:

* Whether the Presidium `status` value appears is determined only by the `config.yml` parameter.

* The situation with the Hugo `draft: true` frontmatter field is more complex: It switches an article to `render` only if a special draft build flag is used, otherwise an article with this frontmatter field is hidden from the `rendered` site. Then if the `show status:true` param is set in the config file then all `status` fields that are populated in any article's frontmatter, are always visible on any `rendered` article. This indicates the maturity status of the article to the reader for any rendered article.