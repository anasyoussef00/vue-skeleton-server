CREATE TABLE IF NOT EXISTS users (
    id SERIAL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    username VARCHAR(20) UNIQUE NOT NULL,
    email_address VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    birthdate DATE NOT NULL,
    gender ENUM('M', 'F') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);