---
title: OpenAPI3
weight: 2
---

Presidium includes a Golang tool ([presidium-oapi3](https://formulae.brew.sh/formula/presidium-oapi3)) for importing your OpenAPI 3 spec into Presidium documentation.

1. Install the [presidium-oapi3](https://github.com/SPANDigital/presidium-oapi3) tool using Homebrew:

   ```sh
   brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git
   brew install presidium-oapi3
```

Example:

```sh
$ presidium-oapi3 convert -f <YOUR_API_SPEC> -o <THE_OUTPUT_DIRECTORY> -r <THE_PRESIDIUM_REFERENCE_URL>
```

The following options are available for `presidium-oapi3`:

| Option                       | Description                                             |
|:-----------------------------|:--------------------------------------------------------|
| -n, \--apiName `string`      | The name under which the generated docs will be grouped |
| -f, \--file `string`         | OpenAPI 3 spec file                                     |
| -o, \--outputDir `string`    | The output directory                                    |
| -r, \--referenceURL `string` | The reference URL (default "reference")                 |
| -h, \--help                  | help for convert                                        |