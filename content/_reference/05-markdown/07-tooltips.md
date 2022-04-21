---
title: Tooltips
---

Tooltips display a short definition of an item. There are two ways to create tooltips:
* Automatically from the Glossary
* Via a link override

# Automatic Tooltips

Automatic tooltips reference Glossary entries. If a Glossary article by the name of "Tooltips" exists, a tooltip will be available for the following item:

Tooltips:

```md
{{< tooltip "Policy" >}} 
```

# Link Override

You can use an internal article as the source of a tooltip. Presidium will use the first paragraph of the article to construct the tooltip, so you should make sure the text will work as a tooltip. Note that the text used for the demarcation of a tooltip does not need to match the article title, like [this,]({{site.baseurl}}/best-practices/#use-article-templates 'presidium-tooltip') which links to an article on templates.

{% raw %}
```md
{{< tooltip "Templates" "/best-practices/#templates” >}}
```
{% endraw %}
