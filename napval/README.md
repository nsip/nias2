# NIAS `napval` module

This is the module of the NIAS toolset used for the validation of NAPLAN registration records

# Installation

## Precondition

* go version 1.8+
* ```go get github.com/pwaller/goupx```


## Binary, DOS
Manually unzip file directory in the zip "go-nias" and put it in c:\
Then run gonias.bat file from that directory

## From source code

[Install golang](https://golang.org/doc/install). Making sure you have a working
`$GOPATH` etc (common mistake is to skip the `src/` directory after `$GOPATH`)

In `$GOPATH/src/github.com/nsip` do:

    git clone https://github.com/nsip/nias2
    ./build_napval.sh

# Running

In Unix-like systems (including OSX), `gonias.sh` launches the required processes, and `shutdown.sh` shuts them down. The PIDs of
the processes are stored in `nias.pid`.

In Windows, `gonias.bat` launches the required processes for an examination year of 2017, and `gonias.bat` launches the required processes for an examination year of 2016.

The web interface to the validator can be accessed at `http://localhost:1325` by default; you can configure the Web Server Port in `nias.toml`


# Code structure
    
`app/napval/nias.toml` 
Configuration file for NIAS run for validation of NAPLAN Registration records:
  * TestYear: the baseline year for date of birth validation 
  * ValidationRoute: the validators to which every incoming message is sent
      * `schema`: schema validation against `core.json`  
      * `schema2`: schema validation against `core_parent2.json` (ensures that if one Parent2 field is present, all such fields are present)
      * `local`: schema validation against `local.json` 
      * `id`: identity validation (detection of duplicates)
      * `dob`: date of birth validation
      * `asl`: check of validity of ASL school identifers
      * `psi`: check of validity of Platform Student Identifier checksums
      * `numericvalid`: does number validation on contents of schema (currently restricted to FTE maximum and minimum for NAPLAN)
  * WebServerPort: the port on which the NIAS web server runs
  * NATSPort: the port on which the NATS Streaming server runs
  * PoolSize: the number of parallel connections run in the microservice distributor

Configuration files:
  * `app/napval/schemas/` : Schemas for validating incoming messages. CSV is converted to JSON, and is validated against JSON Schema:
    * `core.json`: The schema for NAPLAN registration records.
    * `core_parent2.json`: The schema ensuring that if one Parent2 field is present, all such fields are present.
    * `local.json`: Dummy schema for local validation of NAPLAN registration records. Can be used to impose more restrictive conditions on validation, to satisfy local requirements.
  * `app/napval/schoolslist/` : Contains CSV export of the [Australian Schools List](http://asl.acara.edu.au), using in validation
  * `app/napval/templates/` : Contains templates for populating SIF XML
  * `app/napval/var/` : Contains ledis database instance
  * `app/napval/public/` : Contains web server site, including CSS and Javascript
  
`lib/`
* Microservices invoked by NIAS directly via `harness/harness.go`
  * `ledis.go` : Launch the ledis database
  * `aslservice.go` : Validate the ASL school identifiers in a registration record against the ASL data in `harness/schoolslist/`
  * `webserver.go` : Launch web service to deal with RESTful queries for validation. 
  * `distributor.go` : Launch a pool of message handlers to deal with incoming messages, as the microservice bus. 
    * The pool involves instances of NATS Server, NATS Streaming Server, or internal memory channels. 
    * The distributor handles incoming requests as a multiplexer (from "requests" subject to the _distID_ subject: a new GUID)
    * The distributor assigns incoming messages (from the _distID_ subject) for processing by the sequence of microservices named in the message's Route attribute; the output of each named microservice is published to the _srvcID_ subject (a new GUID).
    * The distributor stores incoming messages to ledis (from the _srvcID_ subject).
    * This means that all microservice outputs are output to ledis.
* Microservices invoked by NIAS via the message Route attribute 
  * `dobservice.go` : Date of Birth validator
  * `psiservice.go`: Checksum validator for Platform Student Identifiers
  * `numericvalidservice.go`: Validator of numeric values in NAPLAN registration
  * `idservice.go` : Check each message in a transmission for uniqueness within the transmission. Check involves two keys: (LocalId, ASLSchoolId), and (LocalId, ASLSchoolId, FamilyName, GivenName, BirthDate).
  * `schemaservice.go` : Validate a message against either the core JSON schema or the local JSON schema. The service replaces some JSON Schema error messages with custom messages.
* Support code
  * `nats.go` : Connector code for NATS, involving connectors to storage (store), service handlers (srv), and inbound distributors from the web gateway (dist).
  * `niasmessage.go` : NIAS Message wrapper types
  * `config.go` : read in the NIAS configuration file (`harness/nias.toml`)
  * `vtypes.go` : common types used for validation. Includes the validation error type, and the registration record type (all fields in NAPLAN).
  * `store.go` : code to store messages in ledis. Has mutex support. 
  * `service.go` : interface to handle message requests: request, response, errors
  * `serviceregister.go` : registry of microservices, mapping Route keys to service instances, and with code for processing messages according to their Route attribute


#API

## Supported Queries
* `POST /naplan/reg/validate` : validate the record(s), whether in XML or CSV. This involves publishing the records received onto the microservice bus, with the configured list of validators as the message route. Blank entries are stripped.
* `POST /naplan/reg/convert` : convert the record(s) from CSV to XML, using the templates in  `harness/templates/`. Response is the XML records.
* `GET /naplan/reg/status/:txid` : receive a status report for the validation request with transmission identifier `:txid`
* `GET /naplan/reg/results/:txid` : receive the analysis results for the validation request with transmission identifier `:txid`
* `GET /naplan/reg/results/:txid/:fname` : receive the analysis results for the validation request with transmission identifier `:txid`, as a CSV file to be named `:fname`

## Database
* All messages published to storage as outputs of the distributor are stored in ledis under a list (`rpush`) with the key of `nvr:` followed by the transmission ID. That means that the key for a transmission will access a list of each consecutive microservice output, for each record in that transmission.

## Format
* Ingest Response: response to `POST /naplan/reg/validate`. JSON object:
  * record count (`Records`)
  * transmission identifier (`TxID`). The transmission identifier applies to the specific payload.
* NIAS Message: metadata for any message sent on the microservice bus:
  * Body: message content
  * SeqNo: sequence number of the message within the transmission (corresponding to a single REST payload)
  * TxID: transmission identifier (GUID) for the REST payload
  * MsgID: GUID for the message
  * Target: namespace on ledis under which messages will be stored
  * Route: sequence of microservices that the message is to be passed to. _(Chaining of microservices is not currently supported)_
* NIAS Error Message: as for NIAS Message
  * Body: 
    * Description: description of the validation error
    * Field: field in which the validation error was found
    * OriginalLine: Line of the transmission payload in which the validation error was found
    * Vtype: Validation error type
