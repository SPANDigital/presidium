---
title: Content Versioning
url: content-versioning
slug: content-versioning
draft: true
---

The versioning tool allows the content creator keep up tp 5 versions of the content. These versions live under the `project/versions` folder.  The tool always discards the oldest version number if the user request more than 5 versions. By default this tool is not enabled and only has to be enabled once.

## Workflow

1. Enable to versioning feature ___(once only)___:

   ```sh
   presidium versioning --enable
   ```

2. Start the next version.

   ```sh
   presidium versioning next
   ```

3. Update the last version.

   ```sh
   presidium versioning update
   ```

To go back to any previous version, just copy the version back to the main `/content` folder


