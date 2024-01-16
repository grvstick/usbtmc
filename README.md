# usbtmc
This is a fork from https://github.com/gotmc/usbtmc. The purpose is to renew the repo completely based on gousb, since the feature of selecting between two usb libraries is not working anymore.

## Overview

Refer to [USB TMC Spec](https://www.usb.org/document-library/test-measurement-class-specification) for details.

## Installation

```bash
$ go get github.com/grvstick/usbtmc
```

## Usage

You'll need to install gousb. Please note that libusb is a prerequisite for gousb

```bash
$ go get -v github.com/google/gousb
```

## Documentation
Refer to original repo or see 
- <https://godoc.org/github.com/gotmc/usbtmc>

## Contributing

To contribute, please fork the repository, create a feature branch, and then
submit a [pull request][].

### Testing

Prior to submitting a [pull request][], please run:

```bash
$ make check
```

To update and view the test coverage report:

```bash
$ make cover
```

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
