---
title: Themes
---

Presidium uses [http://bootswatch.com/](http://bootswatch.com/) to manage themes. To change the
theme, navigate to `/media/css/_sass` and edit the `_variables.scss` file:

```scss
// Override any bootstrap or bootswatch variables as you need:
$brand-info: #e49134;
$navbar-default-link-hover-color: #e49033;
$navbar-default-link-hover-bg: white;

// Available Bootswatch Themes
@import 'themes/spacelab';
//@import 'themes/cosmo';
//@import 'themes/darkly';
//@import 'themes/simplex';
//@import 'themes/united';
```
Presidium includes the themes listed above. Uncomment the theme that you want and comment the selected one.
When you call Presidium with `npm start`, it will pick up the change and (the hotloader) will update the 
styles, allowing you view the new theme after refreshing the browser window.

Note that if you want a pure `spacelab` theme for example, you must remove the overrides as shown above (`$brand-info` ... etc).

## Logo

The default logo image is placed and loaded from `/media/images/logo.png`.  To update it, replace
the existing file in the folder. Use the following size ratio:

`260px × 124px`

Additional styling can be added to the `_custom.scss` file to change the position or size of the logo.

## Favicon

A default icon does not exist, but can be added to `/media/images/favicon.ico`. Different favicon sizes are not supported yet.

## Title tag

The browser `<title>` tag is populated with the `name` key found in the configuration file.
