package helper

import (
	"fmt"
	"time"
)

func ValidateDateRange(start_date, end_date string) (err error) {
	start := time.Now()
	end := time.Now()
	if start_date != "" {
		start, err = time.Parse("2006-01-02", start_date)
		if err != nil {
			err = fmt.Errorf("Waktu awal tidak valid.")
			return
		}
	}
	if end_date != "" {
		end, err = time.Parse("2006-01-02", end_date)
		if err != nil {
			err = fmt.Errorf("Waktu akhir tidak valid.")
			return
		}
	}
	if start_date != "" && end_date != "" && start.After(end) {
		err = fmt.Errorf("Waktu awal dan akhir tidak valid.")
		return
	}
	return
}

func ValidateHourRange(startHour, endHour string) (err error) {
	start := time.Now()
	end := time.Now()
	if startHour != "" {
		start, err = time.Parse("15:04", startHour)
		if err != nil {
			err = fmt.Errorf("Waktu awal tidak valid.")
			return
		}
	}
	if endHour != "" {
		end, err = time.Parse("15:04", endHour)
		if err != nil {
			err = fmt.Errorf("Waktu akhir tidak valid.")
			return
		}
	}
	if startHour != "" && endHour != "" && start.After(end) {
		err = fmt.Errorf("Waktu awal dan akhir tidak valid.")
		return
	}
	return
}

var validDay = map[string]bool{
	"monday":    true,
	"tuesday":   true,
	"wednesday": true,
	"thursday":  true,
	"friday":    true,
	"saturday":  true,
	"sunday":    true,
}

func ValdiateDayString(day string) error {
	valid, ok := validDay[day]
	if !valid || !ok {
		return fmt.Errorf("Day %s is not valid.", day)
	}

	return nil
}
