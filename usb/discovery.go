// Copyright (c) 2015-2020 The usbtmc developers. All rights reserved.
// Project site: https://github.com/gotmc/usbtmc
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package usb

import (
	"log"

	"github.com/google/gousb"
)


// NewDevice searches for device matching vid, pid and serial number. Serial number can be omitted by passing empty string
// If serial number is omitted it will look for first device matching vid & pid
// If a device is detected, it will go over configurations to see if thare is a TMC configuration.
func DiscoverUSB() error{
	// Iterate through available devices. Find all devices that match the given
	// Vendor ID and Product ID.
	ctx := gousb.NewContext()
	defer ctx.Close()

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		// This anonymous function is called for every device present. Returning
		// true means the device should be opened.
		return true
	})
	if err != nil {
		log.Println(err.Error())
		for _, d := range devs {
			// I'm ignoring any errors on close at the moment.
			defer d.Close()
		}
	}


	for _, dev := range devs {
		// resource := ""

		log.Printf("%#v", dev.Desc)

		sn, err := dev.SerialNumber()
		if err != nil {
			log.Println(sn)
		}
		activeCfg, err := dev.ActiveConfigNum()
		if err != nil {
			return err
		}
		cfg, err := dev.Config(activeCfg)
		if err != nil {
			return err
		}
		for _, ifDesc := range cfg.Desc.Interfaces {
			for _, alt := range ifDesc.AltSettings {
				isTmc, proto := checkTMC(alt)
				log.Printf("Proto: %v", proto)
				if isTmc {
					intf, err := cfg.Interface(ifDesc.Number, alt.Number)
					if err != nil {
						log.Println(err)
						return err
					}
	
					log.Printf("%#v", intf)
				}
			}
		}
	
	}
	// Switch to configuration #0


	// Loop through the interfaces

	return nil
}
