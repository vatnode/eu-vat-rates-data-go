# eu-vat-rates-data-go

[![Go Reference](https://pkg.go.dev/badge/github.com/vatnode/eu-vat-rates-data-go.svg)](https://pkg.go.dev/github.com/vatnode/eu-vat-rates-data-go)
[![Last updated](https://img.shields.io/github/last-commit/vatnode/eu-vat-rates-data-go?path=eu-vat-rates-data.json&label=last%20updated)](https://github.com/vatnode/eu-vat-rates-data-go/commits/main)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

VAT rates for **45 European countries** — EU-27 plus Norway, Switzerland, UK, and more. EU rates sourced from the European Commission TEDB and checked daily. Non-EU rates maintained manually.

- Standard, reduced, super-reduced, and parking rates
- `EUMember` field on every country — `true` for EU-27, `false` for non-EU
- `VATName` — official name of the VAT tax in the country's primary official language
- `VATAbbr` — short abbreviation used locally (e.g. "ALV", "MwSt", "TVA")
- **`Format` — human-readable VAT number format (e.g. `"ATU + 8 digits"`)** — unique to this package
- **`Pattern` — regex for VAT number validation + built-in `ValidateFormat()` — free, no API key needed** — unique to this package
- Zero dependencies — data embedded with `//go:embed`
- Fully typed — works with Go 1.21+
- EU rates checked daily via GitHub Actions, new version tagged only when rates change

Also available in: [JavaScript/TypeScript (npm)](https://www.npmjs.com/package/eu-vat-rates-data) · [Python (PyPI)](https://pypi.org/project/eu-vat-rates-data/) · [PHP (Packagist)](https://packagist.org/packages/vatnode/eu-vat-rates-data) · [Ruby (RubyGems)](https://rubygems.org/gems/eu_vat_rates_data)

---

## Need live VIES validation?

This package gives you VAT **rates** and **format checks** for free, offline, in your code. It does **not** call VIES — `ValidateFormat()` only checks the shape of a VAT number, not whether it actually exists.

For **live VIES validation** — confirming a VAT ID is real, pulling the registered company name and address, and getting a VIES consultation number (audit-grade proof of validation) — there's **[vatnode](https://vatnode.dev)**:

- Live VIES validation, with national-database fallback when VIES is down
- Registered company name, address, registration date
- VIES consultation number for compliance and audit trails
- Webhooks for VAT status changes
- Official [MCP server](https://www.npmjs.com/package/vatnode-mcp) so AI agents (Claude, Cursor, ChatGPT) can validate VAT IDs directly
- Free tier — no credit card needed

```bash
curl https://api.vatnode.dev/v1/vat/IE6388047V \
  -H "Authorization: Bearer YOUR_API_KEY"
```

[**Get a free API key →**](https://vatnode.dev)

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

    // Dataset membership check (all 45 countries)
    if euvatrates.HasRate(userInput) {
        rate, _ := euvatrates.GetRate(userInput)
        _ = rate
    }

    // All 45 countries
    all := euvatrates.GetAllRates()
    for code, rate := range all {
        fmt.Printf("%s: %.1f%%\n", code, rate.Standard)
    }

    // When were EU rates last fetched?
    fmt.Println(euvatrates.DataVersion()) // "2026-03-27"

    // VAT number format validation — no API key, no network call
    euvatrates.ValidateFormat("ATU12345678") // → true
    euvatrates.ValidateFormat("DE123456789") // → true
    euvatrates.ValidateFormat("INVALID")     // → false

    // Access format metadata directly
    at, _ := euvatrates.GetRate("AT")
    fmt.Println(at.Format)   // "ATU + 8 digits"
    fmt.Println(at.Pattern) // "^ATU\\d{8}$"

    // Flag emoji from a 2-letter country code — no lookup table, computed from regional indicator symbols
    fmt.Println(euvatrates.GetFlag("FI")) // "🇫🇮"
    fmt.Println(euvatrates.GetFlag("DE")) // "🇩🇪"
    fmt.Println(euvatrates.GetFlag("XX")) // "" (empty string for unknown/invalid codes)
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
    Format       string    `json:"format"`   // "ATU + 8 digits"
    Pattern      string    `json:"pattern"`  // "^ATU\\d{8}$" — always present for all 45 countries
}
```

---

## Data source & update frequency

- EU-27 rates: **European Commission TEDB**, refreshed **daily at 07:00 UTC**
- Non-EU rates: maintained manually, updated on official rate changes
- New git tag + pkg.go.dev version published only when rates change

---


## Keeping rates current

Rates are bundled at install time. A new package version is published automatically whenever rates change — but your installed version will not update itself.

**Recommended:** add [Renovate](https://renovatebot.com) or [Dependabot](https://docs.github.com/en/code-security/dependabot) to your repo. They detect new versions and open a PR automatically whenever rates change — no manual update commands needed.

**Need real-time accuracy?** Fetch the always-current JSON directly:

```
https://cdn.jsdelivr.net/gh/vatnode/eu-vat-rates-data@main/data/eu-vat-rates-data.json
```

No package needed — parse it with a single `fetch()` / `http.get()` / `file_get_contents()` call and cache locally.

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

## Changelog

### 2026-04-25
- **fix:** Corrected Sweden (SE) VAT number regex — was `^SE\d{12}$`, now correctly requires the mandatory `01` suffix: `^SE\d{10}01$`.

---

## License

MIT

If you find this useful, a ⭐ on GitHub is appreciated.
