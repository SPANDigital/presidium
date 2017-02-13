---
title: "Data Sources"
category: "Data Sources"
id: "key-concepts" 
author: "VirtualTraveler"
---
Models are composed from data obtained from multiple sources. Data is automatically ingested fro multiple sources and then categorized and loaded to the model.  

- **Directory Services** From on premise directory services systems like LDAP, or similar cloud based services, ADSS gathers information on the enterprise: the groups and divisions that comprise the **Organization** as well as the **People** - developers and managers and their reporting structure. This information is used to track ownership and collaboration within each model. 

- **Development Systems** From source code control systems like github and SVN, continuous integration systems like Jenkins and Travis CI, artifact management systems like Artifactory, and ticket mangement systems like Jira, ADSS gathers information on: the tasks and deliverables of development teams - **Checkins** made to **Source Code** in response to **Work Tickets**, and the **Build** process that creates, **Artifacts** and **Releases** that are **Deployed** to production systems. This information is used to model the logical structure of, and dependencies between, modeled objects. 

- **Configuration Management Systems** From on premis configuration management systems or cloud platforms like AWS ADSS gathers information on **Regions**, **Zones**, **Devices**, **Container Instances**, and **Availability groups**. This information is used to model the physical deployment of the modeled system. 

Container configuration

- **Operational Monitoring Systems** From monitoring systems like App Dynamics and Splunk, ADSS gathers information on calls between modeld objects. This information is mapped and aggregated to produce summary information at the business processes level. 

- **System Designers and Managers** Finally ADSS allows system designers and managers to correct and enhance the automatically collected data allowing them to improve the fidleity of model and guide the system in how t classify the data it collects. The data curation process is ongoing but keeps the workload to the minimum required to maintain data quality. 