# eu-vat-rates-data-go

[![Go Reference](https://pkg.go.dev/badge/github.com/vatnode/eu-vat-rates-data-go.svg)](https://pkg.go.dev/github.com/vatnode/eu-vat-rates-data-go)
[![Last updated](https://img.shields.io/github/last-commit/vatnode/eu-vat-rates-data-go?path=eu-vat-rates.json&label=last%20updated)](https://github.com/vatnode/eu-vat-rates-data-go/commits/main)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

EU VAT rates for all **27 EU member states** plus the **United Kingdom**, sourced from the [European Commission TEDB](https://taxation-customs.ec.europa.eu/tedb/vatRates.html). Checked daily, published automatically when rates change.

- Standard, reduced, super-reduced, and parking rates
- Zero dependencies — data embedded with `//go:embed`
- Fully typed — works with Go 1.21+
- Checked daily via GitHub Actions, new version tagged only when rates change

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
        fmt.Printf("%s: %.1f%%\n", fi.Country, fi.Standard)
        // Finland: 25.5%
    }

    // Just the standard rate
    standard, ok := euvatrates.GetStandardRate("DE")
    fmt.Println(standard) // 19

    // Type guard
    if euvatrates.IsEUMember(userInput) {
        rate, _ := euvatrates.GetRate(userInput)
        _ = rate // guaranteed to be valid
    }

    // All 28 countries
    all := euvatrates.GetAllRates()
    for code, rate := range all {
        fmt.Printf("%s: %.1f%%\n", code, rate.Standard)
    }

    // Data version date
    fmt.Println(euvatrates.DataVersion()) // "2026-02-25"
}
```

---

## Types

```go
type VatRate struct {
    Country      string    `json:"country"`
    Currency     string    `json:"currency"`
    Standard     float64   `json:"standard"`
    Reduced      []float64 `json:"reduced"`
    SuperReduced *float64  `json:"super_reduced"`
    Parking      *float64  `json:"parking"`
}
```

---

## Data source & update frequency

Rates are fetched from the **European Commission Taxes in Europe Database (TEDB)**:

- Canonical data repo: **https://github.com/vatnode/eu-vat-rates-data-js**
- Refreshed: **daily at 08:00 UTC**
- New git tag + pkg.go.dev version published only when rates change

---

## Covered countries

EU-27 member states + United Kingdom (28 countries total):

`AT BE BG CY CZ DE DK EE ES FI FR GB GR HR HU IE IT LT LU LV MT NL PL PT RO SE SI SK`

---

## Need to validate VAT numbers?

This package provides **VAT rates** only. If you also need to **validate EU VAT numbers** against the official VIES database — confirming a business is VAT-registered — check out [vatnode.dev](https://vatnode.dev), a simple REST API with a free tier.

```bash
curl https://api.vatnode.dev/v1/vat/FI17156132 \
  -H "Authorization: Bearer vat_live_..."
# → { "valid": true, "companyName": "Suomen Pehmeä Ikkuna Oy" }
```

## License

MIT

If you find this useful, a ⭐ on GitHub is appreciated.
