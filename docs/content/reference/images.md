---
title: Images
weight: 4
---

 ## Images 
Images used within the content should be placed in a static/images directory
```
/static/images/
    -- static/doc-workflow.png
```

## Logos and favicons

### Logos 

The Logo of the site appears above the Menu Navigation. 

Logos should be placed in the ```static/images``` directory.

**Note**: If a brand module is used the instruction above does not apply.  Ensure the brand module includes the static/images directory with the logo inside.

The menu bar Logo can be configured in the **config** file, under *params*: 
```
params:
    logo: [path to logo] either locally or from the brand module

Examples:

logo: images/logo.png  
logo: images/logo.svg
```

### Favicons
The favicon is the symbol that appears on the tab of the site, next to the title. 

Favicons should be placed in the ```static/images``` directory.

**Note**: If a brand module is used the instruction above does not apply. Ensure the brand module includes the static/images directory with the favicon inside.

The menu bar Logo can be configured in the **config** file, under *params*: 
```
params: 
    favicon: [path to favicon] (either locally or from the brand module)

Examples:

favicon: images/favicon.svg  
favicon: images/favicon.ico
```
