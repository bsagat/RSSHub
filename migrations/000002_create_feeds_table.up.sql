CREATE TABLE feeds (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    description TEXT ,
    url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE articles(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT,
    link TEXT ,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    published_at TIMESTAMP,
    feed_id UUID references feeds (id) ON DELETE CASCADE
);


