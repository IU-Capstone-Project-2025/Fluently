# CSV Import Tool

A comprehensive CLI tool for importing vocabulary data from CSV files into the Fluently backend database.

The tool expects a CSV file with the following columns:

```csv
,Total,word,topic,subtopic,subsubtopic,CEFR_level,translation,sentences
0,462.0,adder,Animals,Animals,amphibians_and_reptiles,c2,гадюка,"[['English sentence', 'Russian translation'], ...]"
```

### Column Details:
- **Column 1**: Index (ignored)
- **Total**: Frequency count (parsed but not stored)
- **word**: The English word (required)
- **topic**: Main topic category
- **subtopic**: Sub-category under topic
- **subsubtopic**: Most specific category
- **CEFR_level**: Language proficiency level (e.g., "c1", "c2")
- **translation**: Translation of the word
- **sentences**: JSON array of sentence pairs `[["English", "Translation"], ...]`

### Basic Usage

```bash
# Build the tool
go build -o import-tool cmd/import/main.go

# Import CSV file
./import-tool csv --file /path/to/your/data.csv
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
