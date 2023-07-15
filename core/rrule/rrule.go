package rrule

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
)

var daystrToDayRrule = map[string]rrule.Weekday{
	"monday":    rrule.MO,
	"tuesday":   rrule.TU,
	"wednesday": rrule.WE,
	"thursday":  rrule.TH,
	"friday":    rrule.FR,
	"saturday":  rrule.SA,
	"sunday":    rrule.SU,
}

func GetDateTime(day, hour string, start, end time.Time) ([]time.Time, error) {
	dr, ok := daystrToDayRrule[day]
	if !ok {
		return nil, fmt.Errorf("day %s is invalid", day)
	}

	hours := strings.Split(hour, ":")
	if len(hours) < 2 {
		return nil, fmt.Errorf("hour %s is invalid", hour)
	}

	hh, err := strconv.Atoi(hours[0])
	if err != nil {
		return nil, fmt.Errorf("hour %s is invalid", hours[0])
	}

	mm, err := strconv.Atoi(hours[1])
	if err != nil {
		return nil, fmt.Errorf("minute %s is invalid", hours[1])
	}

	rr, err := rrule.NewRRule(rrule.ROption{
		Freq:      rrule.DAILY,
		Dtstart:   start,
		Until:     end,
		Interval:  1,
		Byweekday: []rrule.Weekday{dr},
		Byhour:    []int{hh},
		Byminute:  []int{mm},
	})
	if err != nil {
		return nil, err
	}

	return rr.All(), nil
}
