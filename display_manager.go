package seven_seg_display_manager

import (
	"github.com/rs/zerolog/log"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

type DisplayManager struct {
	hltDisplay *NumericDisplay
}

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

	d.hltDisplay, err = NewNumericDisplay(bus, Ht16k33I2CAddr+0)

	return err
}

func (d *DisplayManager) SetHltColon(state bool) error {
	return d.hltDisplay.SetColon(state)
}

func (d *DisplayManager) DisplayStringHlt(dStr string) error {
	if _, err := d.hltDisplay.WriteString(dStr); err != nil {
		log.Error().Err(err).Msg("Failed to display string on HLT")
		return err
	}
	return nil
}
