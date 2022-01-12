package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/lessbutter/habit-tracker-api/graph/generated"
	"github.com/lessbutter/habit-tracker-api/graph/model"
)

func (r *mutationResolver) Login(ctx context.Context, input model.LoginUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateHabit(ctx context.Context, input model.CreateHabitInput) (*model.Habit, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateHabit(ctx context.Context, input model.UpdateHabitInput) (*model.Habit, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateRecord(ctx context.Context, input model.RecordInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Habits(ctx context.Context) ([]*model.Habit, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Histories(ctx context.Context, input model.HistoryQueryInput) (*model.HabitRecordQueryResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
