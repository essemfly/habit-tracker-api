package repository

import (
	"context"
	"log"
	"time"

	"github.com/lessbutter/habit-tracker-api/config"
	"github.com/lessbutter/habit-tracker-api/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HabitStatusEnum string
type WeekDayEnum string
type Hour24Time string

func ChangeDatetoHour24Time(date string) Hour24Time {
	return Hour24Time("13")
}

const (
	IDLE    = HabitStatusEnum("IDLE")
	SKIPPED = HabitStatusEnum("SKIPPED")
	SUCCEED = HabitStatusEnum("SUCCEED")
)

const (
	MONDAY    = WeekDayEnum("MONDAY")
	TUESDAY   = WeekDayEnum("TUESDAY")
	WEDNESDAY = WeekDayEnum("WEDNESDAY")
	THURSDAY  = WeekDayEnum("THURSDAY")
	FRIDAY    = WeekDayEnum("FRIDAY")
	SATURDAY  = WeekDayEnum("SATURDAY")
	SUNDAY    = WeekDayEnum("SUNDAY")
)

type HabitDAO struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	User      *UserDAO
	Title     string
	SkipDays  []WeekDayEnum
	AlertTime Hour24Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (habitDao *HabitDAO) ToDTO() *model.Habit {

	skipdays := []model.WeekDays{}
	for _, day := range habitDao.SkipDays {
		if day == MONDAY {
			skipdays = append(skipdays, model.WeekDaysMonday)
		} else if day == TUESDAY {
			skipdays = append(skipdays, model.WeekDaysTuesday)
		} else if day == WEDNESDAY {
			skipdays = append(skipdays, model.WeekDaysWednesday)
		} else if day == THURSDAY {
			skipdays = append(skipdays, model.WeekDaysThursday)
		} else if day == FRIDAY {
			skipdays = append(skipdays, model.WeekDaysFriday)
		} else if day == SATURDAY {
			skipdays = append(skipdays, model.WeekDaysSaturday)
		} else if day == SUNDAY {
			skipdays = append(skipdays, model.WeekDaysSunday)
		}
	}
	return &model.Habit{
		ID:        habitDao.ID.Hex(),
		Title:     habitDao.Title,
		AlertTime: string(habitDao.AlertTime),
		SkipDays:  skipdays,
	}
}

func GetHabit(ID string) (*HabitDAO, error) {
	c := config.Db.Collection("habits")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var habit *HabitDAO

	habitObjID, _ := primitive.ObjectIDFromHex(ID)
	filter := bson.M{
		"_id": habitObjID,
	}

	if err := c.FindOne(ctx, filter).Decode(&habit); err != nil {
		log.Println(err)
		return nil, err
	}

	return habit, nil
}

func InsertHabit(title string, alertTime Hour24Time, skipDays []WeekDayEnum, user *UserDAO) (*HabitDAO, error) {
	c := config.Db.Collection("habits")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newHabit := HabitDAO{
		ID:        primitive.NewObjectID(),
		User:      user,
		Title:     title,
		SkipDays:  skipDays,
		AlertTime: alertTime,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := c.InsertOne(ctx, newHabit)
	if err != nil {
		log.Println("error on insert habit")
		return nil, err
	}

	return &newHabit, nil
}

func (habitDao *HabitDAO) Update() (*HabitDAO, error) {
	c := config.Db.Collection("habits")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := c.UpdateByID(ctx, habitDao.ID, habitDao)
	if err != nil {
		log.Println("error update habit")
		return nil, err
	}

	return habitDao, nil
}

func ListHabits(email string) ([]*HabitDAO, error) {
	_, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	var habitDaos []*HabitDAO

	c := config.Db.Collection("habits")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user.email": email,
	}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &habitDaos)
	if err != nil {
		return nil, err
	}

	return habitDaos, nil
}

type HabitRecordsDAO struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Habit     *HabitDAO
	Date      WeekDayEnum
	Status    HabitStatusEnum
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (recordDao *HabitRecordsDAO) ToDTO() *model.Habit {
	return nil
}
