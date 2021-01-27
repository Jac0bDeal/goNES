package cpu

import "github.com/Jac0bDeal/goNES/internal/bus"

// Flag is the possible different status flags for the CPU.
type Flag uint8

// Mos6502 Status Flags
const (
	// C is the Carry Bit flag.
	C Flag = 1 << iota

	// Z is the Zero flag.
	Z

	// I is the Disable Interrupts flag.
	I

	// D is the Decimal Mode flag.
	D

	// B is the Break flag.
	B

	// U is the Unused flag.
	U

	// V is the Overflow flag.
	V

	// N is the Negative flag.
	N
)

// Mos6502 represents a Mos 6502 CPU.
type Mos6502 struct {
	// Core registers
	a      uint8
	x      uint8
	y      uint8
	stkp   uint8
	pc     uint16
	status uint8

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
	lookup [256]instruction
}

// NewMos6502 constructs and returns an instance of Mos6502.
func NewMos6502() Mos6502 {
	cpu := Mos6502{}
	cpu.lookup = buildLookupTable(cpu)
	return cpu
}

// ConnectBus connects the CPU to a Bus.
func (c Mos6502) ConnectBus(b *bus.Bus) {
	c.bus = b
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Convenience Methods /////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// GetAccumulator returns the current value of the Accumulator Register.
func (c Mos6502) GetAccumulator() uint8 {
	return c.a
}

// GetX returns the current value of the X Register.
func (c Mos6502) GetX() uint8 {
	return c.x
}

// GetY returns the current value of the Y Register.
func (c Mos6502) GetY() uint8 {
	return c.y
}

// GetStackPointer returns the current value of the Stack Pointer.
func (c Mos6502) GetStackPointer() uint8 {
	return c.stkp
}

// GetProgramCounter returns the current value of the Program Counter.
func (c Mos6502) GetProgramCounter() uint16 {
	return c.pc
}

// GetStatus returns the current CPU status flag.
func (c Mos6502) GetStatus(f Flag) uint8 {
	return c.status
}

// setStatus sets the CPU status flag.
func (c Mos6502) setStatus(f Flag, v bool) {

}

// fetch ..
func (c Mos6502) fetch() uint8 {
	return 0
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// External Event Signals //////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Clock is the clock signal.
func (c Mos6502) Clock() {

}

// Reset is the reset signal.
func (c Mos6502) Reset() {

}

// InterruptRequest is the interrupt request signal.
func (c Mos6502) InterruptRequest() {

}

// NonMaskableInterrupt is the non-maskable interrupt request signal.
func (c Mos6502) NonMaskableInterrupt() {

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Addressing Modes ////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (c Mos6502) imp() {

}

func (c Mos6502) imm() {

}

func (c Mos6502) zp0() {

}

func (c Mos6502) zpx() {

}

func (c Mos6502) zpy() {

}

func (c Mos6502) rel() {

}

func (c Mos6502) abs() {

}

func (c Mos6502) abx() {

}

func (c Mos6502) aby() {

}

func (c Mos6502) ind() {

}

func (c Mos6502) izx() {

}

func (c Mos6502) izy() {

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Opcodes /////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (c Mos6502) adc() uint8 {
	return 0
}

func (c Mos6502) and() uint8 {
	return 0
}

func (c Mos6502) asl() uint8 {
	return 0
}

func (c Mos6502) bcc() uint8 {
	return 0
}

func (c Mos6502) bcs() uint8 {
	return 0
}

func (c Mos6502) beq() uint8 {
	return 0
}

func (c Mos6502) bit() uint8 {
	return 0
}

func (c Mos6502) bmi() uint8 {
	return 0
}

func (c Mos6502) bne() uint8 {
	return 0
}

func (c Mos6502) bpl() uint8 {
	return 0
}

func (c Mos6502) brk() uint8 {
	return 0
}

func (c Mos6502) bvc() uint8 {
	return 0
}

func (c Mos6502) bvs() uint8 {
	return 0
}

func (c Mos6502) clc() uint8 {
	return 0
}

func (c Mos6502) cld() uint8 {
	return 0
}

func (c Mos6502) cli() uint8 {
	return 0
}

func (c Mos6502) clv() uint8 {
	return 0
}

func (c Mos6502) cmp() uint8 {
	return 0
}

func (c Mos6502) cpx() uint8 {
	return 0
}

func (c Mos6502) cpy() uint8 {
	return 0
}

func (c Mos6502) dec() uint8 {
	return 0
}

func (c Mos6502) dex() uint8 {
	return 0
}

func (c Mos6502) dey() uint8 {
	return 0
}

func (c Mos6502) eor() uint8 {
	return 0
}

func (c Mos6502) inc() uint8 {
	return 0
}

func (c Mos6502) inx() uint8 {
	return 0
}

func (c Mos6502) iny() uint8 {
	return 0
}

func (c Mos6502) jmp() uint8 {
	return 0
}

func (c Mos6502) jsr() uint8 {
	return 0
}

func (c Mos6502) lda() uint8 {
	return 0
}

func (c Mos6502) ldx() uint8 {
	return 0
}

func (c Mos6502) ldy() uint8 {
	return 0
}

func (c Mos6502) lsr() uint8 {
	return 0
}

func (c Mos6502) nop() uint8 {
	return 0
}

func (c Mos6502) ora() uint8 {
	return 0
}

func (c Mos6502) pha() uint8 {
	return 0
}

func (c Mos6502) php() uint8 {
	return 0
}

func (c Mos6502) pla() uint8 {
	return 0
}

func (c Mos6502) plp() uint8 {
	return 0
}

func (c Mos6502) rol() uint8 {
	return 0
}

func (c Mos6502) ror() uint8 {
	return 0
}

func (c Mos6502) rti() uint8 {
	return 0
}

func (c Mos6502) rts() uint8 {
	return 0
}

func (c Mos6502) sbc() uint8 {
	return 0
}

func (c Mos6502) sec() uint8 {
	return 0
}

func (c Mos6502) sed() uint8 {
	return 0
}

func (c Mos6502) sei() uint8 {
	return 0
}

func (c Mos6502) sta() uint8 {
	return 0
}

func (c Mos6502) stx() uint8 {
	return 0
}

func (c Mos6502) sty() uint8 {
	return 0
}

func (c Mos6502) tax() uint8 {
	return 0
}

func (c Mos6502) tay() uint8 {
	return 0
}

func (c Mos6502) tsx() uint8 {
	return 0
}

func (c Mos6502) txa() uint8 {
	return 0
}

func (c Mos6502) txs() uint8 {
	return 0
}

func (c Mos6502) tya() uint8 {
	return 0
}

func (c Mos6502) xxx() uint8 {
	return 0
}
