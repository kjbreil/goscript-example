package lights

import (
	"github.com/go-logr/logr"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/goscript-example/circadian"
)

var key = "lights"

type Lights struct {
	Circadian    *circadian.Circadian
	CircadianAll *circadianAll
	MotionLights map[string]motionLight
	service      goscript.ServiceChan
	logger       logr.Logger
}

func Triggers(gs *goscript.GoScript) []*goscript.Trigger {
	inter, err := gs.GetModule(key)
	if err != nil {
		return nil
	}
	l := inter.(*Lights)
	l.service = gs.ServiceChan
	l.logger = gs.Logger()

	l.Circadian.AddGoscript(gs)

	var triggers []*goscript.Trigger
	// Set up the motion lights
	triggers = append(triggers, l.motion()...)
	triggers = append(triggers, l.circadianAll()...)

	return triggers
}
