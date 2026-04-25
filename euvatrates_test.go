package euvatrates

import (
	"regexp"
	"testing"
)

func TestDeIsEuMember(t *testing.T) {
	if !IsEUMember("DE") {
		t.Error("DE should be an EU member")
	}
}

func TestGbIsNotEuMember(t *testing.T) {
	if IsEUMember("GB") {
		t.Error("GB should not be an EU member")
	}
}

func TestNoIsNotEuMember(t *testing.T) {
	if IsEUMember("NO") {
		t.Error("NO should not be an EU member")
	}
}

func TestDatasetSize(t *testing.T) {
	all := GetAllRates()
	if len(all) != 45 {
		t.Errorf("dataset size: got %d, want 45", len(all))
	}
}

func TestAllStandardRatesPositive(t *testing.T) {
	for code, rate := range GetAllRates() {
		if rate.Standard <= 0 {
			t.Errorf("%s: standard rate is %v, want > 0", code, rate.Standard)
		}
	}
}

func TestEuMemberFieldPresent(t *testing.T) {
	de, ok := GetRate("DE")
	if !ok {
		t.Fatal("DE not found in dataset")
	}
	if !de.EUMember {
		t.Error("DE.EUMember should be true")
	}
	no, ok := GetRate("NO")
	if !ok {
		t.Fatal("NO not found in dataset")
	}
	if no.EUMember {
		t.Error("NO.EUMember should be false")
	}
}

func TestAllVatNamesNonEmpty(t *testing.T) {
	for code, rate := range GetAllRates() {
		if rate.VATName == "" {
			t.Errorf("%s: VATName is empty", code)
		}
	}
}

func TestDataVersionFormat(t *testing.T) {
	v := DataVersion()
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, v)
	if !matched {
		t.Errorf("DataVersion() = %q, want YYYY-MM-DD format", v)
	}
}

func TestUnknownCountryNotFound(t *testing.T) {
	_, ok := GetRate("XX")
	if ok {
		t.Error("XX should not be found in dataset")
	}
}
