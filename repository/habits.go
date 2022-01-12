package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HabitStatusEnum string

const (
	IDLE    = HabitStatusEnum("IDLE")
	SKIPPED = HabitStatusEnum("SKIPPED")
	SUCCEED = HabitStatusEnum("SUCCEED")
)

type HabitDAO struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	User      *UserDAO
	Title     string
	AlertTime time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type HabitRecordsDAO struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Habit     *HabitDAO
	Date      time.Time
	Status    HabitStatusEnum
	CreatedAt time.Time
	UpdatedAt time.Time
}
