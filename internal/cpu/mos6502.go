package cpu

import (
	"math/bits"

	"github.com/Jac0bDeal/goNES/internal/bus"
)

type word uint16

// Flag is the possible different status flags for the CPU.
type Flag byte

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
	a      byte // accumulator
	x      byte // x register
	y      byte // y register
	stkp   byte // stack pointer
	pc     word // program counter
	status byte // status register

	// Bus
	bus *bus.Bus

	// Internal Vars
	fetchedData     byte
	temp            word
	addressAbsolute word
	addressRelative word
	opcode          byte
	cycles          byte
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
func (cpu *Mos6502) GetAccumulator() byte {
	return cpu.a
}

// GetX returns the current value of the X Register.
func (cpu *Mos6502) GetX() byte {
	return cpu.x
}

// GetY returns the current value of the Y Register.
func (cpu *Mos6502) GetY() byte {
	return cpu.y
}

// GetStackPointer returns the current value of the Stack Pointer.
func (cpu *Mos6502) GetStackPointer() byte {
	return cpu.stkp
}

// GetProgramCounter returns the current value of the Program Counter.
func (cpu *Mos6502) GetProgramCounter() uint16 {
	return uint16(cpu.pc)
}

// GetStatusFlag returns the current value of specific bit on CPU status register.
func (cpu *Mos6502) GetStatusFlag(f Flag) byte {
	if (cpu.status & byte(f)) > 0 {
		return 1
	}
	return 0
}

// setStatusFlag sets or clears specific bit on the CPU status register.
func (cpu *Mos6502) setStatusFlag(f Flag, v bool) {
	if v {
		cpu.status |= byte(f)
	} else {
		cpu.status &= ^byte(f)
	}
}

// fetch reads the data used by the instruction into the internal
// fetchedData variable. For instructions using the Implied address
// mode, there is no data needed so it is skipped here. It also returns
// the fetched data for convenience.
func (cpu *Mos6502) fetch() byte {
	if !(cpu.lookup[cpu.opcode].addressMode == imp) {
		cpu.fetchedData = cpu.read(cpu.addressAbsolute)
	}
	return cpu.fetchedData
}

// read reads data from the Bus at the passed address.
func (cpu *Mos6502) read(address word) byte {
	return cpu.bus.Read(uint16(address))
}

// write writes data to the Bus at the passed address.
func (cpu *Mos6502) write(address word, data byte) {
	cpu.bus.Write(uint16(address), data)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// External Event Signals //////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Clock is the clock signal.
func (cpu *Mos6502) Clock() {
	// if current instruction complete, read and execute next instruction
	if cpu.cycles == 0 {
		opcode := cpu.read(cpu.pc)
		instruction := cpu.lookup[opcode]

		cpu.setStatusFlag(U, true)
		cpu.pc++

		cpu.cycles = instruction.cycles
		additionalCycleAddr := instruction.setAddressMode()
		additionalCycleOp := instruction.performOp()
		cpu.cycles += additionalCycleAddr & additionalCycleOp

		cpu.setStatusFlag(U, true)
	}

	cpu.cycles--
}

// Reset signals the cpu to reset to a known state.
func (cpu *Mos6502) Reset() {
	cpu.addressAbsolute = 0xfffc
	lowByte := cpu.read(cpu.addressAbsolute)
	highByte := cpu.read(cpu.addressAbsolute + 1)

	cpu.pc = (word(highByte) << 8) | word(lowByte)

	cpu.a = 0x00
	cpu.x = 0x00
	cpu.y = 0x00
	cpu.stkp = 0xfd
	cpu.status = 0x00 | byte(U)

	cpu.addressRelative = 0x0000
	cpu.addressAbsolute = 0x0000
	cpu.fetchedData = 0x00

	cpu.cycles = 8
}

// InterruptRequest is the interrupt request signal. Requires the Interrupt
// Disable (I) flag to be set to 0 or else nothing happens. The currently running
// instruction is allowed to complete before the Interrupt Request does its thing.
func (cpu *Mos6502) InterruptRequest() {
	if cpu.GetStatusFlag(I) == 0 {
		cpu.write(0x0100+word(cpu.stkp), byte((cpu.pc>>8)&0x00ff))
		cpu.stkp--
		cpu.write(0x0100+word(cpu.stkp), byte(cpu.pc&0x00ff))
		cpu.stkp--

		cpu.setStatusFlag(B, false)
		cpu.setStatusFlag(U, true)
		cpu.setStatusFlag(I, true)
		cpu.write(0x0100+word(cpu.stkp), cpu.status)
		cpu.stkp--

		cpu.addressAbsolute = 0xfffe
		lowByte := cpu.read(cpu.addressAbsolute)
		highByte := cpu.read(cpu.addressAbsolute + 1)
		cpu.pc = (word(highByte) << 8) | word(lowByte)

		cpu.cycles = 7
	}
}

// NonMaskableInterrupt is the non-maskable interrupt request signal, which
// cannot be ignored. It has the same behavior as the normal Interrupt Request
// but reads 0xfffa to set the program counter.
func (cpu *Mos6502) NonMaskableInterrupt() {
	cpu.write(0x0100+word(cpu.stkp), byte((cpu.pc>>8)&0x00ff))
	cpu.stkp--
	cpu.write(0x0100+word(cpu.stkp), byte(cpu.pc&0x00ff))
	cpu.stkp--

	cpu.setStatusFlag(B, false)
	cpu.setStatusFlag(U, true)
	cpu.setStatusFlag(I, true)
	cpu.write(0x0100+word(cpu.stkp), cpu.status)
	cpu.stkp--

	cpu.addressAbsolute = 0xfffa
	lowByte := cpu.read(cpu.addressAbsolute)
	highByte := cpu.read(cpu.addressAbsolute + 1)
	cpu.pc = (word(highByte) << 8) | word(lowByte)

	cpu.cycles = 8
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
	b := cpu.read(cpu.pc)
	cpu.pc++
	cpu.addressAbsolute = word(b) & 0x00ff
	return 0
}

// zpx is the Zero Page w/ X Offset address mode.
// Essentially the same as Zero Page addressing, but with an additional
// offset of the read byte by the value in the X register.
func (cpu *Mos6502) zpx() uint8 {
	b := cpu.read(cpu.pc)
	cpu.pc++
	cpu.addressAbsolute = word(b+cpu.x) & 0x00ff
	return 0
}

// zpy is the Zero Page w/ Y Offset address mode.
// The same as Zero Page w/ X offset addressing, but offsetting with the
// Y register instead.
func (cpu *Mos6502) zpy() uint8 {
	b := cpu.read(cpu.pc)
	cpu.pc++
	cpu.addressAbsolute = word(b+cpu.y) & 0x00ff
	return 0
}

// rel is the Relative address mode.
// Exclusive to branch operations, the instruction is addressed within the
// -128 to +127 range of the branched instruction.
func (cpu *Mos6502) rel() uint8 {
	b := cpu.read(cpu.pc)
	cpu.pc++
	if b&0x80 > 0 {
		cpu.addressRelative = word(b) | 0xff00
	} else {
		cpu.addressRelative = word(b)
	}
	return 0
}

// abs is the Absolute address mode.
// Read and load a full 16-bit address.
func (cpu *Mos6502) abs() uint8 {
	lowByte := cpu.read(cpu.pc)
	cpu.pc++

	highByte := cpu.read(cpu.pc)
	cpu.pc++

	cpu.addressAbsolute = (word(highByte) << 8) | word(lowByte)

	return 0
}

// abx is the Absolute w/ X Offset address mode.
// Same as Absolute addressing, but the read address is offset by the
// value of the X register. If this results in a page change, then
// an additional clock cycle is required and returned.
func (cpu *Mos6502) abx() uint8 {
	lowByte := cpu.read(cpu.pc)
	cpu.pc++

	highByte := cpu.read(cpu.pc)
	cpu.pc++

	a := ((word(highByte) << 8) | word(lowByte)) + word(cpu.x)

	cpu.addressAbsolute = a

	if (a & 0xff00) != (word(highByte) << 8) {
		return 1
	}

	return 0
}

// abx is the Absolute w/ Y Offset address mode.
// Same as Absolute w/ X Offset addressing, but offsetting with the
// Y register instead.
func (cpu *Mos6502) aby() uint8 {
	lowByte := cpu.read(cpu.pc)
	cpu.pc++

	highByte := cpu.read(cpu.pc)
	cpu.pc++

	a := ((word(highByte) << 8) | word(lowByte)) + word(cpu.y)

	cpu.addressAbsolute = a

	if (a & 0xff00) != (word(highByte) << 8) {
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
	lowByte := cpu.read(cpu.pc)
	cpu.pc++

	highByte := cpu.read(cpu.pc)
	cpu.pc++

	pointer := (word(highByte) << 8) | word(lowByte)

	var a word
	if lowByte == 0xff {
		a = (word(cpu.read(pointer&0xFF00)) << 8) | word(cpu.read(pointer))
	} else {
		a = (word(cpu.read(pointer+1)) << 8) | word(cpu.read(pointer))
	}
	cpu.addressAbsolute = a

	return 0
}

// izx is the Indirect X address mode.
// The 8-bit address read from the Bus is offset by the byte in the X
// register to reference an address in page 0x00.
func (cpu *Mos6502) izx() uint8 {
	pa := word(cpu.read(cpu.pc) + cpu.x)
	cpu.pc++

	lowByte := cpu.read(pa & 0x00ff)
	highByte := cpu.read((pa + 1) & 0x00ff)

	cpu.addressAbsolute = (word(highByte) << 8) | word(lowByte)

	return 0
}

// izy is the Indirect Y address mode.
// The 8-bit address read from the Bus is used to address a location
// in page 0x00, where a full 16-bit address is read and offset by the
// value in the Y register. If this results in a page change, then an
// additional clock cycle is required and returned.
func (cpu *Mos6502) izy() uint8 {
	pa := word(cpu.read(cpu.pc))
	cpu.pc++

	lowByte := cpu.read(pa & 0x00ff)
	highByte := cpu.read((pa + 1) & 0x00ff)

	a := ((word(highByte) << 8) | word(lowByte)) + word(cpu.y)
	cpu.addressAbsolute = a

	if (a & 0xFF00) != (word(highByte) << 8) {
		return 1
	}

	return 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Opcodes /////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// adc is addition with carry in operation.
func (cpu *Mos6502) adc() uint8 {
	cpu.fetch()

	cpu.temp = word(cpu.a) + word(cpu.fetchedData) + word(cpu.GetStatusFlag(C))

	cpu.setStatusFlag(C, cpu.temp > 255)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0)
	cpu.setStatusFlag(V, ^(word(cpu.a)^word(cpu.fetchedData))&(word(cpu.a)^cpu.temp)&0x0080 > 0)
	cpu.setStatusFlag(N, cpu.temp&0x0080 > 0)

	cpu.a = byte(cpu.temp & 0x00ff)
	return 1
}

// and performs a bitwise AND on the value in the Accumulator and fetched data.
func (cpu *Mos6502) and() uint8 {
	cpu.fetch()
	cpu.a = cpu.a & cpu.fetchedData
	cpu.setStatusFlag(Z, cpu.a == 0x00)
	cpu.setStatusFlag(N, cpu.a&0x80 > 0)
	return 1
}

// asl performs an arithmetic shift left on data at an address in the Bus.
func (cpu *Mos6502) asl() uint8 {
	cpu.fetch()
	cpu.temp = word(cpu.fetchedData) << 1
	cpu.setStatusFlag(C, (cpu.temp&0xff00) > 0)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x00)
	cpu.setStatusFlag(N, cpu.temp&0x80 > 0)
	if cpu.lookup[cpu.opcode].addressMode == imp {
		cpu.a = byte(cpu.temp & 0x00ff)
	} else {
		cpu.write(cpu.addressAbsolute, byte(cpu.temp&0x00ff))
	}
	return 0
}

// branch is a convenience method containing the common branching logic.
func (cpu *Mos6502) branch() {
	cpu.cycles++
	cpu.addressAbsolute = cpu.pc + cpu.addressRelative

	if (cpu.addressAbsolute & 0xff00) != (cpu.pc & 0xff00) {
		cpu.cycles++
	}
	cpu.pc = cpu.addressAbsolute
}

// bcc is the Branch if Carry Clear operation. If a page change occurs as a
// result, an extra cycle is required.
func (cpu *Mos6502) bcc() uint8 {
	if cpu.GetStatusFlag(C) == 0 {
		cpu.branch()
	}
	return 0
}

// bcs is the Branch if Carry set operation. If a page change occurs as a result,
// an extra cycle is required.
func (cpu *Mos6502) bcs() uint8 {
	if cpu.GetStatusFlag(C) == 1 {
		cpu.branch()
	}
	return 0
}

// beq is the Branch if Equal operation. Does nothing if Z flag set to false.
// Adds additional cycle if page change occurs as result.
func (cpu *Mos6502) beq() uint8 {
	if cpu.GetStatusFlag(Z) == 1 {
		cpu.branch()
	}
	return 0
}

// bit is the Test Bits operation.
func (cpu *Mos6502) bit() uint8 {
	cpu.fetch()
	cpu.temp = word(cpu.a & cpu.fetchedData)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x00)
	cpu.setStatusFlag(N, (cpu.fetchedData&(1<<7)) > 0)
	cpu.setStatusFlag(V, (cpu.fetchedData&(1<<6)) > 0)
	return 0
}

// bmi is the Branch if Negative operation. Does nothing if N flag set to false.
// Adds additional cycle if page change occurs as result.
func (cpu *Mos6502) bmi() uint8 {
	if cpu.GetStatusFlag(N) == 1 {
		cpu.branch()
	}
	return 0
}

// bne is the Branch if Not Equal operation. Does nothing if Z flag set to true.
// Adds additional cycle if page change occurs as result.
func (cpu *Mos6502) bne() uint8 {
	if cpu.GetStatusFlag(Z) == 0 {
		cpu.branch()
	}
	return 0
}

// bpl is the Branch if Positive operation. Does nothing if N flag set to true.
// Adds additional cycle if page change occurs as result.
func (cpu *Mos6502) bpl() uint8 {
	if cpu.GetStatusFlag(N) == 0 {
		cpu.branch()
	}
	return 0
}

// brk is the Break operation. It is used to signal an interrupt from the program.
func (cpu *Mos6502) brk() uint8 {
	cpu.pc++

	cpu.setStatusFlag(I, true)
	cpu.write(0x0100+word(cpu.stkp), byte((cpu.pc>>8)&0x00ff))
	cpu.stkp--
	cpu.write(0x0100+word(cpu.stkp), byte(cpu.pc&0x00ff))
	cpu.stkp--

	cpu.setStatusFlag(B, true)
	cpu.write(0x0100+word(cpu.stkp), cpu.status)
	cpu.stkp--
	cpu.setStatusFlag(B, false)

	cpu.pc = word(cpu.read(0xfffe)) | (word(cpu.read(0xffff)) << 8)
	return 0
}

// bvc is the Branch if Overflow Clear operation. Does nothing if V set to true.
// Adds additional cycle if page change occurs as result.
func (cpu *Mos6502) bvc() uint8 {
	if cpu.GetStatusFlag(V) == 0 {
		cpu.branch()
	}
	return 0
}

func (cpu *Mos6502) bvs() uint8 {
	if cpu.GetStatusFlag(V) == 1 {
		cpu.branch()
	}
	return 0
}

// clc is the Clear Carry Flag operation. It sets the Carry flag to false.
func (cpu *Mos6502) clc() uint8 {
	cpu.setStatusFlag(C, false)
	return 0
}

// cld is the Clear Decimal Flag operation. It sets the Decimal flag to false.
func (cpu *Mos6502) cld() uint8 {
	cpu.setStatusFlag(D, false)
	return 0
}

// cli is the Clear Interrupt Flag operation. It sets the Interrupt flag to false.
func (cpu *Mos6502) cli() uint8 {
	cpu.setStatusFlag(I, false)
	return 0
}

// clv is the Clear Overflow Flag operation. It sets the Overflow flag to false.
func (cpu *Mos6502) clv() uint8 {
	cpu.setStatusFlag(V, false)
	return 0
}

// compare is a convenience method containing the common logic used by the compare
// operations.
func (cpu *Mos6502) compare(registerData byte) {
	cpu.fetch()
	cpu.temp = word(registerData) - word(cpu.fetchedData)
	cpu.setStatusFlag(C, registerData >= cpu.fetchedData)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x0000)
	cpu.setStatusFlag(N, (cpu.temp&0x0080) > 0)
}

// cmp is the Compare Accumulator operation. It compares the accumulator to data
// stored on the Bus, setting the C, N, and Z flags accordingly.
func (cpu *Mos6502) cmp() uint8 {
	cpu.compare(cpu.a)
	return 1
}

// cpx is the Compare X Register operation. It compares the X Register to data
// stored on the Bus, setting the C, N, and Z flags accordingly.
func (cpu *Mos6502) cpx() uint8 {
	cpu.compare(cpu.x)
	return 0
}

// cpy is the Compare Y Register operation. It compares the Y Register to data
// stored on the Bus, setting the C, N, and Z flags accordingly.
func (cpu *Mos6502) cpy() uint8 {
	cpu.compare(cpu.y)
	return 0
}

// dec is the Decrement Value at Memory Location operation.
func (cpu *Mos6502) dec() uint8 {
	cpu.fetch()
	cpu.temp = word(cpu.fetchedData) - 1
	cpu.write(cpu.addressAbsolute, byte(cpu.temp&0x00ff))
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x0000)
	cpu.setStatusFlag(N, (cpu.temp&0x0080) > 0)
	return 0
}

// dex is the Decrement X Register operation.
func (cpu *Mos6502) dex() uint8 {
	cpu.x--
	cpu.setStatusFlag(Z, cpu.x == 0x00)
	cpu.setStatusFlag(N, (cpu.x&0x80) > 0)
	return 0
}

// dey is the Decrement Y Register operation.
func (cpu *Mos6502) dey() uint8 {
	cpu.y--
	cpu.setStatusFlag(Z, cpu.y == 0x00)
	cpu.setStatusFlag(N, (cpu.y&0x80) > 0)
	return 0
}

// eor is the Exclusive Or (XOR) operation.
func (cpu *Mos6502) eor() uint8 {
	cpu.fetch()
	cpu.a = cpu.a ^ cpu.fetchedData
	cpu.setStatusFlag(Z, cpu.a == 0x00)
	cpu.setStatusFlag(N, (cpu.a&0x80) > 0)
	return 0
}

// inc is the Increment Value at Memory Location operation.
func (cpu *Mos6502) inc() uint8 {
	cpu.fetch()
	cpu.temp = word(cpu.fetchedData) + 1
	cpu.write(cpu.addressAbsolute, byte(cpu.temp&0x00ff))
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x0000)
	cpu.setStatusFlag(N, (cpu.temp&0x80) > 0)
	return 0
}

// inx is the Increment X Register operation.
func (cpu *Mos6502) inx() uint8 {
	cpu.x++
	cpu.setStatusFlag(Z, cpu.x == 0x00)
	cpu.setStatusFlag(N, (cpu.x&0x80) > 0)
	return 0
}

// iny is the Increment Y Register operation.
func (cpu *Mos6502) iny() uint8 {
	cpu.y++
	cpu.setStatusFlag(Z, cpu.y == 0x00)
	cpu.setStatusFlag(N, (cpu.y&0x80) > 0)
	return 0
}

// jmp is the Jump to Location operation.
func (cpu *Mos6502) jmp() uint8 {
	cpu.pc = cpu.addressAbsolute
	return 0
}

// jsr is the Jump to Sub-Routine operation. Pushes current program counter value
// to stack.
func (cpu *Mos6502) jsr() uint8 {
	cpu.pc--

	cpu.write(0x0100+word(cpu.stkp), byte((cpu.pc>>8)&0x00ff))
	cpu.stkp--
	cpu.write(0x0100+word(cpu.stkp), byte(cpu.pc&0x00ff))
	cpu.stkp--

	cpu.pc = cpu.addressAbsolute
	return 0
}

// lda is the Load Accumulator operation.
func (cpu *Mos6502) lda() uint8 {
	cpu.fetch()
	cpu.a = cpu.fetchedData
	cpu.setStatusFlag(Z, cpu.a == 0x00)
	cpu.setStatusFlag(N, (cpu.a&0x80) > 0)
	return 1
}

// lda is the Load X Register operation.
func (cpu *Mos6502) ldx() uint8 {
	cpu.fetch()
	cpu.x = cpu.fetchedData
	cpu.setStatusFlag(Z, cpu.x == 0x00)
	cpu.setStatusFlag(N, (cpu.x&0x80) > 0)
	return 1
}

// lda is the Load Y Register operation.
func (cpu *Mos6502) ldy() uint8 {
	cpu.fetch()
	cpu.y = cpu.fetchedData
	cpu.setStatusFlag(Z, cpu.y == 0x00)
	cpu.setStatusFlag(N, (cpu.y&0x80) > 0)
	return 1
}

// lsr is the Logical Shift Right operation. Shifts all bits to the right by one,
// shifting original bit 0 into carry flag and carrying to bit 7.
func (cpu *Mos6502) lsr() uint8 {
	cpu.fetch()
	cpu.setStatusFlag(C, (cpu.fetchedData&0x0001) > 0)
	cpu.temp = word(bits.RotateLeft8(cpu.fetchedData, -1))
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x0000)
	cpu.setStatusFlag(N, (cpu.temp&0x0080) > 0)
	if cpu.lookup[cpu.opcode].addressMode == imp {
		cpu.a = byte(cpu.temp & 0x00ff)
	} else {
		cpu.write(cpu.addressAbsolute, byte(cpu.temp&0x00ff))
	}
	return 0
}

// nop is the No Operation operation. Adds additional cycle in some cases.
func (cpu *Mos6502) nop() uint8 {
	switch cpu.opcode {
	case 0x1c, 0x3c, 0x5c, 0x7c, 0xdc, 0xfc:
		return 1
	default:
		return 0
	}
}

// ora is the Bitwise Logic OR operation.
func (cpu *Mos6502) ora() uint8 {
	cpu.fetch()
	cpu.a = cpu.a | cpu.fetchedData
	cpu.setStatusFlag(Z, cpu.a == 0x00)
	cpu.setStatusFlag(N, (cpu.a&0x80) > 0)
	return 1
}

// pha is the Push Accumulator to Stack operation.
func (cpu *Mos6502) pha() uint8 {
	cpu.write(0x0100+word(cpu.stkp), cpu.a)
	cpu.stkp--
	return 0
}

// php is the Push Status Register to Stack operation.
func (cpu *Mos6502) php() uint8 {
	cpu.write(0x0100+word(cpu.stkp), cpu.status|byte(B)|byte(U))
	cpu.setStatusFlag(B, false)
	cpu.setStatusFlag(U, false)
	cpu.stkp--
	return 0
}

// pla is the Pop Accumulator Off Stack operation.
func (cpu *Mos6502) pla() uint8 {
	cpu.stkp++
	cpu.a = cpu.read(0x0100 + word(cpu.stkp))
	cpu.setStatusFlag(Z, cpu.a == 0x00)
	cpu.setStatusFlag(N, (cpu.a&0x80) > 0)
	return 0
}

// plp is the Pop Status Register Off Stack operation.
func (cpu *Mos6502) plp() uint8 {
	cpu.stkp++
	cpu.status = cpu.read(0x0100 + word(cpu.stkp))
	cpu.setStatusFlag(U, true)
	return 0
}

// rol is the Rotate Left operation.
func (cpu *Mos6502) rol() uint8 {
	cpu.fetch()
	cpu.temp = word(cpu.fetchedData)<<1 | word(cpu.GetStatusFlag(C))
	cpu.setStatusFlag(C, (cpu.fetchedData&0x80) > 0)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x0000)
	cpu.setStatusFlag(N, (cpu.temp&0x0080) > 0)
	if cpu.lookup[cpu.opcode].addressMode == imp {
		cpu.a = byte(cpu.temp & 0x00ff)
	} else {
		cpu.write(cpu.addressAbsolute, byte(cpu.temp&0x00ff))
	}
	return 0
}

// rol is the Rotate Right operation.
func (cpu *Mos6502) ror() uint8 {
	cpu.fetch()
	cpu.temp = word(cpu.GetStatusFlag(C))<<7 | word(cpu.fetchedData)>>1
	cpu.setStatusFlag(C, (cpu.fetchedData&0x01) > 0)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x0000)
	cpu.setStatusFlag(N, (cpu.temp&0x0080) > 0)
	if cpu.lookup[cpu.opcode].addressMode == imp {
		cpu.a = byte(cpu.temp & 0x00ff)
	} else {
		cpu.write(cpu.addressAbsolute, byte(cpu.temp&0x00ff))
	}
	return 0
}

// rti is the Return from Interrupt operation.
func (cpu *Mos6502) rti() uint8 {
	cpu.stkp++
	cpu.status = cpu.read(0x0100 + word(cpu.stkp))
	cpu.status &= ^uint8(B)
	cpu.status &= ^uint8(U)

	cpu.stkp++
	cpu.pc = word(cpu.read(0x0100 + word(cpu.stkp)))
	cpu.stkp++
	cpu.pc |= word(cpu.read(0x0100+word(cpu.stkp))) << 8
	return 0
}

// rts is the Return from Subroutine operation.
func (cpu *Mos6502) rts() uint8 {
	cpu.stkp++
	cpu.pc = word(cpu.read(0x0100 + word(cpu.stkp)))
	cpu.stkp++
	cpu.pc |= word(cpu.read(0x0100+word(cpu.stkp))) << 8

	cpu.pc++
	return 0
}

// sbc is the subtract with borrow in operation.
func (cpu *Mos6502) sbc() uint8 {
	cpu.fetch()

	value := word(cpu.fetchedData) ^ 0x00ff
	cpu.temp = word(cpu.a) + value + word(cpu.GetStatusFlag(C))

	cpu.setStatusFlag(C, (cpu.temp&0xff00)>>7 == 1)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0)
	cpu.setStatusFlag(V, (cpu.temp^word(cpu.a))&(cpu.temp^value)&0x0080 > 0)
	cpu.setStatusFlag(N, (cpu.temp&0x0080) > 0)

	cpu.a = byte(cpu.temp & 0x00ff)
	return 1
}

// sec is the Set Carry Flag operation.
func (cpu *Mos6502) sec() uint8 {
	cpu.setStatusFlag(C, true)
	return 0
}

// sed is the Set Decimal Flag operation.
func (cpu *Mos6502) sed() uint8 {
	cpu.setStatusFlag(D, true)
	return 0
}

// sei is the Set Interrupt Flag operation.
func (cpu *Mos6502) sei() uint8 {
	cpu.setStatusFlag(I, true)
	return 0
}

// sta is the Store Accumulator at Address operation.
func (cpu *Mos6502) sta() uint8 {
	cpu.write(cpu.addressAbsolute, cpu.a)
	return 0
}

// stx is the Store X Register at Address operation.
func (cpu *Mos6502) stx() uint8 {
	cpu.write(cpu.addressAbsolute, cpu.x)
	return 0
}

// sty is the Store Y Register at Address operation.
func (cpu *Mos6502) sty() uint8 {
	cpu.write(cpu.addressAbsolute, cpu.y)
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

// xxx represents an unofficial opcode.
func (cpu *Mos6502) xxx() uint8 {
	return 0
}
