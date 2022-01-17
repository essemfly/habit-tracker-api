package repository

import (
	"context"
	"log"
	"time"

	"github.com/lessbutter/habit-tracker-api/config"
	"github.com/lessbutter/habit-tracker-api/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	SkipDays  []WeekDayEnum
	AlertTime Hour24Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (habitDao *HabitDAO) ToDTO() *model.Habit {
	skipdays := []model.WeekDays{}
	for _, day := range habitDao.SkipDays {
		log.Println("my day", day, "AFTER", day.ToDTO())
		skipdays = append(skipdays, day.ToDTO())
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

	_, err := c.UpdateByID(ctx, habitDao.ID, bson.M{"$set": bson.M{
		"title":     habitDao.Title,
		"alerttime": habitDao.AlertTime,
		"skipdays":  habitDao.SkipDays,
	}})
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

type HabitRecordDAO struct {
	ID        primitive.ObjectID `bson:"_id, omitempty"`
	Habit     *HabitDAO
	Date      time.Time
	Status    HabitStatusEnum
	UpdatedAt time.Time
}

func InsertHabitRecord(record *HabitRecordDAO) (*HabitRecordDAO, error) {
	c := config.Db.Collection("records")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	record.ID = primitive.NewObjectID()
	_, err := c.InsertOne(ctx, record)
	if err != nil {
		log.Println("error on insert record", err)
		return nil, err
	}

	return record, nil
}

func UpsertHabitRecord(record *HabitRecordDAO) (*HabitRecordDAO, error) {
	c := config.Db.Collection("records")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"habit._id": record.Habit.ID, "date": record.Date}

	var habbitRecord *HabitRecordDAO
	if err := c.FindOne(ctx, filter).Decode(&habbitRecord); err != nil {
		record.ID = primitive.NewObjectID()
		_, err = c.InsertOne(ctx, record)
		if err != nil {
			return nil, err
		}
		return record, nil
	}

	record.ID = primitive.NilObjectID
	if _, err := c.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"date":   record.Date,
		"habit":  record.Habit,
		"status": record.Status,
	}}); err != nil {
		log.Println("err occured on upsert record", err)
		return nil, err
	}

	return record, nil
}

func ListHabitRecords(habitID string, start, end ExactDate) ([]*HabitRecordDAO, error) {
	c := config.Db.Collection("records")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	habitOID, _ := primitive.ObjectIDFromHex(habitID)
	filter := bson.M{
		"habit._id": habitOID,
		"date": bson.M{
			"$gte": ChangeExactDateToServerDate(start),
			"$lte": ChangeExactDateToServerDate(end),
		},
	}

	options := options.Find()
	options.SetSort(bson.D{{Key: "date", Value: 1}})

	cursor, err := c.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}

	records := []*HabitRecordDAO{}
	err = cursor.All(ctx, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func GetTotalCounts(habitID string) (int, error) {
	c := config.Db.Collection("records")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	habitOID, _ := primitive.ObjectIDFromHex(habitID)
	filter := bson.M{
		"habit._id": habitOID,
		"statis":    "SUCCEED",
	}

	totalCounts, err := c.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(totalCounts), nil
}

func (recordDao *HabitRecordDAO) ToDTO() *model.HabitRecord {
	exactDate := ChangeServerDateToExactDate(recordDao.Date)
	return &model.HabitRecord{
		HabitID: recordDao.Habit.ID.Hex(),
		Status:  model.HabitStatus(recordDao.Status),
		Date:    ChangeExactDateToClientDate(exactDate),
	}
}
