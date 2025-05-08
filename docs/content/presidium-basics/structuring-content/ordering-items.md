---
title: Ordering Items
weight: 20
---
Presidium has two ways of setting the order of folders and articles.
* Alphabetical order by name (filename or folder name). This is the default.
* According to the `weight` attribute. You set this:
  * For top-level sections: in the [config.yml]({{< ref "presidium-basics/structuring-content/#config.yml" >}}) file
  * For all other sections: in the `_index.md` of the folder
  * For articles: in the [frontmatter]({{< ref "presidium-basics/#front-matter" >}})

If `weight` is not specified, ordering is by filename. You can also override ordering by weight by setting the `SortByFilePath` attribute in the `params` section of config.yml to `true`.

> **Note:** The full ordering precedence in Hugo is weight, date, filename, and page title, in that order. So for example if you inadvertently give two articles the same weight, the newer one is ordered first.
