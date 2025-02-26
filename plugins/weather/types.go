package weather

type GeocodeResponse struct {
	Status   string    `json:"status"`
	Info     string    `json:"info"`
	InfoCode string    `json:"infocode"`
	Count    string    `json:"count"`
	GeoCodes []Geocode `json:"geocodes"`
}

type WeatherResponse struct {
	Status   string `json:"status"`
	Count    string `json:"count"`
	Info     string `json:"info"`
	InfoCode string `json:"infocode"`
	Lives    []Live `json:"lives"`
}

type Live struct {
	Province         string `json:"province"`
	City             string `json:"city"`
	Adcode           string `json:"adcode"`
	Weather          string `json:"weather"`
	Temperature      string `json:"temperature"`
	WindDirection    string `json:"winddirection"`
	WindPower        string `json:"windpower"`
	Humidity         string `json:"humidity"`
	ReportTime       string `json:"reporttime"`
	TemperatureFloat string `json:"temperature_float"`
	HumidityFloat    string `json:"humidity_float"`
}

type Geocode struct {
	FormattedAddress string       `json:"formatted_address"`
	Country          string       `json:"country"`
	Province         string       `json:"province"`
	CityCode         string       `json:"citycode"`
	City             string       `json:"city"`
	District         string       `json:"district"`
	Township         []string     `json:"township"`
	Neighborhood     Neighborhood `json:"neighborhood"`
	Building         Building     `json:"building"`
	AdCode           string       `json:"adcode"`
	Street           []string     `json:"street"`
	Number           []string     `json:"number"`
	Location         string       `json:"location"`
	Level            string       `json:"level"`
}

type Neighborhood struct {
	Name []string `json:"name"`
	Type []string `json:"type"`
}

type Building struct {
	Name []string `json:"name"`
	Type []string `json:"type"`
}