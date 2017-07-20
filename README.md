[![Build Status](https://travis-ci.org/EscherAuth/escher.svg?branch=master)](https://travis-ci.org/EscherAuth/escher) Escher - HTTP request signing lib
===================================

Go implementation of the [AWS4](http://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html) compatible [Escher](https://github.com/emartech/escher) HTTP request signing and authentication library.

Prerequisite
------------

this will download the test cases for the escher implementation, and set in the env the required env key(s)

    $ source env.sh

Run the tests
-------------

in 1.9:

    $ go test ./...

in older go versions:

    $ go test $(go list ./... | grep -v /vendor/)

About Escher
------------

More details are available at our [Escher documentation site](http://escherauth.io/).

Install
-------

```bash
# install dep management tool
go get -u github.com/golang/dep/cmd/dep
# dep ensure the missing dependency packages
dep ensure
```

ENV configuration
-----------------

## Configuration

Everything except the Credential scope is optional.

### ESCHER_CONFIG

JSON string that can include the all the configuration parameters:
* credentialScope
* vendorKey
* algoPrefix
* hashAlgo
* authHeaderName
* dateHeaderName

### ESCHER_ALGO_PREFIX
Set the used Algo prefix when config json not includes it

### ESCHER_HASH_ALGO
set the hashAlgo when config json not includes it

### ESCHER_VENDOR_KEY
set the vendorKey when config json not includes it

### ESCHER_AUTH_HEADER_NAME
set the AUTH_HEADER_name when config json not includes it

### ESCHER_DATE_HEADER_NAME
set the DATE_HEADER_name when config json not includes it

### ESCHER_CREDENTIAL_SCOPE
set the credentialScope when config json not includes it

## KeyPool 

### KEY_POOL

JSON serialized array of map that contains credentials with the following keys definitions:
* keyId
* secret
