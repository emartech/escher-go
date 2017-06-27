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
