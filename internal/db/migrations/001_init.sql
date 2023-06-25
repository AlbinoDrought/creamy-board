CREATE TABLE boards (
  board_id SERIAL PRIMARY KEY
, slug VARCHAR(50) UNIQUE NOT NULL
, title VARCHAR(100) NOT NULL
, tagline VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE FUNCTION make_board_post_seq() RETURNS TRIGGER LANGUAGE PLPGSQL AS $$
  BEGIN
    EXECUTE FORMAT('CREATE SEQUENCE IF NOT EXISTS board_post_seq_%s', NEW.board_id);
    RETURN NEW;
  END
$$;

CREATE TRIGGER make_board_post_seq
AFTER INSERT ON boards
FOR EACH ROW EXECUTE PROCEDURE make_board_post_seq();

CREATE FUNCTION board_post_seq_nextval(board_id INTEGER) RETURNS INTEGER LANGUAGE SQL AS $$
  SELECT nextval('board_post_seq_' || board_id);
$$;

CREATE TABLE threads (
  board_id INTEGER NOT NULL
, thread_id BIGINT NOT NULL
, created_at TIMESTAMP NOT NULL DEFAULT NOW()
, bumped_at TIMESTAMP NOT NULL DEFAULT NOW()
, PRIMARY KEY (board_id, thread_id)
, FOREIGN KEY (board_id) REFERENCES boards (board_id) ON DELETE CASCADE
);

CREATE INDEX board_threads ON threads (board_id);

CREATE TABLE posts (
  board_id INTEGER NOT NULL
, thread_id BIGINT NOT NULL
, post_id BIGINT NOT NULL
, created_at TIMESTAMP NOT NULL DEFAULT NOW()
, subject VARCHAR(100) NOT NULL DEFAULT ''
, author VARCHAR(50) NOT NULL DEFAULT ''
, body TEXT NOT NULL DEFAULT ''
, PRIMARY KEY (board_id, thread_id, post_id)
, FOREIGN KEY (board_id) REFERENCES boards (board_id) ON DELETE CASCADE
, FOREIGN KEY (board_id, thread_id) REFERENCES threads (board_id, thread_id) ON DELETE CASCADE
);

CREATE INDEX thread_posts ON posts (board_id, thread_id);

CREATE TABLE files (
  board_id INTEGER NOT NULL
, thread_id BIGINT NOT NULL
, post_id BIGINT NOT NULL
, idx SMALLINT NOT NULL
, path VARCHAR(255) NOT NULL
, extension VARCHAR(10) NOT NULL
, mimetype VARCHAR(255) NOT NULL
, bytes INTEGER NOT NULL
, original_name VARCHAR(255) NOT NULL
, PRIMARY KEY (board_id, thread_id, post_id, idx)
, FOREIGN KEY (board_id) REFERENCES boards (board_id) ON DELETE CASCADE
, FOREIGN KEY (board_id, thread_id) REFERENCES threads (board_id, thread_id) ON DELETE CASCADE
, FOREIGN KEY (board_id, thread_id, post_id) REFERENCES posts (board_id, thread_id, post_id) ON DELETE CASCADE
);

CREATE INDEX thread_files ON files (board_id, thread_id);
CREATE INDEX post_files ON files (board_id, thread_id, post_id);

CREATE TYPE partial_file AS (
  idx SMALLINT
, path VARCHAR(255)
, extension VARCHAR(10)
, mimetype VARCHAR(255)
, bytes INTEGER
, original_name VARCHAR(255)
)
;

INSERT INTO boards (slug, title, tagline) VALUES
  ('cb', 'Creamy Board', 'Welcome Home')
;

WITH
board AS (
  SELECT board_id FROM boards WHERE slug = 'cb'
),
thread AS (
  INSERT INTO threads (board_id, thread_id) VALUES
    ((SELECT board_id FROM board), board_post_seq_nextval((SELECT board_id FROM board)))
  RETURNING thread_id
),
post1 AS (
  INSERT INTO posts (board_id, thread_id, post_id, subject, author, body) VALUES
    ((SELECT board_id FROM board), (SELECT thread_id FROM thread), (SELECT thread_id FROM thread), 'Welcome to Creamy Board', 'Migrator', E'Test thread body please ignore\n\n(pic related)')
  RETURNING post_id
),
post2 AS (
  INSERT INTO posts (board_id, thread_id, post_id, subject, author, body) VALUES
    ((SELECT board_id FROM board), (SELECT thread_id FROM thread), board_post_seq_nextval((SELECT board_id FROM board)), '', 'Migrator', 'chicken chicken dog')
  RETURNING post_id
),
post3 AS (
  INSERT INTO posts (board_id, thread_id, post_id, subject, author, body) VALUES
    ((SELECT board_id FROM board), (SELECT thread_id FROM thread), board_post_seq_nextval((SELECT board_id FROM board)), 'meow', 'Meowgrator', 'cat')
  RETURNING post_id
)
INSERT INTO files (board_id, thread_id, post_id, idx, path, extension, mimetype, bytes, original_name) VALUES
  ((SELECT board_id FROM board), (SELECT thread_id FROM thread), (SELECT post_id FROM post1), 0, 'test_llama.jpg', 'jpg', 'image/jpeg', 101036, 'llama.jpg')
, ((SELECT board_id FROM board), (SELECT thread_id FROM thread), (SELECT post_id FROM post2), 0, 'test_chimkin.jpg', 'jpg', 'image/jpeg', 179262, 'chimkin.jpg')
, ((SELECT board_id FROM board), (SELECT thread_id FROM thread), (SELECT post_id FROM post2), 1, 'test_chimkin2.jpg', 'jpg', 'image/jpeg', 179262, 'chimkin2.jpg')
, ((SELECT board_id FROM board), (SELECT thread_id FROM thread), (SELECT post_id FROM post2), 2, 'test_dog.jpg', 'jpg', 'image/jpeg', 54955, 'notchimkin.jpg')
, ((SELECT board_id FROM board), (SELECT thread_id FROM thread), (SELECT post_id FROM post3), 0, 'test_cat.jpg', 'jpg', 'image/jpeg', 84949, 'cat.jpg')
;
