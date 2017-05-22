---
title: Embed Generated Documentation Manually
---

One way in which to pull in auto generated documentation (Like Javadoc) is to embed it within an iframe.  We'll use
Javadoc as an example.


# Generation

First generate your site documentation by running Javadoc:

`javadoc -sourcepath src -d docs com.spandigital.presidium -notree -noindex -nohelp -nonavbar`

Note the options given, the most important is the `nonavbar` setting which excludes all the top header clutter.

# Serve

In order for the documentation to be statically served by Presidium, place the documentation
somewhere within the `/media` folder:

`/media/javadoc/my-project/`

# Reference

In order to load the documentation, create a folder within the reference section and create a
package / class file:

`01-com-spandigital-presidium.html`:

```markdown

---
title: My Java Package
---

# com.spandigital.presidium

<div>
    <iframe
            src='{{site.baseurl}}/media/my-project/com/spandigital/presidium/package-summary.html'
            sandbox="allow-top-navigation allow-scripts allow-same-origin">
    </iframe>
</div>
```

# Styling

You can override the Javadoc styling to include presidium.css, this will align the styling to
the Presidium side, ensuring that the look and feel remains constant as we've included Javadoc
classes in our css.
