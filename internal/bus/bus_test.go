package bus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBus_Read(t *testing.T) {
	tests := []struct {
		name         string
		ramState     RAM
		address      uint16
		expectedData uint8
	}{
		{
			name:         "correct data read",
			ramState:     RAM{0xff},
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
		ramState     RAM
		address      uint16
		expectedData uint8
	}{
		{
			name:         "correct data read",
			ramState:     RAM{0xff},
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
		initialRAMstate  RAM
		address          uint16
		data             uint8
		expectedRAMstate RAM
	}{
		{
			name:             "write to empty ram succeeds",
			initialRAMstate:  RAM{},
			address:          0x0000,
			data:             0xff,
			expectedRAMstate: RAM{0xff},
		},
		{
			name:             "overwrite of ram address succeeds",
			initialRAMstate:  RAM{0xff},
			address:          0x0000,
			data:             0x11,
			expectedRAMstate: RAM{0x11},
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
