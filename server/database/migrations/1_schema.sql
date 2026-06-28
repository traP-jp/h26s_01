-- +goose Up

CREATE TABLE IF NOT EXISTS users(
    id VARCHAR(32) NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS rooms(
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL
);

CREATE TABLE IF NOT EXISTS room_members(
    room_id UUID NOT NULL,
    user_id VARCHAR(32) NOT NULL,
	is_ready BOOLEAN NOT NULL DEFAULT FALSE,
    guesser_order TINYINT UNSIGNED NOT NULL DEFAULT 0,
    is_connected BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY(room_id, user_id),
    FOREIGN KEY (room_id) REFERENCES rooms(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS games(
    room_id UUID NOT NULL PRIMARY KEY,
    current_round_id UUID DEFAULT NULL,
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE TABLE IF NOT EXISTS game_kanjies(
    id UUID NOT NULL PRIMARY KEY,
    game_id UUID NOT NULL,
    `character` VARCHAR(1) NOT NULL,
    kanji_order TINYINT UNSIGNED NOT NULL DEFAULT 0,
    FOREIGN KEY (game_id) REFERENCES games(room_id)
);

CREATE TABLE IF NOT EXISTS rounds(
    id UUID NOT NULL PRIMARY KEY,
    game_id UUID NOT NULL,
    round_index TINYINT UNSIGNED NOT NULL,
    current_turn_id UUID DEFAULT NULL,
    guesser_id VARCHAR(32) NOT NULL,
    kanji_id UUID NOT NULL,
    started_at TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL,
    FOREIGN KEY (game_id) REFERENCES games(room_id),
    FOREIGN KEY (guesser_id) REFERENCES users(id),
    FOREIGN KEY (kanji_id) REFERENCES game_kanjies(id)
);

CREATE TABLE IF NOT EXISTS turns(
    id UUID NOT NULL PRIMARY KEY,
    round_id UUID NOT NULL,
    turn_index TINYINT UNSIGNED NOT NULL,
    drawer_id VARCHAR(32) NOT NULL,
    FOREIGN KEY (round_id) REFERENCES rounds(id),
    FOREIGN KEY (drawer_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS strokes(
    id UUID NOT NULL PRIMARY KEY,
    turn_id UUID NOT NULL,
    x1 DOUBLE UNSIGNED NOT NULL,
    y1 DOUBLE UNSIGNED NOT NULL,
    x2 DOUBLE UNSIGNED NOT NULL,
    y2 DOUBLE UNSIGNED NOT NULL,
    FOREIGN KEY (turn_id) REFERENCES turns(id)
);

CREATE TABLE IF NOT EXISTS round_results(
    round_id UUID NOT NULL PRIMARY KEY,
    guesser_answer VARCHAR(1) NOT NULL,
    time_ms INT UNSIGNED NOT NULL,
    FOREIGN KEY (round_id) REFERENCES rounds(id)
);
