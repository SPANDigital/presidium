---
title: Images
weight: 7
status: GOOD
---

## Image Links

Put any images you want to include in the <project root>/static/images directory and reference them in the text. 
Note the exclamation point. The image path may be fully resolved:

![Image Name](/images/logo.png)

```md
![Image Name](/images/logo.png)
```

## Image with attributes

{{< img src="/images/logo.png" caption="Sample image" style="width:25%;" >}}

```md
{{</* img src="/images/logo.png" caption="Sample image" style="width:25%;" */>}}
```

### Captions

To include a caption, add *Caption* after an image link. For example:

![Sample Image With Caption](/images/logo.png)
*Sample Image With Caption*

```md
![Sample Image With Caption](/images/logo.png)
*Sample Image With Caption*
```
