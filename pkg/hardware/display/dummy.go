package display

import (
	"context"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	log "github.com/sirupsen/logrus"
)

const TypeDummy = "dummy"

type DummySettings struct{}

type Dummy struct {
	logger   *log.Logger
	settings DummySettings
}

func NewDummy(settings *DummySettings) (*Dummy, error) {
	return &Dummy{
		settings: *settings,
	}, nil
}

func (d *Dummy) DisplayMessage(message display.MessageInfo) error {
	d.logger.WithFields(log.Fields{"message": message}).Info("Displaying message")
	return nil
}

func (d *Dummy) Clear() error {
	d.logger.Info("Clearing display")
	return nil
}

func (d *Dummy) Cleanup(context.Context) error {
	d.logger.Info("Cleaning up display")
	return nil
}

func (d *Dummy) GetType() string {
	return TypeDummy
}
