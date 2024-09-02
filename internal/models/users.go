package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var (
	collectionName = "users"
)

type User struct {
	ID             string    `bson:"_id"`
	Name           string    `bson:"name"`
	Email          string    `bson:"email"`
	HashedPassword []byte    `bson:"hashedPassword"`
	Created        time.Time `bson:"created"`
}

type UserModel struct {
	DB *mongo.Database
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := User{
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
		Created:        time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println(user)
	res, err := m.DB.Collection(collectionName).InsertOne(ctx, &user)
	if err != nil {
		// after adding index check for duplicate email fail - chap 10.3
		return err
	}

	fmt.Println(res)
	/*
		id, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			panic(fmt.Sprintf("insert operation returned unexpected value %v", id))
		}
	*/
	return nil
}

func (m *UserModel) Authenticate(email, password string) (string, error) {
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.Collection(collectionName).FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", ErrInvalidCredentials
		} else {
			return "", err
		}
	}

	return user.ID, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
