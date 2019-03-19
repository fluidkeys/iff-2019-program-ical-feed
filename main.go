package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/arran4/golang-ical"
)

func main() {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetProductId("-//Paul & Ian (www.fluidkeys.com)//Updated hourly from IFF website/")

	for _, dayPage := range dayUrls {
		fmt.Println(dayPage.url)
		page, err := downloadURL(dayPage.url)
		if err != nil {
			panic(err)
		}

		events, err := parsePage(page, dayPage.date)
		if err != nil {
			log.Printf("%v", err)
		}

		for _, event := range events {
			fmt.Printf("%s\n  %s - %s\n  @ %s\n  %s\n\n", event.title, event.startsAt, event.endsAt,
				event.location, event.url)

			icalEvent := cal.AddEvent(event.id)
			icalEvent.SetDtStampTime(time.Now())
			icalEvent.SetStartAt(event.startsAt)
			icalEvent.SetEndAt(event.endsAt)
			icalEvent.SetLocation(event.location)
			icalEvent.SetDescription(event.url)
			icalEvent.SetSummary(event.title)
			icalEvent.SetURL(event.url)
		}

		time.Sleep(time.Duration(1) * time.Second)
	}
	serialized := cal.Serialize()
	bytes := []byte(strings.Replace(serialized, "\n", "\r\n", -1))

	ioutil.WriteFile("iff2019.ics", bytes, 0600)

}

func parsePage(html string, midnight time.Time) (events []eventListing, err error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		panic(err)
	}
	// there are divs called "event_block" and "special_event_block"
	eventBlocks := htmlquery.Find(doc, `//div[contains(@class, "event_block")]`)
	for _, div := range eventBlocks {
		event := eventListing{}

		titleH5 := htmlquery.FindOne(div, `//h5`)
		event.title = strings.TrimSpace(htmlquery.InnerText(titleH5))

		timeDiv := htmlquery.FindOne(div, `//i[contains(@class, "fa-clock-o")]/parent::div`)

		if start, end, err := parseTimeDiv(htmlquery.InnerText(timeDiv), midnight); err != nil {
			log.Printf("%v", err)
			continue
		} else {
			event.startsAt = *start
			event.endsAt = *end
		}

		locationDiv := htmlquery.FindOne(div, `//i[contains(@class, "fa-map-marker")]/parent::div`)
		event.location = strings.TrimSpace(htmlquery.InnerText(locationDiv))

		trackDiv := htmlquery.FindOne(div, `//i[contains(@class, "fa-pencil-square-o")]/parent::div`)
		if trackDiv != nil {
			event.track = strings.TrimSpace(htmlquery.InnerText(trackDiv))
		}

		aTag := htmlquery.FindOne(div, `/parent::a`)
		if aTag != nil {
			event.url = baseURL + htmlquery.SelectAttr(aTag, "href")
			parts := strings.Split(event.url, "/")
			event.id = fmt.Sprintf("2019-%s@internetfreedomfestival.org", parts[len(parts)-1])
		}

		events = append(events, event)
	}
	return events, nil
}

// parseTimeDiv parses two times from a string like this:
// " 09:00 - 11:00AM (2.0h) "
func parseTimeDiv(text string, midnight time.Time) (*time.Time, *time.Time, error) {
	re := regexp.MustCompile(`(\d\d):(\d\d) - (\d\d):(\d\d)([AP])M`)
	match := re.FindStringSubmatch(text)
	if len(match) != 6 {
		return nil, nil, fmt.Errorf("failed to parse `%s`", text)
	}

	startsHour, _ := strconv.Atoi(match[1])
	startsMin, _ := strconv.Atoi(match[2])
	endsHour, _ := strconv.Atoi(match[3])
	endsMin, _ := strconv.Atoi(match[4])

	afternoon := match[5] == "P"

	if afternoon {
		endsHour += 12

		if startsHour+12 < endsHour {
			startsHour += 12
		}
	}

	startsAt := midnight.Add(
		time.Duration(startsHour)*time.Hour + time.Duration(startsMin)*time.Minute,
	)

	endsAt := midnight.Add(
		time.Duration(endsHour)*time.Hour + time.Duration(endsMin)*time.Minute,
	)

	return &startsAt, &endsAt, nil
}

// <div class="col-md-5 event_block">
//   <div class="row">
//     <h5 class="session_titles">The double edge of "Fake News": Disinformation and attacks on the media</h5>
//   </div>
//   <div class="row">
//     <div class="col-md-3 session_info"><i class="fa fa-clock-o" aria-hidden="true"></i> 02:45 - 03:45PM </div>
//   </div>
//   <div class="row">
//     <div class="col-md-3 session_info"><i class="fa fa-pencil-square-o" aria-hidden="true"></i>
//     Journalism, Media and Comms
//     </div>
//   </div>
//   <div class="row">
//     <div class="col-md-3 session_info"><i class="fa fa-map-marker" aria-hidden="true"></i> Theater </div>
//   </div>
// </div>

// another:
// <div class="col-md-11 special_event_block">
//   <div class="row">
//     <h5 class="session_titles">IFF Opening Ceremony</h5>
//   </div>
//   <div class="row">
//     <div class="col-md-2"><i class="fa fa-clock-o" aria-hidden="true"></i> 11:30 - 12:30PM</div>
// 	<div class="col-md-2"><i class="fa fa-map-marker" aria-hidden="true"></i> La Plaza</div>
//   </div>
// </div>

func downloadURL(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: got HTTP %d", url, response.StatusCode)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type eventListing struct {
	url      string
	id       string
	title    string
	startsAt time.Time
	endsAt   time.Time
	track    string
	location string
}

const baseURL = "https://platform.internetfreedomfestival.org"

var valencia, _ = time.LoadLocation("Europe/Madrid")

var dayUrls = []struct {
	url  string
	date time.Time
}{
	{
		baseURL + "/en/IFF2019/public/schedule/custom?day=6",
		time.Date(2019, 4, 1, 0, 0, 0, 0, valencia),
	},
	{
		baseURL + "/en/IFF2019/public/schedule/custom?day=7",
		time.Date(2019, 4, 2, 0, 0, 0, 0, valencia),
	},
	{
		baseURL + "/en/IFF2019/public/schedule/custom?day=8",
		time.Date(2019, 4, 3, 0, 0, 0, 0, valencia),
	},
	{
		baseURL + "/en/IFF2019/public/schedule/custom?day=9",
		time.Date(2019, 4, 4, 0, 0, 0, 0, valencia),
	},
	{
		baseURL + "/en/IFF2019/public/schedule/custom?day=10",
		time.Date(2019, 4, 5, 0, 0, 0, 0, valencia),
	},
}
