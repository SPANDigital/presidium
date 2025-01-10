---
title: Tooltips
weight: 8
status: GOOD
---

Tooltips display a short definition of an item. There are two ways to create tooltips:

- Automatic from the Glossary
- Via Link Override

## Automatic Tooltips

Automatic tooltips reference Glossary entries. If a Glossary article by the name of “Tooltips” exists, a tooltip will 
be available for the following item:

{{< tooltip "Tooltips" >}}

```md
{{</* tooltip "Tooltips" */>}}`
```

## Link Override

You can use an internal article as the source of a tooltip. Presidium will use the first 100 words of the article to 
construct the tooltip, so you should make sure the text will work as a tooltip. Note that the text used for the 
demarcation of a tooltip does not need to match the article title, like 
{{< tooltip "this," "best-practices/plan-content-development.md" >}} which links to an article on plan content 
development.

```md
{{</* tooltip "this," "best-practices/plan-content-development.md" */>}}
```
