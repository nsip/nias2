# nias-go-naplan-registration
NIAS. golang naplan registration modules.

TODO - Update for nias2

Install
* go version 1.6+
* go get github.com/pwaller/goupx

This is a Golang version of the [NIAS](http://github.com/nsip/nias) open-source components, specific to NAPLAN validation. It uses
[NATS](http://nats.io) queues. This module is currently limited to NAPLAN registration validation.

# Installation

## Binary, DOS
Manually unzip file directory in the zip "go-nias" and put it in c:\
Then run gonias.bat file from that directory

## From source code

[Install golang](https://golang.org/doc/install). Making sure you have a working
`$GOPATH` etc (common mistake is to skip the `src/` directory after `$GOPATH`)

In `$GOPATH/src/github.com/nsip` do:
   git clone https://github.com/nsip/nias-go-naplan-registration
   ./build.sh

# Running

In Unix-like systems (including OSX), `gonias.sh` launches the required processes, and `shutdown.sh` shuts them down. The PIDs of
the processes are stored in `nias.pid`.

In Windows, `gonias.bat` launches the required processes for an examination year of 2017, and `gonias.bat` launches the required processes for an examination year of 2016.

# Code Structure

`aggregator`
* Contains the web service coordinating the microservices. Runs on localhost:1234.
* Paths:
  * `/validation`: validation page
  * `/naplan/reg/:stateID`: endpoint to post input CSV files to
  * `/statusfeed/:txID`: SSE endpoint for status/progress updates for transaction txID
  * `/readyfeed/:txID`: SEE endpoint to announce when all messages in a transaction have been processed
  * `/data/:txID`: errors data for a given transaction

`aslvalidator`
* Validates the ACARA School IDs against the ASL schools list for the given state
* Flags:
  * `-s`: The NATS server URLs (comma-delimited) [Default: `nats.DefaultURL`]
  * `-t`: Whether to display timestamps [Default: `false`]
  * `-vtype`: Validation type [Default: `ASL`]
  * `-topic`: Root topic name to subscribe to [Default: `validation`]
  * `-qGroup`: The consumer group to join for parallel processing [Default: `aslvalidation`]
  * `-state`: The state identifier for this service (VIC, SA, NT, WA, ACT, TAS, NSW, QLD) [Default: `naplan`]

`aslvalidator/schoolslist`
* The ASL schools list used for validation

`dobvalidator`
* Validates the Date of Birth of the student for the given examination year
* Flags:
  * `-s`: The NATS server URLs (comma-delimited) [Default: `nats.DefaultURL`]
  * `-t`: Whether to display timestamps [Default: `false`]
  * `-vtype`: Validation type [Default: `date`]
  * `-topic`: Root topic name to subscribe to [Default: `validation`]
  * `-qGroup`: The consumer group to join for parallel processing [Default: `aslvalidation`]
  * `-state`: The state identifier for this service (VIC, SA, NT, WA, ACT, TAS, NSW, QLD) [Default: `naplan`]
  * `-tstyr`: The year in which the test will occur; used to baseline the year/test level age range windows [Default: `2017`]

`idvalidator`
* Validates the identities in the file for possible duplicates
* Flags:
  * `-s`: The NATS server URLs (comma-delimited) [Default: `nats.DefaultURL`]
  * `-t`: Whether to display timestamps [Default: `false`]
  * `-vtype`: Validation type [Default: `identity`]
  * `-topic`: Root topic name to subscribe to [Default: `validation`]
  * `-qGroup`: The consumer group to join for parallel processing [Default: `aslvalidation`]
  * `-state`: The state identifier for this service (VIC, SA, NT, WA, ACT, TAS, NSW, QLD) [Default: `naplan`]

`schemavalidator`
* Validates the data against schemas
* Flags:
  * `-s`: The NATS server URLs (comma-delimited) [Default: `nats.DefaultURL`]
  * `-t`: Whether to display timestamps [Default: `false`]
  * `-vtype`: Validation type [Default: `ASL`]
  * `-topic`: Root topic name to subscribe to [Default: `validation`]
  * `-qGroup`: The consumer group to join for parallel processing [Default: `aslvalidation`]
  * `-state`: The state identifier for this service (VIC, SA, NT, WA, ACT, TAS, NSW, QLD) [Default: `naplan`]
  * `-jsonSchema`: The schema file to be used for validation by this instance of the validator [Default: 'core.json`]

`schemavalidator/schemas`
Schemas used for validation of registration data

`schemavalidator/schemas/core.json`
JSON-Schema of registration data with global applicability

`schemavalidator/schemas/local.json`
JSON-Schema of registration data specific to the jurisdiction


