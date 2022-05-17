package seven_seg_display_manager

import (
	"errors"
	"periph.io/x/conn/v3/i2c"
)

// I2CAddr i2c default address.
const Ht16k33I2CAddr uint16 = 0x70

const (
	cmdRAM        = 0x00
	cmdKeys       = 0x40
	displaySetup  = 0x80
	displayOff    = 0x00
	displayOn     = 0x01
	systemSetup   = 0x20
	oscillatorOff = 0x00
	oscillatorOn  = 0x01
	cmdBrightness = 0xE0
)

// Column positions in memory, colon is considered "column 5"
var pos = []int{0, 2, 6, 8, 4}

type BlinkFrequency byte

const (
	BlinkOff    = 0x00
	Blink2Hz    = 0x02
	Blink1Hz    = 0x04
	BlinkHalfHz = 0x06
)

type Dev struct {
	dev i2c.Dev
}

func NewI2CScreen(bus i2c.Bus, address uint16) (*Dev, error) {
	dev := &Dev{dev: i2c.Dev{Bus: bus, Addr: address}}

	if err := dev.init(); err != nil {
		return nil, err
	}

	return dev, nil
}

func (d *Dev) init() error {
	if _, err := d.dev.Write([]byte{systemSetup | oscillatorOn}); err != nil {
		return err
	}
	if _, err := d.dev.Write([]byte{displaySetup | displayOn}); err != nil {
		return err
	}
	if err := d.SetBlink(BlinkOff); err != nil {
		return err
	}
	return d.SetBrightness(15)
}

func (d *Dev) SetBlink(freq BlinkFrequency) error {
	_, err := d.dev.Write([]byte{displaySetup | displayOn | byte(freq)})
	return err
}

func (d *Dev) SetBrightness(brightness int) error {
	if brightness < 0 || brightness > 15 {
		return errors.New("ht16k33: brightness must be between 0 and 15")
	}
	_, err := d.dev.Write([]byte{cmdBrightness | byte(brightness)})
	return err
}

func (d *Dev) WriteColumn(column int, data uint16) error {
	_, err := d.dev.Write([]byte{byte(pos[column]), byte(data & 0xFF), byte(data >> 8)})
	return err
}

func (d *Dev) Halt() error {
	for i := 0; i < 4; i++ {
		if err := d.WriteColumn(i, 0); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dev) SetColon(state bool) error {
	if state {
		if err := d.WriteColumn(4, 0xFF); err != nil {
			return err
		}
	} else {
		if err := d.WriteColumn(4, 0x00); err != nil {
			return err
		}
	}
	return nil
}
