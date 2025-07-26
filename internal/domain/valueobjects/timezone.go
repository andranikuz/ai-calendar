package valueobjects

import (
	"time"
)

type Timezone struct {
	Name   string
	Offset int // UTC offset in seconds
}

// Common timezones
var (
	TimezoneUTC    = Timezone{Name: "UTC", Offset: 0}
	TimezoneEST    = Timezone{Name: "America/New_York", Offset: -5 * 3600}
	TimezonePST    = Timezone{Name: "America/Los_Angeles", Offset: -8 * 3600}
	TimezoneGMT    = Timezone{Name: "Europe/London", Offset: 0}
	TimezoneCET    = Timezone{Name: "Europe/Berlin", Offset: 1 * 3600}
	TimezoneJST    = Timezone{Name: "Asia/Tokyo", Offset: 9 * 3600}
	TimezoneMSK    = Timezone{Name: "Europe/Moscow", Offset: 3 * 3600}
)

func NewTimezone(name string) (Timezone, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return Timezone{}, err
	}
	
	_, offset := time.Now().In(loc).Zone()
	return Timezone{
		Name:   name,
		Offset: offset,
	}, nil
}

func (tz Timezone) IsValid() bool {
	_, err := time.LoadLocation(tz.Name)
	return err == nil
}

func (tz Timezone) Location() (*time.Location, error) {
	return time.LoadLocation(tz.Name)
}

func (tz Timezone) ConvertTime(t time.Time) (time.Time, error) {
	loc, err := tz.Location()
	if err != nil {
		return time.Time{}, err
	}
	return t.In(loc), nil
}

func (tz Timezone) String() string {
	return tz.Name
}