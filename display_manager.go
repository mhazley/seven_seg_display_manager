package seven_seg_display_manager

import (
	"github.com/rs/zerolog/log"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type DisplayManager struct {
	dispOne *NumericDisplay
	dispTwo *NumericDisplay
}

type Display int8

const (
	DisplayOne Display = iota
	DisplayTwo
)

func NewDisplayManager() (DisplayManager, error) {
	dm := DisplayManager{}
	err := dm.init()
	return dm, err
}

func (d *DisplayManager) init() error {
	if _, err := host.Init(); err != nil {
		log.Fatal().Err(err).Msg("Failed to init host")
	}

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open i2c")
	}

	d.dispOne, err = NewNumericDisplay(bus, Ht16k33I2CAddr+0)
	d.dispTwo, err = NewNumericDisplay(bus, Ht16k33I2CAddr+1)

	return err
}

func (d *DisplayManager) SetDispColon(disp Display, state bool) error {
	switch disp {
	case DisplayOne:
		return d.dispOne.SetColon(state)
	case DisplayTwo:
		return d.dispTwo.SetColon(state)s
	}
	return nil
}

func (d *DisplayManager) DisplayString(disp Display, dStr string) error {
	switch disp {
	case DisplayOne:
		if _, err := d.dispOne.WriteString(dStr); err != nil {
			log.Error().Err(err).Msg("Failed to display string on Display One")
			return err
		}
		break
	case DisplayTwo:
		if _, err := d.dispTwo.WriteString(dStr); err != nil {
			log.Error().Err(err).Msg("Failed to display string on Display Two")
			return err
		}
		break
	}
	return nil
}
