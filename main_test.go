package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestTimeDiv(t *testing.T) {
	midnightUTC := time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		midnight      time.Time
		text          string
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			midnight:      midnightUTC,
			text:          " 09:00 - 11:00AM (2.0h) ", // spaces
			expectedStart: time.Date(2019, 4, 1, 9, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2019, 4, 1, 11, 0, 0, 0, time.UTC),
		},
		{
			midnight:      midnightUTC,
			text:          "09:00 - 11:00AM",
			expectedStart: time.Date(2019, 4, 1, 9, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2019, 4, 1, 11, 0, 0, 0, time.UTC),
		},
		{
			midnight:      midnightUTC,
			text:          "11:00 - 12:00PM",
			expectedStart: time.Date(2019, 4, 1, 11, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2019, 4, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			midnight:      midnightUTC,
			text:          "11:30 - 12:30PM",
			expectedStart: time.Date(2019, 4, 1, 11, 30, 0, 0, time.UTC),
			expectedEnd:   time.Date(2019, 4, 1, 12, 30, 0, 0, time.UTC),
		},
		{
			midnight:      midnightUTC,
			text:          "10:00 - 01:00PM",
			expectedStart: time.Date(2019, 4, 1, 10, 0, 0, 0, time.UTC),
			expectedEnd:   time.Date(2019, 4, 1, 13, 0, 0, 0, time.UTC),
		},
		{
			midnight:      midnightUTC,
			text:          "06:30 - 08:30PM",
			expectedStart: time.Date(2019, 4, 1, 18, 30, 0, 0, time.UTC),
			expectedEnd:   time.Date(2019, 4, 1, 20, 30, 0, 0, time.UTC),
		},
		{
			midnight:      midnightUTC,
			text:          "06:30 - 08:30PM",
			expectedStart: time.Date(2019, 4, 1, 18, 30, 0, 0, time.UTC),
			expectedEnd:   time.Date(2019, 4, 1, 20, 30, 0, 0, time.UTC),
		},
	}

	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			start, end, err := parseTimeDiv(test.text, test.midnight)
			assertNoError(t, err)

			if start == nil {
				t.Fatalf("got nil start")
			} else if end == nil {
				t.Fatalf("got nil end")
			}

			if !test.expectedStart.Equal(*start) {
				t.Errorf("expected start: %s, got %s", test.expectedStart, *start)
			}

			if !test.expectedEnd.Equal(*end) {
				t.Errorf("expected end: %s, got %s", test.expectedEnd, *end)
			}
		})
	}

}

// Equal aims to test equality of any two objects, and call t.Fatalf if they're not equal
func assertEqual(t *testing.T, expected interface{}, got interface{}) {
	t.Helper()

	if !isEqual(expected, got) {
		msg := fmt.Sprintf("expected `%v` got `%v`", expected, got)
		if len(msg) < 50 {
			t.Fatalf(msg)
		} else {
			t.Fatalf("\n--- expected ---\n%v\n--- got ---\n%v\n--- end ---\n", expected, got)
		}
	}
}

func isEqual(a interface{}, b interface{}) bool {
	if aAsError, ok := a.(error); ok {
		if bAsError, ok := b.(error); ok {
			if reflect.TypeOf(a) != reflect.TypeOf(b) {
				return false
			}
			return aAsError.Error() == bAsError.Error() // compare on error string
		}
	}

	return reflect.DeepEqual(a, b)
}

// NoError tests that got is nil and calls t.Fatal if not
func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didnt want one '%s'", got)
	}
}

// GotError tests that got is an error (not nil) and calls t.Fatal if not
func assertGotError(t *testing.T, got error) {
	t.Helper()
	if got == nil {
		t.Fatalf("expected an error, but got none")
	}
}
