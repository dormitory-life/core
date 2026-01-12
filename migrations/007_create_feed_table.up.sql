CREATE TABLE IF NOT EXISTS feed (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    dormitory_id VARCHAR(2) NOT NULL REFERENCES dormitory (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_feed_dormitory_id ON feed (dormitory_id);