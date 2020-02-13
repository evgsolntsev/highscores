package main

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	DatabaseName   = "dbname"
	CollectionName = "scores"
)

type DAO interface {
	Insert(context.Context, Score) error
	GetTop(context.Context, int) ([]Score, error)
	RemoveAll(context.Context) error
}

var _ DAO = (*defaultDAO)(nil)

type defaultDAO struct {
	collection *mgo.Collection
}

func NewDAO(ctx context.Context, session *mgo.Session) *defaultDAO {
	return &defaultDAO{
		collection: session.DB(DatabaseName).C(CollectionName),
	}
}

func (d *defaultDAO) Insert(ctx context.Context, score Score) error {
	return d.collection.Insert(&score)
}

func (d *defaultDAO) RemoveAll(ctx context.Context) error {
	_, err := d.collection.RemoveAll(bson.M{})
	return err
}

func (d *defaultDAO) GetTop(ctx context.Context, quantity int) ([]Score, error) {
	var result []Score
	err := d.collection.Find(nil).Sort("-score").Limit(quantity).All(&result)

	if err != nil {
		return nil, err
	}
	return result, nil
}
