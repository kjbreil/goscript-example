package lights

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/goscript-example/circadian"
	"strings"
	"time"
)

type motionLight struct {
	circadian *circadian.Circadian
	logger    logr.Logger

	triggers  []*goscript.Trigger
	service   goscript.ServiceChan
	TurnOn    bool
	Timeout   int
	BlockIfOn []string
	Detectors []string
	Entities  []string
}

func (l *Lights) motion() []*goscript.Trigger {
	var triggers []*goscript.Trigger

	for _, ml := range l.MotionLights {
		ml.circadian = l.Circadian
		ml.logger = l.logger
		triggers = append(triggers, motionLightsTrigger(ml)...)
	}

	return triggers
}

func motionLightsTrigger(ml motionLight) []*goscript.Trigger {
	ml.addLightsTrigger()
	return ml.triggers
}
func (l *motionLight) addLightsTrigger() {
	l.triggers = append(l.triggers, &goscript.Trigger{
		Triggers: l.Detectors,
		Unique:   &goscript.Unique{},
		States:   append(l.Entities, l.BlockIfOn...),
		Eval:     goscript.Eval(`state == "on"`),
		Func:     l.turnOnLights,
	})
}

func (l *motionLight) turnOnLights(t *goscript.Task) {
	turnOn := l.TurnOn
	l.logger.Info(fmt.Sprintf("Motion Detected: %s", t.Message.DomainEntity()))
	for _, e := range l.BlockIfOn {
		for _, s := range t.States.Slice() {
			if e == s.DomainEntity && s.State == "on" {
				turnOn = false
			}
		}
	}
	if turnOn {
		l.logger.Info(fmt.Sprintf("Turning Lights on: %s", strings.Join(l.Entities, ",")))
		l.circadian.TurnOn(l.Entities...)
	} else {
		allOff := true
		for _, s := range t.States.Slice() {
			if s.Domain == "light" && s.State == "on" {
				allOff = false
			}
		}
		if allOff {
			return
		}
	}
	l.logger.Info(fmt.Sprintf("Waiting for no motion on %s", t.Message.DomainEntity()))
	t.WaitUntil(t.Message.DomainEntity(), goscript.Eval(`state == "off"`), 0)

	l.logger.Info(fmt.Sprintf("Sleeping for %d entity: %s", l.Timeout, t.Message.DomainEntity()))
	t.Sleep(time.Duration(l.Timeout) * time.Second)

	l.logger.Info(fmt.Sprintf("Turning off lights: %s", strings.Join(l.Entities, ",")))
	l.circadian.TurnOff(l.Entities...)
}
