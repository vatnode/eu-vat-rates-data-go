// Package euvatrates provides EU VAT rates for all 27 member states + UK.
//
// Data is sourced from the European Commission TEDB (Taxes in Europe Database)
// and embedded at compile time. Published automatically when rates change.
//
// Usage:
//
//	import euvatrates "github.com/vatnode/eu-vat-rates-go"
//
//	rate, ok := euvatrates.GetRate("FI")
//	// rate.Standard == 25.5, rate.Country == "Finland"
//
//	standard, ok := euvatrates.GetStandardRate("DE")  // 19.0, true
//	euvatrates.IsEUMember("FR")                        // true
//	euvatrates.DataVersion()                           // "2026-02-25"
package euvatrates

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed eu-vat-rates.json
var rawData []byte

// VatRate holds all VAT rates for a single country.
type VatRate struct {
	Country      string    `json:"country"`
	Currency     string    `json:"currency"`
	Standard     float64   `json:"standard"`
	Reduced      []float64 `json:"reduced"`
	SuperReduced *float64  `json:"super_reduced"`
	Parking      *float64  `json:"parking"`
}

// Dataset is the top-level structure of the data file.
type Dataset struct {
	Version string             `json:"version"`
	Source  string             `json:"source"`
	URL     string             `json:"url"`
	Rates   map[string]VatRate `json:"rates"`
}

var data Dataset

func init() {
	if err := json.Unmarshal(rawData, &data); err != nil {
		panic("eu-vat-rates-go: failed to parse embedded data: " + err.Error())
	}
}

// GetRate returns the full VatRate for the given ISO 3166-1 alpha-2 country code.
// The second return value is false if the country is not in the dataset.
func GetRate(countryCode string) (VatRate, bool) {
	rate, ok := data.Rates[strings.ToUpper(countryCode)]
	return rate, ok
}

// GetStandardRate returns the standard VAT rate for the given country code.
// The second return value is false if the country is not in the dataset.
func GetStandardRate(countryCode string) (float64, bool) {
	rate, ok := GetRate(countryCode)
	if !ok {
		return 0, false
	}
	return rate.Standard, true
}

// GetAllRates returns a copy of the full rates map (28 countries).
func GetAllRates() map[string]VatRate {
	out := make(map[string]VatRate, len(data.Rates))
	for k, v := range data.Rates {
		out[k] = v
	}
	return out
}

// IsEUMember returns true if the country code is in the dataset (EU-27 + GB).
func IsEUMember(countryCode string) bool {
	_, ok := data.Rates[strings.ToUpper(countryCode)]
	return ok
}

// DataVersion returns the ISO 8601 date when the data was last fetched from EC TEDB.
func DataVersion() string {
	return data.Version
}

// RawDataset returns the full parsed Dataset struct.
func RawDataset() Dataset {
	return data
}
