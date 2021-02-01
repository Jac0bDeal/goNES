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

func TestMos6502_Clock(t *testing.T) {
	testCases := []struct {
		name              string
		setupInitialState func(*testing.T) *Mos6502
		expectedPC        uint16
		expectedA         uint8
		expectedX         uint8
		expectedY         uint8
		expectedStkp      uint8
		expectedStatus    uint8
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
	testRAM := bus.RAM{}
	testRAM[0xfffc] = 0x20
	testRAM[0xfffd] = 0x04
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
				bus:             bus.NewBus(testRAM),
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
				bus:             bus.NewBus(testRAM),
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
	testRAM := bus.RAM{}
	testRAM[0xfffe] = 0x20
	testRAM[0xffff] = 0x04

	expectedSuccessRAM := bus.RAM{}
	expectedSuccessRAM[0x010f] = 0b00100100
	expectedSuccessRAM[0x0110] = 0x20
	expectedSuccessRAM[0x0111] = 0x04
	expectedSuccessRAM[0xfffe] = 0x20
	expectedSuccessRAM[0xffff] = 0x04

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
				bus:             bus.NewBus(testRAM),
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
				bus:             bus.NewBus(testRAM),
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
				bus:             bus.NewBus(testRAM),
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
				bus:             bus.NewBus(expectedSuccessRAM),
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
	testRAM := bus.RAM{}
	testRAM[0xfffa] = 0x20
	testRAM[0xfffb] = 0x04

	expectedSuccessRAM := bus.RAM{}
	expectedSuccessRAM[0x010f] = 0b00100100
	expectedSuccessRAM[0x0110] = 0x20
	expectedSuccessRAM[0x0111] = 0x04
	expectedSuccessRAM[0xfffa] = 0x20
	expectedSuccessRAM[0xfffb] = 0x04

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
				bus:             bus.NewBus(testRAM),
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
				bus:             bus.NewBus(expectedSuccessRAM),
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
				bus:             bus.NewBus(bus.RAM{0x00}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             bus.NewBus(bus.RAM{0x00}),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0042,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to end of page 0x00 for byte 0xff",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             bus.NewBus(bus.RAM{0xff}),
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
				bus:             bus.NewBus(bus.RAM{0x00}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0x00}),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00 and x=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0042,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to end of page 0x00 for byte 0xff and x=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x00ff,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to start of page 0x00 for byte 0xff and x=0x01",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x01,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x01,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00 and x=0x01",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x01,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x01,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0043,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to byte 0xfe of page 0x00 for byte 0xff and x=0xff",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0xff,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0xff,
				bus:             bus.NewBus(bus.RAM{0xff}),
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
				bus:             bus.NewBus(bus.RAM{0x00}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x00}),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00 and y=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0042,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to end of page 0x00 for byte 0xff and y=0x00",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x00ff,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to start of page 0x00 for byte 0xff and y=0x01",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x01,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x01,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to correct location in page 0x00 and y=0x01",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x01,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x01,
				bus:             bus.NewBus(bus.RAM{0x42}),
				addressAbsolute: 0x0043,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "absolute address set to byte 0xfe of page 0x00 for byte 0xff and y=0xff",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0xff,
				bus:             bus.NewBus(bus.RAM{0xff}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0xff,
				bus:             bus.NewBus(bus.RAM{0xff}),
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
				bus:             bus.NewBus(bus.RAM{0x79}),
				addressRelative: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             bus.NewBus(bus.RAM{0x79}),
				addressRelative: 0x0079,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "relative address within branch range is set correctly",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             bus.NewBus(bus.RAM{0x00}),
				addressRelative: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             bus.NewBus(bus.RAM{0x00}),
				addressRelative: 0x0000,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "relative address at bottom of branch range is set correctly",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             bus.NewBus(bus.RAM{0x80}),
				addressRelative: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				bus:             bus.NewBus(bus.RAM{0x80}),
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
				pc:              0x0000,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0002,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
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
				pc:              0x0000,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0002,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "non-zero x and no page change",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x10,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x02,
				x:               0x0010,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0430,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "x value causing page change",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x42,
				bus:             bus.NewBus(bus.RAM{0xde, 0x03}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0002,
				x:               0x42,
				bus:             bus.NewBus(bus.RAM{0xde, 0x03}),
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
				pc:              0x0000,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0002,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "non-zero y and no page change",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x10,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x02,
				y:               0x0010,
				bus:             bus.NewBus(bus.RAM{0x20, 0x04}),
				addressAbsolute: 0x0430,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "y value causing page change",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x42,
				bus:             bus.NewBus(bus.RAM{0xde, 0x03}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0002,
				y:               0x42,
				bus:             bus.NewBus(bus.RAM{0xde, 0x03}),
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
	bugRAM := bus.RAM{0xff, 0x00}
	bugRAM[0x00ff] = 0x04
	bugRAM[0x0100] = 0x20

	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "reads indirect address correctly with no page change",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             bus.NewBus(bus.RAM{0x02, 0x00, 0x20, 0x4}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0002,
				bus:             bus.NewBus(bus.RAM{0x02, 0x00, 0x20, 0x4}),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "replicates page change bug wrapping around to start of page",
			initialState: &Mos6502{
				pc:              0x0000,
				bus:             bus.NewBus(bugRAM),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0002,
				bus:             bus.NewBus(bugRAM),
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
	largeXram := bus.RAM{0x01}
	largeXram[0x00ff] = 0x20
	largeXram[0x0100] = 0x04

	testCases := []struct {
		name                     string
		initialState             *Mos6502
		expectedState            *Mos6502
		expectedAdditionalCycles uint8
	}{
		{
			name: "indirect assignment of page 0 address when x equals 0",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0x01, 0xff, 0xff, 0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x00,
				bus:             bus.NewBus(bus.RAM{0x01, 0xff, 0xff, 0x20, 0x04}),
				addressAbsolute: 0xffff,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "indirect assignment of page 0x00 is offset by x index value",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0x02,
				bus:             bus.NewBus(bus.RAM{0x01, 0xff, 0xff, 0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0x02,
				bus:             bus.NewBus(bus.RAM{0x01, 0xff, 0xff, 0x20, 0x04}),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "indirect assignment of page 0x00 is offset by large x index value",
			initialState: &Mos6502{
				pc:              0x0000,
				x:               0xfe,
				bus:             bus.NewBus(largeXram),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				x:               0xfe,
				bus:             bus.NewBus(largeXram),
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
				pc:              0x0000,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x03, 0x01, 0x00, 0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x03, 0x01, 0x00, 0x20, 0x04}),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "indirect assignment of page 0x00 address when y equals 0",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x03, 0x01, 0x00, 0x20, 0x04}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0x00,
				bus:             bus.NewBus(bus.RAM{0x03, 0x01, 0x00, 0x20, 0x04}),
				addressAbsolute: 0x0420,
			},
			expectedAdditionalCycles: 0,
		},
		{
			name: "indirect assignment of y shifted address with page change",
			initialState: &Mos6502{
				pc:              0x0000,
				y:               0xff,
				bus:             bus.NewBus(bus.RAM{0x03, 0x01, 0x00, 0x10, 0x00}),
				addressAbsolute: 0x0000,
			},
			expectedState: &Mos6502{
				pc:              0x0001,
				y:               0xff,
				bus:             bus.NewBus(bus.RAM{0x03, 0x01, 0x00, 0x10, 0x00}),
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
			name: "P+N=P C=false cannont overflow",

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
			name: "P+N=P C=true cannont overflow",

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
			name: "P+N=N C=false cannont overflow",

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
			name: "P+N=N C=true cannont overflow",

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
