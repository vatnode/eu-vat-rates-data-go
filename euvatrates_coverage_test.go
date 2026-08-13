package euvatrates

import "testing"

func TestGetRateIsCaseInsensitive(t *testing.T) {
	upper, ok := GetRate("FI")
	if !ok {
		t.Fatal("FI not found in dataset")
	}
	lower, ok := GetRate("fi")
	if !ok {
		t.Fatal("lowercase country code should resolve")
	}
	if upper.Country != lower.Country || upper.Standard != lower.Standard {
		t.Errorf("case-insensitive lookup mismatch: %+v vs %+v", upper, lower)
	}
}

func TestGetStandardRate(t *testing.T) {
	rate, ok := GetStandardRate("DE")
	if !ok {
		t.Fatal("DE not found in dataset")
	}
	if rate != 19 {
		t.Errorf("DE standard rate: got %v, want 19", rate)
	}

	if _, ok := GetStandardRate("XX"); ok {
		t.Error("unknown country should report not found")
	}
}

func TestGetAllRatesReturnsACopy(t *testing.T) {
	all := GetAllRates()
	all["FI"] = VatRate{Country: "Tampered", Standard: 99}

	fresh, ok := GetRate("FI")
	if !ok {
		t.Fatal("FI not found in dataset")
	}
	if fresh.Country == "Tampered" || fresh.Standard == 99 {
		t.Error("mutating the returned map leaked into the package data")
	}
}

func TestHasRate(t *testing.T) {
	for _, code := range []string{"FI", "fi", "GB", "NO"} {
		if !HasRate(code) {
			t.Errorf("HasRate(%q) = false, want true", code)
		}
	}
	for _, code := range []string{"XX", "", "ZZZ"} {
		if HasRate(code) {
			t.Errorf("HasRate(%q) = true, want false", code)
		}
	}
}

func TestValidateFormatAcceptsValidIDs(t *testing.T) {
	cases := []string{
		"ATU12345678",
		"DE123456789",
		"FI12345678",
		"atu12345678", // lowercase input
	}
	for _, id := range cases {
		if !ValidateFormat(id) {
			t.Errorf("ValidateFormat(%q) = false, want true", id)
		}
	}
}

func TestValidateFormatRejectsBadInput(t *testing.T) {
	cases := []string{
		"",             // empty
		"A",            // shorter than a country code
		"INVALID",      // no such country
		"XX123456789",  // unknown country code
		"DE12",         // right country, wrong length
		"ATU1234567",   // one digit short
		"DE1234567890", // one digit too many
	}
	for _, id := range cases {
		if ValidateFormat(id) {
			t.Errorf("ValidateFormat(%q) = true, want false", id)
		}
	}
}

// Greek VAT numbers are issued with the VIES prefix EL while the dataset keys
// Greece under GR. Both spellings have to validate.
func TestValidateFormatHandlesGreece(t *testing.T) {
	if !ValidateFormat("EL123456789") {
		t.Error("ValidateFormat(\"EL123456789\") = false, want true")
	}
	if !ValidateFormat("el123456789") {
		t.Error("lowercase Greek prefix should validate")
	}
	if ValidateFormat("EL12345678") {
		t.Error("a Greek number one digit short should not validate")
	}
}

func TestValidateFormatCoversEveryCountry(t *testing.T) {
	for code, rate := range GetAllRates() {
		if rate.Pattern == "" {
			t.Errorf("%s: pattern is empty", code)
		}
		if rate.Format == "" {
			t.Errorf("%s: format description is empty", code)
		}
	}
}

func TestRawDataset(t *testing.T) {
	ds := RawDataset()
	if ds.Version != DataVersion() {
		t.Errorf("RawDataset().Version = %q, DataVersion() = %q", ds.Version, DataVersion())
	}
	if ds.Source == "" {
		t.Error("RawDataset().Source is empty")
	}
	if len(ds.Rates) != 45 {
		t.Errorf("RawDataset() holds %d countries, want 45", len(ds.Rates))
	}
}

func TestGetFlag(t *testing.T) {
	cases := map[string]string{
		"FI": "🇫🇮",
		"DE": "🇩🇪",
		"GB": "🇬🇧",
		"fi": "🇫🇮",
	}
	for code, want := range cases {
		if got := GetFlag(code); got != want {
			t.Errorf("GetFlag(%q) = %q, want %q", code, got, want)
		}
	}
}

func TestGetFlagRejectsBadInput(t *testing.T) {
	for _, code := range []string{"", "F", "FIN", "F1", "12", "!!"} {
		if got := GetFlag(code); got != "" {
			t.Errorf("GetFlag(%q) = %q, want empty string", code, got)
		}
	}
}

func TestNonEUCountriesArePresent(t *testing.T) {
	for _, code := range []string{"GB", "NO", "CH", "UA", "TR"} {
		rate, ok := GetRate(code)
		if !ok {
			t.Errorf("%s missing from dataset", code)
			continue
		}
		if rate.EUMember {
			t.Errorf("%s: EUMember = true, want false", code)
		}
	}
}

func TestEU27CountAndCurrencies(t *testing.T) {
	euCount := 0
	for code, rate := range GetAllRates() {
		if rate.EUMember {
			euCount++
		}
		if rate.Currency == "" {
			t.Errorf("%s: currency is empty", code)
		}
	}
	if euCount != 27 {
		t.Errorf("EU member count: got %d, want 27", euCount)
	}
}
