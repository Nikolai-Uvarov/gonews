package mongo

import (
	"GoNews/pkg/storage"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "gonews" // имя БД
	collectionName = "posts"  // имя коллекции в БД
)

// Хранилище данных.
type Store struct {
	client *mongo.Client
	ctx    context.Context
}

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {
	var s Store
	var err error
	mongoOpts := options.Client().ApplyURI(constr)
	s.client, err = mongo.Connect(s.ctx, mongoOpts)
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = s.client.Ping(s.ctx, nil)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	return posts, nil
}

func (s *Store) AddPost(p storage.Post) error {
	//выбираем базу и коллекцию
	collection := s.client.Database(databaseName).Collection(collectionName)

	//создаем документ в базе
	_, err := collection.InsertOne(s.ctx, p)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdatePost(p storage.Post) error {
	//выбираем базу и коллекцию
	collection := s.client.Database(databaseName).Collection(collectionName)

	//создаем объекты фильтра и документа-обновления
	filter := bson.D{{"id", p.ID}}
	updatedoc := bson.D{{"$set", p}}

	//обновляем в базе
	_, err := collection.UpdateOne(s.ctx, filter, updatedoc)
	if err != nil {
		return err
	}

	return nil
}
func (s *Store) DeletePost(p storage.Post) error {
	//выбираем базу и коллекцию
	collection := s.client.Database(databaseName).Collection(collectionName)

	//создаем объект фильтра
	filter := bson.D{{"id", p.ID}}

	//удаляем из базы
	_, err := collection.DeleteOne(s.ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

var posts = []storage.Post{
	{
		ID:      1,
		Title:   "Effective Go",
		Content: "Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result—Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.",
	},
	{
		ID:      2,
		Title:   "The Go Memory Model",
		Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
	},
}
