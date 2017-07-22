// Escher - Stateless HTTP Request Signing
//
// Escher is a stateless HTTP request signing spec to allow secure authorization and request validation. It adds an additional security layer and an authentication layer over HTTPS. The algorithm is based on Amazon’s AWS4 authentication. The goal was implementing a protocol where request’s integrity is validated. It is impossible to be modified without knowing the secret.
// Escher is great for creating secure REST API servers, both the client side request signing process, and the server side validation and authentication processes are implemented. The protocol also provides a solution for presigning URLs with expiration time.
// The status is work in progress. We’re working on the documentation and finishing the first implementations.
// Your help will be much welcomed, we are especially interested in code reviews and audits, new library implementations in different languages, writing the documentation and promotion.
// Feel free to join and discuss at Escher’s general mailing list, and follow our Twitter user @escherauth.
//
// Specification
//
// This document defines the Escher HTTP request signing framework.
//
// Abstract
//
// The Escher HTTP request signing framework enables a third-party client to securely access an HTTP service.
// It defines the creation of an authorization header, and includes message integrity checking, and can prevent replaying messages.
// It is designed to be used with the HTTPS protocol.
// The framework is based on Amazon’s AWS4 authentication, and is compatible with Amazon services using their AWS4 protocol.
//
// Copyright Notice
// Copyright (c) 2014 Emarsys. All rights reserved.
//
// 1. Introduction

// The Escher HTTP request signing framework is intended to provide a secure way for clients to sign HTTP requests, and servers to check the integrity of these messages. The goal of the protocol is to introduce an authentication solution for REST API services, that is more secure than the currently available protocols.
// RFC 2617 (HTTP Authentication) defines Basic and Digest Access Authentication. They’re widely used, but Basic Access Authentication doesn’t encrypt the secret and doesn’t add integrity checks to the requests. Digest Access Authentication sends the secret encrypted, but the algorithm with creating a checksum with a nonce and using md5 should not be considered highly secure these days, and as with Basic Access Authentication, there’s no way to check the integrity of the message.
// RFC 6749 (OAuth 2.0 Authorization) enables a third-party application to obtain limited access to an HTTP service, either on behalf of a resource owner by orchestrating an approval interaction between the resource owner and the HTTP service, or by allowing the third-party application to obtain access on its own behalf. This is not helpful for a machine-to-machine communication situation, like a REST API authentication, because typically there’s no third-party user involved. Additionally, after a token is obtained from the authorization endpoint, it is used with no encryption, and doesn’t provide integration checking, or prevent repeating messages. OAuth 2.0 is a stateful protocol which needs a database to store the tokens for client sessions.
// Amazon and other service providers created protocols addressing these issues, however there is no public standard with open source implementations available from them. As Escher is based on a publicly documented, widely, in-the-wild used protocol, the specification does not include novelty techniques.
//
// 2. Signing an HTTP Request
//
// Escher defines a stateless signature generation mechanism. The signature is calculated from the key parts of the HTTP request, a service identifier string called credential scope, and a client key and client secret.
// The signature generation steps are: canonicalizing the request, creating a string to calculate the signature, and adding the signature to the original request.
// Escher supports two hash algorithms: SHA256 and SHA512 designed by the NSA (U.S. National Security Agency).
//
// 2.1. Canonicalizing the Request
//
// In order to calculate a checksum from the key HTTP request parts, the HTTP request method, the request URI, the query parts, the headers, and the request body have to be canonicalized. The output of the canonicalization step will be a string including the request parts separated by LF (line feed, “n”) characters. The string will be used to calculate a checksum for the request.
//
// 2.1.1. The HTTP method
//
// The HTTP method defined by RFC2616 (Hypertext Transfer Protocol) is case sensitive, and must be available in upper case, no transformation has to be applied:
// POST
//
// 2.1.2. The Path
//
// The path is the absolute path of the URL. Starts with a slash (/) character, and does not include the query part (and the question mark).
// Escher follows the rules defined by RFC3986 (Uniform Resource Identifier) to normalize the path. Basically it means:
// Convert relative paths to absolute, remove redundant path components.
// URI-encode each path components:
// the “reserved characters” defined by RFC3986 (Uniform Resource Identifier) have to be kept as they are (no encoding applied)
// all other characters have to be percent encoded, including SPACE (to %20, instead of +)
// non-ASCII, UTF-8 characters should be percent encoded to 2 or more pieces (á to %C3%A1)
// percent encoded hexadecimal numbers have to be upper cased (eg: a%c2%b1b to a%C2%B1b)
// Normalize empty paths to /.
// For example:
//
// 	/path/resource/
//
// 2.1.3. The Query String
//
// RFC3986 (Uniform Resource Identifier) should provide guidance for canonicalization of the query string, but here’s the complete list of the rules to be applied:
// URI-encode each query parameter names and values
// the “reserved characters” defined by RFC3986 (Uniform Resource Identifier) have to be kept as they are (no encoding applied)
// all other characters have to be percent encoded, including SPACE (to %20, instead of +)
// non-ASCII, UTF-8 characters should be percent encoded to 2 or more pieces (á to %C3%A1)
// percent encoded hexadecimal numbers have to be upper cased (eg: a%c2%b1b to a%C2%B1b)
// Normalize empty query strings to empty string.
// Sort query parameters by the encoded parameter names (ASCII order).
// Do not shorten parameter values if their parameter name is the same (key=B&key=A is a valid output), the order of parameters in a URL may be significant (this is not defined by the HTTP standard).
// Separate parameter names and values by = signs, include = for empty values, too
// Separate parameters by &
// For example:
//
//	foo=bar&abc=efg
//	2.1.4. The Headers
//
// To canonicalize the headers, the following rules have to be followed:
// Lower case the header names.
// Separate header names and values by a :, with no spaces.
// Sort header names to alphabetical order (ASCII).
// Group headers with the same names into a header, and separate their values by commas, without sorting.
// Trim header values, keep all the spaces between quote characters (").
// For example:
//
// 	accept:*/*
// 	user-agent:example-client
// 	connection:close
// 	content-type:application/x-www-form-urlencoded
// 	content-length:21
// 	host:example.com
//
// 2.1.5. Signed Headers
//
// The list of headers to include when calculating the signature. Lower cased value of header names, separated by ;, like this:
// date;host
//
// 2.1.6. Body Checksum
//
// A checksum for the request body, aka the payload has to be calculated. Escher supports SHA-256 and SHA-512 algorithms for checksum calculation. If the request contains no body, an empty string has to be used as the input for the hash algorithm.
// The selected algorithm will be added later to the authorization header, so the server will be able to use the same algorithm for validation.
// The checksum of the body has to be presented as a lower cased hexadecimal string,
// for example:
// 	fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
//
// 2.1.7. Concatenating the Canonicalized Parts
//
// All the steps above produce a row of data, except the headers canonicalization, as it creates one row per headers. These have to be concatenated with LF (line feed, “n”) characters into a string.
// An example:
//
// 	POST
// 	/path/resource/
// 	foo=bar&abc=efg
// 	accept:*/*
// 	user-agent:example-client
// 	connection:close
// 	content-type:application/x-www-form-urlencoded
// 	content-length:21
// 	host:example.com
// 	date;host
// 	fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210
//
// 2.2. Creating the Signature
//
// The next step is creating another string which will be directly used to calculate the signature.
//
// 2.2.1. Algorithm ID
//
// The algorithm ID comes from the algo_prefix (default value is ESR) and the algorithm used to calculate checksums during the signing process. The string algo_prefix, “HMAC”, and the algorithm name should be concatenated with dashes, like this:
//
// 	ESR-HMAC-SHA256
//
// 2.2.2. Long Date
//
// The long date is the request date in the ISO 8601 basic format, like YYYYMMDD + T + HHMMSS + Z. Note that the basic format uses no punctuation. Example is:
//
// 	20141022T120000Z
//
// This date has to be added later, too, as a date header (default header name is X-Escher-Date).
//
// 2.2.3. Date and Credential Scope
//
// Next information is the short date, and the credential scope concatenated with a / character. The short date is the request date’s date part, an ISO 8601 basic formatted representation, the credential scope is defined by the service. Example:
//
// 	20141022/eu-vienna/yourproductname/escher_request
//
// This will be added later, too, as part of the authorization header (default header name is X-Escher-Auth).
//
// 2.2.4. Checksum of the Canonicalized Request
//
// Take the output of step 2.1.7., and create a checksum from the canonicalized checksum string. This checksum has to be represented as a lower cased hexadecimal string, too. Something like this will be an output:
//
// 	0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
//
// 2.2.5. Concatenating the Signing String
//
// Concatenate the outputs of steps 2.2. with LF characters. Example output:
//
// 	ESR-HMAC-SHA256
// 	20141022T120000Z
// 	20141022/eu-vienna/yourproductname/escher_request
// 	0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
//
// 2.2.6. The Signing Key
//
// The signing key is based on the algo_prefix, the client secret, the parts of the credential scope, and the request date.
// Take the algo_prefix, concatenate the client secret to it. First apply the HMAC algorithm to the request date, then apply the actual value on each of the credential scope parts (splitted at /). The end result is a binary signing key.
// Pseudo code:
//
// 	signing_key = algo_prefix + client_secret
// 	signing_key = HMAC.Digest(short_request_date, signing_key)
// 	foreach credential_scope.split('/') as scope_part
// 	  signing_key = HMAC.Digest(scope_part, signing_key)
// 	end_foreach
// 	return signing_key
//
// 2.2.7. Create the Signature
//
// The signature is created from the output of steps 2.2.5. (Signing String) and 2.2.6. (Signing Key). With the selected algorithm, create a checksum. It has to be represented as a lower cased hexadecimal string. Something like this will be an output:
//
// 	abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd
//
// 2.3. Adding the Signature to the Request
//
// The final step of the Escher signing process is adding the Signature to the request. Escher adds a new header to the request, by default, the header name is X-Escher-Auth. The header value will include the algorithm ID (see 2.2.1.), the client key, the short date and the credential scope (see 2.2.3.), the signed headers string (see 2.1.5.) and finally the signature (see 2.2.7.).
// The values of this inputs have to be concatenated like this:
//
// 	ESR-HMAC-SHA256 Credential=CLIENT_KEY/20141022/eu-vienna/yourproductname/escher_request,
// 	SignedHeaders=date;host, Signature=abcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcdabcd
//
// 3. Presigning a URL
//
// The URL presigning process is very similar to the request signing procedure. But for a URL, there are no headers, no request body, so the calculation of the Signature is different. Also, the Signature cannot be added to the headers, but is included as query parameters.
// A significant difference is that the presigning allows defining an expiration time. By default, it is 86400 secs, 24 hours. The current time and the expiration time will be included in the URL, and the server has to check if the URL is expired.
//
// 3.1. Canonicalizing the URL to Presign
//
// The canonicalization for URL presigning is the same process as for HTTP requests, in this section we will cover the differences only.
//
// 3.1.1. The HTTP method
//
// The HTTP method for presigned URLs is fixed to:
//
// 	GET
// 	3.1.2. The Path
// 	The path is coming from the URL, and the same canonicalization process has to be applied to them as for HTTP requests.
//
// For example:
//
// 	/path/resource/
//
// 3.1.3. The Query String
//
// The query is coming from the URL, but the algorithm, credentials, date, expiration time, and signed headers have to be added to the query parts.
//
// 	foo=bar&abc=efg
//
// 3.1.4. The Headers
//
// A URL has no headers, Escher creates the Host header based on the URL’s domain information, and adds it to the canonicalized request.
// For example:
//
// 	host:example.com
//
// 3.1.5. Signed Headers
//
// It will be host, as that’s the only header included. Example:
//
// 	host
package escher
