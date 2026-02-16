# vanilladbc-csv

CSV conversion plugin for [vanilladbc-cli](https://github.com/suprsokr/vanilladbc-cli).

Provides bidirectional conversion between World of Warcraft Vanilla DBC files and CSV format.

## Features

- **DBC to CSV** - Export DBC files to CSV format
- **CSV to DBC** - Import CSV files back to DBC format
- **Header Row** - Automatic column name headers
- **Type Preservation** - Maintains data types during round-trip conversion
- **Quoted Strings** - Properly handles string values with commas

## Installation

```bash
go get github.com/suprsokr/vanilladbc-csv
```

## Usage

### With vanilladbc-cli

```bash
# Convert DBC to CSV
vanilladbc convert Spell.dbc Spell.dbd 1.12.1.5875 --plugin csv --output spell.csv

# Import CSV back to DBC
vanilladbc import spell.csv Spell.dbd 1.12.1.5875 --plugin csv --output Spell.dbc
```

### As a Library

```go
package main

import (
    "os"
    csvplugin "github.com/suprsokr/vanilladbc-csv"
    "github.com/suprsokr/vanilladbc/pkg/dbd"
    "github.com/suprsokr/vanilladbc/pkg/dbc"
)

func main() {
    // Writing to CSV
    writer, _ := os.Create("output.csv")
    plugin := csvplugin.New(writer)
    
    // ... use plugin.WriteHeader(), WriteRecord(), WriteFooter()
    
    // Reading from CSV
    reader, _ := os.Open("input.csv")
    readerPlugin := csvplugin.NewReader(reader)
    
    // ... use readerPlugin.ReadHeader(), ReadRecord()
}
```

## CSV Format

The CSV output includes a header row with column names from the DBD definition:

```csv
ID,School,Category,CastUI,DispelType,Mechanic
1,2,0,0,0,0
133,4,76,76,0,0
```

## Dependencies

- [vanilladbc](https://github.com/suprsokr/vanilladbc) - Core DBC/DBD library

## Related Projects

- [vanilladbc-cli](https://github.com/suprsokr/vanilladbc-cli) - Command-line tool
- [vanilladbc-json](https://github.com/suprsokr/vanilladbc-json) - JSON plugin
- [vanilladbc-mysql](https://github.com/suprsokr/vanilladbc-mysql) - MySQL plugin
- [VanillaDBDefs](https://github.com/suprsokr/VanillaDBDefs) - Database definitions

## License

MIT License - See LICENSE file for details.
