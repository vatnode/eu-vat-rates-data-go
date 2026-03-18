package euvatrates

import (
	"regexp"
	"testing"
)

func TestDeStandardRate(t *testing.T) {
	rate, ok := GetRate("DE")
	if !ok {
		t.Fatal("DE not found in dataset")
	}
	if rate.Standard != 19.0 {
		t.Errorf("DE standard rate: got %v, want 19.0", rate.Standard)
	}
}

func TestEeStandardRate(t *testing.T) {
	rate, ok := GetRate("EE")
	if !ok {
		t.Fatal("EE not found in dataset")
	}
	if rate.Standard != 24.0 {
		t.Errorf("EE standard rate: got %v, want 24.0", rate.Standard)
	}
}

func TestFrIsEuMember(t *testing.T) {
	if !IsEUMember("FR") {
		t.Error("FR should be an EU member")
	}
}

func TestGbIsNotEuMember(t *testing.T) {
	if IsEUMember("GB") {
		t.Error("GB should not be an EU member")
	}
}

func TestEuMemberField(t *testing.T) {
	de, _ := GetRate("DE")
	if !de.EUMember {
		t.Error("DE.EUMember should be true")
	}
	no, _ := GetRate("NO")
	if no.EUMember {
		t.Error("NO.EUMember should be false")
	}
}

func TestDatasetSize(t *testing.T) {
	all := GetAllRates()
	if len(all) != 44 {
		t.Errorf("dataset size: got %d, want 44", len(all))
	}
}

func TestDataVersionFormat(t *testing.T) {
	v := DataVersion()
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, v)
	if !matched {
		t.Errorf("DataVersion() = %q, want YYYY-MM-DD format", v)
	}
}
