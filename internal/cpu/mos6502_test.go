package cpu

import (
	"testing"

	"github.com/Jac0bDeal/goNES/internal/bus"
	"github.com/stretchr/testify/assert"
)

func newTestMos6502() *Mos6502 {
	return &Mos6502{
		bus: bus.NewBus(bus.RAM{}),
		lookup: mos6502LookupTable{
			instruction{
				operation:   "TOP",
				addressMode: "TAM",
				performOp: func() uint8 {
					return 0
				},
				setAddressMode: func() uint8 {
					return 0
				},
				cycles: 0,
			},
		},
	}
}

// busBuilder is a test utility for creating Bus instances inline
// test case definitions.
type busBuilder struct {
	bus *bus.Bus
}

// newBusBuilder creates and returns a pointer to a busBuilder instance.
func newBusBuilder() *busBuilder {
	return &busBuilder{
		bus.NewBus(bus.RAM{}),
	}
}

// build returns the built Bus instance.
func (b *busBuilder) build() *bus.Bus {
	return b.bus
}

// write assigns the passed data to the given address on the Bus being built.
func (b *busBuilder) write(address word, data byte) *busBuilder {
	b.bus.Write(uint16(address), data)
	return b
}

func TestMos6502_Clock(t *testing.T) {
	testCases := []struct {
		name              string
		setupInitialState func(*testing.T) *Mos6502
		expectedPC        word
		expectedA         byte
		expectedX         byte
		expectedY         byte
		expectedStkp      byte
		expectedStatus    byte
		expectedCycles    uint8
	}{
		{
			name: "cpu only decrements cycle count if non-zero",
			setupInitialState: func(*testing.T) *Mos6502 {
				return &Mos6502{
					pc:     0x0000,
					a:      0x11,
					x:      0x11,
					y:      0x11,
					stkp:   0x11,
					status: 0b00000000,
					cycles: 1,
					lookup: mos6502LookupTable{
						{
							operation:   "TOP",
							addressMode: "TAM",
							performOp: func() uint8 {
								return 0
							},
							setAddressMode: func() uint8 {
								return 0
							},
							cycles: 0,
						},
					},
					bus: bus.NewBus(bus.RAM{}),
				}
			},
			expectedPC:     0x0000,
			expectedA:      0x11,
			expectedX:      0x11,
			expectedY:      0x11,
			expectedStkp:   0x11,
			expectedStatus: 0b00000000,
			expectedCycles: 0,
		},
		{
			name: "cpu performs clock cycle correctly with no additional cycles",
			setupInitialState: func(t *testing.T) *Mos6502 {
				t.Helper()
				cpu := &Mos6502{
					pc:     0x0000,
					a:      0x11,
					x:      0x11,
					y:      0x11,
					stkp:   0x11,
					status: 0b00000000,
					cycles: 0,
					bus: bus.NewBus(bus.RAM{
						0x00,
					}),
				}
				cpu.lookup = mos6502LookupTable{
					{
						operation:   "TOP",
						addressMode: "TAM",
						performOp: func() uint8 {
							// assert that U flag is set to true first
							assert.Equal(t, uint8(0x01), cpu.GetStatusFlag(U))
							// set U flag to false for final state assertion
							cpu.setStatusFlag(U, false)
							return 0
						},
						setAddressMode: func() uint8 {
							return 0
						},
						cycles: 2,
					},
				}
				return cpu
			},
			expectedPC:     0x0001,
			expectedA:      0x11,
			expectedX:      0x11,
			expectedY:      0x11,
			expectedStkp:   0x11,
			expectedStatus: 0b00100000,
			expectedCycles: 1,
		},
		{
			name: "cpu performs clock cycle with additional operation cycle only",
			setupInitialState: func(t *testing.T) *Mos6502 {
				t.Helper()
				cpu := &Mos6502{
					pc:     0x0000,
					a:      0x11,
					x:      0x11,
					y:      0x11,
					stkp:   0x11,
					status: 0b00000000,
					cycles: 0,
					bus: bus.NewBus(bus.RAM{
						0x00,
					}),
				}
				cpu.lookup = mos6502LookupTable{
					{
						operation:   "TOP",
						addressMode: "TAM",
						performOp: func() uint8 {
							return 1
						},
						setAddressMode: func() uint8 {
							return 0
						},
						cycles: 2,
					},
				}
				return cpu
			},
			expectedPC:     0x0001,
			expectedA:      0x11,
			expectedX:      0x11,
			expectedY:      0x11,
			expectedStkp:   0x11,
			expectedStatus: 0b00100000,
			expectedCycles: 1,
		},
		{
			name: "cpu performs clock cycle with additional address mode cycle only",
			setupInitialState: func(t *testing.T) *Mos6502 {
				t.Helper()
				cpu := &Mos6502{
					pc:     0x0000,
					a:      0x11,
					x:      0x11,
					y:      0x11,
					stkp:   0x11,
					status: 0b00000000,
					cycles: 0,
					bus: bus.NewBus(bus.RAM{
						0x00,
					}),
				}
				cpu.lookup = mos6502LookupTable{
					{
						operation:   "TOP",
						addressMode: "TAM",
						performOp: func() uint8 {
							return 0
						},
						setAddressMode: func() uint8 {
							return 1
						},
						cycles: 2,
					},
				}
				return cpu
			},
			expectedPC:     0x0001,
			expectedA:      0x11,
			expectedX:      0x11,
			expectedY:      0x11,
			expectedStkp:   0x11,
			expectedStatus: 0b00100000,
			expectedCycles: 1,
		},
		{
			name: "cpu performs clock cycle with additional address mode cycle only",
			setupInitialState: func(t *testing.T) *Mos6502 {
				t.Helper()
				cpu := &Mos6502{
					pc:     0x0000,
					a:      0x11,
					x:      0x11,
					y:      0x11,
					stkp:   0x11,
					status: 0b00000000,
					cycles: 0,
					bus: bus.NewBus(bus.RAM{
						0x00,
					}),
				}
				cpu.lookup = mos6502LookupTable{
					{
						operation:   "TOP",
						addressMode: "TAM",
						performOp: func() uint8 {
							return 1
						},
						setAddressMode: func() uint8 {
							return 1
						},
						cycles: 2,
					},
				}
				return cpu
			},
			expectedPC:     0x0001,
			expectedA:      0x11,
			expectedX:      0x11,
			expectedY:      0x11,
			expectedStkp:   0x11,
			expectedStatus: 0b00100000,
			expectedCycles: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.setupInitialState(t)
			cpu.Clock()

			assert.Equal(t, tc.expectedPC, cpu.pc)
			assert.Equal(t, tc.expectedA, cpu.a)
			assert.Equal(t, tc.expectedX, cpu.x)
			assert.Equal(t, tc.expectedY, cpu.y)
			assert.Equal(t, tc.expectedStkp, cpu.stkp)
			assert.Equal(t, tc.expectedStatus, cpu.status)
			assert.Equal(t, tc.expectedCycles, cpu.cycles)
		})
	}
}

func TestMos6502_Reset(t *testing.T) {
	testCases := []struct {
		name          string
		initialState  *Mos6502
		expectedState *Mos6502
	}{
		{
			name: "cpu resets correctly",
			initialState: &Mos6502{
				pc:              0x1111,
				a:               0x11,
				x:               0x11,
				y:               0x11,
				stkp:            0x11,
				status:          0b11111111,
				addressAbsolute: 0x1111,
				addressRelative: 0x1111,
				fetchedData:     0x11,
				cycles:          1,
				bus: newBusBuilder().
					write(0xfffc, 0x20).
					write(0xfffd, 0x04).
					build(),
			},
			expectedState: &Mos6502{
				pc:              0x0420,
				a:               0x00,
				x:               0x00,
				y:               0x00,
				stkp:            0xfd,
				status:          0b00100000,
				addressAbsolute: 0x0000,
				addressRelative: 0x0000,
				fetchedData:     0x00,
				cycles:          8,
				bus: newBusBuilder().
					write(0xfffc, 0x20).
					write(0xfffd, 0x04).
					build(),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			cpu.Reset()

			assert.Equal(t, tc.expectedState, cpu)
		})
	}
}

func TestMos6502_InterruptRequest(t *testing.T) {
	testCases := []struct {
		name          string
		initialState  *Mos6502
		expectedState *Mos6502
	}{
		{
			name: "cpu Interrupt Request does nothing with I true",
			initialState: &Mos6502{
				pc:              0x1111,
				a:               0x11,
				x:               0x11,
				y:               0x11,
				stkp:            0x11,
				status:          0b00100100,
				addressAbsolute: 0x1111,
				addressRelative: 0x1111,
				fetchedData:     0x11,
				cycles:          1,
				bus: newBusBuilder().
					write(0xfffe, 0x20).
					write(0xffff, 0x04).
					build(),
			},
			expectedState: &Mos6502{
				pc:              0x1111,
				a:               0x11,
				x:               0x11,
				y:               0x11,
				stkp:            0x11,
				status:          0b00100100,
				addressAbsolute: 0x1111,
				addressRelative: 0x1111,
				fetchedData:     0x11,
				cycles:          1,
				bus: newBusBuilder().
					write(0xfffe, 0x20).
					write(0xffff, 0x04).
					build(),
			},
		},
		{
			name: "cpu Interrupt Request executes correctly",
			initialState: &Mos6502{
				pc:              0x0420,
				a:               0x11,
				x:               0x11,
				y:               0x11,
				stkp:            0x11,
				status:          0b00010000,
				addressAbsolute: 0x1111,
				addressRelative: 0x1111,
				fetchedData:     0x11,
				cycles:          1,
				bus: newBusBuilder().
					write(0xfffe, 0x20).
					write(0xffff, 0x04).
					build(),
			},
			expectedState: &Mos6502{
				pc:              0x0420,
				a:               0x11,
				x:               0x11,
				y:               0x11,
				stkp:            0x0e,
				status:          0b00100100,
				addressAbsolute: 0xfffe,
				addressRelative: 0x1111,
				fetchedData:     0x11,
				cycles:          7,
				bus: newBusBuilder().
					write(0x010f, 0b00100100).
					write(0x0110, 0x20).
					write(0x0111, 0x04).
					write(0xfffe, 0x20).
					write(0xffff, 0x04).
					build(),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			cpu.InterruptRequest()

			assert.Equal(t, tc.expectedState, cpu)
		})
	}
}

func TestMos6502_NonMaskableInterrupt(t *testing.T) {
	testCases := []struct {
		name          string
		initialState  *Mos6502
		expectedState *Mos6502
	}{
		{
			name: "cpu NMI request executes correctly",
			initialState: &Mos6502{
				pc:              0x0420,
				a:               0x11,
				x:               0x11,
				y:               0x11,
				stkp:            0x11,
				status:          0b00010000,
				addressAbsolute: 0x1111,
				addressRelative: 0x1111,
				fetchedData:     0x11,
				cycles:          1,
				bus: newBusBuilder().
					write(0xfffa, 0x20).
					write(0xfffb, 0x04).
					build(),
			},
			expectedState: &Mos6502{
				pc:              0x0420,
				a:               0x11,
				x:               0x11,
				y:               0x11,
				stkp:            0x0e,
				status:          0b00100100,
				addressAbsolute: 0xfffa,
				addressRelative: 0x1111,
				fetchedData:     0x11,
				cycles:          8,
				bus: newBusBuilder().
					write(0x010f, 0b00100100).
					write(0x0110, 0x20).
					write(0x0111, 0x04).
					write(0xfffa, 0x20).
					write(0xfffb, 0x04).
					build(),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			cpu.NonMaskableInterrupt()

			assert.Equal(t, tc.expectedState, cpu)
		})
	}
}

func TestMos6502_GetStatusFlag(t *testing.T) {
	cpu := newTestMos6502()
	cpu.status = 0b10101010
	testCases := []struct {
		name          string
		flag          Flag
		expectedValue uint8
	}{
		{
			name:          "C flag correct",
			flag:          C,
			expectedValue: 0,
		},
		{
			name:          "Z flag correct",
			flag:          Z,
			expectedValue: 1,
		},
		{
			name:          "I flag correct",
			flag:          I,
			expectedValue: 0,
		},
		{
			name:          "D flag correct",
			flag:          D,
			expectedValue: 1,
		},
		{
			name:          "B flag correct",
			flag:          B,
			expectedValue: 0,
		},
		{
			name:          "U flag correct",
			flag:          U,
			expectedValue: 1,
		},
		{
			name:          "V flag correct",
			flag:          V,
			expectedValue: 0,
		},
		{
			name:          "N flag correct",
			flag:          N,
			expectedValue: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := cpu.GetStatusFlag(tc.flag)
			assert.Equal(t, tc.expectedValue, value)
		})
	}
}

func TestMos6502_setStatusFlag(t *testing.T) {
	allFlags := []Flag{
		C,
		Z,
		I,
		D,
		B,
		U,
		V,
		N,
	}
	testCases := []struct {
		name           string
		initialStatus  uint8
		value          bool
		expectedStatus uint8
	}{
		{
			name:           "initially false sets to true successfully",
			initialStatus:  0b00000000,
			value:          true,
			expectedStatus: 0b11111111,
		},
		{
			name:           "initially true sets to false successfully",
			initialStatus:  0b11111111,
			value:          false,
			expectedStatus: 0b00000000,
		},
		{
			name:           "initially false set to false does nothing",
			initialStatus:  0b0000000,
			value:          false,
			expectedStatus: 0b00000000,
		},
		{
			name:           "initially true set to true does nothing",
			initialStatus:  0b11111111,
			value:          true,
			expectedStatus: 0b11111111,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.status = tc.initialStatus

			for _, flag := range allFlags {
				cpu.setStatusFlag(flag, tc.value)
			}

			assert.Equal(t, tc.expectedStatus, cpu.status)
		})
	}
}

func TestMos6502_imp(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "correct data fetched from accumulator",
			initialState: &Mos6502{
				a:           42,
				fetchedData: 0,
			},
			expectedState: &Mos6502{
				a:           42,
				fetchedData: 42,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.imp()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_imm(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "absolute address points to address immediately after program counter",
			initialState: &Mos6502{
				pc:              0,
				addressAbsolute: 0,
			},
			expectedState: &Mos6502{
				pc:              0,
				addressAbsolute: 1,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.imm()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_zp0(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "absolute address set to start of page 0x00 for byte 0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             newBusBuilder().write(0x0000, 0x00).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             newBusBuilder().write(0x0000, 0x00).build(),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0042,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to end of page 0x00 for byte 0xff",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x00ff,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.zp0()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_zpx(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "absolute address set to start of page 0x00 for byte 0x00 and x=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x00,
				bus:             newBusBuilder().write(0x0000, 0x00).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x00,
				bus:             newBusBuilder().write(0x0000, 0x00).build(),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00 and x=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x00,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x00,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0042,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to end of page 0x00 for byte 0xff and x=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x00,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x00,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x00ff,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to start of page 0x00 for byte 0xff and x=0x01",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x01,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x01,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00 and x=0x01",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x01,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x01,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0043,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to byte 0xfe of page 0x00 for byte 0xff and x=0xff",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0xff,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0xff,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x00fe,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.zpx()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_zpy(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "absolute address set to start of page 0x00 for byte 0x00 and y=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x00,
				bus:             newBusBuilder().write(0x0000, 0x00).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x00,
				bus:             newBusBuilder().write(0x0000, 0x00).build(),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00 and y=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x00,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x00,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0042,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to end of page 0x00 for byte 0xff and y=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x00,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x00,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x00ff,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to start of page 0x00 for byte 0xff and y=0x01",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x01,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x01,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00 and y=0x01",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x01,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x01,
				bus:             newBusBuilder().write(0x0000, 0x42).build(),
				addressAbsolute: 0x0043,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to byte 0xfe of page 0x00 for byte 0xff and y=0xff",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0xff,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0xff,
				bus:             newBusBuilder().write(0x0000, 0xff).build(),
				addressAbsolute: 0x00fe,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.zpy()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_rel(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "relative address at top of branch range is set correctly",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             newBusBuilder().write(0x0000, 0x79).build(),
				addressRelative: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             newBusBuilder().write(0x0000, 0x79).build(),
				addressRelative: 0x0079,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "relative address within branch range is set correctly",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             newBusBuilder().write(0x0000, 0x00).build(),
				addressRelative: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             newBusBuilder().write(0x0000, 0x00).build(),
				addressRelative: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "relative address at bottom of branch range is set correctly",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             newBusBuilder().write(0x0000, 0x80).build(),
				addressRelative: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             newBusBuilder().write(0x0000, 0x80).build(),
				addressRelative: 0xff80,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.rel()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_abs(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "absolute address read correctly",
			initialState: &Mos6502{
				pc: 0x0000,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0002,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.abs()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_abx(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "x equals zero",
			initialState: &Mos6502{
				pc: 0x0000,
				x:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0002,
				x:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "non-zero x and no page change",
			initialState: &Mos6502{
				pc: 0x0000,
				x:  0x10,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x02,
				x:  0x0010,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0430,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "x value causing page change",
			initialState: &Mos6502{
				pc: 0x0000,
				x:  0x42,
				bus: newBusBuilder().
					write(0x0000, 0xde).
					write(0x0001, 0x03).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0002,
				x:  0x42,
				bus: newBusBuilder().
					write(0x0000, 0xde).
					write(0x0001, 0x03).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.abx()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_aby(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "y equals zero",
			initialState: &Mos6502{
				pc: 0x0000,
				y:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0002,
				y:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "non-zero y and no page change",
			initialState: &Mos6502{
				pc: 0x0000,
				y:  0x10,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x02,
				y:  0x0010,
				bus: newBusBuilder().
					write(0x0000, 0x20).
					write(0x0001, 0x04).
					build(),
				addressAbsolute: 0x0430,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "y value causing page change",
			initialState: &Mos6502{
				pc: 0x0000,
				y:  0x42,
				bus: newBusBuilder().
					write(0x0000, 0xde).
					write(0x0001, 0x03).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0002,
				y:  0x42,
				bus: newBusBuilder().
					write(0x0000, 0xde).
					write(0x0001, 0x03).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.aby()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_ind(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "reads indirect address correctly with no page change",
			initialState: &Mos6502{
				pc: 0x0000,
				bus: newBusBuilder().
					write(0x0000, 0x02).
					write(0x0001, 0x00).
					write(0x0002, 0x20).
					write(0x0003, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0002,
				bus: newBusBuilder().
					write(0x0000, 0x02).
					write(0x0001, 0x00).
					write(0x0002, 0x20).
					write(0x0003, 0x04).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "replicates page change bug wrapping around to start of page",
			initialState: &Mos6502{
				pc: 0x0000,
				bus: newBusBuilder().
					write(0x0000, 0xff).
					write(0x0001, 0x00).
					write(0x00ff, 0x04).
					write(0x0100, 0x20).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0002,
				bus: newBusBuilder().
					write(0x0000, 0xff).
					write(0x0001, 0x00).
					write(0x00ff, 0x04).
					write(0x0100, 0x20).
					build(),
				addressAbsolute: 0xff04,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.ind()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_izx(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "indirect assignment of page 0 address when x equals 0",
			initialState: &Mos6502{
				pc: 0x0000,
				x:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x01).
					write(0x0001, 0xff).
					write(0x0002, 0xff).
					write(0x0003, 0x20).
					write(0x0004, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0001,
				x:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x01).
					write(0x0001, 0xff).
					write(0x0002, 0xff).
					write(0x0003, 0x20).
					write(0x0004, 0x04).
					build(),
				addressAbsolute: 0xffff,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "indirect assignment of page 0x00 is offset by x index value",
			initialState: &Mos6502{
				pc: 0x0000,
				x:  0x02,
				bus: newBusBuilder().
					write(0x0000, 0x01).
					write(0x0001, 0xff).
					write(0x0002, 0xff).
					write(0x0003, 0x20).
					write(0x0004, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0001,
				x:  0x02,
				bus: newBusBuilder().
					write(0x0000, 0x01).
					write(0x0001, 0xff).
					write(0x0002, 0xff).
					write(0x0003, 0x20).
					write(0x0004, 0x04).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "indirect assignment of page 0x00 is offset by large x index value",
			initialState: &Mos6502{
				pc: 0x0000,
				x:  0xfe,
				bus: newBusBuilder().
					write(0x0000, 0x01).
					write(0x00ff, 0x20).
					write(0x0100, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0001,
				x:  0xfe,
				bus: newBusBuilder().
					write(0x0000, 0x01).
					write(0x00ff, 0x20).
					write(0x0100, 0x04).
					build(),
				addressAbsolute: 0x0120,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.izx()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_izy(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "indirect assignment of page 0x00 address when y equals 0",
			initialState: &Mos6502{
				pc: 0x0000,
				y:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x03).
					write(0x0001, 0x01).
					write(0x0002, 0x00).
					write(0x0003, 0x20).
					write(0x0004, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0001,
				y:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x03).
					write(0x0001, 0x01).
					write(0x0002, 0x00).
					write(0x0003, 0x20).
					write(0x0004, 0x04).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "indirect assignment of page 0x00 address when y equals 0",
			initialState: &Mos6502{
				pc: 0x0000,
				y:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x03).
					write(0x0001, 0x01).
					write(0x0002, 0x00).
					write(0x0003, 0x20).
					write(0x0004, 0x04).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0001,
				y:  0x00,
				bus: newBusBuilder().
					write(0x0000, 0x03).
					write(0x0001, 0x01).
					write(0x0002, 0x00).
					write(0x0003, 0x20).
					write(0x0004, 0x04).
					build(),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "indirect assignment of y shifted address with page change",
			initialState: &Mos6502{
				pc: 0x0000,
				y:  0xff,
				bus: newBusBuilder().
					write(0x0000, 0x03).
					write(0x0001, 0x01).
					write(0x0002, 0x00).
					write(0x0003, 0x10).
					write(0x0004, 0x00).
					build(),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc: 0x0001,
				y:  0xff,
				bus: newBusBuilder().
					write(0x0000, 0x03).
					write(0x0001, 0x01).
					write(0x0002, 0x00).
					write(0x0003, 0x10).
					write(0x0004, 0x00).
					build(),
				addressAbsolute: 0x010f,
			},
			expectedAdditionalCycles: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.izy()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_adc(t *testing.T) {
	testCases := []struct {
		name string

		busValue      uint8
		initialAvalue uint8
		initialCflag  bool

		expectedAvalue           uint8
		expectedCflag            uint8
		expectedZflag            uint8
		expectedVflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "0+0=0 sets Z=true",

			busValue:      0x00,
			initialAvalue: 0x00,
			initialCflag:  false,

			expectedAvalue:           0x00,
			expectedCflag:            0,
			expectedZflag:            1,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "x+(-x)=0 C=false sets Z=true",

			busValue:      0x81,
			initialAvalue: 0x7f,
			initialCflag:  false,

			expectedAvalue:           0x00,
			expectedCflag:            1,
			expectedZflag:            1,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "x+(-x)+1=0 C=true sets Z=true",

			busValue:      0x80,
			initialAvalue: 0x7f,
			initialCflag:  true,

			expectedAvalue:           0x00,
			expectedCflag:            1,
			expectedZflag:            1,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P+P=P C=false no overflow",

			busValue:      0x01,
			initialAvalue: 0x01,
			initialCflag:  false,

			expectedAvalue:           0x02,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P+P=P C=true no overflow",

			busValue:      0x01,
			initialAvalue: 0x01,
			initialCflag:  true,

			expectedAvalue:           0x03,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N+N=N C=false no overflow",

			busValue:      0xff,
			initialAvalue: 0xff,
			initialCflag:  false,

			expectedAvalue:           0xfe,
			expectedCflag:            1,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N+N=N C=true no overflow",

			busValue:      0xff,
			initialAvalue: 0xff,
			initialCflag:  true,

			expectedAvalue:           0xff,
			expectedCflag:            1,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P+P=N C=false causes overflow",

			busValue:      0x7f,
			initialAvalue: 0x01,
			initialCflag:  false,

			expectedAvalue:           0x80,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            1,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P+P=P C=true causes overflow",

			busValue:      0x7e,
			initialAvalue: 0x01,
			initialCflag:  true,

			expectedAvalue:           0x80,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            1,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N+N=P C=false causes overflow",

			busValue:      0x81,
			initialAvalue: 0x81,
			initialCflag:  false,

			expectedAvalue:           0x02,
			expectedCflag:            1,
			expectedZflag:            0,
			expectedVflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N+N=P C=true causes overflow",

			busValue:      0x80,
			initialAvalue: 0x80,
			initialCflag:  true,

			expectedAvalue:           0x01,
			expectedCflag:            1,
			expectedZflag:            0,
			expectedVflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P+N=P C=false cannot overflow",

			busValue:      0x7f,
			initialAvalue: 0xff,
			initialCflag:  false,

			expectedAvalue:           0x7e,
			expectedCflag:            1,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P+N=P C=true cannot overflow",

			busValue:      0x7f,
			initialAvalue: 0xff,
			initialCflag:  true,

			expectedAvalue:           0x7f,
			expectedCflag:            1,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P+N=N C=false cannot overflow",

			busValue:      0x1,
			initialAvalue: 0x80,
			initialCflag:  false,

			expectedAvalue:           0x81,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P+N=N C=true cannot overflow",

			busValue:      0x1,
			initialAvalue: 0x81,
			initialCflag:  true,

			expectedAvalue:           0x83,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := &Mos6502{
				a:   tc.initialAvalue,
				bus: bus.NewBus(bus.RAM{tc.busValue}),
			}
			cpu.lookup = mos6502LookupTable{
				{
					operation:      adc,
					addressMode:    "TST",
					performOp:      cpu.adc,
					setAddressMode: func() uint8 { return 0 },
				},
			}
			cpu.setStatusFlag(C, tc.initialCflag)

			additionalCycles := cpu.adc()

			assert.Equal(t, tc.expectedAvalue, cpu.a, "incorrect A value")
			assert.Equal(t, tc.expectedCflag, cpu.GetStatusFlag(C), "incorrect C flag")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedVflag, cpu.GetStatusFlag(V), "incorrect V flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_and(t *testing.T) {
	testCases := []struct {
		name string

		aValue    uint8
		dataValue uint8

		expectedAvalue           uint8
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name:      "and operation performed correctly",
			aValue:    0b00100101,
			dataValue: 0b00100001,

			expectedAvalue:           0b00100001,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name:      "and operation results in 0 sets Z true",
			aValue:    0b00100101,
			dataValue: 0b10000010,

			expectedAvalue:           0b00000000,
			expectedZflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name:      "and operation results in negative result sets N true",
			aValue:    0b10100101,
			dataValue: 0b10000100,

			expectedAvalue:           0b10000100,
			expectedZflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := &Mos6502{
				a:               tc.aValue,
				addressAbsolute: 0x0000,
				bus:             bus.NewBus(bus.RAM{tc.dataValue}),
			}
			additionalCycles := cpu.and()

			assert.Equal(t, tc.expectedAvalue, cpu.a, "incorrect A value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_asl(t *testing.T) {
	testCases := []struct {
		name string

		dataValue   uint8
		instruction instruction

		expectedAvalue           uint8
		expectedBusValue         uint8
		expectedCflag            uint8
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "asl operation in implied mode performed correctly",

			dataValue: 0b00000001,
			instruction: instruction{
				operation:   asl,
				addressMode: imp,
			},

			expectedAvalue:           0b00000010,
			expectedBusValue:         0b00000000,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "asl operation in non-implied mode performed correctly",

			dataValue: 0b00000001,
			instruction: instruction{
				operation:   asl,
				addressMode: "TST",
			},

			expectedAvalue:           0b00000000,
			expectedBusValue:         0b00000010,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "asl operation resulting in 0 sets Z true",

			dataValue: 0b00000000,
			instruction: instruction{
				operation:   asl,
				addressMode: imp,
			},

			expectedAvalue:           0b00000000,
			expectedBusValue:         0b00000000,
			expectedCflag:            0,
			expectedZflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "asl operation resulting in negative result sets N true",

			dataValue: 0b01000010,
			instruction: instruction{
				operation:   asl,
				addressMode: imp,
			},

			expectedAvalue:           0b10000100,
			expectedBusValue:         0b00000000,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "asl operation resulting in carry sets C true",

			dataValue: 0b10000001,
			instruction: instruction{
				operation:   asl,
				addressMode: imp,
			},

			expectedAvalue:           0b00000010,
			expectedBusValue:         0b00000000,
			expectedCflag:            1,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := &Mos6502{
				addressAbsolute: 0x0000,
				bus:             bus.NewBus(bus.RAM{}),
				lookup:          mos6502LookupTable{tc.instruction},
			}
			if tc.instruction.addressMode == imp {
				cpu.fetchedData = tc.dataValue
			} else {
				cpu.write(cpu.addressAbsolute, tc.dataValue)
			}
			additionalCycles := cpu.asl()

			assert.Equal(t, tc.expectedAvalue, cpu.a, "incorrect A value")
			assert.Equal(t, tc.expectedBusValue, cpu.read(cpu.addressAbsolute), "incorrect Bus value")
			assert.Equal(t, tc.expectedCflag, cpu.GetStatusFlag(C), "incorrect C flag")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles value")
		})
	}
}

func TestMos6502_bcc(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "C=true nothing happens",
			initialState: &Mos6502{
				status: 0b00000001,
				cycles: 2,
			},
			expectedState: &Mos6502{
				status: 0b00000001,
				cycles: 2,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "C=false assigns pc correctly",
			initialState: &Mos6502{
				status:          0b11111110,
				addressAbsolute: 0x00000,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b11111110,
				addressAbsolute: 0x0011,
				pc:              0x0011,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "C=false relative addressing assigns pc correctly",
			initialState: &Mos6502{
				status:          0b11111110,
				addressAbsolute: 0x0000,
				addressRelative: 0x0011,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b11111110,
				addressAbsolute: 0x00022,
				addressRelative: 0x0011,
				pc:              0x0022,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "C=false page change causes extra cycle",
			initialState: &Mos6502{
				status:          0b11111110,
				addressAbsolute: 0x0000,
				addressRelative: 0x1100,
				pc:              0x1111,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b11111110,
				addressAbsolute: 0x2211,
				addressRelative: 0x1100,
				pc:              0x2211,
				cycles:          4,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.bcc()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_bcs(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "C=false nothing happens",
			initialState: &Mos6502{
				status: 0b11111110,
				cycles: 2,
			},
			expectedState: &Mos6502{
				status: 0b11111110,
				cycles: 2,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "C=true assigns pc correctly",
			initialState: &Mos6502{
				status:          0b00000001,
				addressAbsolute: 0x00000,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b00000001,
				addressAbsolute: 0x0011,
				pc:              0x0011,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "C=true relative addressing assigns pc correctly",
			initialState: &Mos6502{
				status:          0b00000001,
				addressAbsolute: 0x0000,
				addressRelative: 0x0011,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b00000001,
				addressAbsolute: 0x00022,
				addressRelative: 0x0011,
				pc:              0x0022,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "C=true page change causes extra cycle",
			initialState: &Mos6502{
				status:          0b00000001,
				addressAbsolute: 0x0000,
				addressRelative: 0x1100,
				pc:              0x1111,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b00000001,
				addressAbsolute: 0x2211,
				addressRelative: 0x1100,
				pc:              0x2211,
				cycles:          4,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.bcs()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_beq(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "Z=false nothing happens",
			initialState: &Mos6502{
				status: 0b11111101,
				cycles: 2,
			},
			expectedState: &Mos6502{
				status: 0b11111101,
				cycles: 2,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "Z=true assigns pc correctly",
			initialState: &Mos6502{
				status:          0b00000010,
				addressAbsolute: 0x00000,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b00000010,
				addressAbsolute: 0x0011,
				pc:              0x0011,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "Z=true relative addressing assigns pc correctly",
			initialState: &Mos6502{
				status:          0b00000010,
				addressAbsolute: 0x0000,
				addressRelative: 0x0011,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b00000010,
				addressAbsolute: 0x00022,
				addressRelative: 0x0011,
				pc:              0x0022,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "Z=true page change causes extra cycle",
			initialState: &Mos6502{
				status:          0b00000010,
				addressAbsolute: 0x0000,
				addressRelative: 0x1100,
				pc:              0x1111,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b00000010,
				addressAbsolute: 0x2211,
				addressRelative: 0x1100,
				pc:              0x2211,
				cycles:          4,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.beq()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles)
		})
	}
}

func TestMos6502_bit(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedZflag            uint8
		expectedNflag            uint8
		expectedVflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "nothing gets set",
			initialState: &Mos6502{
				a: 0x11,
				bus: bus.NewBus(bus.RAM{
					0x11,
				}),
			},
			expectedZflag:            0,
			expectedNflag:            0,
			expectedVflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "result of mask is zero",
			initialState: &Mos6502{
				a: 0xff,
				bus: bus.NewBus(bus.RAM{
					0x00,
				}),
			},
			expectedZflag:            1,
			expectedNflag:            0,
			expectedVflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "addressed data is negative",
			initialState: &Mos6502{
				a: 0xff,
				bus: bus.NewBus(bus.RAM{
					0x80,
				}),
			},
			expectedZflag:            0,
			expectedNflag:            1,
			expectedVflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "addressed data is overflow value",
			initialState: &Mos6502{
				a: 0xff,
				bus: bus.NewBus(bus.RAM{
					0x7f,
				}),
			},
			expectedZflag:            0,
			expectedNflag:            0,
			expectedVflag:            1,
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.bit()

			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedVflag, cpu.GetStatusFlag(V), "incorrect V flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_bmi(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "N=false nothing happens",
			initialState: &Mos6502{
				status: 0b01111111,
				cycles: 2,
			},
			expectedState: &Mos6502{
				status: 0b01111111,
				cycles: 2,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "N=true assigns pc correctly",
			initialState: &Mos6502{
				status:          0b10000000,
				addressAbsolute: 0x00000,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b10000000,
				addressAbsolute: 0x0011,
				pc:              0x0011,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "N=true relative addressing assigns pc correctly",
			initialState: &Mos6502{
				status:          0b10000000,
				addressAbsolute: 0x0000,
				addressRelative: 0x0011,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b10000000,
				addressAbsolute: 0x00022,
				addressRelative: 0x0011,
				pc:              0x0022,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "N=true page change causes extra cycle",
			initialState: &Mos6502{
				status:          0b10000000,
				addressAbsolute: 0x0000,
				addressRelative: 0x1100,
				pc:              0x1111,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b10000000,
				addressAbsolute: 0x2211,
				addressRelative: 0x1100,
				pc:              0x2211,
				cycles:          4,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.bmi()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_bne(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "Z=true nothing happens",
			initialState: &Mos6502{
				status: 0b00000010,
				cycles: 2,
			},
			expectedState: &Mos6502{
				status: 0b00000010,
				cycles: 2,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "Z=false assigns pc correctly",
			initialState: &Mos6502{
				status:          0b11111101,
				addressAbsolute: 0x00000,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b11111101,
				addressAbsolute: 0x0011,
				pc:              0x0011,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "Z=false relative addressing assigns pc correctly",
			initialState: &Mos6502{
				status:          0b11111101,
				addressAbsolute: 0x0000,
				addressRelative: 0x0011,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b11111101,
				addressAbsolute: 0x00022,
				addressRelative: 0x0011,
				pc:              0x0022,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "Z=false page change causes extra cycle",
			initialState: &Mos6502{
				status:          0b11111101,
				addressAbsolute: 0x0000,
				addressRelative: 0x1100,
				pc:              0x1111,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b11111101,
				addressAbsolute: 0x2211,
				addressRelative: 0x1100,
				pc:              0x2211,
				cycles:          4,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.bne()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_bpl(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "N=true nothing happens",
			initialState: &Mos6502{
				status: 0b10000000,
				cycles: 2,
			},
			expectedState: &Mos6502{
				status: 0b10000000,
				cycles: 2,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "N=false assigns pc correctly",
			initialState: &Mos6502{
				status:          0b01111111,
				addressAbsolute: 0x00000,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b01111111,
				addressAbsolute: 0x0011,
				pc:              0x0011,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "N=false relative addressing assigns pc correctly",
			initialState: &Mos6502{
				status:          0b01111111,
				addressAbsolute: 0x0000,
				addressRelative: 0x0011,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b01111111,
				addressAbsolute: 0x00022,
				addressRelative: 0x0011,
				pc:              0x0022,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "N=false page change causes extra cycle",
			initialState: &Mos6502{
				status:          0b01111111,
				addressAbsolute: 0x0000,
				addressRelative: 0x1100,
				pc:              0x1111,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b01111111,
				addressAbsolute: 0x2211,
				addressRelative: 0x1100,
				pc:              0x2211,
				cycles:          4,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.bpl()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_brk(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "break is performed correctly",
			initialState: &Mos6502{
				pc:     0x11fe,
				stkp:   0x42,
				status: 0b11101011,
				bus: newBusBuilder().
					write(0xfffe, 0x20).
					write(0xffff, 0x04).
					build(),
			},
			expectedState: &Mos6502{
				pc:     0x0420,
				status: 0b11101111,
				stkp:   0x3f,
				bus: newBusBuilder().
					write(0x0140, 0b11111111).
					write(0x0141, 0xff).
					write(0x0142, 0x11).
					write(0xfffe, 0x20).
					write(0xffff, 0x04).
					build(),
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.brk()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_bvc(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "V=true nothing happens",
			initialState: &Mos6502{
				status: 0b01000000,
				cycles: 2,
			},
			expectedState: &Mos6502{
				status: 0b01000000,
				cycles: 2,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "V=false assigns pc correctly",
			initialState: &Mos6502{
				status:          0b10111111,
				addressAbsolute: 0x00000,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b10111111,
				addressAbsolute: 0x0011,
				pc:              0x0011,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "V=false relative addressing assigns pc correctly",
			initialState: &Mos6502{
				status:          0b10111111,
				addressAbsolute: 0x0000,
				addressRelative: 0x0011,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b10111111,
				addressAbsolute: 0x00022,
				addressRelative: 0x0011,
				pc:              0x0022,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "V=false page change causes extra cycle",
			initialState: &Mos6502{
				status:          0b10111111,
				addressAbsolute: 0x0000,
				addressRelative: 0x1100,
				pc:              0x1111,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b10111111,
				addressAbsolute: 0x2211,
				addressRelative: 0x1100,
				pc:              0x2211,
				cycles:          4,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.bvc()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_bvs(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "V=false nothing happens",
			initialState: &Mos6502{
				status: 0b10111111,
				cycles: 2,
			},
			expectedState: &Mos6502{
				status: 0b10111111,
				cycles: 2,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "V=true assigns pc correctly",
			initialState: &Mos6502{
				status:          0b01000000,
				addressAbsolute: 0x00000,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b01000000,
				addressAbsolute: 0x0011,
				pc:              0x0011,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "V=true relative addressing assigns pc correctly",
			initialState: &Mos6502{
				status:          0b01000000,
				addressAbsolute: 0x0000,
				addressRelative: 0x0011,
				pc:              0x0011,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b01000000,
				addressAbsolute: 0x00022,
				addressRelative: 0x0011,
				pc:              0x0022,
				cycles:          3,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "V=true page change causes extra cycle",
			initialState: &Mos6502{
				status:          0b01000000,
				addressAbsolute: 0x0000,
				addressRelative: 0x1100,
				pc:              0x1111,
				cycles:          2,
			},
			expectedState: &Mos6502{
				status:          0b01000000,
				addressAbsolute: 0x2211,
				addressRelative: 0x1100,
				pc:              0x2211,
				cycles:          4,
			},
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.bvs()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_clc(t *testing.T) {
	testCases := []struct {
		name                     string
		initialCflag             bool
		expectedCflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name:                     "sets C=true to false",
			initialCflag:             true,
			expectedCflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name:                     "sets C=false to false",
			initialCflag:             false,
			expectedCflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.setStatusFlag(C, tc.initialCflag)
			additionalCycles := cpu.clc()

			assert.Equal(t, tc.expectedCflag, cpu.GetStatusFlag(C), "incorrect C flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_cld(t *testing.T) {
	testCases := []struct {
		name                     string
		initialDflag             bool
		expectedDflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name:                     "sets D=true to false",
			initialDflag:             true,
			expectedDflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name:                     "sets D=false to false",
			initialDflag:             false,
			expectedDflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.setStatusFlag(D, tc.initialDflag)
			additionalCycles := cpu.cld()

			assert.Equal(t, tc.expectedDflag, cpu.GetStatusFlag(D), "incorrect D flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_cli(t *testing.T) {
	testCases := []struct {
		name                     string
		initialIflag             bool
		expectedIflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name:                     "sets I=true to false",
			initialIflag:             true,
			expectedIflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name:                     "sets I=false to false",
			initialIflag:             false,
			expectedIflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.setStatusFlag(I, tc.initialIflag)
			additionalCycles := cpu.cli()

			assert.Equal(t, tc.expectedIflag, cpu.GetStatusFlag(I), "incorrect I flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_clv(t *testing.T) {
	testCases := []struct {
		name                     string
		initialVflag             bool
		expectedVflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name:                     "sets V=true to false",
			initialVflag:             true,
			expectedVflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name:                     "sets V=false to false",
			initialVflag:             false,
			expectedVflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.setStatusFlag(V, tc.initialVflag)
			additionalCycles := cpu.clv()

			assert.Equal(t, tc.expectedVflag, cpu.GetStatusFlag(V), "incorrect V flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_cmp(t *testing.T) {
	testCases := []struct {
		name string

		aValue  byte
		busData byte

		expectedCflag            uint8
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "A=M results in Z=true",

			aValue:  0x42,
			busData: 0x42,

			expectedCflag:            1,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "A>=M results in C=true",

			aValue:  0x42,
			busData: 0x01,

			expectedCflag:            1,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "A<=M results in N=true",

			aValue:  0x01,
			busData: 0x42,

			expectedCflag:            0,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.a = tc.aValue
			cpu.write(0x0000, tc.busData)
			additionalCycles := cpu.cmp()

			assert.Equal(t, tc.expectedCflag, cpu.GetStatusFlag(C), "incorrect C flag")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_cpx(t *testing.T) {
	testCases := []struct {
		name string

		xValue  byte
		busData byte

		expectedCflag            uint8
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "X=M results in Z=true",

			xValue:  0x42,
			busData: 0x42,

			expectedCflag:            1,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "X>=M results in C=true",

			xValue:  0x42,
			busData: 0x01,

			expectedCflag:            1,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "X<=M results in N=true",

			xValue:  0x01,
			busData: 0x42,

			expectedCflag:            0,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.x = tc.xValue
			cpu.write(0x0000, tc.busData)
			additionalCycles := cpu.cpx()

			assert.Equal(t, tc.expectedCflag, cpu.GetStatusFlag(C), "incorrect C flag")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_cpy(t *testing.T) {
	testCases := []struct {
		name string

		yValue  byte
		busData byte

		expectedCflag            uint8
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "Y=M results in Z=true",

			yValue:  0x42,
			busData: 0x42,

			expectedCflag:            1,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "Y>=M results in C=true",

			yValue:  0x42,
			busData: 0x01,

			expectedCflag:            1,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "Y<=M results in N=true",

			yValue:  0x01,
			busData: 0x42,

			expectedCflag:            0,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.y = tc.yValue
			cpu.write(0x0000, tc.busData)
			additionalCycles := cpu.cpy()

			assert.Equal(t, tc.expectedCflag, cpu.GetStatusFlag(C), "incorrect C flag")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_dec(t *testing.T) {
	testCases := []struct {
		name string

		address      word
		initialValue byte

		expectedValue            byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "value is decremented",

			address:      0x0420,
			initialValue: 0x42,

			expectedValue:            0x41,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "decrement resulting in 0 sets Z=true",

			address:      0x0420,
			initialValue: 0x01,

			expectedValue:            0x00,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "decrement resulting in -1 sets N=true",

			address:      0x0420,
			initialValue: 0x00,

			expectedValue:            0xff,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.addressAbsolute = tc.address
			cpu.write(cpu.addressAbsolute, tc.initialValue)
			additionalCycles := cpu.dec()

			assert.Equal(t, tc.expectedValue, cpu.read(tc.address), "incorrect value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_dex(t *testing.T) {
	testCases := []struct {
		name string

		initialXvalue byte

		expectedXvalue           byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "value is decremented",

			initialXvalue: 0x42,

			expectedXvalue:           0x41,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "decrement resulting in 0 sets Z=true",

			initialXvalue: 0x01,

			expectedXvalue:           0x00,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "decrement resulting in -1 sets N=true",

			initialXvalue: 0x00,

			expectedXvalue:           0xff,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.x = tc.initialXvalue
			additionalCycles := cpu.dex()

			assert.Equal(t, tc.expectedXvalue, cpu.x, "incorrect X value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_dey(t *testing.T) {
	testCases := []struct {
		name string

		initialYvalue byte

		expectedYvalue           byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "value is decremented",

			initialYvalue: 0x42,

			expectedYvalue:           0x41,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "decrement resulting in 0 sets Z=true",

			initialYvalue: 0x01,

			expectedYvalue:           0x00,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "decrement resulting in -1 sets N=true",

			initialYvalue: 0x00,

			expectedYvalue:           0xff,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.y = tc.initialYvalue
			additionalCycles := cpu.dey()

			assert.Equal(t, tc.expectedYvalue, cpu.y, "incorrect Y value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_eor(t *testing.T) {
	testCases := []struct {
		name string

		address       word
		dataValue     byte
		initialAvalue byte

		expectedAvalue           byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "result is positive",

			address:       0x0420,
			dataValue:     0b01010101,
			initialAvalue: 0b00101011,

			expectedAvalue:           0b01111110,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "result is zero sets Z=true",

			address:       0x0420,
			dataValue:     0b01010101,
			initialAvalue: 0b01010101,

			expectedAvalue:           0b00000000,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "result is negative sets N=true",

			address:       0x0420,
			dataValue:     0b11010101,
			initialAvalue: 0b00101010,

			expectedAvalue:           0b11111111,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.a = tc.initialAvalue
			cpu.addressAbsolute = tc.address
			cpu.write(cpu.addressAbsolute, tc.dataValue)
			additionalCycles := cpu.eor()

			assert.Equal(t, tc.expectedAvalue, cpu.a, "incorrect A value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_inc(t *testing.T) {
	testCases := []struct {
		name string

		address      word
		initialValue byte

		expectedValue            byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "value is incremented",

			address:      0x0420,
			initialValue: 0x42,

			expectedValue:            0x43,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "increment resulting in 0 sets Z=true",

			address:      0x0420,
			initialValue: 0xff,

			expectedValue:            0x00,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "increment resulting in negative sets N=true",

			address:      0x0420,
			initialValue: 0xf0,

			expectedValue:            0xf1,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.addressAbsolute = tc.address
			cpu.write(cpu.addressAbsolute, tc.initialValue)
			additionalCycles := cpu.inc()

			assert.Equal(t, tc.expectedValue, cpu.read(tc.address), "incorrect value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_inx(t *testing.T) {
	testCases := []struct {
		name string

		initialXvalue byte

		expectedXvalue           byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "value is incremented",

			initialXvalue: 0x42,

			expectedXvalue:           0x43,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "increment resulting in 0 sets Z=true",

			initialXvalue: 0xff,

			expectedXvalue:           0x00,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "increment resulting in negative sets N=true",

			initialXvalue: 0xf0,

			expectedXvalue:           0xf1,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.x = tc.initialXvalue
			additionalCycles := cpu.inx()

			assert.Equal(t, tc.expectedXvalue, cpu.x, "incorrect X value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_iny(t *testing.T) {
	testCases := []struct {
		name string

		initialYvalue byte

		expectedYvalue           byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "value is incremented",

			initialYvalue: 0x42,

			expectedYvalue:           0x43,
			expectedNflag:            0,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "increment resulting in 0 sets Z=true",

			initialYvalue: 0xff,

			expectedYvalue:           0x00,
			expectedNflag:            0,
			expectedZflag:            1,
			expectedAdditionalCycles: 0,
		},
		{
			name: "increment resulting in negative sets N=true",

			initialYvalue: 0xf0,

			expectedYvalue:           0xf1,
			expectedNflag:            1,
			expectedZflag:            0,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.y = tc.initialYvalue
			additionalCycles := cpu.iny()

			assert.Equal(t, tc.expectedYvalue, cpu.y, "incorrect Y value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_jmp(t *testing.T) {
	testCases := []struct {
		name string

		address word

		expectedPC               word
		expectedAdditionalCycles uint8
	}{
		{
			name: "program counter assigned correct address",

			address: 0x0420,

			expectedPC:               0x0420,
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.addressAbsolute = tc.address
			additionalCycles := cpu.jmp()

			assert.Equal(t, tc.expectedPC, cpu.pc, "incorrect PC value")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_jsr(t *testing.T) {
	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "program counter written to stack and then assigned with correct address",
			initialState: &Mos6502{
				pc:              0x1200,
				stkp:            0x42,
				addressAbsolute: 0x0420,
				bus:             newBusBuilder().build(),
			},
			expectedState: &Mos6502{
				pc:              0x0420,
				stkp:            0x40,
				addressAbsolute: 0x0420,
				bus: newBusBuilder().
					write(0x0141, 0xff).
					write(0x0142, 0x11).
					build(),
			},
			expectedAdditionalCycles: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := tc.initialState
			additionalCycles := cpu.jsr()

			assert.Equal(t, tc.expectedState, cpu)
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_lda(t *testing.T) {
	testCases := []struct {
		name string

		busData byte

		expectedAvalue           byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "accumulator assigned correct value",

			busData: 0x42,

			expectedAvalue:           0x42,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "accumulator assigned zero value and Z set true",

			busData: 0x00,

			expectedAvalue:           0x00,
			expectedZflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "accumulator assigned negative value and N set true",

			busData: 0x80,

			expectedAvalue:           0x80,
			expectedZflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.addressAbsolute = 0x0000
			cpu.bus = newBusBuilder().write(cpu.addressAbsolute, tc.busData).build()
			additionalCycles := cpu.lda()

			assert.Equal(t, tc.expectedAvalue, cpu.a, "incorrect A value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_ldx(t *testing.T) {
	testCases := []struct {
		name string

		busData byte

		expectedXvalue           byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "x register assigned correct value",

			busData: 0x42,

			expectedXvalue:           0x42,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "x register assigned zero value and Z set true",

			busData: 0x00,

			expectedXvalue:           0x00,
			expectedZflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "x register assigned negative value and N set true",

			busData: 0x80,

			expectedXvalue:           0x80,
			expectedZflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.addressAbsolute = 0x0000
			cpu.bus = newBusBuilder().write(cpu.addressAbsolute, tc.busData).build()
			additionalCycles := cpu.ldx()

			assert.Equal(t, tc.expectedXvalue, cpu.x, "incorrect X value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_ldy(t *testing.T) {
	testCases := []struct {
		name string

		busData byte

		expectedYvalue           byte
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "y register assigned correct value",

			busData: 0x42,

			expectedYvalue:           0x42,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "y register assigned zero value and Z set true",

			busData: 0x00,

			expectedYvalue:           0x00,
			expectedZflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "y register assigned negative value and N set true",

			busData: 0x80,

			expectedYvalue:           0x80,
			expectedZflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := newTestMos6502()
			cpu.addressAbsolute = 0x0000
			cpu.bus = newBusBuilder().write(cpu.addressAbsolute, tc.busData).build()
			additionalCycles := cpu.ldy()

			assert.Equal(t, tc.expectedYvalue, cpu.y, "incorrect Y value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_lsr(t *testing.T) {
	testCases := []struct {
		name string

		dataValue   uint8
		instruction instruction

		expectedAvalue           uint8
		expectedBusValue         uint8
		expectedCflag            uint8
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "lsr operation in implied mode performed correctly",

			dataValue: 0x42,
			instruction: instruction{
				operation:   lsr,
				addressMode: imp,
			},

			expectedAvalue:           0x21,
			expectedBusValue:         0x00,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "lsr operation in non-implied mode performed correctly",

			dataValue: 0x42,
			instruction: instruction{
				operation:   lsr,
				addressMode: "TST",
			},

			expectedAvalue:           0x00,
			expectedBusValue:         0x21,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "lsr operation resulting in 0 sets Z true",

			dataValue: 0x00,
			instruction: instruction{
				operation:   lsr,
				addressMode: imp,
			},

			expectedAvalue:           0x00,
			expectedBusValue:         0x00,
			expectedCflag:            0,
			expectedZflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 0,
		},
		{
			name: "lsr operation resulting in negative result sets N and C true",

			dataValue: 0x05,
			instruction: instruction{
				operation:   lsr,
				addressMode: imp,
			},

			expectedAvalue:           0x82,
			expectedBusValue:         0x00,
			expectedCflag:            1,
			expectedZflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := &Mos6502{
				addressAbsolute: 0x0000,
				bus:             bus.NewBus(bus.RAM{}),
				lookup:          mos6502LookupTable{tc.instruction},
			}
			if tc.instruction.addressMode == imp {
				cpu.fetchedData = tc.dataValue
			} else {
				cpu.write(cpu.addressAbsolute, tc.dataValue)
			}
			additionalCycles := cpu.lsr()

			assert.Equal(t, tc.expectedAvalue, cpu.a, "incorrect A value")
			assert.Equal(t, tc.expectedBusValue, cpu.read(cpu.addressAbsolute), "incorrect Bus value")
			assert.Equal(t, tc.expectedCflag, cpu.GetStatusFlag(C), "incorrect C flag")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles value")
		})
	}
}

func TestMos6502_nop(t *testing.T) {
	additionalCycleOpcodes := []byte{
		0x1c,
		0x3c,
		0x5c,
		0x7c,
		0xdc,
		0xfc,
	}
	noAdditionalCycleOpcodes := make([]byte, 0, 256-len(additionalCycleOpcodes))
	for opcode := 0x00; opcode < 0x100; opcode++ {
		var match bool
		for _, o := range additionalCycleOpcodes {
			if byte(opcode) == o {
				match = true
			}
		}
		if !match {
			noAdditionalCycleOpcodes = append(noAdditionalCycleOpcodes, byte(opcode))
		}
	}

	testCases := []struct {
		name string

		opcodes []byte

		expectedAdditionalCycles uint8
	}{
		{
			name: "opcodes requiring no additional cycls",

			opcodes: noAdditionalCycleOpcodes,

			expectedAdditionalCycles: 0,
		},
		{
			name: "opcodes requiring 1 additional cycle",

			opcodes: additionalCycleOpcodes,

			expectedAdditionalCycles: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, opcode := range tc.opcodes {
				cpu := newTestMos6502()
				cpu.opcode = opcode
				additionalCycles := cpu.nop()

				assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
			}
		})
	}
}

func TestMos6502_ora(t *testing.T) {
	testCases := []struct {
		name string

		aValue    uint8
		dataValue uint8

		expectedAvalue           uint8
		expectedZflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name:      "ora operation performed correctly",
			aValue:    0b01010101,
			dataValue: 0b00101010,

			expectedAvalue:           0b01111111,
			expectedZflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name:      "ora operation results in 0 sets Z true",
			aValue:    0b00000000,
			dataValue: 0b00000000,

			expectedAvalue:           0b00000000,
			expectedZflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name:      "ora operation results in negative result sets N true",
			aValue:    0b10000000,
			dataValue: 0b01111111,

			expectedAvalue:           0b11111111,
			expectedZflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := &Mos6502{
				a:               tc.aValue,
				addressAbsolute: 0x0000,
				bus:             bus.NewBus(bus.RAM{tc.dataValue}),
			}
			additionalCycles := cpu.ora()

			assert.Equal(t, tc.expectedAvalue, cpu.a, "incorrect A value")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}

func TestMos6502_sbc(t *testing.T) {
	testCases := []struct {
		name string

		busValue      byte
		initialAvalue byte
		initialCflag  bool

		expectedAvalue           byte
		expectedCflag            uint8
		expectedZflag            uint8
		expectedVflag            uint8
		expectedNflag            uint8
		expectedAdditionalCycles uint8
	}{
		{
			name: "0-0=-1 C=false sets Z=true",

			busValue:      0x00,
			initialAvalue: 0x00,
			initialCflag:  false,

			expectedAvalue:           0xff,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "0-0=0 C=true sets Z=true",

			busValue:      0x00,
			initialAvalue: 0x00,
			initialCflag:  true,

			expectedAvalue:           0x00,
			expectedCflag:            0,
			expectedZflag:            1,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P-P=-1 C=false sets Z=true",

			busValue:      0x7f,
			initialAvalue: 0x7f,
			initialCflag:  false,

			expectedAvalue:           0xff,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P-P=0 C=true sets Z=true",

			busValue:      0x7f,
			initialAvalue: 0x7f,
			initialCflag:  true,

			expectedAvalue:           0x00,
			expectedCflag:            0,
			expectedZflag:            1,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N-N=-1 C=false cannot overflow",

			busValue:      0x80,
			initialAvalue: 0x80,
			initialCflag:  false,

			expectedAvalue:           0xff,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N-N=0 C=true cannot overflow",

			busValue:      0x80,
			initialAvalue: 0x80,
			initialCflag:  true,

			expectedAvalue:           0x00,
			expectedCflag:            0,
			expectedZflag:            1,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N-N=P C=false cannot overflow",

			busValue:      0x80,
			initialAvalue: 0xff,
			initialCflag:  false,

			expectedAvalue:           0x7e,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N-N=P C=true cannot overflow",

			busValue:      0x80,
			initialAvalue: 0xff,
			initialCflag:  true,

			expectedAvalue:           0x7f,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "0-N=N C=true overflows",

			busValue:      0x80,
			initialAvalue: 0x00,
			initialCflag:  true,

			expectedAvalue:           0x80,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            1,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P-P=P C=false cannot overflow",

			busValue:      0x01,
			initialAvalue: 0x7f,
			initialCflag:  false,

			expectedAvalue:           0x7d,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P-P=P C=true cannot overflow",

			busValue:      0x7e,
			initialAvalue: 0x7f,
			initialCflag:  true,

			expectedAvalue:           0x01,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "0-P=N C=false cannot overflow",

			busValue:      0x7f,
			initialAvalue: 0x00,
			initialCflag:  false,

			expectedAvalue:           0x80,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "0-P=N C=true cannot overflow",

			busValue:      0x7f,
			initialAvalue: 0x00,
			initialCflag:  true,

			expectedAvalue:           0x81,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P-N=N C=false causes overflow",

			busValue:      0x80,
			initialAvalue: 0x01,
			initialCflag:  false,

			expectedAvalue:           0x80,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            1,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "P-N=N C=true causes overflow",

			busValue:      0x80,
			initialAvalue: 0x01,
			initialCflag:  true,

			expectedAvalue:           0x81,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            1,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N-P=P C=false causes overflow",

			busValue:      0x7f,
			initialAvalue: 0xff,
			initialCflag:  false,

			expectedAvalue:           0x7f,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            1,
			expectedNflag:            0,
			expectedAdditionalCycles: 1,
		},
		{
			name: "N-P=P C=true cannot overflow",

			busValue:      0x7f,
			initialAvalue: 0xff,
			initialCflag:  true,

			expectedAvalue:           0x80,
			expectedCflag:            0,
			expectedZflag:            0,
			expectedVflag:            0,
			expectedNflag:            1,
			expectedAdditionalCycles: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := &Mos6502{
				a:   tc.initialAvalue,
				bus: bus.NewBus(bus.RAM{tc.busValue}),
			}
			cpu.lookup = mos6502LookupTable{
				{
					operation:      sbc,
					addressMode:    "TST",
					performOp:      cpu.sbc,
					setAddressMode: func() uint8 { return 0 },
				},
			}
			cpu.setStatusFlag(C, tc.initialCflag)

			additionalCycles := cpu.sbc()

			assert.Equal(t, tc.expectedAvalue, cpu.a, "incorrect A value")
			assert.Equal(t, tc.expectedCflag, cpu.GetStatusFlag(C), "incorrect C flag")
			assert.Equal(t, tc.expectedZflag, cpu.GetStatusFlag(Z), "incorrect Z flag")
			assert.Equal(t, tc.expectedVflag, cpu.GetStatusFlag(V), "incorrect V flag")
			assert.Equal(t, tc.expectedNflag, cpu.GetStatusFlag(N), "incorrect N flag")
			assert.Equal(t, tc.expectedAdditionalCycles, additionalCycles, "incorrect additional cycles")
		})
	}
}
