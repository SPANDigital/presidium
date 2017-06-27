---
title: Images
---

Images may be placed with your content or Place images you wish to include in the `/media/images` directory.

The image path may be fully resolved or templated using a custom property defined in your site config:
 
{% raw %}
```markdown 
![Sample Image](/media/images/logo.png)
![Sample Image]({{site.baseurl}}/media/images/logo.png)
![Sample Image]({{site.images}}/logo.png)
```
{% endraw %}

![Sample Image](/media/images/logo.png)
![Sample Image]({{site.baseurl}}/media/images/logo.png)
![Sample Image]({{site.images}}/logo.png)

To simplify links, you can add the base path to your images in your site config:
site:
images: ${site.baseurl}//media/images
{% raw %}
```markdown 
![Sample Image]({{site.images}}/logo.png)
```
{% endraw %}

![Sample Image]({{site.images}}/logo.png)

To include a caption, follow an image link with a`*Caption*`
```markdown 
![Sample Image With Caption]({{site.baseurl}}/media/images/logo.png)
*Sample Image With Caption*
```
![Sample Image With Caption]({{site.baseurl}}/media/images/logo.png)
*Sample Image With Caption*