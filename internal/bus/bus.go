package bus

const ramSize = 64 * 1024

// Bus represents the bus used by the CPU to communicate with other components. It can be
// read from and written to.
type Bus struct {
	ram [ramSize]uint8
}

// NewBus constructs and returns a Bus instance.
func NewBus() *Bus {
	return &Bus{}
}

// Read reads a byte at a given address on the Bus.
func (b *Bus) Read(address uint16) uint8 {
	return b.ram[address]
}

// ReadByteOnly ...
func (b *Bus) ReadByteOnly(address uint16) uint8 {
	return 0
}

// Write writes a byte of data to an address on the Bus.
func (b *Bus) Write(address uint16, data uint8) {
	b.ram[address] = data
}
