# eu-vat-rates-data-go

[![Go Reference](https://pkg.go.dev/badge/github.com/vatnode/eu-vat-rates-data-go.svg)](https://pkg.go.dev/github.com/vatnode/eu-vat-rates-data-go)
[![Last updated](https://img.shields.io/github/last-commit/vatnode/eu-vat-rates-data-go?path=eu-vat-rates-data.json&label=last%20updated)](https://github.com/vatnode/eu-vat-rates-data-go/commits/main)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

VAT rates for **44 European countries** — EU-27 plus Norway, Switzerland, UK, and more. EU rates sourced from the European Commission TEDB and checked daily. Non-EU rates maintained manually.

- Standard, reduced, super-reduced, and parking rates
- `EUMember` field on every country — `true` for EU-27, `false` for non-EU
- `VATName` — official name of the VAT tax in the country's primary official language
- `VATAbbr` — short abbreviation used locally (e.g. "ALV", "MwSt", "TVA")
- Zero dependencies — data embedded with `//go:embed`
- Fully typed — works with Go 1.21+
- EU rates checked daily via GitHub Actions, new version tagged only when rates change

Also available in: [JavaScript/TypeScript (npm)](https://www.npmjs.com/package/eu-vat-rates-data) · [Python (PyPI)](https://pypi.org/project/eu-vat-rates-data/) · [PHP (Packagist)](https://packagist.org/packages/vatnode/eu-vat-rates-data) · [Ruby (RubyGems)](https://rubygems.org/gems/eu_vat_rates_data)

---

## Installation

```bash
go get github.com/vatnode/eu-vat-rates-data-go
```

---

## Usage

```go
package main

import (
    "fmt"
    euvatrates "github.com/vatnode/eu-vat-rates-data-go"
)

func main() {
    // Full rate struct for a country
    fi, ok := euvatrates.GetRate("FI")
    if ok {
        fmt.Printf("%s: %.1f%% (EU member: %v, tax: %s / %s)\n", fi.Country, fi.Standard, fi.EUMember, fi.VATName, fi.VATAbbr)
        // Finland: 25.5% (EU member: true, tax: Arvonlisävero / ALV)
    }

    // Just the standard rate
    standard, ok := euvatrates.GetStandardRate("DE")
    fmt.Println(standard) // 19

    // EU membership check — false for non-EU countries (GB, NO, CH, ...)
    if euvatrates.IsEUMember(userInput) {
        rate, _ := euvatrates.GetRate(userInput)
        _ = rate
    }

    // Dataset membership check (all 44 countries)
    if euvatrates.HasRate(userInput) {
        rate, _ := euvatrates.GetRate(userInput)
        _ = rate
    }

    // All 44 countries
    all := euvatrates.GetAllRates()
    for code, rate := range all {
        fmt.Printf("%s: %.1f%%\n", code, rate.Standard)
    }

    // When were EU rates last fetched?
    fmt.Println(euvatrates.DataVersion()) // "2026-03-27"
}
```

---

## Types

```go
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
}
```

---

## Data source & update frequency

- EU-27 rates: **European Commission TEDB**, refreshed **daily at 07:00 UTC**
- Non-EU rates: maintained manually, updated on official rate changes
- New git tag + pkg.go.dev version published only when rates change

---

## Covered countries

**EU-27** (daily auto-updates via EC TEDB):

`AT` `BE` `BG` `CY` `CZ` `DE` `DK` `EE` `ES` `FI` `FR` `GR` `HR` `HU` `IE` `IT` `LT` `LU` `LV` `MT` `NL` `PL` `PT` `RO` `SE` `SI` `SK`

**Non-EU Europe** (manually maintained):

`AD` `AL` `BA` `CH` `GB` `GE` `IS` `LI` `MC` `MD` `ME` `MK` `NO` `RS` `TR` `UA` `XK`

---

## Need to validate VAT numbers?

This package provides **VAT rates** only. If you also need to **validate EU VAT numbers** against the official VIES database — confirming a business is VAT-registered — check out [vatnode.dev](https://vatnode.dev), a simple REST API with a free tier.

```bash
curl https://api.vatnode.dev/v1/vat/FI17156132 \
  -H "Authorization: Bearer vat_live_..."
# → { "valid": true, "companyName": "Suomen Pehmeä Ikkuna Oy" }
```

---

## License

MIT

If you find this useful, a ⭐ on GitHub is appreciated.
