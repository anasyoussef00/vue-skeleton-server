CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    user_id BIGINT UNSIGNED,
    title VARCHAR(255) NOT NULL,
    alternative_title VARCHAR(255) NULL,
    description TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX user_ind(user_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
)