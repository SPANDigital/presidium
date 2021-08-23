---
title: Themes
slug: themes
url: reference/themes
weight: 2
---

Presidium provides some themes when spinning up your Presidium Site through the command line. You will see something
like this when running `presidium-hugo init`:

```
Select a theme
  ✔ Presidium Theme - (Presidium's default theme)
```

If you already have a presidium site, you can enable your theme by adding a module in the configuration file:

```
module:
  imports:
    - path: <Your theme>
```

## Logo

The default logo image is placed and loaded from the pre-configured template. 
[TODO] // Please provide information on how to update the Logo for a Presidium Site
Use `260px × 124px` as the size ratio.

## Favicon

[TODO] // Please provide documentation on how to update a favicon

## Title tag

The browser `<title>` tag is populated with the `title` provided when spinning up your Presidium Site throught the
command line. You can modify it by changing the `title` key in the configuration file.
