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
	temp            uint16
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

// fetch reads the data used by the instruction into the internal
// fetchedData variable. For instructions using the Implied address
// mode, there is no data needed so it is skipped here. It also returns
// the fetched data for convenience.
func (cpu *Mos6502) fetch() uint8 {
	if !(cpu.lookup[cpu.opcode].addressMode == imp) {
		cpu.fetchedData = cpu.bus.Read(cpu.addressAbsolute)
	}
	return cpu.fetchedData
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
	lowByte := cpu.bus.Read(cpu.addressAbsolute)
	highByte := cpu.bus.Read(cpu.addressAbsolute + 1)

	cpu.pc = (uint16(highByte) << 8) | uint16(lowByte)

	cpu.a = 0x00
	cpu.x = 0x00
	cpu.y = 0x00
	cpu.stkp = 0xfd
	cpu.status = 0x00 | uint8(U)

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
		cpu.bus.Write(0x0100+uint16(cpu.stkp), uint8((cpu.pc>>8)&0x00ff))
		cpu.stkp--
		cpu.bus.Write(0x0100+uint16(cpu.stkp), uint8(cpu.pc&0x00ff))
		cpu.stkp--

		cpu.setStatusFlag(B, false)
		cpu.setStatusFlag(U, true)
		cpu.setStatusFlag(I, true)
		cpu.bus.Write(0x0100+uint16(cpu.stkp), cpu.status)
		cpu.stkp--

		cpu.addressAbsolute = 0xfffe
		lowByte := cpu.bus.Read(cpu.addressAbsolute)
		highByte := cpu.bus.Read(cpu.addressAbsolute + 1)
		cpu.pc = (uint16(highByte) << 8) | uint16(lowByte)

		cpu.cycles = 7
	}
}

// NonMaskableInterrupt is the non-maskable interrupt request signal, which
// cannot be ignored. It has the same behavior as the normal Interrupt Request
// but reads 0xfffa to set the program counter.
func (cpu *Mos6502) NonMaskableInterrupt() {
	cpu.bus.Write(0x0100+uint16(cpu.stkp), uint8((cpu.pc>>8)&0x00ff))
	cpu.stkp--
	cpu.bus.Write(0x0100+uint16(cpu.stkp), uint8(cpu.pc&0x00ff))
	cpu.stkp--

	cpu.setStatusFlag(B, false)
	cpu.setStatusFlag(U, true)
	cpu.setStatusFlag(I, true)
	cpu.bus.Write(0x0100+uint16(cpu.stkp), cpu.status)
	cpu.stkp--

	cpu.addressAbsolute = 0xfffa
	lowByte := cpu.bus.Read(cpu.addressAbsolute)
	highByte := cpu.bus.Read(cpu.addressAbsolute + 1)
	cpu.pc = (uint16(highByte) << 8) | uint16(lowByte)

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
	if b&0x80 > 0 {
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

// adc is addition with carry in operation.
func (cpu *Mos6502) adc() uint8 {
	cpu.fetch()

	cpu.temp = uint16(cpu.a) + uint16(cpu.fetchedData) + uint16(cpu.GetStatusFlag(C))

	cpu.setStatusFlag(C, cpu.temp > 255)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0)
	cpu.setStatusFlag(V, ^(uint16(cpu.a)^uint16(cpu.fetchedData))&(uint16(cpu.a)^cpu.temp)&0x0080 > 0)
	cpu.setStatusFlag(N, cpu.temp&0x0080 > 0)

	cpu.a = uint8(cpu.temp & 0x00ff)
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
	cpu.temp = uint16(cpu.fetchedData) << 1
	cpu.setStatusFlag(C, (cpu.temp&0xff00) > 0)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0x00)
	cpu.setStatusFlag(N, cpu.temp&0x80 > 0)
	if cpu.lookup[cpu.opcode].addressMode == imp {
		cpu.a = uint8(cpu.temp & 0x00ff)
	} else {
		cpu.bus.Write(cpu.addressAbsolute, uint8(cpu.temp&0x00ff))
	}
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

// sbc is the subtract with borrow in operation.
func (cpu *Mos6502) sbc() uint8 {
	cpu.fetch()

	value := uint16(cpu.fetchedData) ^ 0x00ff
	cpu.temp = uint16(cpu.a) + value + uint16(cpu.GetStatusFlag(C))

	cpu.setStatusFlag(C, cpu.temp&0xff00 > 0)
	cpu.setStatusFlag(Z, (cpu.temp&0x00ff) == 0)
	cpu.setStatusFlag(V, (cpu.temp^uint16(cpu.a))&(cpu.temp^value)&0x0080 > 0)
	cpu.setStatusFlag(N, cpu.temp&0x0080 > 0)

	cpu.a = uint8(cpu.temp & 0x00ff)
	return 1
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
