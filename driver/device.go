// Copyright (c) 2015-2020 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package driver

import (
	"fmt"
	"log"

	"github.com/google/gousb"
)

// BareUsbDevice represents a USB device not a USBMTC device.
type BareUsbDevice struct {
	dev                 *gousb.Device
	intf                *gousb.Interface
	cfg                 *gousb.Config
	BulkInEndpoint      *gousb.InEndpoint
	BulkOutEndpoint     *gousb.OutEndpoint
	InterruptInEndpoint *gousb.InEndpoint
}

// Close closes the Device.
func (d *BareUsbDevice) Close() error {
	d.intf.Close()
	err := d.cfg.Close()
	if err != nil {
		return err
	}
	return d.dev.Close()
}

// String providers the Stringer interface method for Device.
func (d *BareUsbDevice) String() string {
	return fmt.Sprintf("%#v", d)
}

// Write writes to the USB device's bulk out endpoint.
func (d *BareUsbDevice) Write(p []byte) (n int, err error) {
	log.Printf("Writing Data: %x", p)
	log.Printf("Writing Data: %s", p)
	return d.BulkOutEndpoint.Write(p)
}

// WriteString writes the given string to the Device and returns the number
// of bytes written along with an error code.
func (d *BareUsbDevice) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// Read reads from the USB device's bulk in endpoint.
func (d *BareUsbDevice) Read(p []byte) (n int, err error) {
	return d.BulkInEndpoint.Read(p)
}
