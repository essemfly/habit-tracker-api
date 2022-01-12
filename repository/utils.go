package repository

import (
	"log"
	"strings"
	"time"

	"github.com/lessbutter/habit-tracker-api/graph/model"
)

type WeekDayEnum string
type Hour24Time string
type ExactDate struct {
	Date string
	Day  string
}

const (
	MONDAY    = WeekDayEnum("Monday")
	TUESDAY   = WeekDayEnum("Tuesday")
	WEDNESDAY = WeekDayEnum("Wednesday")
	THURSDAY  = WeekDayEnum("Thursday")
	FRIDAY    = WeekDayEnum("Friday")
	SATURDAY  = WeekDayEnum("Saturday")
	SUNDAY    = WeekDayEnum("Sunday")
)

// time.RFC3339
func ChangeDatetoHour24Time(hours string) Hour24Time {
	return Hour24Time(hours)
}

// RFC3339에서 parsing하기
func ChangeClientDateToExactDate(date string) ExactDate {
	RFC3339local := "2006-01-02T15:04:05Z"
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatal(err)
	}

	timeObj, _ := time.ParseInLocation(RFC3339local, date, loc)

	return ChangeServerDateToExactDate(timeObj)
}

// RFC3339로 보내기
func ChangeExactDateToClientDate(date ExactDate) string {
	timeObj := ChangeExactDateToServerDate(date)
	return timeObj.Format(time.RFC3339)
}

func ChangeServerDateToExactDate(dbtime time.Time) ExactDate {
	dbtimeString := dbtime.Format("2006-01-02 Monday")
	timeSlice := strings.Split(dbtimeString, " ")
	return ExactDate{
		Date: timeSlice[0],
		Day:  timeSlice[1],
	}
}

func ChangeExactDateToServerDate(datetime ExactDate) time.Time {
	t, _ := time.Parse("2006-01-02", datetime.Date)
	return t
}

func ChangeWeekDayToWeekDayEnum(day model.WeekDays) WeekDayEnum {
	if day == model.WeekDaysMonday {
		return MONDAY
	} else if day == model.WeekDaysTuesday {
		return TUESDAY
	} else if day == model.WeekDaysWednesday {
		return WEDNESDAY
	} else if day == model.WeekDaysThursday {
		return THURSDAY
	} else if day == model.WeekDaysFriday {
		return FRIDAY
	} else if day == model.WeekDaysSaturday {
		return SATURDAY
	} else if day == model.WeekDaysSunday {
		return SUNDAY
	}
	return ""
}

func (weekdayEnum WeekDayEnum) ToDTO() model.WeekDays {
	if weekdayEnum == MONDAY {
		return model.WeekDaysMonday
	} else if weekdayEnum == TUESDAY {
		return model.WeekDaysTuesday
	} else if weekdayEnum == WEDNESDAY {
		return model.WeekDaysWednesday
	} else if weekdayEnum == THURSDAY {
		return model.WeekDaysThursday
	} else if weekdayEnum == FRIDAY {
		return model.WeekDaysFriday
	} else if weekdayEnum == SATURDAY {
		return model.WeekDaysSaturday
	} else if weekdayEnum == SUNDAY {
		return model.WeekDaysSunday
	}
	return ""
}
