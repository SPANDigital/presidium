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

For a full list of possible frontmatter fields, see https://gohugo.io/content-management/front-matter/. Following are details on the variables in the above example:
* `title`: String; must be enclosed in double quotes if it contains special characters such as colon or parentheses.
* `slug`: String; the slug for deep linking the article. Note the slug only updates the leaf node, and does not update the section slug. This overwrites the default, which is the last segment of the URL.
* `url`: String; the URL for the article. Note that this is an absolute URL, and not a relative URL. Overwrites the default URL.
* `weight`: Integer; sets the order of this article relative to others in its section. Higher numbers are ordered later.
* `author`: String; generally the email address of the main author of the article. 
* `status`: String; indicates the status of the article. Possible values are draft, review, published, and deprecated. These article status indicators display next to the author's name:
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