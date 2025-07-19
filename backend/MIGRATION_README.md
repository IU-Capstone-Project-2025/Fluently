# Database Migration Guide

## Overview

This migration fixes the relationship between `users` and `user_preferences` tables and removes redundant refresh token handling from the users table.

## Changes Made

### 1. User-Preferences Relationship
- **Before**: Users had a `PrefID` field pointing to preferences
- **After**: Preferences have a `UserID` field pointing to users (one-to-one relationship)
- This is the correct GORM pattern for one-to-one relationships

### 2. Refresh Token Handling
- **Before**: Refresh tokens were stored as a string field in the users table
- **After**: Refresh tokens are stored in a separate `refresh_tokens` table with proper expiration and revocation handling

### 3. Code Changes
- Updated `UserRepository` methods to be backward compatible but deprecated
- Created proper `RefreshTokenRepository` for handling refresh tokens
- Fixed GORM model relationships

## Running the Migration

### Option 1: Using the provided script
```bash
cd backend
./run_migration.sh
```

### Option 2: Manual execution
```bash
# Set your database connection details
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=fluently
export DB_USER=postgres
export DB_PASSWORD=your_password

# Run the migration
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migration_fix.sql
```

### Option 3: Using Docker (if using docker-compose)
```bash
docker-compose exec postgres psql -U postgres -d fluently -f /path/to/migration_fix.sql
```

## What the Migration Does

1. **Removes redundant columns**: Drops `refresh_token` and `pref_id` columns from users table
2. **Ensures proper table structure**: Creates/updates `user_preferences` and `refresh_tokens` tables
3. **Adds foreign key constraints**: Properly links tables with CASCADE delete
4. **Creates indexes**: For better query performance
5. **Creates default preferences**: For users who don't have any preferences yet
6. **Verifies migration**: Shows counts of users, preferences, and refresh tokens

## After Migration

1. **Restart your application**: The new GORM models will work correctly
2. **Update your code**: Gradually migrate from deprecated methods to new `RefreshTokenRepository`
3. **Test thoroughly**: Ensure all authentication and preference features work correctly

## Rollback (if needed)

If you need to rollback, you can restore from a backup taken before the migration. The migration is designed to be safe and idempotent, but it's always good practice to have a backup.

## Code Migration Guide

### Old way (deprecated):
```go
// Update refresh token
userRepo.UpdateRefreshToken(ctx, userID, token)

// Get user by refresh token
user, err := userRepo.GetByRefreshToken(ctx, token)
```

### New way:
```go
// Create refresh token
refreshToken := &models.RefreshToken{
    UserID:    userID,
    Token:     token,
    ExpiresAt: time.Now().Add(24 * time.Hour),
}
refreshTokenRepo.Create(ctx, refreshToken)

// Get user by refresh token
refreshToken, err := refreshTokenRepo.GetByToken(ctx, token)
if err != nil {
    return nil, err
}
user, err := userRepo.GetByID(ctx, refreshToken.UserID)
```

## Troubleshooting

### Common Issues

1. **Permission denied**: Make sure your database user has ALTER TABLE permissions
2. **Connection refused**: Check your database connection settings
3. **Column already exists**: The migration uses `IF EXISTS` and `IF NOT EXISTS` to handle this safely

### Verification

After migration, you can verify it worked by checking:
```sql
-- Check table structure
\d users
\d user_preferences
\d refresh_tokens

-- Check data
SELECT COUNT(*) FROM users;
SELECT COUNT(*) FROM user_preferences;
SELECT COUNT(*) FROM refresh_tokens;
```

## Support

If you encounter any issues during migration, please:
1. Check the database logs for detailed error messages
2. Verify your database connection settings
3. Ensure you have proper permissions on the database 