---
title: Themes
role: Developer
---

Presidium uses [http://bootswatch.com/](http://bootswatch.com/) for managing its themes.  If you wish to change the 
current theme of your documentation, site, simply navigate to `/media/css/_sass` and edit the `_variables.scss` file:

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
Presidium gives you the supported themes listed above. Uncomment the theme that you want (and comment the existing one).
Then, if you called Presidium with `npm start`, it will pick up the change and (the hotloader) will update the 
styles allowing you to refresh the browser window and view the new theme.

Note that if you want a *pure* `spacelab` theme for example, you must remove the overrides as displayed above.

## Logo

The default logo image is placed and loaded from `/media/images/logo.png`.  If you wish to update it, simply override 
the existing file in the folder.  Its size ratio should be:

`260px × 124px`