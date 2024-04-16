package main_test

// import (
// 	"log"
// 	"strings"

// 	"testing"

// 	"github.com/grvstick/usbtmc/usbtmc"
// 	"github.com/grvstick/usbtmc/visa"
// )

// // const addr string = "USB0::0xF4EC::0x1631::SDL13GCX4R0117::INSTR"
// // const addr string = "USB0::1155::30016::SPD3ECAD2R1470::0::INSTR"

// const serialNumber = "MY"

// func TestDAQ970(t *testing.T) {
// 	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
 
// 	var inst *usbtmc.UsbTmc
// 	var err error
// 	log.Println("Listing Resources")

// 	for _, resource := range visa.ListResources() {
// 		log.Println(resource)
// 		if strings.Contains(resource, serialNumber){
// 			inst, err = visa.OpenResource(resource, '\n')

// 			if err != nil {
// 				log.Fatal(err.Error())
// 			}
// 		}
// 	}

// 	if inst == nil {
// 		log.Fatal("no devices match the serial number")
// 	}

// 	defer inst.Close()
// 	// inst.WriteString("*RST")
// 	// time.Sleep(3 * time.Second)

// 	// initCmds := []string {
// 	// 	"INST:DMM ON",
// 	// 	"CONF:VOLT:DC DEF,(@101:120)",
// 	// 	"VOLT:DC:NPLC 10,(@101:120)",
// 	// 	"ROUT:SCAN (@101:120)",
// 	// }

// 	// for _, cmd := range initCmds {
// 	// 	_, err := inst.WriteString(cmd)
// 	// 	if err != nil {
// 	// 		log.Printf("error during init cmds %s", err)
// 	// 	}
// 	// }
// 	for i := 0; i < 5; i++ {
// 		_, err := inst.Write([]byte("*IDN?"))
// 		if err != nil {
// 			log.Printf("error during writing query %s", err)
// 		}
// 		rsp, err := inst.Read()
// 		if err != nil {
// 			log.Printf("error during reading query response %s", err)
// 		}
// 		log.Println(string(rsp))
// 	}
// }
