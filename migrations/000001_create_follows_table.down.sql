-- Drop indexes
DROP INDEX IF EXISTS idx_follows_followee;
DROP INDEX IF EXISTS idx_follows_follower;

-- Drop follows table
DROP TABLE IF EXISTS follows;
