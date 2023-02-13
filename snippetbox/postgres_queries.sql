CREATE TABLE snippets (
    id SERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);
CREATE INDEX idx_snippets_created ON snippets(created);
INSERT INTO snippets (title, content, created, expires)
VALUES (
        'An old silent pond',
        E'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
        NOW(),
        NOW() + INTERVAL '365 DAY'
    );
INSERT INTO snippets (title, content, created, expires)
VALUES (
        'Over the wintry forest',
        E'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
        NOW(),
        NOW() + INTERVAL '365 DAY'
    );
INSERT INTO snippets (title, content, created, expires)
VALUES (
        'First autumn morning',
        E'First autumn morning\nthe mirror I stare into\nshows my father' 's face.\n\n– Murakami Kijo',
        NOW(),
        NOW() + INTERVAL '7 DAY'
    );
SELECT *
FROM snippets;