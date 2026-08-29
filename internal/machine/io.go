package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
)

type IOController struct {
	inputValue int32
	inputReady bool

	irqPending bool
	irqVector  uint8

	inputEvents []InputEvent
	nextEvent   int

	output []byte
}

var ioOpcodes = map[isa.Opcode]ControlOperation{
	isa.IN:  (*CPU).execIn,
	isa.OUT: (*CPU).execOut,
}

func newIOController() *IOController {
	return newIOControllerWithEvents(nil)
}

func newIOControllerWithEvents(inputEvents []InputEvent) *IOController {
	events := make([]InputEvent, len(inputEvents))
	copy(events, inputEvents)

	return &IOController{
		irqVector:   isa.PortsTable[isa.PortInput].InterruptVector,
		inputEvents: events,
		output:      make([]byte, 0),
	}
}

func (io *IOController) Advance(tick int) error {
	if tick < 0 {
		return nil
	}

	for io.nextEvent < len(io.inputEvents) {
		event := io.inputEvents[io.nextEvent]
		if event.Tick > uint64(tick) {
			return nil
		}

		if io.inputReady {
			return fmt.Errorf("input overrun at tick %d", tick)
		}

		io.inputValue = event.Value
		io.inputReady = true
		io.irqPending = true
		io.irqVector = isa.PortsTable[isa.PortInput].InterruptVector
		io.nextEvent++
	}

	return nil
}

func (io *IOController) IRQPending() bool {
	return io.irqPending
}

func (io *IOController) IRQVector() uint8 {
	return io.irqVector
}

func (io *IOController) ClearIRQ() {
	io.irqPending = false
}

func (io *IOController) InputReady() bool {
	return io.inputReady
}

func (c *CPU) execOut() error {
	return c.io.WritePort(uint8(c.IR.Operand), c.ACC)
}

func (c *CPU) execIn() error {
	value, err := c.io.ReadPort(uint8(c.IR.Operand))
	if err != nil {
		return err
	}

	c.ACC = value
	c.updateZN(c.ACC)
	return nil
}

func (io *IOController) ReadPort(port uint8) (int32, error) {
	if port != isa.PortInput {
		return 0, fmt.Errorf("unknown input port: %v", port)
	}

	if !io.inputReady {
		return 0, fmt.Errorf("input port is not ready")
	}

	value := io.inputValue
	io.inputValue = 0
	io.inputReady = false
	io.irqPending = false
	return value, nil
}

func (io *IOController) WritePort(port uint8, value int32) error {
	if port != isa.PortOutput {
		return fmt.Errorf("unknown output port: %v", port)
	}

	if value < 0 || value > 255 {
		return fmt.Errorf("output value out of byte range: %d", value)
	}

	io.output = append(io.output, byte(value))
	return nil
}

func (io *IOController) Output() string {
	return string(io.output)
}

func (io *IOController) OutputBytes() []byte {
	output := make([]byte, len(io.output))
	copy(output, io.output)
	return output
}
