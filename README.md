# Ancre :anchor: [WIP]

![License][license-img]

A golang implementation of [OpenTimestamps](https://opentimestamps.org/) client

Ancre anchors checksums in time.

## Install

Install [the go programming language](https://golang.org).

```bash
go get github.com/edouardparis/ancre
go install ./...
```

## Test

```bash
go test ./...
```
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg

## Usage

Only reading `.ots` file is good enough implemented

```
ancre info examples/hello-world.txt.ots
```
