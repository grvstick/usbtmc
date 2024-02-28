// Copyright (c) 2015-2020 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package usbtmc

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/grvstick/usbtmc/usb"
)

// UsbTmc models a USBTMC device, which includes a USB device and the required
// USBTMC attributes and methods.
type UsbTmc struct {
	UsbDevice       usb.UsbDevice
	BTag            byte
	TermChar        byte
	TermCharEnabled bool
}

func (d *UsbTmc) Write(data []byte) (int, error) {
	log.Printf("Sending: %s", data)
	d.BTag = (d.BTag % 255) + 1
	if d.TermCharEnabled {
		data = append(data, d.TermChar)
	}
	header := encodeBulkOutHeader(d.BTag, uint32(len(data)), true)
	packet := append(header[:], data...)

	if len(packet)%4 != 0 {
		packet = append(packet, make([]byte, 4-len(packet)%4)...)
	}

	return d.UsbDevice.BulkOutEndpoint.Write(packet)
}

// Read creates and sends the header on the bulk out endpoint and then reads
// from the bulk in endpoint per USBTMC standard.
func (d *UsbTmc) Read() ([]byte, error) {
	received := []byte{}
	eom := false
	for !eom {
		buf := make([]byte, d.UsbDevice.BulkInMaxPktSize*12)
		d.BTag = nextbTag(d.BTag)
		reqInMsg := encodeMsgInBulkOutHeader(d.BTag, uint32(len(buf)), d.TermCharEnabled, d.TermChar)
		if _, err := d.UsbDevice.BulkOutEndpoint.Write(reqInMsg[:]); err != nil {
			return []byte{}, err
		}

		n, err := d.UsbDevice.BulkInEndpoint.Read(buf)
		if err != nil || n < 12 {
			return []byte{}, err
		}
		header, buf := buf[:12], buf[12:]
		if header[0] != byte(MsgIdDevDepMsgIn) {
			return []byte{}, fmt.Errorf("not a valid read response: %x", header)
		}
		transferSize := int(binary.LittleEndian.Uint32(header[4:8]))
		received = append(received, buf[:transferSize]...)

		if n >= headerSize+transferSize {
			eom = header[8]&1 == 1
		}

	}
	return received, nil
}

// Close closes the underlying USB device.
func (d *UsbTmc) Close() error {
	return d.UsbDevice.Close()
}

// WriteString writes a string using the underlying USB device. A newline
// terminator is not automatically added.
func (d *UsbTmc) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}

// Query writes the given string to the USBTMC device and returns the returned
// value as a string. A newline character is automatically added to the query
// command sent to the instrument.
func (d *UsbTmc) Query(s string) (string, error) {
	_, err := d.Write([]byte(s))
	if err != nil {
		return "", err
	}
	buf, err := d.Read()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
