package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	DB  *pgxpool.Pool
	ctx context.Context
}

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {

	var s Store
	s.ctx = context.Background()

	var err error
	s.DB, err = pgxpool.Connect(s.ctx, constr)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	var posts []storage.Post

	rows, err := s.DB.Query(s.ctx,
		`SELECT p.id, p.author_id, p.title, p.content, p.created_at, p.published_at, a.name AS author_name 
		FROM posts AS p 
		JOIN authors AS a ON p.author_id = a.id 
		ORDER BY p.id;`)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p storage.Post

		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.PublishedAt,
			&p.AuthorName,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	//проверить rows.Err()
	return posts, rows.Err()
}

func (s *Store) addAuthor(name string) (int64, error) {

	//проверить, существует ли автор по имени
	//если нет, то вставить нового автора
	//узнать id существовавшего или вновь созданного автора

	_, err := s.DB.Exec(s.ctx,
		`INSERT INTO authors (name) 
		SELECT ($1)
		WHERE NOT EXISTS (SELECT * FROM authors WHERE name =($1) LIMIT 1);`, name)

	if err != nil {
		return 0, err
	}

	rows, err := s.DB.Query(s.ctx,
		`SELECT id FROM authors WHERE name = ($1);`, name)

	if err != nil {
		return 0, err
	}

	var id []int64
	for rows.Next() {

		var ci int64
		err = rows.Scan(&ci)

		if err != nil {
			return 0, err
		}
		id = append(id, ci)
	}

	return id[0], rows.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	//получаем из БД id созданного или уже существовавшего автора
	id, err := s.addAuthor(p.AuthorName)

	if err != nil {
		return err
	}
	//записываем в базу новый пост
	_, err = s.DB.Exec(s.ctx,
		`INSERT INTO posts(author_id,title, content, created_at, published_at) 
		VALUES (($1), ($2), ($3), ($4),($5));`, id, p.Title, p.Content, p.CreatedAt, time.Now().Unix())

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdatePost(p storage.Post) error {
	//получаем из БД id созданного или уже существовавшего автора
	id, err := s.addAuthor(p.AuthorName)

	if err != nil {
		return err
	}
	//обновляем в базе пост
	_, err = s.DB.Exec(s.ctx,
		`UPDATE posts 
		SET author_id = ($1), title = ($2), content = ($3), published_at = ($4) 
		WHERE id = ($5);`,
		id, p.Title, p.Content, time.Now().Unix(), p.ID)

	if err != nil {
		return err
	}

	return nil
}
func (s *Store) DeletePost(storage.Post) error {
	return nil
}
