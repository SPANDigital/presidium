---
title: "Track Implementation Technologies"
id: "track-implementation-technologies"
author: "Virtual Traveler"
---
It is common to want to track all the containers that use a particular implementation technology. For example all the Data Containers implemented using Redis, or all the Message Containers using Kafka, or all the Functional Containers implemented in Java. This is helpful when doing security audits, license reviews, or impact analysis for technology upgrades. 

Modeling implementation Technologies can be achieved by creating a containment relationship between a container and the system it is implemented on 

## Discussion 

1. Create a System object in ADSS, name it after the implementation technology you wish to track (Eg. Oracle, Redis, Apache Solr, Kafka, etc.)  Do not assign a layer to the system - this will prevent it from being drawn on the diagram.  
2. [Edit the Container](/references/screens/object-management/) whose implementation technology you want to track. Create a contains relationship from the Container to the Sytem it is implemented on. The system must be contained by the container not the other way around. 

## See Also

* [REPLACE WITH a link description](http://www.google.com) 
