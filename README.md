# Authful

A lightweight OAuth 2.0 server in Go with accompanying CLI

[![Build Status](https://circleci.com/gh/danielsomerfield/authful.svg?style=svg)](https://circleci.com/gh/danielsomerfield/authful)

## Requirements

Tested with go 1.8.3

## Building

To run the tests and build the authful binary run:

    make
    
## Running

Once built, you can run the server. Currently it only runs on port 8080.
Configuration and command line options are on the way.

    ➜ ./authful &
    2017/07/19 05:28:26 Starting server up http server on port 8080
    
    ➜ curl http://localhost:8080/health
    {"status":"ok"}                                                              ➜  authful git:(master) ✗ curl -v http://localhost:8080/health
    
## Using
Coming soon...