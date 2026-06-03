CREATE TABLE IF NOT EXISTS game_library (
    id             SERIAL PRIMARY KEY,
    rawg_id        INTEGER NOT NULL UNIQUE,
    title          VARCHAR(255) NOT NULL,
    genre          VARCHAR(100),
    platform       VARCHAR(100),
    cover_url      TEXT,
    personal_note  TEXT,
    personal_score INTEGER,
    status         VARCHAR(20),
    added_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_rawg_id ON game_library(rawg_id);
CREATE INDEX IF NOT EXISTS idx_status ON game_library(status);
