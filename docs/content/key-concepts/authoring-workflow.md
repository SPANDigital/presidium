---
title: Authoring Workflow
weight: "4"
---

The following workflow is recommended for authoring content using Presidium.

# Identify the Author

Whether you're creating a single article or an entire section, make sure to include the author tag in the front matter.

# Work Locally

Presidium encourages you to write and review documentation on your local machine, rather than making changes directly on Github. The advantages to this workflow are:

* Changes are immediately viewable
* You can leverage Presidium's validation tools
* Git commit history will have less noise

# Protect Master

If you own the documentation project, make sure you prevent others from making commits or pushing to the Master branch. You will be publishing from the Master branch, so you want to maintain its integrity.

# Work On a Branch

If you're adding new content, start out on a branch based off of Develop.

# Review

Changes should be submitted as Pull Requests. Review all changes and merge into your Develop branch, then tag the release as a release for verification.

# Tag, Release & Publish

After you've finalized the new content, merge it into Master, tag it as a release, then publish to Github pages.
