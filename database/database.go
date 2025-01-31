package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	client *mongo.Client
}

type ServerSync struct {
	ServerId     string `bson:"server_id"`
	ChannelId    string `bson:"channel_id"`
	WebhookId    string `bson:"webhook_id"`
	WebhookToken string `bson:"webhook_token"`
}

type StickerSubscription struct {
	UserId    string `bson:"user_id"`
	Sticker   string `bson:"sticker"`
	TradeType string `bson:"trade_type"`
}

type User struct {
	UserId       string             `bson:"user_id"`
	PremiumLevel int                `bson:"premium_level"`
	PremiumSince primitive.DateTime `bson:"premium_since"`
	Bookmarks    []string           `bson:"bookmarks"`
	GlobalName   string             `bson:"global_name"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	Image        string             `bson:"image"`
}

type WebsitePost struct {
	UserId         string             `bson:"user_id"`
	UserName       string             `bson:"usǒr_name"`
	UserGlobalName string             `bson:"user_global_name"`
	UserAvatar     string             `bson:"user_avatar"`
	CreatedAt      primitive.DateTime `bson:"created_at"`
	ExpireTime     primitive.DateTime `bson:"expire_time"`
	ServerSync     bool               `bson:"server_sync"`
	Locked         bool               `bson:"locked"`
	Trade          bson.D             `bson:"trade"`
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Collection(coll string) *mongo.Collection {
	return d.client.Database("trade-builder").Collection(coll)
}

func (d *Database) Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	d.client = client
	return client, nil
}
