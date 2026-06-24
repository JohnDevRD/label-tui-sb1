package printers

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

func ListUSBPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("listing ports: %w", err)
	}
	if len(ports) == 0 {
		return nil, nil
	}
	return ports, nil
}

func PrintZPL(portName string, zplData string) error {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		return fmt.Errorf("opening port %s: %w", portName, err)
	}
	defer port.Close()

	port.SetReadTimeout(2 * time.Second)

	n, err := port.Write([]byte(zplData))
	if err != nil {
		return fmt.Errorf("writing to port %s: %w", portName, err)
	}

	if n != len(zplData) {
		return fmt.Errorf("wrote %d bytes but expected %d", n, len(zplData))
	}

	return nil
}
