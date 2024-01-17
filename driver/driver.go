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
			log.Printf("Bulk out: %#v", bulkOut.Desc)
		}
		if ep.Direction == gousb.EndpointDirectionIn && ep.TransferType == gousb.TransferTypeBulk {
			bulkIn, err = intf.InEndpoint(ep.Number)
			if err != nil {
				return nil, err
			}
			log.Printf("Bulk in: %#v", bulkIn.Desc)
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
		BulkInMaxPktSize:    bulkIn.Desc.MaxPacketSize,
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
