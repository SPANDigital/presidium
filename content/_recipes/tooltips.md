---
title: Tooltips
---

Tooltips provide a quick definition for an item. There are two ways of creating tooltip: automatically from the glossary, and via a link override.

# Automatic Tooltips

Automatic tooltips reference glossary entries. If a glossary article by the name of "Tooltips" exists, then a tooltip will be available for the following item:

[Tooltips](# 'presidium-tooltip')

```md
[Tooltips](# 'presidium-tooltip')
```

# Link Override

You may also supply an internal article as a source for a tooltip. Presidium will use the article's first paragraph to construct the tooltip. You are required to ensure, however, that the first paragraph is semantically sufficient to be used as a tooltip. Note that the text used for the demarcation of a tooltip need not match the article title.

[Article Templates]({{site.baseurl}}/best-practices/#use-article-templates 'presidium-tooltip')