package bus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBus_Read(t *testing.T) {
	tests := []struct {
		name         string
		ramState     [ramSize]uint8
		address      uint16
		expectedData uint8
	}{
		{
			name:         "correct data read",
			ramState:     [ramSize]uint8{0xff},
			address:      0x0000,
			expectedData: 0xff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBus(tt.ramState)
			data := b.Read(tt.address)

			assert.Equal(t, tt.expectedData, data)
		})
	}
}

func TestBus_ReadByteOnly(t *testing.T) {
	tests := []struct {
		name         string
		ramState     [ramSize]uint8
		address      uint16
		expectedData uint8
	}{
		{
			name:         "correct data read",
			ramState:     [ramSize]uint8{0xff},
			address:      0x0000,
			expectedData: 0xff,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBus(tt.ramState)
			data := b.ReadByteOnly(tt.address)

			assert.Equal(t, tt.expectedData, data)
		})
	}
}

func TestBus_Write(t *testing.T) {
	tests := []struct {
		name             string
		initialRAMstate  [ramSize]uint8
		address          uint16
		data             uint8
		expectedRAMstate [ramSize]uint8
	}{
		{
			name:             "write to empty ram succeeds",
			initialRAMstate:  [ramSize]uint8{},
			address:          0x0000,
			data:             0xff,
			expectedRAMstate: [ramSize]uint8{0xff},
		},
		{
			name:             "overwrite of ram address succeeds",
			initialRAMstate:  [ramSize]uint8{0xff},
			address:          0x0000,
			data:             0x11,
			expectedRAMstate: [ramSize]uint8{0x11},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBus(tt.initialRAMstate)
			b.Write(tt.address, tt.data)

			assert.Equal(t, tt.expectedRAMstate, b.ram)
		})
	}
}
