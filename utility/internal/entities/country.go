package entities

// IsoParam represents parameters for ISO codes.
type IsoParam struct {
	Iso string `form:"iso"`
}

// CountryExists represents the existence of a country and missing country codes.
type CountryExists struct {
	Exists              bool     `json:"exists"`
	MissingCountryCodes []string `json:"missing_country_codes,omitempty"`
}

type IsoList struct {
	Iso []string `json:"iso"`
}
