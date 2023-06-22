CREATE TABLE boards (
  board_id SERIAL PRIMARY KEY
, slug VARCHAR(50) UNIQUE NOT NULL
, title VARCHAR(100) NOT NULL
, tagline VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE board_post_counters (
  board_id INTEGER NOT NULL PRIMARY KEY
, next_post_id BIGINT NOT NULL DEFAULT 1
, FOREIGN KEY (board_id) REFERENCES boards (board_id) ON DELETE CASCADE
);

CREATE TABLE threads (
  board_id INTEGER NOT NULL
, thread_id BIGINT NOT NULL
, created_at TIMESTAMP NOT NULL DEFAULT NOW()
, bumped_at TIMESTAMP NOT NULL DEFAULT NOW()
, subject VARCHAR(100) NOT NULL DEFAULT ''
, PRIMARY KEY (board_id, thread_id)
, FOREIGN KEY (board_id) REFERENCES boards (board_id) ON DELETE CASCADE
);

CREATE TABLE posts (
  board_id INTEGER NOT NULL
, thread_id BIGINT NOT NULL
, post_id BIGINT NOT NULL
, created_at TIMESTAMP NOT NULL DEFAULT NOW()
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

CREATE INDEX post_files ON files (board_id, thread_id, post_id);

INSERT INTO boards (slug, title, tagline) VALUES
  ('cb', 'Creamy Board', 'Welcome Home')
;
INSERT INTO board_post_counters (board_id, next_post_id) VALUES
  (1, 3)
;
INSERT INTO threads (board_id, thread_id, subject) VALUES
  (1, 1, 'Welcome to Creamy Board')
;
INSERT INTO posts (board_id, thread_id, post_id, author, body) VALUES
  (1, 1, 1, 'Migrator', 'Test thread body please ignore')
, (1, 1, 2, 'Migrator', 'Test post body please ignore')
;
