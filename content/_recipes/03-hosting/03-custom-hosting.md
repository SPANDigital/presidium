---
title: Custom Hosting
---

When you run `npm start` or `npm run build` a static website is built in the `dist/site` directory. This may be copied over and served by any web server capable of hosting static content.

A common use case is to build your Presidium docs by running `npm run build` on your continuous integration server and copying the generated site to a hosting web server of your choice. This may include, Apache HTTP, Nginx or an AWS S3 bucket where private and public access may be controlled.