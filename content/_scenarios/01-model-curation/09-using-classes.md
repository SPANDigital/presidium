---
title: "Using Classes"
category: "model-curation"
id: "using-classes"
author: "Virtual Traveler"
---
ADSS supports the defintion of classes. In general it is best to avoid classes. They are implementation details tht tend to change rapidly and are best encapsulated in components. However there are a few cases when classes can be useful.

## Discussion 

**Data Classes** are equivalent to individual entities or files. As such it is sometimes worth identifying where they are used throughout the system. My creating containment relationships between data componenets and the classes they contain it is possible to identify all usage of a particular data class throughout the system. 

**Messaging Classes** are equivalent to specific message formats. Similar to data classes, above, it is possible to identify when containers create or consume massages of various formats.  

## See Also

* [REPLACE WITH a link description](http://www.google.com) 
