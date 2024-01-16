# USB TMC
This is a fork from https://github.com/gotmc/usbtmc. 

## Overview

The purpose is to renew the repo completely based on gousb, since the feature of selecting between two usb libraries is not working anymore. This is targeted specifically for controlling USB test instruments without visa library.

## USB TMC Specs
Refer to [USB TMC Spec](https://www.usb.org/document-library/test-measurement-class-specification) for details. If you are not familiar with USB protocol, which I'm also, take a tour on [USB in a nutshell](http://www.beyondlogic.org/usbnutshell/)

## Installation

```bash
$ go get github.com/grvstick/usbtmc
```

## Usage
You'll need to install gousb. Please note that libusb is a prerequisite for gousb

```bash
$ go get -v github.com/google/gousb
```
```device_test.go``` contains the usage of the library

## Documentation
Refer to original repo or see 
- <https://godoc.org/github.com/gotmc/usbtmc>


### Disclosure and Call for Help

While this package works, it does not fully implement the [USBTMC][]
specification.  Please submit pull requests as needed to increase
functionality, maintainability, or reliability.

## License

[usbtmc][gousbtmc] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[godoc badge]: https://godoc.org/github.com/gotmc/usbtmc?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/usbtmc
[golibusb]: https://github.com/gotmc/libusb
[gousb]: https://github.com/google/gousb
[libusb]: http://libusb.info
[LICENSE.txt]: https://github.com/gotmc/libusb/blob/master/LICENSE.txt
[license badge]: https://img.shields.io/badge/license-MIT-blue.svg
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/gotmc/usbtmc
[report card]: https://goreportcard.com/report/github.com/gotmc/usbtmc
[Scott Chacon]: http://scottchacon.com/about.html
[travis badge]: http://img.shields.io/travis/gotmc/usbtmc/master.svg
[travis link]: https://travis-ci.org/gotmc/usbtmc
[usbtmc]: http://www.usb.org/developers/docs/devclass_docs/
[gousbtmc]: https://github.com/gotmc/usbtmc
