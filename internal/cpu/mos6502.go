package cpu

import (
	"github.com/Jac0bDeal/goNES/internal/bus"
)

// Flag is the possible different status flags for the CPU.
type Flag uint8

// Mos6502 Status Flags
const (
	C Flag = 1 << iota // C is the Carry Bit flag.
	Z                  // Z is the Zero flag.
	I                  // I is the Disable Interrupts flag.
	D                  // D is the Decimal Mode flag.
	B                  // B is the Break flag.
	U                  // U is the Unused flag.
	V                  // V is the Overflow flag.
	N                  // N is the Negative flag.
)

// Mos6502 represents a Mos 6502 CPU.
type Mos6502 struct {
	// Core registers
	a      uint8  // accumulator
	x      uint8  // x register
	y      uint8  // y register
	stkp   uint8  // stack pointer
	pc     uint16 // program counter
	status uint8  // status register

	// Bus
	bus *bus.Bus

	// Internal Vars
	fetchedData     uint8
	temp            uint8
	addressAbsolute uint16
	addressRelative uint16
	opcode          uint8
	cycles          uint8
	clockCount      uint32

	// OpCode Lookup Table
	lookup mos6502LookupTable
}

// NewMos6502 constructs and returns a pointer to an instance of Mos6502.
func NewMos6502() *Mos6502 {
	cpu := &Mos6502{}
	cpu.lookup = buildMos502LookupTable(cpu)
	return cpu
}

// ConnectBus connects the CPU to a Bus.
func (cpu *Mos6502) ConnectBus(b *bus.Bus) {
	cpu.bus = b
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Convenience Methods /////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetAccumulator returns the current value of the Accumulator Register.
func (cpu *Mos6502) GetAccumulator() uint8 {
	return cpu.a
}

// GetX returns the current value of the X Register.
func (cpu *Mos6502) GetX() uint8 {
	return cpu.x
}

// GetY returns the current value of the Y Register.
func (cpu *Mos6502) GetY() uint8 {
	return cpu.y
}

// GetStackPointer returns the current value of the Stack Pointer.
func (cpu *Mos6502) GetStackPointer() uint8 {
	return cpu.stkp
}

// GetProgramCounter returns the current value of the Program Counter.
func (cpu *Mos6502) GetProgramCounter() uint16 {
	return cpu.pc
}

// GetStatusFlag returns the current value of specific bit on CPU status register.
func (cpu *Mos6502) GetStatusFlag(f Flag) uint8 {
	if (cpu.status & uint8(f)) > 0 {
		return 1
	}
	return 0
}

// setStatusFlag sets or clears specific bit on the CPU status register.
func (cpu *Mos6502) setStatusFlag(f Flag, v bool) {
	if v {
		cpu.status |= uint8(f)
	} else {
		cpu.status &= ^uint8(f)
	}
}

// fetch ..
func (cpu *Mos6502) fetch() uint8 {
	return 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// External Event Signals //////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Clock is the clock signal.
func (cpu *Mos6502) Clock() {
	// if current instruction complete, read and execute next instruction
	if cpu.cycles == 0 {
		opcode := cpu.bus.Read(cpu.pc)
		instruction := cpu.lookup[opcode]

		cpu.setStatusFlag(U, true)
		cpu.pc++

		cpu.cycles = instruction.cycles
		additionalCycleAddr := instruction.addressMode()
		additionalCycleOp := instruction.operate()
		cpu.cycles += additionalCycleAddr & additionalCycleOp

		cpu.setStatusFlag(U, true)
	}

	cpu.cycles--
}

// Reset is the reset signal.
func (cpu *Mos6502) Reset() {

}

// InterruptRequest is the interrupt request signal.
func (cpu *Mos6502) InterruptRequest() {

}

// NonMaskableInterrupt is the non-maskable interrupt request signal.
func (cpu *Mos6502) NonMaskableInterrupt() {

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Addressing Modes ////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// imp is the Implied address mode.
// On original hardware, signals there is no additional data needed for instruction.
// Code implementation will be used to fetch data from accumulator into
// internal variable.
func (cpu *Mos6502) imp() uint8 {
	cpu.fetchedData = cpu.a
	return 0
}

// imm is the Immediate address mode.
// Signals the instruction needs the next byte after the program counter
// as a value, so we point the point the read address at the next byte.
func (cpu *Mos6502) imm() uint8 {
	cpu.addressAbsolute = cpu.pc + 1
	return 0
}

// zp0 is the Zero Page address mode.
// Allows for reduction in program byte usage by absolutely addressing a
// location in the first address block (0x00 - 0xff), saving a byte of
// usage in the process.
func (cpu *Mos6502) zp0() uint8 {
	b := cpu.bus.Read(cpu.pc)
	cpu.pc++
	cpu.addressAbsolute = uint16(b) & 0x00ff
	return 0
}

// zpx is the Zero Page w/ X Offset address mode.
// Essentially the same as Zero Page addressing, but with an additional
// offset of the read byte by the value in the X register.
func (cpu *Mos6502) zpx() uint8 {
	b := cpu.bus.Read(cpu.pc)
	cpu.pc++
	cpu.addressAbsolute = uint16(b+cpu.x) & 0x00ff
	return 0
}

// zpy is the Zero Page w/ Y Offset address mode.
// The same as Zero Page w/ X offset addressing, but offsetting with the
// Y register instead.
func (cpu *Mos6502) zpy() uint8 {
	b := cpu.bus.Read(cpu.pc)
	cpu.pc++
	cpu.addressAbsolute = uint16(b+cpu.y) & 0x00ff
	return 0
}

// rel is the Relative address mode.
// Exclusive to branch operations, the instruction is addressed within the
// -128 to +127 range of the branched instruction.
func (cpu *Mos6502) rel() uint8 {
	b := cpu.bus.Read(cpu.pc)
	cpu.pc++
	if (b&0x80)>>7 == 1 {
		cpu.addressRelative = uint16(b) | 0xff00
	} else {
		cpu.addressRelative = uint16(b)
	}
	return 0
}

// abs is the Absolute address mode.
// Read and load a full 16-bit address.
func (cpu *Mos6502) abs() uint8 {
	lowByte := cpu.bus.Read(cpu.pc)
	cpu.pc++

	highByte := cpu.bus.Read(cpu.pc)
	cpu.pc++

	cpu.addressAbsolute = (uint16(highByte) << 8) | uint16(lowByte)

	return 0
}

// abx is the Absolute w/ X Offset address mode.
// Same as Absolute addressing, but the read address is offset by the
// value of the X register. If this results in a page change, then
// an additional clock cycle is required and returned.
func (cpu *Mos6502) abx() uint8 {
	lowByte := cpu.bus.Read(cpu.pc)
	cpu.pc++

	highByte := cpu.bus.Read(cpu.pc)
	cpu.pc++

	a := ((uint16(highByte) << 8) | uint16(lowByte)) + uint16(cpu.x)

	cpu.addressAbsolute = a

	if (a & 0xff00) != (uint16(highByte) << 8) {
		return 1
	}

	return 0
}

// abx is the Absolute w/ Y Offset address mode.
// Same as Absolute w/ X Offset addressing, but offsetting with the
// Y register instead.
func (cpu *Mos6502) aby() uint8 {
	lowByte := cpu.bus.Read(cpu.pc)
	cpu.pc++

	highByte := cpu.bus.Read(cpu.pc)
	cpu.pc++

	a := ((uint16(highByte) << 8) | uint16(lowByte)) + uint16(cpu.y)

	cpu.addressAbsolute = a

	if (a & 0xff00) != (uint16(highByte) << 8) {
		return 1
	}

	return 0
}

// ind is the Indirect address mode.
// The address in the program counter is used to read an address from the
// Bus and references that address instead (kinda like a pointer). A bug in
// the original hardware is included here for accuracy, where a page boundary
// is crossed if the low byte read is 0xff. But instead of reading the high
// byte from the next page, the bug causes the start of the same page to be
// read instead and resulting in an invalid address.
func (cpu *Mos6502) ind() uint8 {
	lowByte := cpu.bus.Read(cpu.pc)
	cpu.pc++

	highByte := cpu.bus.Read(cpu.pc)
	cpu.pc++

	pointer := (uint16(highByte) << 8) | uint16(lowByte)

	var a uint16
	if lowByte == 0xff {
		a = (uint16(cpu.bus.Read(pointer&0xFF00)) << 8) | uint16(cpu.bus.Read(pointer))
	} else {
		a = (uint16(cpu.bus.Read(pointer+1)) << 8) | uint16(cpu.bus.Read(pointer))
	}
	cpu.addressAbsolute = a

	return 0
}

// izx is the Indirect X address mode.
// The 8-bit address read from the Bus is offset by the byte in the X
// register to reference an address in page 0x00.
func (cpu *Mos6502) izx() uint8 {
	pa := uint16(cpu.bus.Read(cpu.pc) + cpu.x)
	cpu.pc++

	lowByte := cpu.bus.Read(pa & 0x00ff)
	highByte := cpu.bus.Read((pa + 1) & 0x00ff)

	cpu.addressAbsolute = (uint16(highByte) << 8) | uint16(lowByte)

	return 0
}

// izy is the Indirect Y address mode.
// The 8-bit address read from the Bus is used to address a location
// in page 0x00, where a full 16-bit address is read and offset by the
// value in the Y register. If this results in a page change, then an
// additional clock cycle is required and returned.
func (cpu *Mos6502) izy() uint8 {
	pa := uint16(cpu.bus.Read(cpu.pc))
	cpu.pc++

	lowByte := cpu.bus.Read(pa & 0x00ff)
	highByte := cpu.bus.Read((pa + 1) & 0x00ff)

	a := ((uint16(highByte) << 8) | uint16(lowByte)) + uint16(cpu.y)
	cpu.addressAbsolute = a

	if (a & 0xFF00) != (uint16(highByte) << 8) {
		return 1
	}

	return 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Opcodes /////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (cpu *Mos6502) adc() uint8 {
	return 0
}

func (cpu *Mos6502) and() uint8 {
	return 0
}

func (cpu *Mos6502) asl() uint8 {
	return 0
}

func (cpu *Mos6502) bcc() uint8 {
	return 0
}

func (cpu *Mos6502) bcs() uint8 {
	return 0
}

func (cpu *Mos6502) beq() uint8 {
	return 0
}

func (cpu *Mos6502) bit() uint8 {
	return 0
}

func (cpu *Mos6502) bmi() uint8 {
	return 0
}

func (cpu *Mos6502) bne() uint8 {
	return 0
}

func (cpu *Mos6502) bpl() uint8 {
	return 0
}

func (cpu *Mos6502) brk() uint8 {
	return 0
}

func (cpu *Mos6502) bvc() uint8 {
	return 0
}

func (cpu *Mos6502) bvs() uint8 {
	return 0
}

func (cpu *Mos6502) clc() uint8 {
	return 0
}

func (cpu *Mos6502) cld() uint8 {
	return 0
}

func (cpu *Mos6502) cli() uint8 {
	return 0
}

func (cpu *Mos6502) clv() uint8 {
	return 0
}

func (cpu *Mos6502) cmp() uint8 {
	return 0
}

func (cpu *Mos6502) cpx() uint8 {
	return 0
}

func (cpu *Mos6502) cpy() uint8 {
	return 0
}

func (cpu *Mos6502) dec() uint8 {
	return 0
}

func (cpu *Mos6502) dex() uint8 {
	return 0
}

func (cpu *Mos6502) dey() uint8 {
	return 0
}

func (cpu *Mos6502) eor() uint8 {
	return 0
}

func (cpu *Mos6502) inc() uint8 {
	return 0
}

func (cpu *Mos6502) inx() uint8 {
	return 0
}

func (cpu *Mos6502) iny() uint8 {
	return 0
}

func (cpu *Mos6502) jmp() uint8 {
	return 0
}

func (cpu *Mos6502) jsr() uint8 {
	return 0
}

func (cpu *Mos6502) lda() uint8 {
	return 0
}

func (cpu *Mos6502) ldx() uint8 {
	return 0
}

func (cpu *Mos6502) ldy() uint8 {
	return 0
}

func (cpu *Mos6502) lsr() uint8 {
	return 0
}

func (cpu *Mos6502) nop() uint8 {
	return 0
}

func (cpu *Mos6502) ora() uint8 {
	return 0
}

func (cpu *Mos6502) pha() uint8 {
	return 0
}

func (cpu *Mos6502) php() uint8 {
	return 0
}

func (cpu *Mos6502) pla() uint8 {
	return 0
}

func (cpu *Mos6502) plp() uint8 {
	return 0
}

func (cpu *Mos6502) rol() uint8 {
	return 0
}

func (cpu *Mos6502) ror() uint8 {
	return 0
}

func (cpu *Mos6502) rti() uint8 {
	return 0
}

func (cpu *Mos6502) rts() uint8 {
	return 0
}

func (cpu *Mos6502) sbc() uint8 {
	return 0
}

func (cpu *Mos6502) sec() uint8 {
	return 0
}

func (cpu *Mos6502) sed() uint8 {
	return 0
}

func (cpu *Mos6502) sei() uint8 {
	return 0
}

func (cpu *Mos6502) sta() uint8 {
	return 0
}

func (cpu *Mos6502) stx() uint8 {
	return 0
}

func (cpu *Mos6502) sty() uint8 {
	return 0
}

func (cpu *Mos6502) tax() uint8 {
	return 0
}

func (cpu *Mos6502) tay() uint8 {
	return 0
}

func (cpu *Mos6502) tsx() uint8 {
	return 0
}

func (cpu *Mos6502) txa() uint8 {
	return 0
}

func (cpu *Mos6502) txs() uint8 {
	return 0
}

func (cpu *Mos6502) tya() uint8 {
	return 0
}

func (cpu *Mos6502) xxx() uint8 {
	return 0
}
