---
title: Links
weight: 4
---

You can link to internal articles in your repository, to external articles, or even to other semantically significant 
text. Any text enclosed in angle brackets will be interpreted as a link. If you want to add a description, 
use square brackets for the description and parenthesis for the link.

- Internal Page link: [Presidium Authors]({{< ref "reference/front-matter/authors.md" >}})
- Internal Anchor link: [Presidium Authors]({{< ref "reference/front-matter/#authors" >}})
- External link: <https://github.com/SPANDigital/presidium>
- Alternative: [Presidium on Github](https://github.com/SPANDigital/presidium)


```md
- Internal Page link: [Presidium Authors]({{</* ref "reference/front-matter/authors.md" */>}})
- Internal Anchor link: [Presidium Authors]({{</* ref "reference/front-matter/#authors" */>}})
- External link: <https://github.com/SPANDigital/presidium>
- Alternative: [Presidium on Github](https://github.com/SPANDigital/presidium)
```
