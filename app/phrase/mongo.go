package phrase

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Model struct {
	collection *mongo.Collection
}

func NewModel(client *mongo.Client, db string) (*Model, error) {
	err := client.Connect(context.Background())

	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Ping(ctxPing, nil)

	defer cancel()

	if err != nil {
		return nil, err
	}

	collection := client.Database(db).Collection("phrase")

	return &Model{collection: collection}, nil
}

func (model *Model) CreatePhrase(userId, incNumber int, textEnglish, textRussian string) error {
	_, err := model.collection.InsertOne(context.TODO(), Phrase{
		UserId:      userId,
		IncNumber:   incNumber,
		EnglishText: textEnglish,
		RussianText: textRussian,
		IsMemorized: false,
	})
	return err
}

func (model *Model) AllPhrases(userId int) ([]Phrase, error) {
	var phrases []Phrase

	cur, err := model.collection.Find(context.TODO(), bson.M{"isMemorized": false, "userId": userId})

	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem Phrase
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("[ERROR] Fail to decode phrase: %v", err)
		} else {
			phrases = append(phrases, elem)
		}
	}

	return phrases, nil
}

func (model *Model) NextIncNumber(userId int) (int, error) {
	var phrase Phrase

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "incNumber", Value: -1}})

	err := model.collection.FindOne(context.TODO(), bson.M{"isMemorized": false, "userId": userId}, findOptions).Decode(&phrase)

	if err == mongo.ErrNoDocuments {
		return 1, nil
	}

	if err != nil {
		return 0, err
	}

	return phrase.IncNumber + 1, nil
}

func (model *Model) FindPhraseByIncNumber(userId, incNumber int) (*Phrase, error) {
	var phrase Phrase

	err := model.collection.FindOne(context.TODO(), bson.M{"incNumber": incNumber, "userId": userId}).Decode(&phrase)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &phrase, nil
}
