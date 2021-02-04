package main

import (
	"fmt"

	"github.com/Jac0bDeal/goNES/internal/bus"
	"github.com/Jac0bDeal/goNES/internal/cpu"
)

func main() {
	var (
		programStart  uint16 = 0x8000
		programEnd    uint16
		programOffset = programStart
	)

	// load program
	program := []byte{0xA2, 0x0A, 0x8E, 0x00, 0x00, 0xA2, 0x03, 0x8E, 0x01, 0x00, 0xAC, 0x00, 0x00, 0xA9, 0x00, 0x18, 0x6D, 0x01, 0x00, 0x88, 0xD0, 0xFA, 0x8D, 0x02, 0x00, 0xEA, 0xEA, 0xEA}
	r := bus.RAM{}
	for _, b := range program {
		r[programOffset] = b
		programOffset++
	}
	programEnd = programOffset

	// set reset pointer
	r[0xFFFC] = byte(programStart & 0x00ff)
	r[0xFFFD] = byte((programStart & 0xff00) >> 8)

	// initialize bus and cpu
	b := bus.NewBus(r)
	c := cpu.NewMos6502()
	c.ConnectBus(b)

	// disassemble program and print
	asmMap := c.Disassemble(programStart, programEnd)
	for a := programStart; a < programEnd; a++ {
		asmLine, exists := asmMap[a]
		if exists {
			fmt.Println(asmLine)
		}
	}

	// reset cpu
	c.Reset()
}
