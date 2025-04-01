package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

// Format represents an output format
type Format string

const (
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
	FormatTable Format = "table"
)

// Formatter handles formatting output data
type Formatter struct {
	format Format
	output io.Writer
}

// NewFormatter creates a new formatter with the specified format
func NewFormatter(format Format, output io.Writer) *Formatter {
	return &Formatter{
		format: format,
		output: output,
	}
}

// Format formats the provided data according to the formatter's format
func (f *Formatter) Format(data interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.formatJSON(data)
	case FormatCSV:
		return f.formatCSV(data)
	case FormatTable:
		return f.formatTable(data)
	default:
		return fmt.Errorf("unsupported format: %s", f.format)
	}
}

// formatJSON formats data as JSON
func (f *Formatter) formatJSON(data interface{}) error {
	encoder := json.NewEncoder(f.output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// formatCSV formats data as CSV
func (f *Formatter) formatCSV(data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Handle slices of structs
	if val.Kind() == reflect.Slice {
		if val.Len() == 0 {
			return nil
		}

		// Get the first element to determine field names
		firstElem := val.Index(0)
		if firstElem.Kind() == reflect.Ptr {
			firstElem = firstElem.Elem()
		}

		if firstElem.Kind() != reflect.Struct {
			return fmt.Errorf("CSV format only supports slices of structs")
		}

		writer := csv.NewWriter(f.output)
		defer writer.Flush()

		// Extract field names for header
		firstElemType := firstElem.Type()
		headers := make([]string, firstElemType.NumField())
		for i := 0; i < firstElemType.NumField(); i++ {
			headers[i] = firstElemType.Field(i).Name
		}

		// Write header
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write CSV header: %w", err)
		}

		// Write each row
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}

			row := make([]string, elem.NumField())
			for j := 0; j < elem.NumField(); j++ {
				row[j] = fmt.Sprintf("%v", elem.Field(j).Interface())
			}

			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write CSV row: %w", err)
			}
		}

		return nil
	}

	return fmt.Errorf("CSV format only supports slices of structs")
}

// formatTable formats data as a table
func (f *Formatter) formatTable(data interface{}) error {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Handle slices of structs
	if val.Kind() == reflect.Slice {
		if val.Len() == 0 {
			fmt.Fprintln(f.output, "No data")
			return nil
		}

		// Get the first element to determine field names
		firstElem := val.Index(0)
		if firstElem.Kind() == reflect.Ptr {
			firstElem = firstElem.Elem()
		}

		if firstElem.Kind() != reflect.Struct {
			return fmt.Errorf("table format only supports slices of structs")
		}

		// Create tabwriter for aligned output
		w := tabwriter.NewWriter(f.output, 0, 0, 2, ' ', 0)

		// Extract field names for header
		firstElemType := firstElem.Type()
		headers := make([]string, firstElemType.NumField())
		for i := 0; i < firstElemType.NumField(); i++ {
			headers[i] = firstElemType.Field(i).Name
		}

		// Write header
		fmt.Fprintln(w, strings.Join(headers, "\t"))

		// Write separator line
		seps := make([]string, len(headers))
		for i := range seps {
			seps[i] = strings.Repeat("-", len(headers[i]))
		}
		fmt.Fprintln(w, strings.Join(seps, "\t"))

		// Write each row
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}

			row := make([]string, elem.NumField())
			for j := 0; j < elem.NumField(); j++ {
				row[j] = fmt.Sprintf("%v", elem.Field(j).Interface())
			}

			fmt.Fprintln(w, strings.Join(row, "\t"))
		}

		return w.Flush()
	}

	return fmt.Errorf("table format only supports slices of structs")
}
