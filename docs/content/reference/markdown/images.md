---
title: Images
weight: "7"
---

# Image Links

## Fully Resolved

Put any images you want to include in the `<project root>/static/images` directory and reference them in the text. Note the exclamation point.
The image path may be fully resolved:


```md
![Image Name]({{% baseurl %}}/images/image.png)
```


## Baseurl

**Or** templated using the `baseurl` shortcode:


```md
![Image Name]({{% baseurl %}}/images/image.png)
```


## Path

**Or** templated using the `path` shortcode if the image is in the same directory as the article.


```md
![Image Name]({{% baseurl %}}/image.png)
```


## Captions

To include a caption use the img shortcode. For example:
{{< img src="photo.jpg" caption="My Name is Jeff">}}

