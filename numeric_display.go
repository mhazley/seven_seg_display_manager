package seven_seg_display_manager

import (
	"github.com/rs/zerolog/log"
	"periph.io/x/conn/v3/i2c"
)

var digitValues = map[rune]uint16{
	' ': 0x0,
	'.': 0x80,
	'0': 0x3F,
	'1': 0x06,
	'2': 0x5B,
	'3': 0x4F,
	'4': 0x66,
	'5': 0x6D,
	'6': 0x7D,
	'7': 0x07,
	'8': 0x7F,
	'9': 0x6F,
	'A': 0x77,
	'B': 0x7C,
	'C': 0x39,
	'D': 0x5E,
	'E': 0x79,
	'F': 0x71,
}

type NumericDisplay struct {
	dev *Dev
}

func NewNumericDisplay(bus i2c.Bus, address uint16) (*NumericDisplay, error) {
	dev, err := NewI2CScreen(bus, address)
	if err != nil {
		return nil, err
	}
	return &NumericDisplay{dev: dev}, nil
}

func (d *NumericDisplay) SetDigit(pos int, digit rune, decimal bool) error {
	val := digitValues[digit]
	if decimal {
		val |= digitValues['.']
	}
	return d.dev.WriteColumn(pos, val)
}

func (d *NumericDisplay) WriteString(s string) (int, error) {
	log.Debug().Str("NumericDisplay->WriteString", s).Send()
	if err := d.dev.Halt(); err != nil {
		return 0, err
	}

	pos := 4 - len(s)
	if pos < 0 {
		pos = 0
	}

	for _, ch := range s {
		if ch == '.' {
			// Print decimal points on the previous digit.
			c := rune(s[pos-1])
			if err := d.SetDigit(pos-1, c, true); err != nil {
				return pos, err
			}
		} else {
			if err := d.SetDigit(pos, ch, false); err != nil {
				return pos, err
			}
			pos++
		}
	}
	return pos, nil
}

func (d *NumericDisplay) Halt() error {
	return d.dev.Halt()
}

func (d *NumericDisplay) SetColon(state bool) error {
	return d.dev.SetColon(state)
}
