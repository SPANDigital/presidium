---
title: Landing Page
---

# Configure

If you want your documentation to have a landing page, edit your `config.yml`, so that no section has `/` as its url.

Change this:

```yml
   identifier: Introduction
    name: Introduction
    url: /
    weight: 1
```

to this:

```yml
   -identifier: Introduction
    name: Introduction
    url: /introduction/
    weight: 1
```

# Create a Landing Page

In the root of `./content`, add a file called `index.md`. The example below is an excerpt from Presidium's landing page. Note the front matter at the top which is required for paths to be rendered correctly.

{% raw %}
```html
---
---

<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    .
    .
    .
```
{% endraw %}
