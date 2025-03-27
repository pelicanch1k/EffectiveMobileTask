-- Удаляем индексы
DROP INDEX IF EXISTS idx_songs_genre;
DROP INDEX IF EXISTS idx_songs_release_date;
DROP INDEX IF EXISTS idx_songs_group_id;

-- Удаляем внешний ключ и колонку
ALTER TABLE songs DROP COLUMN IF EXISTS group_id;

-- Удаляем таблицы
DROP TABLE IF EXISTS songs;
DROP TABLE IF EXISTS groups;