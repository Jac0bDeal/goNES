package bus

// RAMsize is the size of the bus RAM space.
const RAMsize = 64 * 1024

// RAM represents the addressable RAM space on the Bus.
type RAM [RAMsize]uint8

// Bus represents the bus used by the CPU to communicate with other components. It can be
// read from and written to.
type Bus struct {
	ram RAM
}

// NewBus constructs and returns a Bus instance.
func NewBus(ram RAM) *Bus {
	return &Bus{
		ram: ram,
	}
}

// Read reads a byte at a given address on the Bus.
func (b *Bus) Read(address uint16) uint8 {
	return b.ram[address]
}

// ReadByteOnly will be used by diassaembler to read address with mutating state.
// currently unused.
func (b *Bus) ReadByteOnly(address uint16) uint8 {
	return b.Read(address)
}

// Write writes a byte of data to an address on the Bus.
func (b *Bus) Write(address uint16, data uint8) {
	b.ram[address] = data
}
