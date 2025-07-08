# Import CLI Tool

A command-line tool for managing data in the Fluently backend database.

## Features

- **CSV Import**: Import words, topics, and sentences from CSV files
- **Database Clear**: Safely clear all learning data from the database

## Setup

1. Make sure you have the database running
2. Create a `.env` file with the required environment variables:

```env
# Database connection
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=fluently

# Required for clear command (set a strong password)
CLEAR_PASSWORD=your_secure_clear_password
```

## Commands

### CSV Import

Import words, topics, and sentences from a CSV file:

```bash
go run main.go csv --file path/to/your/file.csv
```

**CSV Format:**
The CSV should have the following columns:
- Index (empty or number)
- Total (empty or number)
- word (required)
- topic
- subtopic  
- subsubtopic
- CEFR_level
- translation
- sentences (JSON format: [["English sentence", "Russian translation"], ...])

### Database Clear

**⚠️ WARNING: This permanently deletes ALL learning data!**

Clear all words, sentences, topics, and learned words from the database:

```bash
go run main.go clear
```

This command will:
1. Prompt for the password set in `CLEAR_PASSWORD` environment variable
2. Delete data in the correct order to maintain referential integrity:
   - learned_words (user progress)
   - sentences 
   - words
   - topics (in hierarchy order)
3. Show statistics of deleted records

**What gets deleted:**
- ✅ learned_words table (user learning progress)
- ✅ sentences table 
- ✅ words table
- ✅ topics table

**What is preserved:**
- ❌ users table
- ❌ refresh_tokens table  
- ❌ link_tokens table
- ❌ preferences table

### Word Enrichment

Enrich existing words in the database with phonetic transcription, audio URLs, and part of speech from the dictionary API:

```bash
go run main.go enrich --limit 1000 --delay 500
```

**Options:**
- `--limit` (default: 100000): Maximum number of words to process
- `--delay` (default: 500): Delay between API calls in milliseconds

**Features:**
- ✅ Real-time database updates (each word saved immediately)
- ✅ Rate limiting to prevent API failures
- ✅ Comprehensive progress logging
- ✅ Handles API failures gracefully
- ✅ Marks failed enrichments to avoid reprocessing

### Sentence Enrichment

Enrich existing sentences with distractor options using the ML distractor service:

```bash
go run main.go enrich-sentences --limit 1000 --delay 10
```

**Options:**
- `--limit` (default: 1000000): Maximum number of sentences to process
- `--delay` (default: 10): Delay between API calls in milliseconds

### Reset Failed Enrichments

Reset words that failed enrichment back to "unknown" status to allow retry:

```bash
go run main.go reset-enrichment
```

This command:
- Finds words marked as "unknown_processed" (failed enrichments)
- Resets them back to "unknown" status
- Allows them to be retried in future enrichment runs
- Shows count of reset words

## Safety Features

- **Password Protection**: Clear command requires environment variable `CLEAR_PASSWORD`
- **Transaction Safety**: All operations run in database transactions
- **Confirmation Prompt**: Interactive password confirmation
- **Detailed Logging**: Full audit trail of operations
- **Statistics**: Detailed reports of imported/deleted records

## Examples

### Import CSV data:
```bash
# Set environment variables
export DB_USER=fluently_user
export DB_PASSWORD=your_password
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=fluently_dev

# Import data
go run main.go csv --file words_data.csv
```

### Clear database:
```bash
# Set clear password
export CLEAR_PASSWORD=super_secret_clear_password

# Clear all data
go run main.go clear
# Enter password when prompted: super_secret_clear_password
```

## Build

To build the executable:

```bash
go build -o import-tool main.go
```

Then use:
```bash
./import-tool csv --file data.csv
./import-tool clear
```

## How It Works

### 1. Topic Hierarchy Creation
- Creates hierarchical topic structure: `Animals` → `Animals` → `amphibians_and_reptiles`
- Each level can have the same name (creates separate records)
- Handles missing subtopics/subsubtopics gracefully

### 2. Word Processing
- Groups CSV rows by unique word+topic+CEFR combination
- Merges multiple translations with comma separation
- Deduplicates identical translations
- Sets PartOfSpeech to "unknown" by default
- Updates existing words or creates new ones

### 3. Sentence Processing
- Parses JSON sentence arrays
- Creates individual Sentence records
- Links sentences to their respective words
- Skips duplicate sentences

### 4. Batch Processing
- Processes data in batches of 100 records
- Each batch runs in a database transaction
- Continues processing even if individual batches fail

## Output

The tool provides comprehensive statistics upon completion:

```
==================================================
📊 IMPORT STATISTICS
==================================================
⏱️  Duration: 2m30s
🚀 Speed: 125.5 rows/second
📄 Total rows processed: 18750
📝 Words imported: 12500
🏷️  Topics created: 450
💬 Sentences added: 75000
❌ Errors encountered: 23
==================================================
```

## Error Handling

The tool implements robust error handling:

- **CSV Format Errors**: Logs error and continues with next row
- **Duplicate Words**: Merges translations, logs info message
- **Invalid JSON Sentences**: Logs error and skips sentences for that word
- **Database Errors**: Logs error and continues with next batch
- **Missing Required Fields**: Logs warning and skips row

## Technical Details

### Dependencies
- **Cobra CLI**: Command-line interface framework
- **GORM**: Database ORM
- **Progressbar**: Real-time progress reporting
- **Zap**: Structured logging

### Database Models
- **Topic**: Hierarchical topic structure with parent-child relationships
- **Word**: Vocabulary words with translations, CEFR level, and topic association
- **Sentence**: Example sentences linked to words

### Performance Optimizations
- Batch processing reduces database load
- Transaction-based processing ensures data integrity
- Progress reporting provides user feedback
- Memory-efficient CSV streaming

## Troubleshooting

### Common Issues

1. **"CSV file does not exist"**
   - Verify the file path is correct
   - Use absolute paths if relative paths don't work

2. **Database connection errors**
   - Check your `.env` file configuration
   - Ensure the database is running
   - Verify network connectivity

3. **"Invalid CSV format"**
   - Ensure CSV has exactly 9 columns
   - Check for missing headers
   - Verify CSV encoding (UTF-8 recommended)

4. **High error counts**
   - Check CSV data quality
   - Validate JSON format in sentences column
   - Review log output for specific error patterns

### Logs

The tool uses structured logging. Check the console output for:
- **INFO**: General progress information
- **WARN**: Non-critical issues (skipped rows, etc.)
- **ERROR**: Serious issues that don't stop processing

## Development

### Adding New Features

The tool is designed to be extensible:

1. **New Data Sources**: Add new commands under the root `import` command
2. **New CSV Formats**: Modify the `CSVRecord` struct and parsing logic
3. **Additional Processing**: Add new functions to the processing pipeline
