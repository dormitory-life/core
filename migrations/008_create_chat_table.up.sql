CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    dormitory_id VARCHAR(2) NOT NULL REFERENCES dormitory (id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    text TEXT NOT NULL CHECK (
        LENGTH(text) BETWEEN 1 AND 2000
    ),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_chat_messages_dormitory_id ON chat_messages (dormitory_id);

CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages (created_at DESC);