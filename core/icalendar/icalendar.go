package icalendar

import (
	"fmt"
	"time"
)

const (
	HeaderCalscale        = "CALSCALE"
	HeaderMethod          = "METHOD"
	HeaderProductId       = "PRODID"
	HeaderVersion         = "VERSION"
	HeaderRefreshInterval = "REFRESH-INTERVAL;VALUE=DURATION"
	HeaderXWRCalName      = "X-WR-CALNAME"
	HeaderXWRTimezone     = "X-WR-TIMEZONE"
)

type Calendar struct {
	headers   []string
	timezones []Timezone
	events    []Event
}

type Timezone struct {
	id   string
	name string
	tz   string
}

type Event struct {
	uid   string
	name  string
	start time.Time
	end   time.Time
	tz    string
}

func New() *Calendar {
	calendar := &Calendar{}
	calendar.AddHeader(HeaderMethod, "PUBLISH")
	calendar.AddHeader(HeaderCalscale, "GREGORIAN")
	calendar.AddHeader(HeaderVersion, "1.0")
	return calendar
}

func (c *Calendar) Serialize() string {
	serial := "BEGIN:VCALENDAR"

	for _, header := range c.headers {
		serial += "\n" + header
	}
	for _, tz := range c.timezones {
		serial += "\n" + tz.Serialize()
	}
	for _, ev := range c.events {
		serial += "\n" + ev.Serialize()
	}

	serial += "\nEND:VCALENDAR"
	return serial
}

func (c *Calendar) AddHeader(key, value string) {
	c.headers = append(c.headers, fmt.Sprintf("%s:%s", key, value))
}

func (c *Calendar) AddTimezone(id, name, tz string) {
	c.timezones = append(c.timezones, Timezone{id, name, tz})
}

func (c *Calendar) AddEvent(id, name string, start, end time.Time, tz string) {
	c.events = append(c.events, Event{id, name, start, end, tz})
}

const timezoneSection = `BEGIN:VTIMEZONE
TZID:%s
X-LIC-LOCATION:%s
BEGIN:STANDARD
TZOFFSETFROM:%s
TZOFFSETTO:%s
TZNAME:%s
DTSTART:19700101T000000
END:STANDARD
END:VTIMEZONE`

func (tz Timezone) Serialize() string {
	return fmt.Sprintf(timezoneSection, tz.id, tz.id, tz.tz, tz.tz, tz.name)
}

const eventSection = `BEGIN:VEVENT
DTSTART;TZID=%s:%s
DTEND;TZID=%s:%s
UID:%s
SEQUENCE:0
STATUS:CONFIRMED
SUMMARY:%s
TRANSP:OPAQUE
END:VEVENT`

func (ev Event) Serialize() string {
	start := ev.start.Format("20060102T150405")
	end := ev.end.Format("20060102T150405")
	return fmt.Sprintf(eventSection, ev.tz, start, ev.tz, end, ev.uid, ev.name)
}
