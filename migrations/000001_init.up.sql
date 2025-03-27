CREATE TABLE songs (
    id SERIAL PRIMARY KEY,

    genre VARCHAR(30),
    song VARCHAR(30),

    releaseDate DATE,
    text TEXT,
    link VARCHAR(255)
);

CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

ALTER TABLE songs ADD COLUMN group_id INT REFERENCES groups(id);

CREATE INDEX idx_songs_genre ON songs(genre);
CREATE INDEX idx_songs_release_date ON songs(releaseDate);
CREATE INDEX idx_songs_group_id ON songs(group_id);

ALTER TABLE songs ALTER COLUMN releaseDate TYPE DATE USING releaseDate::DATE;

