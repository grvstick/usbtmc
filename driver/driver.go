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

// Driver models libusb context and implements the driver.Driver interface.
type Driver struct {
	ctx *gousb.Context
}

// NewDriver creates a new libusb session/context.
func NewDriver() (*Driver, error) {
	d := Driver{
		ctx: gousb.NewContext(),
	}
	return &d, nil
}

func (drv *Driver) Close() {
	drv.ctx.Close()
}

func (drv *Driver) SetDbgLevel(lv int) {
	drv.ctx.Debug(lv)
}

// NewDevice searches for device matching vid, pid and serial number. Serial number can be omitted by passing empty string
// If serial number is omitted it will look for first device matching vid & pid
// If a device is detected, it will go over configurations to see if thare is a TMC configuration.
func (drv *Driver) NewDevice(vid, pid int, sn string) (*BareUsbDevice, error) {
	// Iterate through available devices. Find all devices that match the given
	// Vendor ID and Product ID.
	gVID, gPID := gousb.ID(vid), gousb.ID(pid)
	devs, err := drv.ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// This anonymous function is called for every device present. Returning
		// true means the device should be opened.
		return desc.Vendor == gVID && desc.Product == gPID
	})
	if err != nil {
		// Close all devices and return error.
		for _, d := range devs {
			// I'm ignoring any errors on close at the moment.
			d.Close()
		}
		return nil, err
	}
	// if len(devs) == 0 && vidRaw == 0x0957 {
	// 	// Didn't find a device matching the given vendor ID and product ID. Need
	// 	// to make sure we're not trying to find one of the Agilent/Keysight (VID
	// 	// 2391 = 0x0957) USB modular test equipment that powers up in a firmware
	// 	// update mode. In order to get out of boot mode and into the normal USBTMC
	// 	// mode, some control commands need to be sent. Once in normal USBTMC mode,
	// 	// the Product ID will change.
	// 	bootPIDs := map[gousb.ID]gousb.ID{
	// 		// usbtmcPID: bootPID
	// 		0x2818: 0x2918, // U2702A 200 MHz Oscilloscope
	// 		0x3D18: 0x3E18, // U2751A 4x8 2-wire Switch Matrix
	// 		0x4118: 0x4218, // U2722A Source Measure Unit
	// 		0x4318: 0x4418, // U2723A Source Measure Unit
	// 	}
	// 	if bootPID, ok := bootPIDs[pid]; ok {
	// 		// Iterate through available USB devices. Find all devices that match the
	// 		// Keysight USB modular boot PID.
	// 		devs, err = drv.ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
	// 			return desc.Vendor == vid && desc.Product == bootPID
	// 		})
	// 		if err != nil {
	// 			// Close all devices and return error.
	// 			for _, d := range devs {
	// 				// I'm ignoring any errors on close at the moment.
	// 				d.Close()
	// 			}
	// 			return nil, err
	// 		}
	// 		if len(devs) == 0 {
	// 			return nil, fmt.Errorf("no devices found matching VID %s and PID %s", vid, pid)
	// 		}
	// 		// Found a Keysight USB modular device, so exit boot mode.
	// 		err = exitBootMode(devs[0], bootPID)

	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		// Now find the normal USBTMC mode.
	// 		devs, err = drv.ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
	// 			// This anonymous function is called for every device present. Returning
	// 			// true means the device should be opened.
	// 			return desc.Vendor == vid && desc.Product == pid
	// 		})
	// 		if err != nil {
	// 			// Close all devices and return error.
	// 			for _, d := range devs {
	// 				// I'm ignoring any errors on close at the moment.
	// 				d.Close()
	// 			}
	// 			return nil, err
	// 		}
	// 		if len(devs) == 0 {
	// 			return nil, fmt.Errorf("no devices found after reboot matching VID %s and PID %s", vid, pid)
	// 		}
	// 	}
	// } else if len(devs) == 0 {
	// 	return nil, fmt.Errorf("no devices found matching VID %s and PID %s", vid, pid)
	// }

	// Close all except the first returned device.
	// for i, d := range devs {
	// 	if i != 0 {
	// 		d.Close()
	// 	}
	// }

	// pick the device with matching serial number
	var dev *gousb.Device

	// empty string implies picking the first one
	if len(sn) == 0 {
		for i, d := range devs {
			if i == 0 {
				dev = d
			} else {
				d.Close()
			}
		}
	} else {
		for _, d := range devs {
			serial, err := d.SerialNumber()
			if err == nil && serial == sn {
				dev = d
			} else {
				d.Close()
			}
		}

	}
	// There are cases not being able to claim USB device without this code
	dev.SetAutoDetach(true)

	if dev == nil {
		return nil, fmt.Errorf("no devices found matching vid(%v), pid(%v), sn(%v)", vid, pid, sn)
	}

	// Switch to configuration #0
	activeCfg, err := dev.ActiveConfigNum()
	if err != nil {
		return nil, err
	}
	cfg, err := dev.Config(activeCfg)
	if err != nil {
		return nil, err
	}

	// Loop through the interfaces
	for _, ifDesc := range cfg.Desc.Interfaces {
		for _, alt := range ifDesc.AltSettings {
			isTmc, proto := checkTMC(alt)
			if isTmc {
				intf, err := cfg.Interface(ifDesc.Number, alt.Number)
				if err != nil {
					log.Println(err)
					return nil, err
				}

				return tryGetUsbDevice(cfg, dev, intf, proto)
			}
		}
	}

	return nil, fmt.Errorf("target device has no TMC class")
}

func tryGetUsbDevice(cfg *gousb.Config, dev *gousb.Device, intf *gousb.Interface, proto gousb.Protocol) (*BareUsbDevice, error) {
	var bulkIn *gousb.InEndpoint
	var bulkOut *gousb.OutEndpoint
	var intIn *gousb.InEndpoint
	var err error

	for _, ep := range intf.Setting.Endpoints {
		if ep.Direction == gousb.EndpointDirectionOut && ep.TransferType == gousb.TransferTypeBulk {
			bulkOut, err = intf.OutEndpoint(ep.Number)
			if err != nil {
				return nil, err
			}
		}
		if ep.Direction == gousb.EndpointDirectionIn && ep.TransferType == gousb.TransferTypeBulk {
			bulkIn, err = intf.InEndpoint(ep.Number)
			if err != nil {
				return nil, err
			}
		}
		if ep.Direction == gousb.EndpointDirectionIn && ep.TransferType == gousb.TransferTypeInterrupt && proto == prot488 {
			intIn, err = intf.InEndpoint(ep.Number)
			if err != nil {
				return nil, err
			}
		}
	}

	return &BareUsbDevice{
		dev:                 dev,
		intf:                intf,
		cfg:                 cfg,
		BulkInEndpoint:      bulkIn,
		BulkOutEndpoint:     bulkOut,
		InterruptInEndpoint: intIn,
	}, nil
}

func checkTMC(val gousb.InterfaceSetting) (bool, gousb.Protocol) {
	//based on http://www.linux-usb.org/usb.ids
	var (
		class, sub gousb.Class
		proto      gousb.Protocol
	)
	class, sub, proto = val.Class, val.SubClass, val.Protocol

	isTmc := class == baseClsAppSpecific && sub == subClsTmc

	return isTmc, proto
}

func (drv *Driver) NewDeviceFromVisaString(addr string) (*BareUsbDevice, error) {
	v, err := NewVisaResource(addr)
	if err != nil {
		return nil, err
	}
	return drv.NewDevice(v.manufacturerID, v.modelCode, v.serialNumber)
}

// func exitBootMode(dev *gousb.Device, bootPID gousb.ID) error {
// 	thirdIndex := uint16(0x0487)
// 	if bootPID == 0x2818 || bootPID == 0x3E18 {
// 		thirdIndex = 0x0484
// 	}
// 	bRequest := uint8(0x0C)
// 	value := uint16(0x0000)
// 	packets := []struct {
// 		bmRequestType uint8
// 		index         uint16
// 		data          []byte
// 	}{
// 		{0xC0, 0x047E, make([]byte, 0x01)},
// 		{0xC0, 0x047D, make([]byte, 0x06)},
// 		{0xC0, thirdIndex, make([]byte, 0x05)},
// 		{0xC0, 0x0472, make([]byte, 0x0C)},
// 		{0xC0, 0x047A, make([]byte, 0x01)},
// 		{0x40, 0x0475, []byte{0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x08, 0x01}},
// 	}
// 	for i, packet := range packets {
// 		_, err := dev.Control(
// 			packet.bmRequestType,
// 			bRequest,
// 			value,
// 			packet.index,
// 			packet.data,
// 		)
// 		if err != nil {
// 			return fmt.Errorf("error sending control transfer #%d: %s", i+1, err)
// 		}
// 	}

// 	// We need to wait for the USB device to exit boot mode and reboot in normal
// 	// mode.
// 	rebootDelay := time.Second * 7
// 	time.Sleep(rebootDelay)
// 	return nil
// }
