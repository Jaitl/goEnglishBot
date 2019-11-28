package category

import (
	"context"
	"errors"
	"fmt"
	"github.com/jaitl/goEnglishBot/app/phrase"
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
	collection := client.Database(db).Collection("category")

	return &Model{collection: collection}, nil
}

func (model *Model) CreateCategory(userId, incNumber int, name string) (*Category, error) {
	cat := Category{
		UserId:     userId,
		Active:     true,
		IncNumber:  incNumber,
		Name:       name,
		CreateDate: time.Now().In(time.UTC).Unix(),
		Phrases:    make([]phrase.Phrase, 0),
	}
	_, err := model.collection.InsertOne(context.TODO(), cat)

	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (model *Model) FindCategoryByIncNumber(userId, incNumber int) (*Category, error) {
	var cat Category

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "phrases.incNumber", Value: 1}})

	err := model.collection.FindOne(context.TODO(), bson.M{"incNumber": incNumber, "userId": userId}, findOptions).Decode(&cat)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("category with incNumber: %d not found", incNumber)
	}

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (model *Model) FindActiveCategory(userId int) (*Category, error) {
	var cat Category

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "phrases.incNumber", Value: 1}})

	err := model.collection.FindOne(context.TODO(), bson.M{"userId": userId, "active": true}, findOptions).Decode(&cat)

	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("active category not found")
	}

	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (model *Model) AllCategories(userId int) ([]Category, error) {
	var cats []Category

	cur, err := model.collection.Find(context.TODO(), bson.M{"userId": userId})

	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var elem Category
		err := cur.Decode(&elem)
		if err != nil {
			log.Printf("[ERROR] Fail to decode phrase: %v", err)
		} else {
			cats = append(cats, elem)
		}
	}

	if err := cur.Err(); err != nil {
		log.Printf("[ERROR] Fail during work with coursor: %v", err)
	}

	return cats, nil
}

func (model *Model) NextIncNumberCategory(userId int) (int, error) {
	var cat Category

	findOptions := options.FindOne()
	findOptions.SetSort(bson.D{{Key: "incNumber", Value: -1}})

	err := model.collection.FindOne(context.TODO(), bson.M{"userId": userId}, findOptions).Decode(&cat)

	if err == mongo.ErrNoDocuments {
		return 1, nil
	}

	if err != nil {
		return 0, err
	}

	return cat.IncNumber + 1, nil
}

func (model *Model) DeactivateAllCategories(userId int) error {
	filter := bson.M{"userId": userId, "active": true}
	update := bson.D{{Key: "$set", Value: bson.M{"active": false}}}

	_, err := model.collection.UpdateMany(context.TODO(), filter, update)

	return err
}

func (model *Model) ActivateCategory(userId, incNumber int) error {
	filter := bson.M{"userId": userId, "incNumber": incNumber}
	update := bson.D{{Key: "$set", Value: bson.M{"active": true}}}
	_, err := model.collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func (model *Model) RemoveCategory(userId, incNumber int) (int64, error) {
	cat, err := model.FindCategoryByIncNumber(userId, incNumber)

	if err != nil {
		return 0, err
	}

	if cat.Active {
		return 0, errors.New("unable to remove active category")
	}

	delRes, err := model.collection.DeleteOne(context.TODO(), bson.M{"userId": userId, "incNumber": incNumber})

	if err != nil {
		return 0, err
	}

	return delRes.DeletedCount, nil
}

func (model *Model) NextIncNumberPhrase(userId int) (int, error) {
	cat, err := model.FindActiveCategory(userId)

	if err != nil {
		return 0, err
	}

	if len(cat.Phrases) == 0 {
		return 1, nil
	}

	lastPhrase := cat.Phrases[len(cat.Phrases)-1]

	return lastPhrase.IncNumber + 1, nil
}

func (model *Model) CreatePhrase(userId, incNumber int, textEnglish, textRussian string) (*phrase.Phrase, *Category, error) {
	cat, err := model.FindActiveCategory(userId)

	if err != nil {
		return nil, nil, err
	}

	ph := phrase.Phrase{
		UserId:      userId,
		IncNumber:   incNumber,
		EnglishText: textEnglish,
		RussianText: textRussian,
		AudioId:     "",
	}

	_, err = model.collection.UpdateOne(context.TODO(), bson.M{"_id": cat.Id}, bson.M{"$push": bson.M{"phrases": ph}})

	if err != nil {
		return nil, nil, err
	}

	cat.Phrases = append(cat.Phrases, ph)

	return &ph, cat, nil
}

func (model *Model) RemovePhrase(userId, incNumber int) (int64, error) {
	cat, err := model.FindActiveCategory(userId)

	if err != nil {
		return 0, err
	}

	filter := bson.M{"_id": cat.Id}
	update := bson.M{"$pull": bson.M{"phrases": bson.M{"incNumber": incNumber}}}

	updateRes, err := model.collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return 0, err
	}

	return updateRes.ModifiedCount, nil
}

func (model *Model) UpdatePhraseAudioId(userId, incNumber int, audioId string) error {
	cat, err := model.FindActiveCategory(userId)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": cat.Id, "phrases.incNumber": incNumber}
	update := bson.D{{Key: "$set", Value: bson.M{"phrases.$.audioId": audioId}}}

	_, err = model.collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func (model *Model) FindPhraseByIncNumber(userId, incNumber int) (*phrase.Phrase, error) {
	cat, err := model.FindActiveCategory(userId)

	if err != nil {
		return nil, err
	}

	for _, ph := range cat.Phrases {
		if ph.IncNumber == incNumber {
			return &ph, nil
		}
	}

	return nil, fmt.Errorf("phrase with incNumber: %d not found", incNumber)
}

func (model *Model) SmartFindByRange(userId int, from, to *int) ([]phrase.Phrase, error) {
	cat, err := model.FindActiveCategory(userId)

	if err != nil {
		return nil, err
	}

	var phrases []phrase.Phrase

	if from == nil && to == nil {
		phrases = cat.Phrases
	} else if from != nil && to != nil {
		for _, ph := range cat.Phrases {
			if ph.IncNumber >= *from && ph.IncNumber <= *to {
				phrases = append(phrases, ph)
			}
		}
	} else if from != nil {
		for _, ph := range cat.Phrases {
			if ph.IncNumber == *from {
				phrases = append(phrases, ph)
				break
			}
		}
	} else {
		return nil, errors.New("params not correct")
	}

	return phrases, nil
}
