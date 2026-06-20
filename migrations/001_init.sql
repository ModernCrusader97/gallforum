CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS channels (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    owner_id INTEGER REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    channel_id INTEGER NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id),
    guest_name TEXT,
    guest_password TEXT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_urls TEXT DEFAULT '[]',
    likes INTEGER DEFAULT 0,
    dislikes INTEGER DEFAULT 0,
    comment_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post_votes (
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id),
    guest_ip TEXT,
    vote INTEGER NOT NULL CHECK(vote IN (1, -1)),
    PRIMARY KEY (post_id, user_id, guest_ip)
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    parent_id INTEGER REFERENCES comments(id),
    user_id INTEGER REFERENCES users(id),
    guest_name TEXT,
    guest_password TEXT,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT OR IGNORE INTO channels (slug, name, description) VALUES
    ('free', '자유게시판', '자유롭게 이야기하는 공간'),
    ('humor', '유머게시판', '웃긴 게시물을 공유하세요'),
    ('news', '뉴스', '시사 뉴스 및 이슈');
