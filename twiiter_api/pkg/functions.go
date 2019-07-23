package functions

import "time"

//CalculateTime function to get time from tweet then convert and calculate time to usable
func CalculateTime(layout string, timeValue string) (int, string) {
	t, _ := time.Parse(layout, timeValue)
	l, _ := time.LoadLocation("Local")
	timeWithFormat := t.In(l).Format("15:04")
	diff := time.Now().Sub(t).Minutes()
	return int(diff), timeWithFormat
}
