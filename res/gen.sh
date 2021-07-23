#!/bin/sh
go get github.com/akavel/rsrc

rsrc -ico icon.ico -arch 386
rsrc -ico icon.ico -arch amd64
