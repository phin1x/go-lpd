# go-lpd

[comment]: [![Version](https://img.shields.io/github/release-pre/phin1x/go-lpd.svg)](https://github.com/phin1x/go-ipp/releases/tag/v1.0.0)
[![Licence](https://img.shields.io/github/license/phin1x/go-lpd.svg)](https://github.com/phin1x/go-lpd/blob/master/LICENSE)

## Go Get

To get the package, execute:
```
go get -u github.com/phin1x/go-lpd
```

## Features

* rfc1179 compatible client
* create custom lpd requests
* parse control files

## Examples

Print a file
```go
client := lpd.NewClient("printserver", 515)
client.PrintFile("/path/to/file", "my-printer", nil)
```
## TODO's

* test print document / file, print jobs and delete jobs method
* parse result in get queue state methods
* add usernmae and job number filter in queue state and delete job methods
* write some tests

## Licence

Apache Licence Version 2.0

