DROP TABLE IF EXISTS authors, posts;

CREATE TABLE IF NOT EXISTS authors (
    id BIGSERIAL PRIMARY KEY, 
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY, 
	author_id BIGINT NOT NULL REFERENCES authors(id),
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created_at BIGINT NOT NULL DEFAULT 0,
	published_at BIGINT NOT NULL DEFAULT 0
);

-- Очистка всех таблиц перед наполнением тестовыми данными.
TRUNCATE TABLE authors, posts;

-- Наполнение таблиц тестовыми данными.
INSERT INTO authors(name) VALUES
    ('Маск'), ('Цукерберг');
	
INSERT INTO posts(author_id,title, content, created_at, published_at) VALUES
    (1,'Особенности сражений в Колизее', 'Главное быть готовым к дерзости противника...', 1687180530, 1687180550),
	(2,'Моя борьба и победа', 'Все начиналось как шоу...', 1687180550, 1687180590),
	(1,'Держи карман шире', 'Некоторые самоуверенные личности...', 1687180570, 1687180580);


