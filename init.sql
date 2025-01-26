CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'user', 'moderator')),
  password_hash VARCHAR(255) NOT NULL,
  refresh_token text DEFAULT null,
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS snippets (
  id SERIAL PRIMARY KEY,
  userid INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  language VARCHAR(50) NOT NULL,
  title VARCHAR(255) NOT NULL UNIQUE,
  code BYTEA NOT NULL,
  description TEXT,
  private boolean NOT NULL,
  tags VARCHAR(50)[],
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);