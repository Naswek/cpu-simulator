package machine

import (
	"cpu-simulator/internal/isa"
	"fmt"
	"sort"
)

type InputEvent struct {
	Tick  uint64
	Value int32
}

type IOController struct {
	inputQueue []int32

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
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Tick < events[j].Tick
	})

	return &IOController{
		irqVector:   isa.PortsTable[isa.PortInput].InterruptVector,
		inputEvents: events,
		inputQueue:  make([]int32, 0),
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

		if event.Value < 0 || event.Value > 255 {
			return fmt.Errorf("input value out of byte range at tick %d: %d", tick, event.Value)
		}

		io.inputQueue = append(io.inputQueue, event.Value)
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
	return len(io.inputQueue) > 0
}

func (c *CPU) execOut() error {
	c.appendStages((*CPU).stageWritePort)
	return nil
}

func (c *CPU) execIn() error {
	c.appendStages(
		(*CPU).stageReadPort,
	)
	return nil
}

func (io *IOController) ReadPort(port uint8) (int32, error) {
	if port != isa.PortInput {
		return 0, fmt.Errorf("unknown input port: %v", port)
	}

	if len(io.inputQueue) == 0 {
		return 0, fmt.Errorf("input port is not ready")
	}

	value := io.inputQueue[0]
	io.inputQueue = io.inputQueue[1:]
	io.irqPending = len(io.inputQueue) > 0
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

func (c *CPU) stageReadPort() error {
	value, err := c.io.ReadPort(uint8(c.IR.Operand))
	if err != nil {
		return err
	}

	c.ACC = value
	c.SR = statusWithZN(c.SR, c.ACC)
	c.tick()
	return nil
}

func (c *CPU) stageWritePort() error {
	if err := c.io.WritePort(uint8(c.IR.Operand), c.ACC); err != nil {
		return err
	}

	c.tick()
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

func (c *CPU) Output() string {
	return c.io.Output()
}

func (c *CPU) OutputBytes() []byte {
	return c.io.OutputBytes()
}
