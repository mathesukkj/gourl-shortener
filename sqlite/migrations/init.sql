CREATE TABLE IF NOT EXISTS urls (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  initial_url TEXT NOT NULL,
  shortened_url TEXT NOT NULL
); 

CREATE INDEX url_index ON urls (shortened_url);
