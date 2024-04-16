// Copyright (c) 2015-2020 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"encoding/binary"
	"fmt"
)

// UsbTmc models a USBTMC device, which includes a USB device and the required
// USBTMC attributes and methods.
type UsbTmc struct {
	UsbDevice       *Device
	bTag            byte
	termChar        byte
	termCharEnabled bool
}

func (d *UsbTmc) Write(data []byte) (int, error) {
	// log.Printf("Sending: %s", data)
	d.rotateBtag()
	if d.termCharEnabled {
		data = append(data, d.termChar)
	}
	header := encodeBulkOutHeader(d.bTag, uint32(len(data)), true)
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
		d.rotateBtag()
		reqInMsg := encodeMsgInBulkOutHeader(d.bTag, uint32(len(buf)), d.termCharEnabled, d.termChar)
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

func (d *UsbTmc) rotateBtag() {
	d.bTag = (d.bTag % 255) + 1
}

func NewUsbTmc(dev *Device, termchar byte) *UsbTmc {
	return &UsbTmc{
		UsbDevice:       dev,
		bTag:            0,
		termChar:        termchar,
		termCharEnabled: true,
	}
}
