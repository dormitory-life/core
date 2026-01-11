CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    owner_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    dormitory_id VARCHAR(2) NOT NULL REFERENCES dormitory (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reviews_dormitory_id ON reviews (dormitory_id);

CREATE INDEX idx_reviews_created_at ON reviews (created_at)