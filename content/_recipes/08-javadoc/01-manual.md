---
title: Manually Embed Generated Documentation
---

One way in which to pull in auto generated documentation (like Javadoc) is to embed it within an iframe.  We'll use
Javadoc as an example.


# Generation

First generate your site documentation by running Javadoc:

`javadoc -sourcepath src -d docs com.spandigital.presidium -notree -noindex -nohelp -nonavbar`

Note that, of the options given, the most important is the `nonavbar` setting which excludes all the top header clutter.

# Serve

In order for the documentation to be statically served by Presidium, place the documentation
somewhere within the `/media` folder, the Presidium convention is to place it under `/media/import`:

`/media/import/javadoc/my-project/`

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
            src='{{site.baseurl}}/media/import/javadoc/my-project/com/spandigital/presidium/package-summary.html'
    </iframe>
</div>
```

You are also able to *frame* content either at the package level (as above) or class level - here, you can create 
multiple markdown files for each class, and structure them appropriately in the left-hand side navigation.

# Styling

The default styles that you provide to your javadoc will be used when rendering the content in the iframe.

You are able to override those styles by passing the presidium stylesheet as an argument to the javadoc command:

`javadoc ... -stylesheetfile ./dist/site/media/css/presidium.css`

The Presidium css file is available in ./dist after a build or serve.