---
title: Images
---

# Image Links

## Fully Resolved

Place any images you wish to include in `<project root>/media/images` directory and reference them in the text. Note the exclamation point.
The image path may be fully resolved:

{% raw %}
```md
![Image Name](/media/images/image.png)
```
{% endraw %}

## Baseurl

**Or** templated using the `baseurl` property defined in your site config:

{% raw %}
```md
![Image Name]({{site.baseurl}}/media/images/image.png)
```
{% endraw %}

## Custom Property

**Or** templated using a custom property defined in your site config:

```yml
images: ${site.baseurl}/media/images
```

{% raw %}
```md
![Image Name]({{site.images}}/image.png)
```
{% endraw %}

## Captions

To include a caption, follow an image link with a `*Caption*`, e.g.:

{% raw %}
```markdown 
![Sample Image With Caption]({{site.baseurl}}/media/images/logo.png)
*Sample Image With Caption*
```
{% endraw %}

![Sample Image With Caption]({{site.baseurl}}/media/images/logo.png)
*Sample Image With Caption*