---
title: General Performance Guidelines
weight: 7
---

For Presidium to support the functionality it has built on top of Hugo, there are some performance considerations to take into account so that build times can remain reasonable. Below are some steps that can be taken to help with overall build performance:

- Keep total markdown content to 50MB or less. This should result in a build time of around one minute, although this cannot be guaranteed due to other factors.

- Content spread across multiple files is better than having fewer files with more content.

- Less content per markdown file is encouraged. This is also in line with Presidium's philosophy of micro-articles.
