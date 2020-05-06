# Command servedir

![CI](https://github.com/foxygoat/servedir/workflows/ci/badge.svg?branch=master)
[![Godoc](https://img.shields.io/badge/godoc-ref-blue)](https://pkg.go.dev/foxygo.at/servedir)
[![Slack chat](https://img.shields.io/badge/slack-gophers-795679?logo=slack)](https://gophers.slack.com/messages/foxygoat)

Simple HTTP server inspired by `python -mSimpleHTTPServer`.

Serves files from given directory on specified or next free port.

Install and run with

    > go get foxygo.at/servedir
    > servedir
    Starting HTTP server at http://localhost:52537

There are options for port, listen interface and directory, see `servedir --help`.
