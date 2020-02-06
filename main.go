package main

import (
	"log"
	"time"

	MCP "github.com/ardnew/mcp2221a"
)

const (
	mcpPinFETGate = 0
	mcpPinLEDI2C  = 3
)

var (
	mcp *MCP.MCP2221A
)

func main() {

	// open the MCP2221A interface
	if m, err := MCP.New(0, MCP.VID, MCP.PID); nil != err {
		log.Fatalf("MCP.New(): %v", err)
	} else {
		mcp = m
	}
	defer mcp.Close()

	// configure the MOSFET gate pin as GPIO output, default to OFF
	if err := mcp.GPIO.FlashConfig(mcpPinFETGate, 0, MCP.ModeGPIO, MCP.DirOutput); nil != err {
		log.Fatalf("GPIO.FlashConfig(): %v", err)
	}

	// configure the I2C LED pin for dedicated function
	if err := mcp.Alt.LEDI2C.FlashConfig(false); nil != err {
		log.Fatalf("LEDI2C.FlashConfig(): %v", err)
	}

	if err := mcp.Reset(time.Second * 5); nil != err {
		log.Fatalf("Reset(): %v", err)
	}

	// configure the I2C interface
	if err := mcp.I2C.SetConfig(MCP.I2CBaudRate); nil != err {
		log.Fatalf("I2C.SetConfig(): %v", err)
	}

	for i := 0; i < 100; i++ {
		// read the 16-bit data from device ID register (0xFF) from an INA260 power
		// sensor at default slave address (0x40)
		if buf, err := mcp.I2C.ReadReg(0x40, 0xFF, 2); nil != err {
			log.Fatalf("I2C.ReadReg(): %v", err)
		} else {

			// parse the data received, packing it into a 16-bit unsigned int. the
			// INA260 returns data MSB-first.
			var ub uint16 = (uint16(buf[0]) << 8) | uint16(buf[1])

			rev := ub & 0x0F // revision is 4 bits (LSB)
			die := ub >> 4   // device ID is remaining 12 bits

			log.Printf("Revision  = %3d {0x%4X} [0b%16b]", rev, rev, rev)
			log.Printf("Device ID = %3d {0x%4X} [0b%16b]", die, die, die)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
