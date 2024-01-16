package driver

import "github.com/google/gousb"

// The base class codes are part of the USB class codes "used to identify a
// device's functionality and to nominally load a device driver based on that
// functionality." [Source: http://www.usb.org/developers/defined_class]
// The USBTMC standard refers to these as the bInterfaceClass. The only base
// class code required for USBTMC is the Application Specific Base Class 0xFE.

const (
	baseClsAppSpecific gousb.Class    = 0xfe
	subClsTmc          gousb.Class    = 0x03
	protTmc            gousb.Protocol = 0x01
	prot488            gousb.Protocol = 0x02
)
