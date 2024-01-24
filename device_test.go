package usbtmc

import (
	"log"

	"testing"

	"github.com/grvstick/usbtmc/driver"
)

// const addr string = "USB0::0xF4EC::0x1631::SDL13GCX4R0117::INSTR"
// const addr string = "USB0::1155::30016::SPD3ECAD2R1470::0::INSTR"

func TestDevice(t *testing.T) {
	drv, err := driver.NewDriver()
	drv.SetDbgLevel(0)

	if err != nil {
		log.Fatal(err)
	}
	defer drv.Close()
	inst := UsbTmc{
		TermChar:        '\n',
		BTag:            5,
		TermCharEnabled: true,
	}

	dev, err := drv.NewDevice(int(0xF4EC), int(0x1631), "SDL13GCX4R0117")
	// dev, err := drv.NewDevice(int(1155), int(30016), "SPD3ECAD2R1470")
	if err != nil {
		log.Fatal(err)
	}
	defer drv.Close()

	inst.BareUsbDev = *dev
	// inst.Command("*RST")
	// resp, err := inst.Query("*OPC?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(resp)
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
		log.Printf("%x", rsp)
	}
	// inst.WriteString("*IDN?")
	// buf := make([]byte, 100)
	// _, err = inst.Read(buf)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(buf))
}