package circadian

import (
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/hass-ws/services"
	"github.com/sixdouglas/suncalc"
	"math"
	"time"
)

type Circadian struct {
	Lat              float64
	Long             float64
	MinTemperature   float64
	MaxTemperature   float64
	MinBrightnessPct float64
	MaxBrightnessPct float64

	transition float64

	currentTemperature   float64
	currentBrightnessPct float64

	service goscript.ServiceChan
}

func (c *Circadian) ChangeEnough(temperature, brightness float64) bool {
	brightnessStepSize := (c.MaxBrightnessPct - c.MinBrightnessPct) / 100
	temperatureStepSize := (c.MaxTemperature - c.MinTemperature) / 100

	return math.Abs(temperature-c.currentTemperature) <= temperatureStepSize &&
		math.Abs(brightness-c.currentBrightnessPct) <= brightnessStepSize

}

func (c *Circadian) AddGoscript(gs *goscript.GoScript) {
	c.service = gs.ServiceChan
	c.transition = 0.5
}

func (c *Circadian) Calculate() (float64, float64) {
	now := time.Now()
	times := suncalc.GetTimes(now, c.Lat, c.Long)

	currentPos := suncalc.GetPosition(now, c.Lat, c.Long)
	noonPos := suncalc.GetPosition(times[suncalc.SolarNoon].Value, c.Lat, c.Long)

	//l.gs.CallService(services.NewInputBooleanToggle([]string{"input_boolean.test_toggle"}))

	//azPct := (currentPos.Azimuth * 180 / math.Pi) / (noonPos.Azimuth * 180 / math.Pi)
	//altPct := (currentPos.Altitude / noonPos.Altitude) * 100
	altPct := math.Round((currentPos.Altitude*180/math.Pi)/(noonPos.Altitude*180/math.Pi)*10000) / 100

	c.currentTemperature = mapRange(0, 100, c.MaxTemperature, c.MinTemperature, altPct)
	c.currentBrightnessPct = mapRange(0, 100, c.MinBrightnessPct, c.MaxBrightnessPct, altPct)
	c.within()

	return c.currentTemperature, c.currentBrightnessPct
}

func (c *Circadian) TurnOn(entities ...string) {
	if len(entities) == 0 {
		return
	}
	c.Calculate()

	c.service <- services.NewLightTurnOn(services.Targets(entities...), &services.LightTurnOnParams{
		BrightnessPct: &c.currentBrightnessPct,
		ColorTemp:     &c.currentTemperature,
		Transition:    &c.transition,
	})
}

func (c *Circadian) TurnOnTemperature(entities ...string) {
	if len(entities) == 0 {
		return
	}
	c.Calculate()

	c.service <- services.NewLightTurnOn(services.Targets(entities...), &services.LightTurnOnParams{
		ColorTemp:  &c.currentTemperature,
		Transition: &c.transition,
	})
}

func (c *Circadian) TurnOff(entities ...string) {
	c.service <- services.NewLightTurnOff(services.Targets(entities...), &services.LightTurnOffParams{
		Transition: &c.transition,
	})
}

func mapRange(rangeLow, rangeHigh, mapLow, mapHigh, value float64) float64 {
	return mapLow + ((value - rangeLow) / (rangeHigh - (rangeLow)) * (mapHigh - (mapLow)))
}

func (c *Circadian) within() {
	if c.currentTemperature > c.MaxTemperature {
		c.currentTemperature = c.MaxTemperature
	}
	if c.currentTemperature < c.MinTemperature {
		c.currentTemperature = c.MinTemperature
	}
	if c.currentBrightnessPct > c.MaxBrightnessPct {
		c.currentBrightnessPct = c.MaxBrightnessPct
	}
	if c.currentBrightnessPct < c.MinBrightnessPct {
		c.currentBrightnessPct = c.MinBrightnessPct
	}
}
