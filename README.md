# USB TMC
This is a fork from https://github.com/gotmc/usbtmc. 

## Overview

The purpose is to renew the repo completely based on gousb, since the feature of selecting between two usb libraries is not working anymore(or I failed miserably). This is targeted specifically for controlling USB test instruments without visa library. The goal is to provide similar experience as in pyvisa package from python.

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

Refer to  ```tests``` directory for example usage of the library

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
