*All files in this repository are licensed under Apache 2.0. For a full description of the license [click this link](LICENSE).*


# nias2
NIAS: NSIP Integration As A Service.

This is a Golang version of the [NIAS](http://github.com/nsip/nias) open-source components. This version incorporates NAPLAN validation and post-processing functionality, as well as the generic  functionality of NIAS.

# 1. Overview

NIAS is a suite of open-source components designed to enable as many different users as possible to quickly and easily solve issues of system integration using the Australian [SIF Data Model](http://specification.sifassociation.org/Implementation/AU/1.4/html/) for school education.

The product was developed by harnessing existing open source middleware components, including:
* [NATS](http://nats.io) streaming queues 
* [Echo](https://github.com/labstack/echo) web framework

Over these components, two main modules have been built:
* The __SIF Store & Forward (SSF)__ is an opinionated message queueing system, which ingests very large quantities of data and stores them for delivery to clients. XML messages on the system are assumed by default to be in SIF. The SSF service builds an education-standards aware REST interface on top of the NATS message queues, and provides a number of utility services to ease SIF-based integrations.
* The __SIF Memory Store (SMS)__ is a database that builds its internal structures from the data it receives, using RefIds both as keys to access stored messages, and to map out a network graph for SIF objects.

The software also uses these components as architecture to support Test Administration Authorities' interaction with NAPLAN Online:
* __napval__ validates NAPLAN registration records, in either SIF/XML or CSV format.
* __naprr__ post-processes the NAPLAN results & reporting dataset, including generating local reports, and aligning Year 3 Writing results to the codeframe.

The code is open source, and is released through the NSIP Github repository. 

## 1.1. Scope

This product delivers the following high level functions: 

1. Support for persistent and ordered queues of SIF messages, which can be reread multiple times.
2. Support for asynchronous queues in both clients and servers.
3. Support for format-agnostic messaging infrastructure.
4. Support for data exchange through an event/subscribe model (in brokered environments)
~~5. Support for message validation.~~
6. Support for extracting arbitrary relations between object types within SIF (bypassing need to configure service path queries, and simplifying the query API for objects).
~~7. Support for extracting arbitrary relations between object types from different standards (allowing multiple data standards to coexist in an integration, referring to the same entities).~~
8. Support for privacy filtering in middleware (which releases object providers from having to do privacy filtering internally).
~~9. Support for simple and extensible interactive analytics.~~
~~10. Support for the ODDESSA data dictionary as a service.~~
11. Support for data format conversions, including CSV to SIF, and SIF 2 to SIF 3.

This product only acts as middleware. It does not provide integration with the back ends of products (although this can be provided by combining NIAS with the [SIF Framework](https://github.com/nsip/sif3-framework-java)). It is not intended to deliver business value to end consumers, or to compete with existing market offerings.

The product delivers only exemplar analytics, and the SIF team is not committing to developing analytics and queries for all product users. Users that do develop their own analytics and queries are encouraged to contribute these back as open source.

~The product delivers only exemplar integrations between multiple standards (SIF/XML and [IMS OneRoster](https://www.imsglobal.org/lis/index.html)/CSV), and the SIF team is not committing to developing standards integrations for all product users. Users that do develop their own standards integrations are encouraged to contribute these back as open source.~

The product does not incorporate authentication or authorisation.


## 1.2. Constraints

* The product is released with the SIF-AU 3.4 XSD schema, and validates against it. Other schemas can be used, but may require re-coding of some modules.


# 2. Installation

## 2.1. Precondition

* go version 1.8+
* ```go get github.com/pwaller/goupx```


## 2.2. Binary, DOS
Manually unzip file directory in the zip `go-nias` and put it in c:\
Then run `gonias.bat` file from the `nias` subdirectory

## 2.3. From source code

[Install golang](https://golang.org/doc/install). Making sure you have a working
`$GOPATH` etc (common mistake is to skip the `src/` directory after `$GOPATH`)

In `$GOPATH/src/github.com/nsip` do:

    git clone https://github.com/nsip/nias2
    ./build.sh


## 2.4. Running NAPLAN Results & Reporting modules

Separate executables are run to process NAPLAN data; see [NAPVAL readme](./napval/README.md), [NAPRR readme](./naprr/README.md)

# 3. Code Structure

See also [NAPRR readme](./naprr/README.md), [NAPLVAL readme](./napval/README.md)

NIAS2 relies on the following infrastructure:
* [ledis](http://ledisdb.com), a NoSQL database similar to [redis](http://redis.io)
* xmllint, an XML validator included in the [libxml](http://xmlsoft.org) package

`unit_test_files/`
* Contains files used in unit/integration testing of the code. Currently restricted to CSV files input into the validation module.

`build.sh`, `build/`, `release.sh`
* `build.sh` is the script to build all NIAS2 executables for the various supported platforms. 
* `build_sms.sh` is the script to build the executables specific to the SSF and SMS. The builds for each platform are built in `build/{PLATFORM}/nias/`. 
* `build_napval.sh` is the script to build NIAS2 executables speficic to NAPLAN registration for the various supported platforms. The builds for each platform are built in `build/{PLATFORM}/napval/`. 
* `build_naprr.sh` is the script to build NIAS2 executables specific to NAPLAN results & reporting for the various supported platforms. The builds for each platform are built in `build/{PLATFORM}/naprr/`. 
* The `release_naprr.sh` and `release_napval.sh` script zips the builds at the top level of the `build` folder, and creates a new release of the NIAS2 code on github, including the binary files specific to NAPLAN results & reporting and NAPLAN registration, respectively. 
* The `release_napval.sh` script zips the builds at the top level of the `build` folder, and creates a new release of the NIAS2 code on github, including the binary files specific to both NAPLAN results & reporting and NAPLAN registration.
* The supported platforms are: 
  * Mac OSX
  * Windows 32 bit
  * Windows 64 bit
  * Linux 32 bit
  * Linux 64 bit

`bin/`
Contains the scripts and batch files to start and stop running NIAS. These are copied into the builds for each platform:
  * `gonias.sh`: launch NIAS (OSX, Linux)
  * `gonias.bat`: launch NIAS (Windows)
  * `stopnias.sh`: stop NIAS (OSX, Linux)
  * `stopnias.bat`: stop NIAS (Windows)

`tools/`
Contains utilities for managing NIAS2
  * `release.go` creates a new release of the NIAS2 code on github.

`app/`
Contains the code to run executables within NIAS as single pieces of software, along with necessary configuration files, and test executables. The configuration fields are copied into the binary distributions of NIAS.
  * `sms/` : the SIF Memory Store functionality of NIAS (ported from Ruby)
  * `napval/` : NAPLAN Registration records validation
  * `naprr/` : NAPLAN Results and Reporting post-processing
  * `napcomp/` : Comparison of students between NAPLAN Registration and NAPLAN Results and Reporting (included in `naprr` distribution)
  
`napval/`
Code for NAPLAN Registration records validation

`naprr/`
Code for NAPLAN Results and Reporting post-processing

`sms/`
SIF Memory Store functionality of NIAS

`xml/`
Golang structs corresponding to SIF XML objects relevant to executables. Currently limited to NAPLAN-specific objects.

`lib/`
Library code shared between all executables:
* `config.go`: read configuration files (in [toml](github.com/BurntSushi/toml) format)
* `encoding.go`: encode NIAS messages
* `nats.go`: create NATS connections and process chains
* `niasmessage.go`: message wrapper types
* `server_connections.go`: standardised NATS server access module
* `service.go`: service interface to handle message requests
* `transactiontracker.go`: transaction status reporting structure
    

