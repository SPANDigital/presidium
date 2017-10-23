---
title: Authoring Workflow
---

The suggested workflow for authoring content using Presidium is as follows:

# Take Ownership

Whether you're creating a single article or an entire section, make sure you include the author tag in the front matter.

# Work Locally

Presidium encourages you to write documentation and review it locally on your machine. There are many benefits to this as opposed to making changes directly on Github:

* Changes are immediately viewable.
* You can leverage Presidium's validation tools.
* Git commit history has less noise.

# Protect Master

If you own the documentation project, ensure that you prevent other members from making commits, or pushing to, the Master branch.
This is the branch you should be publishing from so you'll want to maintain its integrity.

# Work On a Branch

If you're adding new content, start out on a branch based off of develop.

# Review

Changes should be submitted as Pull Requests. Review all changes and merge into your develop branch, tag the release as a release for verification.

# Tag, Release & Publish

Once you're happy with the new content, merge it into Master, tag it as a release, and then publish to github pages.