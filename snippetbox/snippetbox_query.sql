CREATE TABLE snippets (
	id SERIAL PRIMARY KEY NOT NULL,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL,
	created TIMESTAMP NOT NULL,
	expires TIMESTAMP NOT NULL
);

CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond 2',
    E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    NOW(),
    NOW() + INTERVAL '365 DAY'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    E'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    NOW(),
    NOW() + INTERVAL '365 DAY'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    E'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    NOW(),
    NOW() + INTERVAL '7 DAY'
);

SELECT * FROM snippets;

-- DELETE FROM snippets WHERE id=3

CREATE TABLE sessions (
	token VARCHAR(43) PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMP NOT NULL
)

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

SELECT * FROM sessions

CREATE TABLE users (
	id SERIAL PRIMARY KEY NOT NULL,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	hashed_password CHAR(60) NOT NULL,
	created TIMESTAMP NOT NULL
)
ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);