package usbtmc_test

import (
	"log"

	"testing"

	"github.com/grvstick/usbtmc"
	"github.com/grvstick/usbtmc/driver"
)

// const addr string = "USB0::0xF4EC::0x1631::SDL13GCX4R0117::INSTR"
// const addr string = "USB0::1155::30016::SPD3ECAD2R1470::0::INSTR"

func TestSiglentLoad(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	drv, err := driver.NewDriver()
	drv.SetDbgLevel(0)

	if err != nil {
		log.Fatal(err)
	}
	defer drv.Close()
	inst := usbtmc.UsbTmc{
		TermChar:        '\n',
		BTag:            5,
		TermCharEnabled: true,
	}

	dev, err := drv.NewDevice(int(0xF4EC), int(0x1631), "SDL13GCX4R0117")
	if err != nil {
		log.Fatal(err)
	}
	defer drv.Close()

	inst.BareUsbDev = *dev
	inst.WriteString("*CLS")
	resp, err := inst.Query("*OPC?")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
	log.Printf("%x", resp)

	for i := 0; i < 5; i++ {
		_, err := inst.Write([]byte("*IDN?"))
		if err != nil {
			log.Fatal(err)
		}
		rsp, err := inst.Read()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(rsp))
	}
}
