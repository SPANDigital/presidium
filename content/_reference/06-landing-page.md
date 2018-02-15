---
title: Landing Page
---

# Configure

If you want your documentation to have a landing page, edit your `_config.yml`, so that no section has `/` as its url.

Change this:

```yml
  - title: "Overview"
    url: "/"
    collection: overview
    collapsed: true
```

to this:

```yml
  - title: "Overview"
    url: "/overview"
    collection: overview
    collapsed: true
```

# Create a Landing Page

In the root of `./content`, add a file called `index.html`. The example below is an excerpt from Presidium's landing page. Note the front matter at the top which is required for paths to be rendered correctly.

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

# Create & Resolve Links

Create links to content and resolve them by including the `baseurl` in the path:

{% raw %}
```html
    .
    .
    .
    <link rel="shortcut icon" href="{{site.baseurl}}/media/images/favicon.ico" type="image/X-icon">
    <title>Presidium</title>

    <link rel="stylesheet" href="{{site.baseurl}}/media/css/presidium.css" type="text/css" media="screen"/>
</head>
<body>
<div class="presidium-landing-container container-fluid">
    <nav class="presidium-landing-header">
        <div class="logo"><img src="{{site.baseurl}}/media/images/landing/presidium-header.png" /></div>
        <ul class="links">
            <li class="link"><a href="{{site.baseurl}}/overview"><img src="{{site.baseurl}}/media/images/landing/presidium-icon.png" /><span>DOCUMENTATION</span></a></li>
            <li class="link"><a href="http://www.github.com/SPANDigital/presidium"><img src="{{site.baseurl}}/media/images/landing/github-icon.png" /><span>ON GITHUB</span></a></li>
            <li class="link" style="display: none"><a href="#"><img src="{{site.baseurl}}/media/images/landing/slack-icon.png" /><span>ON SLACK</span></a></li>
        </ul>
    </nav>
    ...
    <div class="presidium-landing-footer">
        <div class="credit">Brought to you by </div><a href="https://www.spandigital.com"><img src="{{site.baseurl}}/media/images/landing/span-footer.png" title="SPAN Digital"/></a>
        <div class="copyright"><span>©2017 SPAN Digital LLC</span></div>
    </div>
</div>

</body>
</html>
```
{% endraw %}
