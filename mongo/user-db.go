package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (userDB *UserDB) GetUsers() ([]User, error) {
	collection := userDB.DB.Database("my-database").Collection("user-database")

	filter := bson.D{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []User
	for cursor.Next(context.TODO()) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (userDB *UserDB) GetUserByID(id string) (*User, error) {
	collection := userDB.DB.Database("my-database").Collection("user-database")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.D{{"_id", objectID}}
	var user User

	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (userDB *UserDB) GetUserByUserName(username string) (*User, error) {
	collection := userDB.DB.Database("my-database").Collection("user-database")

	filter := bson.D{{"username", username}}
	var user User

	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (userDB *UserDB) InsertUser(user User) error {
	collection := userDB.DB.Database("my-database").Collection("user-database")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if _, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
		return nil
	}

	return nil
}

func (userDB *UserDB) UpdateUser(id string, updatedUser User) error {
	collection := userDB.DB.Database("my-database").Collection("user-database")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}
	filter := bson.D{{"_id", objectID}}

	update := bson.M{"$set": bson.M{}}

	if updatedUser.Name != "" {
		update["$set"].(bson.M)["name"] = updatedUser.Name
	}
	if updatedUser.Lastname != "" {
		update["$set"].(bson.M)["lastname"] = updatedUser.Lastname
	}
	if updatedUser.Username != "" {
		update["$set"].(bson.M)["username"] = updatedUser.Username
	}
	if updatedUser.Password != "" {
		update["$set"].(bson.M)["password"] = updatedUser.Password
	}
	if updatedUser.Tell != "" {
		update["$set"].(bson.M)["tell"] = updatedUser.Tell
	}

	updateResult, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("no user found with ID: %s", id)
	}
	if _, ok := updateResult.UpsertedID.(primitive.ObjectID); ok {
		return nil
	}

	return nil
}

func (userDB *UserDB) DeleteUser(id string) error {
	collection := userDB.DB.Database("my-database").Collection("user-database")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}
	filter := bson.D{{"_id", objectID}}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("no user found with ID: %s", id)
	}
	return nil
}

func (userDB *UserDB) UserExist(username string) (bool, error) {
	collection := userDB.DB.Database("my-database").Collection("user-database")

	filter := bson.D{
		{"username", username},
	}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No user found
			return false, nil
		}
		// Some other error occurred
		return false, fmt.Errorf("error checking user existence: %v", err)
	}
	return true, nil
}
