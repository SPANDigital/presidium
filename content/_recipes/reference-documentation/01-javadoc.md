---
title: Javadoc
---

Java source code comments may be imported to Presidium using the Presidium Doclet implementation. 
Imported documentation will be included in the menu and sitemap.

To import javadoc, you will need to use the [presidium-javadoc](https://www.npmjs.com/package/presidium-javadoc) package.

1. Add the [presidium-javadoc](https://www.npmjs.com/package/presidium-javadoc) dependency to your site's `package.json`.
1. Add a generation script that parses the provided `<src-path>` directory and `<packages>` and generates markdown in your `content/_reference` section.
1. Run `npm run import-javadoc-api` whenever you need to update your source documentation

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
| -d,--directory `path`                     | The destination directory to save the generated documentation to. default: './docs'
| -h,--help                                 | Shows this help.
| -p,--subpackages `package1:package2:...`  | Packages to generate documentation from. default: all
| -s,--sourcepath `path`                    | Java source code directory.
| -t,--title `string`                        | Reference title. default: 'javadoc'
| -u,--url `foo/bar/{title-slug}`            | Section url. default: 'reference/javadoc'
| -d,--directory `path`                      | The destination directory to save the generated documentation to. default: './docs'
| -h,--help                                  | Shows this help.
| -p,--subpackages `package1:package2:...`   | Packages to generate documentation from. default: all
| -s,--sourcepath `path`                     | Java source code directory.
| -t,--title `string`                        | Reference title. default: 'javadoc'
| -u,--url `foo/bar/{title-slug}`            | Section url. default: 'reference/javadoc'

Alternatively you can use the Doclet within an existing build workflow such as `gradle` using the [javadoc-plugin](https://docs.gradle.org/current/dsl/org.gradle.api.tasks.javadoc.Javadoc.html).

