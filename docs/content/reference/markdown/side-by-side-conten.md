---
title: Side By Side Content
weight: 10
---

A shortcode is available (`{{%/* sidebyside */%}}`) in order to have content, within a single file, in two different columns that are side by side.

Create side-by-side content by using the following:

```
{{% sidebyside %}}

Content you want on the left hand side of the site

--split-content--

Content you want on the right hand side of the site

{{% /sidebyside %}}
```

> **Note**: `--split-content--` is where the content is to be split between left and right. Without this keyword the site will fail to render the two columns correctly
