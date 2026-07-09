-- tweets table 
CREATE TABLE IF NOT EXISTS tweets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL,
    content TEXT NOT NULL,
    reply_to_tweet_id UUID NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

-- likes table
CREATE TABLE IF NOT EXISTS likes (
    tweet_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tweet_id, user_id)
);

-- retweets table
CREATE TABLE IF NOT EXISTS retweets (
    tweet_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tweet_id, user_id)
);

-- indexes for performance optimization 
CREATE INDEX IF NOT EXISTS idx_tweets_author_created_at ON tweets (author_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_tweets_reply ON tweets (reply_to_tweet_id);
CREATE INDEX IF NOT EXISTS idx_likes_tweet ON likes (tweet_id);
CREATE INDEX IF NOT EXISTS idx_retweets_tweet ON retweets (tweet_id);
