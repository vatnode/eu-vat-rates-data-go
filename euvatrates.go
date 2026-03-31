// Package euvatrates provides VAT rates for 44 European countries (EU-27 + 17 non-EU).
//
// EU rates are sourced from the European Commission TEDB (Taxes in Europe Database)
// and embedded at compile time. Non-EU rates are maintained manually.
//
// Usage:
//
//	import euvatrates "github.com/vatnode/eu-vat-rates-data-go"
//
//	rate, ok := euvatrates.GetRate("FI")
//	// rate.Standard == 25.5, rate.Country == "Finland", rate.EUMember == true
//
//	standard, ok := euvatrates.GetStandardRate("DE")  // 19.0, true
//	euvatrates.IsEUMember("NO")                        // false
//	euvatrates.IsEUMember("FR")                        // true
//	euvatrates.DataVersion()                           // "2026-03-18"
package euvatrates

import (
	_ "embed"
	"encoding/json"
	"regexp"
	"strings"
)

//go:embed eu-vat-rates-data.json
var rawData []byte

// VatRate holds all VAT rates for a single country.
type VatRate struct {
	Country      string    `json:"country"`
	Currency     string    `json:"currency"`
	EUMember     bool      `json:"eu_member"`
	VATName      string    `json:"vat_name"`
	VATAbbr      string    `json:"vat_abbr"`
	Standard     float64   `json:"standard"`
	Reduced      []float64 `json:"reduced"`
	SuperReduced *float64  `json:"super_reduced"`
	Parking      *float64  `json:"parking"`
	Format       string    `json:"format"`
	Pattern      *string   `json:"pattern"`
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

// GetAllRates returns a copy of the full rates map (44 countries).
func GetAllRates() map[string]VatRate {
	out := make(map[string]VatRate, len(data.Rates))
	for k, v := range data.Rates {
		out[k] = v
	}
	return out
}

// IsEUMember returns true if the country is an EU-27 member state.
// Returns false for non-EU countries in the dataset (GB, NO, CH, etc.)
// and for unknown country codes.
func IsEUMember(countryCode string) bool {
	rate, ok := data.Rates[strings.ToUpper(countryCode)]
	return ok && rate.EUMember
}

// HasRate returns true if the given country code is present in the dataset (all 44 countries).
// Use IsEUMember to check EU membership specifically.
func HasRate(countryCode string) bool {
	_, ok := data.Rates[strings.ToUpper(countryCode)]
	return ok
}

// DataVersion returns the ISO 8601 date when the EU data was last fetched from EC TEDB.
func DataVersion() string {
	return data.Version
}

// ValidateFormat returns true if vatID matches the expected format for its country.
// Input must include the country code prefix (e.g. "ATU12345678").
// Returns false when the country has no standardised format or the ID does not match.
// Note: Greece uses the "EL" prefix, not "GR".
func ValidateFormat(vatID string) bool {
	if len(vatID) < 2 {
		return false
	}
	code := strings.ToUpper(vatID[:2])
	rate, ok := data.Rates[code]
	if !ok || rate.Pattern == nil {
		return false
	}
	matched, err := regexp.MatchString(*rate.Pattern, strings.ToUpper(vatID))
	return err == nil && matched
}

// RawDataset returns the full parsed Dataset struct.
func RawDataset() Dataset {
	return data
}

// GetFlag returns the flag emoji for a 2-letter ISO 3166-1 alpha-2 country code.
// Computed from regional indicator symbols — no lookup table needed.
// Returns an empty string if the input is not exactly 2 ASCII letters.
//
// Example:
//
//	GetFlag("FI") // "🇫🇮"
//	GetFlag("DE") // "🇩🇪"
//	GetFlag("GB") // "🇬🇧"
func GetFlag(countryCode string) string {
	code := strings.ToUpper(countryCode)
	if len(code) != 2 {
		return ""
	}
	a, b := rune(code[0]), rune(code[1])
	if a < 'A' || a > 'Z' || b < 'A' || b > 'Z' {
		return ""
	}
	const base = 0x1F1E6
	return string([]rune{base + a - 'A', base + b - 'A'})
}
