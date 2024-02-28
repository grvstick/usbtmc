package usb_test

import (
	"log"
	"testing"

	"github.com/grvstick/usbtmc/usb"
)

func TestDiscovery(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	err := usb.DiscoverUSB()
	log.Fatal(err)
}
