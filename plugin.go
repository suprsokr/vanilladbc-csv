package csvplugin

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/suprsokr/vanilladbc/pkg/dbc"
	"github.com/suprsokr/vanilladbc/pkg/dbd"
)

// Plugin implements both plugin.Writer and plugin.Reader interfaces for CSV
type Plugin struct {
	reader io.Reader
	writer io.Writer
	csvReader *csv.Reader
	csvWriter *csv.Writer
	
	// For writing
	columnOrder []string
	versionDef  *dbd.VersionDefinition
	columns     map[string]dbd.ColumnDefinition
	
	// For reading
	headers []string
	index   int
}

// New creates a new CSV plugin that writes to the given writer
func New(writer io.Writer) *Plugin {
	return &Plugin{
		writer:    writer,
		csvWriter: csv.NewWriter(writer),
	}
}

// NewReader creates a new CSV plugin that reads from the given reader
func NewReader(reader io.Reader) *Plugin {
	return &Plugin{
		reader:    reader,
		csvReader: csv.NewReader(reader),
	}
}

// NewConverter creates a new CSV plugin that can both read and write
func NewConverter(reader io.Reader, writer io.Writer) *Plugin {
	return &Plugin{
		reader:    reader,
		writer:    writer,
		csvReader: csv.NewReader(reader),
		csvWriter: csv.NewWriter(writer),
	}
}

// WriteHeader is called once before any records are written
func (p *Plugin) WriteHeader(versionDef *dbd.VersionDefinition, columns map[string]dbd.ColumnDefinition) error {
	p.versionDef = versionDef
	p.columns = columns
	
	// Build column order from version definition
	p.columnOrder = make([]string, 0, len(versionDef.Definitions))
	for _, def := range versionDef.Definitions {
		p.columnOrder = append(p.columnOrder, def.Column)
	}
	
	// Write header row
	return p.csvWriter.Write(p.columnOrder)
}

// WriteRecord is called for each record in the DBC file
func (p *Plugin) WriteRecord(record dbc.Record) error {
	row := make([]string, len(p.columnOrder))
	
	for i, colName := range p.columnOrder {
		value := record[colName]
		row[i] = fmt.Sprintf("%v", value)
	}
	
	return p.csvWriter.Write(row)
}

// WriteFooter is called once after all records are written
func (p *Plugin) WriteFooter() error {
	p.csvWriter.Flush()
	return p.csvWriter.Error()
}

// ReadHeader is called once before reading records
func (p *Plugin) ReadHeader() (*dbd.VersionDefinition, map[string]dbd.ColumnDefinition, error) {
	// Read the header row
	headers, err := p.csvReader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV header: %w", err)
	}
	
	p.headers = headers
	p.index = 0
	
	// CSV doesn't store schema info, so we return what was set
	return p.versionDef, p.columns, nil
}

// ReadRecord is called repeatedly to read records
func (p *Plugin) ReadRecord() (dbc.Record, error) {
	row, err := p.csvReader.Read()
	if err == io.EOF {
		return nil, nil // No more records
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV record: %w", err)
	}
	
	// Convert row to record
	record := make(dbc.Record)
	for i, value := range row {
		if i >= len(p.headers) {
			break
		}
		
		colName := p.headers[i]
		
		// Try to parse the value based on column definition
		if colDef, ok := p.columns[colName]; ok {
			parsedValue, err := parseValue(value, colDef)
			if err != nil {
				return nil, fmt.Errorf("failed to parse value for column %s: %w", colName, err)
			}
			record[colName] = parsedValue
		} else {
			// Fallback: store as string
			record[colName] = value
		}
	}
	
	p.index++
	return record, nil
}

// Close is called to cleanup resources
func (p *Plugin) Close() error {
	// Nothing to cleanup for CSV
	return nil
}

// SetSchema allows the caller to set the schema information needed for reading
func (p *Plugin) SetSchema(versionDef *dbd.VersionDefinition, columns map[string]dbd.ColumnDefinition) {
	p.versionDef = versionDef
	p.columns = columns
}

// parseValue parses a string value based on column definition
func parseValue(value string, colDef dbd.ColumnDefinition) (interface{}, error) {
	value = strings.TrimSpace(value)
	
	switch colDef.Type {
	case dbd.TypeInt:
		// Default to int32 for integers
		v, err := strconv.ParseInt(value, 10, 32)
		return int32(v), err
	case dbd.TypeUInt:
		// Default to uint32 for unsigned integers
		v, err := strconv.ParseUint(value, 10, 32)
		return uint32(v), err
	case dbd.TypeFloat:
		// Default to float32
		v, err := strconv.ParseFloat(value, 32)
		return float32(v), err
	case dbd.TypeString, dbd.TypeLocString:
		return value, nil
	}
	
	// Fallback: return as string
	return value, nil
}
