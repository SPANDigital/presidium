---
title: Establish Enterprise Standards
id: "Establish-Enterprise-Standards"
author: "VirtualTraveler"
---
You will get the most out of ADSS if you and your peers use it consistently. ADSS uses the C4 Architectural model as it's foundation. With only four object types; systems, containers, components and classes, and three relationship types; calls, contains and replaces ADSS Models are deliberately simple. Despite  this simple structure ADSS has been used to model even the most complex of systems.  

## Discussion

The following convensions have been found to work well.

**Systems** Agree on what you will use systems for. The following three  types of system have been found to work well 
- **External Systems** are  generally outside your control - You may know the name of the system and maybe the name of the interface you interact with but the internal details are not known, or not important, to your area of concern. External systems are "black boxes". 
- **Internal Systems** are not common. An internal system has a very well defined system boundary and is typically deployed many times across many systems, not just yours. Just because a container or group of containers is well defined does not make it a system. Ask yourself, if the candidate system is supported by a dedicated team? do they publish a service level agreement? If no, it's probably not a system. 
- **Implementation Technologies** are a very common way to use systems, and highly recommended. Implementation technologies are external systems and frameworks that are embedded in, or otherwise, used to build containers in your area of concern. Implementation technologies include; languages like C++, Java, C#, Python, datastores like Postgres and Redis, or messgaing systems like Kafka and RabbitMQ. Follow these [best practices for using implementation technologies](example.com) that you should follow. 

**Containers** Agree on the naming convension for Containers, be consistent. There are sometimes different names for the same container in different stages of the build process. Sometimes the same container is deployed multiple times with slightly different configuration and called by a different name. Agree on how you will handle these situations.

**Components**
- **Internal Libraries** and Frameworks
- **External Libraries** are included in the build of many Containers. Tracking these is often required to understand licenses complience. 

**Classes** Are not drawn on teh various views. It is best to stay away from classes until you have your model well established. If you are ready to use classes see [best practices for using classes](example.com).

## See Also

- [Visualise, document and explore your software architecture
Software Architecture for Developers - Volume 2](https://leanpub.com/visualising-software-architecture) 

