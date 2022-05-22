---
title: Presidium's Learning Model
weight: 1
---

Presidium is more than a framework for producing documentation. It asks you to think about, structure, author, and
publish content based on sensible best practices so you can write effective documentation, rather than producing 
documentation for documentation’s sake.

## Supporting Learning Objectives & Efficiency

Presidium recommends creating small micro-articles that convey enough information to fulfill a learning objective. 
Small articles are advantageous because they:

- Usually only require a single author
- Force the author to focus on effectively describe a single concept
- Enable authors to quickly add documentation while also working on a product or service.

Using this method, creating high quality content is part of the engineering process, not a deferred or protracted task.

## Categorization

Often, either only a few people hold all the critical business and technical knowledge or information is spread across 
documents without organization or cohesion. Presidium suggests a way to categorize information that helps the writer 
compartmentalize and categorize knowledge, leading to more cohesive and logical documentation.

### Overview & Key Concepts

Almost any knowledge base requires an overview to orient the reader in the correct domain. The Overview section should 
support and lead to a discussion of key concepts that **prime** the reader to use the rest of the content. 
These two sections are naturally presented first, and one should follow the another.

### Prerequisites & Getting Started

For software projects, APIs, and similar documentation, most readers will require some initial preparation. 
This may include installing libraries or binaries and reading information before proceeding.

The Prerequisites section should support the information in Getting Started. The Getting Started section should be a 
quick deep-dive than enables a user to get from nowhere to running in minutes.

### Best Practices

The Best Practices section should contain the ideal path to follow to get the most out of your product or service and 
reduce the potential of encountering problems. They can be a combination of anecdotal knowledge, RFCs, white papers etc.

### Reference

References should enhance or support the information in other sections (for example, imported API documentation, 
deep-dives beyond the needs of most readers, etc.)

### Recipes

Recipes are clear sets of steps that explain how to do something useful or unusual with your product or service.

### Glossary

Glossary entries are bite-sized pieces of information that explain a concept, resolve jargon, or list synonyms for a 
particular term. Because all content in the Glossary can be used for automatic tooltip generation, a good rule of thumb 
is to make sure a glossary entry can be condensed into a tooltip without being verbose or introducing confusion.

### More

Tools, Uses Cases, Updates, and Support are other sections you can use. You can also define your own sections in Presidium.
