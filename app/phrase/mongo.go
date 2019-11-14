package phrase

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (model *Model) CreatePhrase(userId, incNumber int, textEnglish, textRussian string) (error, *Phrase) {
	ph := Phrase{
		UserId:      userId,
		IncNumber:   incNumber,
		EnglishText: textEnglish,
		RussianText: textRussian,
		IsMemorized: false,
		AudioId:     "",
	}
	_, err := model.collection.InsertOne(context.TODO(), ph)

	if err != nil {
		return err, nil
	}
	return nil, &ph
}

func (model *Model) AllPhrases(userId int) ([]Phrase, error) {
	var phrases []Phrase

	cur, err := model.collection.Find(context.TODO(), bson.M{"isMemorized": false, "userId": userId})

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var elem Phrase
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("[ERROR] Fail to decode phrase: %v", err)
		} else {
			phrases = append(phrases, elem)
		}
	}

	if err := cur.Err(); err != nil {
		log.Printf("[ERROR] Fail during work with coursor: %v", err)
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
		return nil, fmt.Errorf("phrase with incNumber: %d not found", incNumber)
	}

	if err != nil {
		return nil, err
	}

	return &phrase, nil
}

func (model *Model) UpdateAudioId(id primitive.ObjectID, audioId string) error {
	filter := bson.M{"_id": id}
	update := bson.D{{Key: "$set", Value: bson.M{"audioId": audioId}}}

	_, err := model.collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func (model *Model) RemovePhrase(userId, incNumber int) (int64, error) {
	delRes, err := model.collection.DeleteOne(context.TODO(), bson.M{"incNumber": incNumber, "userId": userId})

	if err != nil {
		return 0, err
	}

	return delRes.DeletedCount, nil
}

func (model *Model) SmartFindByRange(userId int, from, to *int) ([]Phrase, error) {
	var phrases []Phrase
	var filter bson.M

	if from == nil && to == nil {
		filter = bson.M{"isMemorized": false, "userId": userId}
	} else if from != nil && to != nil {
		filter = bson.M{"isMemorized": false, "userId": userId, "incNumber": bson.M{"$gte": *from, "$lte": *to}}
	} else if from != nil {
		filter = bson.M{"isMemorized": false, "userId": userId, "incNumber": *from}
	} else {
		return nil, errors.New("params not correct")
	}

	cur, err := model.collection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var elem Phrase
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("[ERROR] Fail to decode phrase: %v", err)
		} else {
			phrases = append(phrases, elem)
		}
	}

	if err := cur.Err(); err != nil {
		log.Printf("[ERROR] Fail during work with coursor: %v", err)
	}

	return phrases, nil
}
