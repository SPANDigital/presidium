---
title: Themes
weight: 2
status: OUTDATED
---

Presidium provides some themes when spinning up your Presidium Site through the command line. You will see something
like this when running `presidium init`:

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

The default Presidium logo image is placed and loaded from the pre-configured template. 

To update the site's logo place the logo image into the static/images folder and update the site's ```config ``` logo variable as seen below
```logo: [path to logo]```

Ideally the logo should have the following specifications `260px × 124px` and can be of type: png, jpg or svg

## Favicon

The default Presidium favicon image is placed and 
loaded from the pre-configured template. 

To update the site's favicon place the favicon image into the static/images folder and update the site's ```config ``` favIcon variable as seen below
```favIcon: [path to favicon]```

Ideally the favIcon can be of type: png, jpg or svg, ico.

## Title tag

The browser `<title>` tag is populated with the `title` provided when spinning up your Presidium Site throught the
command line. You can modify it by changing the `title` key in the configuration file.
