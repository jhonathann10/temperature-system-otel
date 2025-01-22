package entity

type Weather struct {
	Localidade string
	Celsius    float64
	Fahrenheit float64
	Kelvin     float64
}

func NewWeather(localidade string, celsius float64) (*Weather, error) {
	weather := &Weather{
		Localidade: localidade,
		Celsius:    celsius,
	}
	weather.CalculateFahrenheit()
	weather.CalculateKelvin()

	return weather, nil
}

func (w *Weather) CalculateFahrenheit() {
	w.Fahrenheit = (w.Celsius * 1.8) + 32
}

func (w *Weather) CalculateKelvin() {
	w.Kelvin = w.Celsius + 273.15
}
