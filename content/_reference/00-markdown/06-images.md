---
title: Images
---

Place any images you wish to include in `<project root>/media/images` directory and reference them in the text. Note the exclamation point.

```
![Image Name]({{site.baseurl}}/media/images/image.png)
```

## Presidium Image Path

If you've defined a helper variable in `_config.yml` such as:

```yml
images: ${site.baseurl}/media/images
```

you can instead write:

```
![Image Name]({{site.images}}/image.png)
```