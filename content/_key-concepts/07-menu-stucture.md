---
title: Menu Structure & Behavior
---

# Structure

For simplicity and transparency, the menu structure and associated directory structure are the same.

# Behavior

A significant feature of the menu, is that it has a tree structure: each **section or subsection is a node**, and all content is simply a traversal from this node through its children. A convention that Presidium applies, is that the node ‘name’ or title then becomes the main header of the page and everything outside of that sub-tree is not presented on the page. In this example, Content Structure has been clicked on, the content shown is everything that falls under that section:

![Content Structure]({{site.baseurl}}/media/images/content_structure.png)

Note the URL: `/recipes/content-structure/`. When clicking on a section/subsection title, the url will reflect that information in the path. If an article is clicked in a subsection, e.g. 'Article Concatenation', the URL becomes `/recipes/content-structure/#article-concatenation`.
