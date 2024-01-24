// Copyright (c) 2015-2020 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package usbtmc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/grvstick/usbtmc/driver"
)

// UsbTmc models a USBTMC device, which includes a USB device and the required
// USBTMC attributes and methods.
type UsbTmc struct {
	BareUsbDev      driver.BareUsbDevice
	BTag            byte
	TermChar        byte
	TermCharEnabled bool
}

// Write creates the appropriate USBMTC header, writes the header and data on
// the bulk out endpoint, and returns the number of bytes written and any
// errors.
// func (d *UsbTmc) Write(data []byte) (n int, err error) {
// 	maxTransferSize := d.usbDevice.BulkOutEndpoint.Desc.MaxPacketSize
// 	for pos := 0; pos < len(data); {
// 		d.bTag = nextbTag(d.bTag)
// 		dataLen := len(data[pos:])
// 		if dataLen > maxTransferSize-bulkOutHeaderSize {
// 			dataLen = maxTransferSize - bulkOutHeaderSize
// 		}
// 		header := encodeBulkOutHeader(d.bTag, uint32(dataLen), true)
// 		data := append(header[:], data[pos:pos+dataLen]...)
// 		if moduloFour := len(data) % 4; moduloFour > 0 {
// 			numAlignment := 4 - moduloFour
// 			alignment := bytes.Repeat([]byte{0x00}, numAlignment)
// 			data = append(data, alignment...)
// 		}
// 		_, err := d.usbDevice.Write(data)
// 		if err != nil {
// 			return pos, err
// 		}
// 		pos += dataLen
// 	}
// 	return len(data), nil
// }

func (d *UsbTmc) Write(data []byte) (int, error) {
	d.BTag = (d.BTag % 255) + 1
	header := encodeBulkOutHeader(d.BTag, uint32(len(data)), true)
	packet := append(header[:], data...)

	if len(data)%4 != 0 {
		packet = append(packet, make([]byte, 4-len(data)%4)...)
	}

	return d.BareUsbDev.BulkOutEndpoint.Write(packet)
}

// Read creates and sends the header on the bulk out endpoint and then reads
// from the bulk in endpoint per USBTMC standard.
func (d *UsbTmc) Read() ([]byte, error) {
	buf := make([]byte, d.BareUsbDev.BulkInMaxPktSize*8)
	d.BTag = nextbTag(d.BTag)
	reqInMsg := encodeMsgInBulkOutHeader(d.BTag, uint32(len(buf)), d.TermCharEnabled, d.TermChar)
	if _, err := d.BareUsbDev.Write(reqInMsg[:]); err != nil {
		return []byte{}, err
	}

	n, err := d.BareUsbDev.BulkInEndpoint.Read(buf)
	if err != nil || n < 12 {
		return []byte{}, err
	}

	header, buf := buf[:12], buf[12:]
	// log.Printf("header: %x", header)
	if header[0] != byte(MsgIdDevDepMsgIn) {
		return []byte{}, fmt.Errorf("not a valid read response: %x", header)
	}

	transferSize := int(binary.LittleEndian.Uint32(header[4:8]))
	if transferSize > len(buf)-headerSize {
		reqBytes := transferSize - (len(buf) - headerSize)
		reqPktSize := reqBytes/d.BareUsbDev.BulkInMaxPktSize + 1
		extraBuf := make([]byte, reqPktSize*d.BareUsbDev.BulkInMaxPktSize)
		n, err := d.BareUsbDev.BulkInEndpoint.Read(extraBuf)
		if err != nil {
			return []byte{}, err
		}
		return append(buf, extraBuf[:n]...), nil
	}

	if d.TermCharEnabled {
		termPos := bytes.Index(buf, []byte{d.TermChar})
		if termPos < n && termPos > 0 {
			return buf[:termPos], nil
		}
	}
	return buf[:n], nil
}

// func (d *UsbTmc) readRemoveHeader(p []byte) (n int, transfer int, err error) {
// 	// FIXME(mdr): Seems like I shouldn't use 512 as a magic number or as a hard
// 	// size limit. I should grab the max size of the bulk in endpoint.
// 	usbtmcHeaderLen := 12
// 	temp := make([]byte, 512)
// 	n, err = d.bareUsbDev.Read(temp)
// 	log.Printf("ReadRemoveHeader %x\n", temp[:n])
// 	log.Printf("ReadRemoveHeader %d\n", n)
// 	// Remove the USBMTC Bulk-IN Header from the data and the number of bytes
// 	if n < usbtmcHeaderLen {
// 		return 0, 0, err
// 	}
// 	t32 := binary.LittleEndian.Uint32(temp[4:8])
// 	transfer = int(t32)
// 	reader := bytes.NewReader(temp)
// 	_, err = reader.ReadAt(p, int64(usbtmcHeaderLen))

// 	if err != nil && err != io.EOF {
// 		return n - usbtmcHeaderLen, transfer, err
// 	}
// 	return n - usbtmcHeaderLen, transfer, nil
// }

// func (d *UsbTmc) readKeepHeader(p []byte) (n int, err error) {
// 	return d.bareUsbDev.Read(p)
// }

// Close closes the underlying USB device.
func (d *UsbTmc) Close() error {
	return d.BareUsbDev.Close()
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
