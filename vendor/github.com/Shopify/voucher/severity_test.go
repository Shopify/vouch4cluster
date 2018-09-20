package voucher

import "testing"

var testSeverities = map[string]Severity{
	negligibleSeverityString: NegligibleSeverity,
	lowSeverityString:        LowSeverity,
	mediumSeverityString:     MediumSeverity,
	unknownSeverityString:    UnknownSeverity,
	highSeverityString:       HighSeverity,
	criticalSeverityString:   CriticalSeverity,
	"whatever":               UnknownSeverity,
}

func TestSeverityToString(t *testing.T) {
	for expected, severity := range testSeverities {
		if "whatever" == expected {
			continue
		}
		value := severity.String()
		if value != expected {
			t.Errorf("Severity.String() returned the wrong output, should be: %v, was %v", expected, value)
		}
	}

}

func TestStringToSeverity(t *testing.T) {
	for name, expected := range testSeverities {
		value, err := StringToSeverity(name)
		if nil != err && "whatever" != name {
			t.Errorf("got error converting severities: %s", err)
			continue
		}
		if value != expected {
			t.Errorf("StringToSeverity returned the wrong Severity, should be: %v, was %v", expected, value)
		}
		if "whatever" == name && err == nil {
			t.Errorf("should have gotten an error due to unknown severity")
		}
	}
}
