package mongo

import (
	_ "time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	Lastname string             `bson:"lastname" json:"lastname"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Tell     string             `bson:"tell" json:"tell"`
}
