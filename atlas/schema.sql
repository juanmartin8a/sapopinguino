CREATE TYPE available_languages AS ENUM (
    'English',
    'Spanish',
    'French',
    'German',
    'Arabic',
    'Mandarin',
    'Portuguese',
    'Russian',
);

CREATE TABLE languages (
    id SERIAL PRIMARY KEY,
    language available_languages NOT NULL,
    words INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
);

CREATE TABLE words (
    id SERIAL PRIMARY KEY,
    language_id INT NOT NULL,
    word TEXT NOT NULL,
    UNIQUE (language_id, word),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_language FOREIGN KEY (language_id) REFERENCES languages (id) ON DELETE CASCADE,
);

CREATE TABLE phonetic_mappings (
    id SERIAL PRIMARY KEY,
    from_language_id INT NOT NULL,
    to_language_id INT NOT NULL,
    word_id INT NOT NULL,
    translation TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (from_language_id, to_language_id, word_id),
    CONSTRAINT fk_from_language FOREIGN KEY (from_language_id) REFERENCES languages (id) ON DELETE CASCADE,
    CONSTRAINT fk_to_language FOREIGN KEY (to_language_id) REFERENCES languages (id) ON DELETE CASCADE,
    CONSTRAINT fk_word FOREIGN KEY (word_id) REFERENCES words (id) ON DELETE CASCADE
);

CREATE TABLE ipa_transcriptions (
    id SERIAL PRIMARY KEY,
    word_id INT NOT NULL,
    ipa_transcription TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE (word_id),
    CONSTRAINT fk_word FOREIGN KEY (word_id) REFERENCES words (id) ON DELETE CASCADE,
);

