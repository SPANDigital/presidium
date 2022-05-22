---
title: Javadoc
weight: "3"
---

Presidium can import Java source code comments using the Presidium Doclet implementation.
Imported documentation will be included in the menu and sitemap.

To import javadoc, use the [presidium-javadoc](https://www.npmjs.com/package/presidium-javadoc) package.

1. Add the [presidium-javadoc](https://www.npmjs.com/package/presidium-javadoc) dependency to your site's `package.json`.
1. Add a generation script that parses the provided `<src-path>` directory and `<packages>` and generates Markdown in your `content/_reference` section.
1. Run `npm run import-javadoc-api` when you need to update your source documentation.

```json
{
    "scripts" : {
        "import-javadoc-api" : "presidium-javadoc -s <src-path>"
    },
    "devDependencies": {
        "presidium-core" : "#.#.#",
        "presidium-javadoc" : "#.#.#"
    }
}
```

The following options are available to `presidium-javadoc`:

| Option | Description
|:---|:---
| -d,--directory `path`                     | The destination directory in which to save the generated documentation. (Default is './docs'.)
| -h,--help                                 | Displays this help.
| -p,--subpackages `package1:package2:...`  | Packages to generate documentation from. (Default is all.)
| -s,--sourcepath `path`                    | Java source code directory.
| -t,--title `string`                        | Reference title. (Default is 'javadoc'.)
| -u,--url `foo/bar/{title-slug}`            | Section url. (Default is 'reference/javadoc'.)


Alternatively, you can use the Doclet within an existing build workflow such as `gradle` using the [javadoc-plugin](https://docs.gradle.org/current/dsl/org.gradle.api.tasks.javadoc.Javadoc.html).
