CREATE TABLE feeds (
    ID UUID PRIMARY KEY,
    Name TEXT UNIQUE NOT NULL,
    URL TEXT NOT NULL,
    Created_at TIMESTAMP DEFAULT NOW(),
    Updated_at TIMESTAMP
);

CREATE TABLE articles(
    ID UUID PRIMARY KEY,
    Title TEXT,
    Link TEXT ,
    Description TEXT,
    Created_at TIMESTAMP DEFAULT NOW(),
    Updated_at TIMESTAMP,
    Published_at TIMESTAMP,
    feed_id UUID references feeds (ID) ON DELETE CASCADE
);

