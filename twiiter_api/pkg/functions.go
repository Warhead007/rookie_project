package functions

import (
	"log"
	"time"
)

//CalculateTime function to get time from tweet then convert and calculate time
//between feed_time and time now
func CalculateTime(layout string, timeValue string) (int, string) {
	t, _ := time.Parse(layout, timeValue)
	l, _ := time.LoadLocation("Local")
	timeWithFormat := t.In(l).Format("15:04")
	diff := time.Now().Sub(t).Minutes()
	return int(diff), timeWithFormat
}

//FailOnError function to handle with error by show message
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
