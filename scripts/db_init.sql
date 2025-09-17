-- CREATE DATABASE social;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, -- Will store bcrypt hash
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Posts table
CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    title VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags TEXT[], -- PostgreSQL array type
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Trigger to automatically update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_posts_updated_at 
    BEFORE UPDATE ON posts 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insert some test data
INSERT INTO users (username, email, password) VALUES 
('john_doe', 'john@example.com', '$2a$10$dummy_hash_for_testing'),
('jane_smith', 'jane@example.com', '$2a$10$dummy_hash_for_testing')
ON CONFLICT (email) DO NOTHING;

INSERT INTO posts (title, content, user_id, tags) VALUES 
('My First Post', 'Hello world! This is my first post.', 1, ARRAY['hello', 'first']),
('Learning Go', 'Go is an amazing language for backend development.', 1, ARRAY['go', 'programming']),
('Database Design', 'Designing good database schemas is crucial.', 2, ARRAY['database', 'design'])
ON CONFLICT DO NOTHING;