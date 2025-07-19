-- Migration script to fix user and preferences relationship
-- This script should be run to fix the database schema issues

-- Step 1: Remove the redundant refresh_token column from users table
-- (since we have a separate refresh_tokens table)
ALTER TABLE users DROP COLUMN IF EXISTS refresh_token;

-- Step 2: Ensure user_preferences table has the correct structure
-- First, let's check if the table exists and has the right columns
DO $$
BEGIN
    -- Check if user_preferences table exists
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_preferences') THEN
        -- Create the table if it doesn't exist
        CREATE TABLE user_preferences (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID NOT NULL,
            cefr_level VARCHAR(2) NOT NULL,
            fact_everyday BOOLEAN DEFAULT FALSE,
            notifications BOOLEAN DEFAULT FALSE,
            notifications_at TIMESTAMP DEFAULT NULL,
            words_per_day INTEGER DEFAULT 10,
            goal VARCHAR(255),
            subscribed BOOLEAN DEFAULT FALSE,
            avatar_image_url TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        -- Add foreign key constraint
        ALTER TABLE user_preferences 
        ADD CONSTRAINT fk_user_preferences_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        
        -- Add index on user_id for better performance
        CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
    ELSE
        -- Table exists, ensure it has the correct structure
        -- Add user_id column if it doesn't exist
        IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'user_id') THEN
            ALTER TABLE user_preferences ADD COLUMN user_id UUID;
        END IF;
        
        -- Add foreign key constraint if it doesn't exist
        IF NOT EXISTS (
            SELECT FROM information_schema.table_constraints 
            WHERE constraint_name = 'fk_user_preferences_user_id'
        ) THEN
            ALTER TABLE user_preferences 
            ADD CONSTRAINT fk_user_preferences_user_id 
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        END IF;
        
        -- Add index if it doesn't exist
        IF NOT EXISTS (
            SELECT FROM pg_indexes 
            WHERE indexname = 'idx_user_preferences_user_id'
        ) THEN
            CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
        END IF;
    END IF;
END $$;

-- Step 3: Update any existing preferences to link them to users
-- This assumes that if there are any orphaned preferences, they should be cleaned up
-- or linked to the appropriate users based on your business logic

-- Step 4: Remove any PrefID column from users table if it exists
ALTER TABLE users DROP COLUMN IF EXISTS pref_id;

-- Step 5: Ensure refresh_tokens table exists with correct structure
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'refresh_tokens') THEN
        CREATE TABLE refresh_tokens (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID NOT NULL,
            token TEXT NOT NULL UNIQUE,
            revoked BOOLEAN DEFAULT FALSE,
            expires_at TIMESTAMP NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
        
        -- Add foreign key constraint
        ALTER TABLE refresh_tokens 
        ADD CONSTRAINT fk_refresh_tokens_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        
        -- Add indexes
        CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
        CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
    END IF;
END $$;

-- Step 6: Add any missing columns to user_preferences table
DO $$
BEGIN
    -- Add missing columns if they don't exist
    IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'cefr_level') THEN
        ALTER TABLE user_preferences ADD COLUMN cefr_level VARCHAR(2) NOT NULL DEFAULT 'A1';
    END IF;
    
    IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'fact_everyday') THEN
        ALTER TABLE user_preferences ADD COLUMN fact_everyday BOOLEAN DEFAULT FALSE;
    END IF;
    
    IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'notifications') THEN
        ALTER TABLE user_preferences ADD COLUMN notifications BOOLEAN DEFAULT FALSE;
    END IF;
    
    IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'notifications_at') THEN
        ALTER TABLE user_preferences ADD COLUMN notifications_at TIMESTAMP DEFAULT NULL;
    END IF;
    
    IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'words_per_day') THEN
        ALTER TABLE user_preferences ADD COLUMN words_per_day INTEGER DEFAULT 10;
    END IF;
    
    IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'goal') THEN
        ALTER TABLE user_preferences ADD COLUMN goal VARCHAR(255);
    END IF;
    
    IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'subscribed') THEN
        ALTER TABLE user_preferences ADD COLUMN subscribed BOOLEAN DEFAULT FALSE;
    END IF;
    
    IF NOT EXISTS (SELECT FROM information_schema.columns WHERE table_name = 'user_preferences' AND column_name = 'avatar_image_url') THEN
        ALTER TABLE user_preferences ADD COLUMN avatar_image_url TEXT;
    END IF;
END $$;

-- Step 7: Create default preferences for users who don't have any
INSERT INTO user_preferences (user_id, cefr_level, fact_everyday, notifications, words_per_day, subscribed)
SELECT 
    u.id,
    'A1',
    FALSE,
    FALSE,
    10,
    FALSE
FROM users u
WHERE NOT EXISTS (
    SELECT 1 FROM user_preferences up WHERE up.user_id = u.id
);

-- Step 8: Verify the migration
SELECT 
    'Migration completed successfully' as status,
    COUNT(*) as total_users,
    (SELECT COUNT(*) FROM user_preferences) as total_preferences,
    (SELECT COUNT(*) FROM refresh_tokens) as total_refresh_tokens
FROM users; 