---
title: Images
---

# Image Links

## Fully Resolved

Put any images you want to include in the `<project root>/static/images` directory and reference them in the text. Note the exclamation point.
The image path may be fully resolved:

{% raw %}
```md
![Image Name](/images/image.png)
```
{% endraw %}

## Baseurl

**Or** templated using the `baseurl` shortcode:

{% raw %}
```md
![Image Name]({{site.baseurl}}/images/image.png)
```
{% endraw %}

## Path

**Or** templated using the `path` shortcode if the image is in the same directory as the article.

{% raw %}
```md
![Image Name]({{%path%}}/image.png)
```
{% endraw %}

## Captions

To include a caption use the img shortcode. For example:
{{< img src="photo.jpg" caption="My Name is Jeff">}}

