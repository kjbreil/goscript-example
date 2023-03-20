package lights

import (
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/goscript-example/circadian"
)

type circadianAll struct {
	Enabled bool
	c       *circadian.Circadian
}

func (l *Lights) circadianAll() []*goscript.Trigger {
	if l.CircadianAll == nil {
		return nil
	}

	var triggers []*goscript.Trigger

	var ca circadianAll

	ca.c = l.Circadian

	triggers = append(triggers, &goscript.Trigger{
		Unique:        &goscript.Unique{},
		DomainTrigger: nil,
		Periodic:      goscript.Periodics("* * * * *", ""),
		States:        nil,
		DomainStates:  []string{"light"},
		Eval:          nil,
		Func:          ca.changeLightTemperature,
	})
	triggers = append(triggers, &goscript.Trigger{
		Unique:        nil,
		DomainTrigger: []string{"light"},
		Periodic:      nil,
		States:        nil,
		Eval:          goscript.Eval(`state == "on"`),
		Func:          ca.lightTurnedOn,
	})

	return triggers
}

func (ca *circadianAll) changeLightTemperature(t *goscript.Task) {
	ca.c.Calculate()

	var entitiesToChangeBrightness []string
	var entitiesToChangeTemperature []string
	for k, s := range t.States {
		var brightness float64
		if b, ok := s.Attributes["brightness"]; ok {
			brightness = b.(float64)
		}
		if s.State == "on" {
			if brightness > 60 {
				entitiesToChangeBrightness = append(entitiesToChangeBrightness, k)
			} else {
				entitiesToChangeTemperature = append(entitiesToChangeTemperature, k)
			}
		}
	}
	ca.c.TurnOn(entitiesToChangeBrightness...)
	ca.c.TurnOnTemperature(entitiesToChangeTemperature...)
}

func (ca *circadianAll) lightTurnedOn(t *goscript.Task) {
	if oldState := t.Message.OldState(); oldState != nil {
		currentState := t.Message.State()
		if *oldState != currentState {
			var brightness float64

			if b, ok := t.Message.Attributes()["brightness"]; ok {
				brightness = b.(float64)
			}
			if brightness > 60 {
				ca.c.TurnOn(t.Message.DomainEntity())
			} else {
				ca.c.TurnOnTemperature(t.Message.DomainEntity())
			}
		}
	}
}
