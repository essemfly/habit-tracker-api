package repository

import (
	"context"
	"log"
	"time"

	"github.com/lessbutter/habit-tracker-api/config"
	"github.com/lessbutter/habit-tracker-api/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserDAO struct {
	ID       primitive.ObjectID `bson:"_id, omitempty"`
	Name     string
	Email    string
	Password string
}

func GetUserByEmail(email string) (*UserDAO, error) {
	c := config.Db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user *UserDAO
	filter := bson.M{
		"email": email,
	}

	if err := c.FindOne(ctx, filter).Decode(&user); err != nil {
		log.Println(err)
		return nil, err
	}
	return user, nil
}

func InsertUser(email, password, name string) (*UserDAO, error) {
	c := config.Db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error on hashing password")
		return nil, err
	}

	newUser := UserDAO{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}

	_, err = c.InsertOne(ctx, newUser)
	if err != nil {
		log.Println("error on insert one")
		return nil, err
	}

	return &newUser, nil
}

func (userDao *UserDAO) CheckPassword(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(userDao.Password), []byte(password))

	if err != nil {
		log.Println("hash error on password: ", err)
		return false, err
	}
	return true, nil
}

func (userDao *UserDAO) ToDTO() *model.User {
	return &model.User{
		ID:    userDao.ID.Hex(),
		Email: userDao.Email,
		Name:  userDao.Name,
	}
}
