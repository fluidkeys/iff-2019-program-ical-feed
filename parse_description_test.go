package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDescription(t *testing.T) {
	html := loadHTMLTestData(t, "session.html")

	gotFields, err := parseSessionFields(html)
	assertNoError(t, err)

	expectedFields := []sessionField{
		{
			key:   "Date",
			value: "Wednesday, April  3, 2019",
		},
		{
			key:   "Time",
			value: "09:45 - 10:45 AM",
		},
		{
			key:   "Room",
			value: "Theater",
		},
		{
			key:   "Notes",
			value: "pad.internetfreedomfestival.org/p/1356"},
		{
			key:   "Duration",
			value: "1 hour(s)",
		},
		{
			key:   "Format",
			value: "Collaborative Talk",
		},
		{
			key:   "Theme",
			value: "Hacking the Net",
		},
		{
			key:   "Presenter",
			value: "Quyen",
		},
		{
			key:   "Other Presenters",
			value: "Elodie Vialle",
		},

		{
			key:   "Description",
			value: "Vietnam’s cyber army has become more creative with its attacks. From abusing the Facebook report button to more sophisticated phishing attacks, the cyber army is determined to stop any form of free expression on Facebook. Learn more about these attacks so we can come up with ways to counter them.",
		},
		{
			key:   "Target Audience",
			value: "Advocacy/policy professionals, organisations under online harassment, journalists, security trainers, communications professionals, front line activists",
		},
		{
			key:   "Goal of the session",
			value: "Learn about the new tactics troll armies are implementing to take down activist accounts and formulate methods to counter these tactics.",
		},
		{
			key:   "Desired Outcome",
			value: "To share best practices and forms of support we can provide to Vietnamese activists and netizens who are victims of cyberbullying.\nTo create a list of methods to tackle these various issues from the discussion",
		},
	}

	assertFieldsEqual(t, expectedFields, gotFields)

	gotDescription := formatDescription(gotFields)

	expectedDescription := `Vietnam’s cyber army has become more creative with its attacks. From abusing the Facebook report button to more sophisticated phishing attacks, the cyber army is determined to stop any form of free expression on Facebook. Learn more about these attacks so we can come up with ways to counter them.

* Notes: https://pad.internetfreedomfestival.org/p/1356
* Format: Collaborative Talk
* Theme: Hacking the Net
* Presenter: Quyen
* Other Presenters: Elodie Vialle

# Target Audience
Advocacy/policy professionals, organisations under online harassment, journalists, security trainers, communications professionals, front line activists

# Goal of the session
Learn about the new tactics troll armies are implementing to take down activist accounts and formulate methods to counter these tactics.

# Desired Outcome
To share best practices and forms of support we can provide to Vietnamese activists and netizens who are victims of cyberbullying.
To create a list of methods to tackle these various issues from the discussion`

	assertEqual(t, expectedDescription, gotDescription)
}

func loadHTMLTestData(t *testing.T, filename string) string {
	f, err := os.Open(filepath.Join("testdata", "session.html"))
	assertNoError(t, err)

	data, err := ioutil.ReadAll(f)
	assertNoError(t, err)
	return string(data)
}

func assertFieldsEqual(t *testing.T, expected []sessionField, got []sessionField) {
	for i := range expected {
		assertEqual(t, expected[i], got[i])
	}
}
