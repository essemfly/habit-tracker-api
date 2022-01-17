package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/lessbutter/habit-tracker-api/auth"
	"github.com/lessbutter/habit-tracker-api/graph/generated"
	"github.com/lessbutter/habit-tracker-api/graph/model"
	"github.com/lessbutter/habit-tracker-api/repository"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginUserInput) (string, error) {
	userDao, err := repository.GetUserByEmail(input.Email)
	if err != nil {
		return "", errors.New("no user found by email")
	}

	loginCorrect, err := userDao.CheckPassword(input.Password)
	if !loginCorrect {
		return "", err
	}

	token, err := auth.GenerateToken(input.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	_, err := repository.GetUserByEmail(input.Email)
	if err == nil {
		return nil, errors.New("duplicate email")
	}

	newUser, err := repository.InsertUser(input.Email, input.Password, input.Name)
	if err != nil {
		return nil, err
	}
	return newUser.ToDTO(), nil
}

func (r *mutationResolver) CreateHabit(ctx context.Context, input model.CreateHabitInput) (*model.Habit, error) {
	userDao := auth.ForContext(ctx)
	if userDao == nil {
		return nil, errors.New("invalid token")
	}

	newAlertTime := repository.Hour24Time(*input.AlertTime)
	if input.AlertTime != nil {
		newAlertTime = repository.ChangeDatetoHour24Time(*input.AlertTime)
	}

	skipDays := []repository.WeekDayEnum{}
	for _, day := range input.SkipDays {
		skipDays = append(skipDays, repository.WeekDayEnum(*day))
	}

	habit, err := repository.InsertHabit(input.Title, newAlertTime, skipDays, userDao)
	if err != nil {
		return nil, err
	}

	return habit.ToDTO(), nil
}

func (r *mutationResolver) UpdateHabit(ctx context.Context, input model.UpdateHabitInput) (*model.Habit, error) {
	userDao := auth.ForContext(ctx)
	if userDao == nil {
		return nil, errors.New("invalid token")
	}

	habitDao, err := repository.GetHabit(input.ID)
	if err != nil {
		return nil, err
	}

	newAlertTime := repository.Hour24Time("")
	if input.AlertTime != nil {
		newAlertTime = repository.ChangeDatetoHour24Time(*input.AlertTime)
	}

	skipDays := []repository.WeekDayEnum{}
	for _, day := range input.SkipDays {
		skipDays = append(skipDays, repository.WeekDayEnum(*day))
	}

	habitDao.AlertTime = newAlertTime
	habitDao.SkipDays = skipDays
	habitDao.Title = input.Title

	newHabitDao, err := habitDao.Update()
	if err != nil {
		return nil, err
	}

	return newHabitDao.ToDTO(), nil
}

func (r *mutationResolver) UpsertRecord(ctx context.Context, input model.RecordInput) (bool, error) {
	userDao := auth.ForContext(ctx)
	if userDao == nil {
		return false, errors.New("invalid token")
	}

	habitDao, _ := repository.GetHabit(input.HabitID)

	exactDate := repository.ChangeClientDateToExactDate(input.Date)
	newRecord := &repository.HabitRecordDAO{
		Habit:     habitDao,
		Date:      repository.ChangeExactDateToServerDate(exactDate),
		Status:    repository.HabitStatusEnum(input.Status),
		UpdatedAt: time.Now(),
	}

	newRecord, err := repository.UpsertHabitRecord(newRecord)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *queryResolver) Habits(ctx context.Context) ([]*model.Habit, error) {
	userDao := auth.ForContext(ctx)
	if userDao == nil {
		return nil, errors.New("invalid token")
	}
	habitDaos, err := repository.ListHabits(userDao.Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	newHabits := []*model.Habit{}
	for _, habitDao := range habitDaos {
		newHabits = append(newHabits, habitDao.ToDTO())
	}

	return newHabits, nil
}

func (r *queryResolver) Histories(ctx context.Context, input model.HistoryQueryInput) (*model.HabitRecordQueryResult, error) {
	userDao := auth.ForContext(ctx)
	if userDao == nil {
		return nil, errors.New("invalid token")
	}

	habitDao, err := repository.GetHabit(input.HabitID)
	if err != nil {
		return nil, errors.New("invalid habit id")
	}
	if habitDao.User.ID != userDao.ID {
		return nil, errors.New("not your habit")
	}

	startExactDate := repository.ChangeClientDateToExactDate(input.StartDate)
	endExactDate := repository.ChangeClientDateToExactDate(input.EndDate)
	recordDaos, err := repository.ListHabitRecords(input.HabitID, startExactDate, endExactDate)
	if err != nil {
		return nil, err
	}

	totalCounts, _ := repository.GetTotalCounts(input.HabitID)
	records := []*model.HabitRecord{}
	for _, dao := range recordDaos {
		records = append(records, dao.ToDTO())
	}

	ret := model.HabitRecordQueryResult{
		HabitID:    input.HabitID,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		Records:    records,
		TotalCount: totalCounts,
	}

	return &ret, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
