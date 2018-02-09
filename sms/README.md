# SMS: SIF Memory Store

The core components of this module are:

### 2.2.1. SIF Store & Forward (SSF)

This is the core messaging infrastructure of NIAS. Based on the Apache Kafka message broker it simply ingests very large quantities of data and stores them for delivery to clients. The transport is utterly format agnostic and will support XML (SIF/EdFi), CSV (including IMS OneRoster), and JSON.

The defining feature of Kafka message queues is that they are persistent and ordered, meaning they can be re-read multiple times without consequence. For integration tasks, this means that both clients and servers can be asynchronous in their behavior; the inability for providers to behave in this way, or for one party in an integration to be forced to play the role of always-available-provider, is a current issue in market adoption, in particular for interfaces between solution providers.

Since the Kafka protocol is entirely message driven, services to work with its data can be written in any development language.
The SSF service that is part of NIAS simply builds an education-standards aware REST interface on top of Kafka, and also provides a number of utility services to ease SIF-based integrations.

### 2.2.2. SIF Memory Store (SMS)

The SMS is a database that builds its internal structures from the data it receives, rather than from the imposition of a schema.
When a SIF data object is provided to the SMS (typically via an SSF queue), the SMS creates a graph of all references to and from that object, and adds each object to a collection based on its type. As more data is added the number of collections can grow, but producers and consumers of data are not required to implement any more collections than those they wish to work with. In an integration scenario, if parties wish to exchange only invoice data and student personal data, then they can. 

Allowing the data to drive the structure lowers the effort for data providers considerably. The net result is that users no longer have to know the relationships between the different parts of the SIF model:, by providing data in the correct format the relationships will build themselves. The simplest possible input API for data, then, becomes achievable: data providers simply need to know how to represent the entities in their own systems as the appropriate SIF object. Entire schools-worth of data can be ingested in a single operation, with no API needed for the producer.

When it comes to retrieving data from the SMS, the query service exploits the graph of references to find any relationship between the requested items. Queries can all be expressed in the form of two parameters: the ID of an item in the database, and the name of a collection. For instance, providing the ID of a school as the item and *students* as the collection will find all students who have a relationship with the school. A teaching-group ID and the collection *attendance time lists* will return the attendance information for that teaching group.

One important consequence of the SMS traversing all relationships is that user queries no longer need to be aware of intermediate join objects, such as the StudentSchoolEnrollment objects that currently join students to schools. This radically simplifies the necessary understanding of the SIF data model when undertaking integration tasks. Users can focus on core entities such as students and class groups, without having to handle the wider complexities of the data model.

The combination of SSF and SMS achieves a highly simplified interface for integration based on SIF. In effect the inbound API is simply a stream of SIF messages, and the outbound API is a single query requiring only two parameters to fulfil any service path available in the provided data.

### 2.2.3. Non-SIF data

A significant side-effect of using non-SQL tools to construct the data store is that non-SIF data can also be easily accommodated for integration purposes. For instance IMS OneRoster data can be ingested, and if provided at ingestion time with a linking object identity, will be inserted into the relationship graph at that point. The students imported via OneRoster become equal citizens as far as querying all further relationships are concerned. Hence the rest of the data model is now linked to the OneRoster students. The SMS can now be queried, for example, to retrieve all invoices for a OneRoster student, or all timetable subjects that they undertake.

The key point is that the receiving systems do not have to build out the whole of their data model to understand or implement IMS OneRoster, and it becomes an ongoing choice as to whether they ever need to ingest that data back into their core systems.

The goal of NIAS is to potentially allow multiple open standards with specialist areas of expertise to co-exist in the most productive way for end users. This removes the need for a single standard to cover all uses, and means that each open standard can focus on adding its particular value. There is no need to pick a winner in order to build out a comprehensive model that covers all possible activities.

## 2.3. NIAS Servces

The core components provide a lightweight architecture that works through pure exchange of data messages. To this foundation we can add a number of services, each of which helps to solve particular integration concerns that NSIP has identified through its work with stakeholders.

### 2.3.1. SIF Privacy Service (SPS)

This service allows users to attach privacy filters to any outbound stream of SIF messages. Filters are held independently of the data and can be edited or specialized for any particular purpose. An editing UI is provided.

The current filters implement the NSIP Privacy Framework constraints for profiles of SIF data against the APP ratings; the default profiles are, therefore, Low, Medium, High and Extreme.

All data transformed via the service is then exposed as SSF endpoints for consumption, meaning that privacy control for clients is achieved simply by pointing the client to the correct endpoint rather than managing any data access. This approach also means that data producers do not have to be concerned with implementing privacy policy in their own solutions, and that policy is applied consistently to all data passing through NIAS.

### 2.3.2. IMS Ingest 

This is a specialized input to receive IMS OneRoster information, with an additional object id parameter that allows the data to be connected to the main data model at an insertion point of the users choosing—thus linking the supplied data to all other queries available through the SMS.

IMS data in its original form can also be consumed from the endpoint. Thus, if the integrations scenario is focused on linking the data, IMS OneRoster clients can produce and consume their data from the service, but with no need for onward systems to ever ingest the data unless they choose to.

### 2.3.3. Lightweight Analytics Service (LAS)
 
This service extracts SIF data from the SMS via query (all attendance records for a school, for example), and creates data arrays suitable for presentation in the family of visualisation tools based on the D3 specifications (a specialized json infrastructure for visualising data).

D3-based clients are then easily instantiated in html pages for lightweight dashboarding and basic reporting. These same services can be used to provide interactive data analysis support to NAPLAN results users where no systemic BI capability is available.

### 2.3.4. CSV-SIF Conversion

This is a simple service to support NAPLAN Online integrations. Validating CSV files is significantly harder than validating XML, but there will be a strong preference in schools and jurisdictions initially to produce registration data in CSV format. This service converts CSV input to the relevant SIF objects for onward transmission to the National Assessment Platform.
