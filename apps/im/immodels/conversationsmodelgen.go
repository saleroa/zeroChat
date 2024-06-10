// Code generated by goctl. DO NOT EDIT!
package immodels

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type conversationsModel interface {
	Insert(ctx context.Context, data *Conversations) error
	FindOne(ctx context.Context, id string) (*Conversations, error)
	Update(ctx context.Context, data *Conversations) error
	Delete(ctx context.Context, id string) error
	FindByUserId(ctx context.Context, uid string) (*Conversations, error)
}

type defaultConversationsModel struct {
	conn *mon.Model
}

func newDefaultConversationsModel(conn *mon.Model) *defaultConversationsModel {
	return &defaultConversationsModel{conn: conn}
}

func (m *defaultConversationsModel) Insert(ctx context.Context, data *Conversations) error {
	if !data.ID.IsZero() {
		data.ID = primitive.NewObjectID()
		data.CreateAt = time.Now()
		data.UpdateAt = time.Now()
	}

	_, err := m.conn.InsertOne(ctx, data)
	return err
}

func (m *defaultConversationsModel) FindOne(ctx context.Context, id string) (*Conversations, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidObjectId
	}

	var data Conversations

	err = m.conn.FindOne(ctx, &data, bson.M{"_id": oid})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultConversationsModel) Update(ctx context.Context, data *Conversations) error {
	data.UpdateAt = time.Now()

	_, err := m.conn.UpdateOne(ctx, bson.M{"_id": data.ID}, bson.M{
		"$set": data,
	}, options.Update().SetUpsert(true))
	return err
}

func (m *defaultConversationsModel) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidObjectId
	}

	_, err = m.conn.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

// 根据用户 id 查询他的对应的的 conversation 
func (m *defaultConversationsModel) FindByUserId(ctx context.Context, uid string) (*Conversations, error) {
	var data Conversations

	err := m.conn.FindOne(ctx, &data, bson.M{"userId": uid})
	switch err {
	case nil:
		return &data, nil
	case mon.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
