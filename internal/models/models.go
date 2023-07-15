package models

import "time"

type EmptyBodyReq struct {
	OrgID    int64
	ActionBy int64
}

type Pagination struct {
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
	Limit int   `json:"limit"`
}

type SendVerifCodeReq struct {
	OrgID    int64
	ActionBy int64
	Timezone string

	UserID      int64
	Password    string
	Destination string
	RequestType int
}

type EmptyDataResp struct {
	Errors `json:"-"`
}

const (
	TimeFormatMon02Jan       = "Mon 02 Jan"
	TimeFormatJan02Comma2006 = "Jan 02, 2006"
	TimeZoneCodeDefault      = "Asia/Jakarta"
	TimestampFormat          = "2006-01-02 15:04:05"
)

// GetTimeFormat and GetAppTime used only in data timestamp only (with hours precision)
// currently only used on admission_events, last_action_time on user, admission_histories timestamp
// be careful when compare time.Now with GetAppTime value, because time.Now use server timezone
func GetTimeFormat(timestamp time.Time, zone, format string) string {
	loc, _ := time.LoadLocation(zone)

	return timestamp.In(loc).Format(format)
}
func GetAppTime(timestampFE string, zone string) string {
	format := TimestampFormat
	// get time on user timezone
	userLoc, _ := time.LoadLocation(zone)
	timestamp, _ := time.ParseInLocation(format, timestampFE, userLoc)
	// parse time on user timezone to utc timezone
	utcLoc, _ := time.LoadLocation("UTC")
	result := timestamp.In(utcLoc).Format(format)
	return result
}
func GetUserLoc(zone string) *time.Location {
	// get time on user timezone
	userLoc, _ := time.LoadLocation(zone)
	return userLoc
}
