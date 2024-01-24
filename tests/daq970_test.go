package usbtmc_test

import (
	"log"

	"testing"

	"github.com/grvstick/usbtmc"
	"github.com/grvstick/usbtmc/driver"
)

// const addr string = "USB0::0xF4EC::0x1631::SDL13GCX4R0117::INSTR"
// const addr string = "USB0::1155::30016::SPD3ECAD2R1470::0::INSTR"

func TestDAQ970(t *testing.T) {
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

	dev, err := drv.NewDevice(10893, 20737, "MY58014078")
	if err != nil {
		log.Fatal(err)
	}
	defer drv.Close()

	inst.BareUsbDev = *dev
	// inst.WriteString("*RST")
	// time.Sleep(3 * time.Second)

	// initCmds := []string {
	// 	"INST:DMM ON",
	// 	"CONF:VOLT:DC DEF,(@101:120)",
	// 	"VOLT:DC:NPLC 10,(@101:120)",
	// 	"ROUT:SCAN (@101:120)",
	// }

	// for _, cmd := range initCmds {
	// 	_, err := inst.WriteString(cmd)
	// 	if err != nil {
	// 		log.Printf("error during init cmds %s", err)
	// 	}
	// }
	for i := 0; i < 5; i++ {
		_, err := inst.Write([]byte("*IDN?"))
		if err != nil {
			log.Printf("error during writing query %s", err)
		}
		rsp, err := inst.Read()
		if err != nil {
			log.Printf("error during reading query response %s", err)
		}
		log.Println(string(rsp))
	}
}
