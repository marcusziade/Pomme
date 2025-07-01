package sales

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/marcusziade/pomme/internal/models"
)

// Parser handles CSV parsing for sales reports
type Parser struct {
	dateFormats []string
}

// NewParser creates a new CSV parser
func NewParser() *Parser {
	return &Parser{
		dateFormats: []string{
			"01/02/2006", // US format
			"2006-01-02", // ISO format
			"02/01/2006", // EU format
		},
	}
}

// ParseCSV parses sales report CSV data
func (p *Parser) ParseCSV(data []byte) ([]models.SalesRecord, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	reader.Comma = '\t' // Apple uses tab-separated values
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create field index map
	fieldMap := p.createFieldMap(header)
	
	// Validate required fields
	if err := p.validateFields(fieldMap); err != nil {
		return nil, err
	}

	// Read records
	var records []models.SalesRecord
	lineNum := 1
	
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row %d: %w", lineNum, err)
		}
		lineNum++

		record, err := p.parseRecord(row, fieldMap)
		if err != nil {
			// Log warning but continue processing
			// fmt.Printf("Warning: failed to parse row %d: %v\n", lineNum, err)
			continue
		}

		records = append(records, record)
	}

	return records, nil
}

// createFieldMap creates a map of field names to column indices
func (p *Parser) createFieldMap(header []string) map[string]int {
	fieldMap := make(map[string]int)
	
	for i, field := range header {
		// Normalize field names
		normalized := strings.TrimSpace(field)
		fieldMap[normalized] = i
		
		// Also map common variations
		switch normalized {
		case "Provider":
			fieldMap["provider"] = i
		case "Provider Country":
			fieldMap["provider_country"] = i
		case "SKU":
			fieldMap["sku"] = i
		case "Developer":
			fieldMap["developer"] = i
		case "Title":
			fieldMap["title"] = i
			fieldMap["app_name"] = i
		case "Version":
			fieldMap["version"] = i
		case "Product Type Identifier":
			fieldMap["product_type_id"] = i
		case "Units":
			fieldMap["units"] = i
		case "Developer Proceeds":
			fieldMap["developer_proceeds"] = i
			fieldMap["proceeds"] = i
		case "Begin Date":
			fieldMap["begin_date"] = i
			fieldMap["start_date"] = i
		case "End Date":
			fieldMap["end_date"] = i
		case "Customer Currency":
			fieldMap["customer_currency"] = i
		case "Country Code":
			fieldMap["country_code"] = i
			fieldMap["country"] = i
		case "Currency of Proceeds":
			fieldMap["currency_of_proceeds"] = i
			fieldMap["proceeds_currency"] = i
		case "Apple Identifier":
			fieldMap["apple_id"] = i
			fieldMap["apple_identifier"] = i
		case "Customer Price":
			fieldMap["customer_price"] = i
			fieldMap["price"] = i
		case "Promo Code":
			fieldMap["promo_code"] = i
		case "Parent Identifier":
			fieldMap["parent_id"] = i
		case "Subscription":
			fieldMap["subscription"] = i
		case "Period":
			fieldMap["period"] = i
		case "Category":
			fieldMap["category"] = i
		case "CMB":
			fieldMap["cmb"] = i
		case "Device":
			fieldMap["device"] = i
			fieldMap["device_type"] = i
		case "Supported Platforms":
			fieldMap["supported_platforms"] = i
			fieldMap["platforms"] = i
		case "Proceeds Reason":
			fieldMap["proceeds_reason"] = i
		case "Preserved Pricing":
			fieldMap["preserved_pricing"] = i
		case "Client":
			fieldMap["client"] = i
		case "Order Type":
			fieldMap["order_type"] = i
		}
	}
	
	return fieldMap
}

// validateFields ensures all required fields are present
func (p *Parser) validateFields(fieldMap map[string]int) error {
	required := []string{
		"Title",
		"Units",
		"Developer Proceeds",
		"Apple Identifier",
	}
	
	var missing []string
	for _, field := range required {
		if _, ok := fieldMap[field]; !ok {
			missing = append(missing, field)
		}
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}
	
	return nil
}

// parseRecord parses a single CSV row into a SalesRecord
func (p *Parser) parseRecord(row []string, fieldMap map[string]int) (models.SalesRecord, error) {
	record := models.SalesRecord{}
	
	// Helper function to safely get field value
	getValue := func(field string) string {
		if idx, ok := fieldMap[field]; ok && idx < len(row) {
			return strings.TrimSpace(row[idx])
		}
		return ""
	}
	
	// Parse string fields
	record.Provider = getValue("Provider")
	record.ProviderCountry = getValue("Provider Country")
	record.SKU = getValue("SKU")
	record.Developer = getValue("Developer")
	record.Title = getValue("Title")
	record.Version = getValue("Version")
	record.ProductTypeID = getValue("Product Type Identifier")
	record.AppleID = getValue("Apple Identifier")
	record.CountryCode = getValue("Country Code")
	record.CustomerCurrency = getValue("Customer Currency")
	record.CurrencyOfProceeds = getValue("Currency of Proceeds")
	record.BeginDate = getValue("Begin Date")
	record.EndDate = getValue("End Date")
	record.PromoCode = getValue("Promo Code")
	record.ParentID = getValue("Parent Identifier")
	record.Subscription = getValue("Subscription")
	record.Period = getValue("Period")
	record.Category = getValue("Category")
	record.CMB = getValue("CMB")
	record.DeviceType = getValue("Device")
	record.SupportedPlatforms = getValue("Supported Platforms")
	record.ProceedsReason = getValue("Proceeds Reason")
	record.PreservedPricing = getValue("Preserved Pricing")
	record.Client = getValue("Client")
	record.OrderType = getValue("Order Type")
	
	// Parse numeric fields
	var err error
	
	// Units
	unitsStr := getValue("Units")
	if unitsStr != "" {
		record.Units, err = strconv.Atoi(unitsStr)
		if err != nil {
			return record, fmt.Errorf("invalid units value: %s", unitsStr)
		}
	}
	
	// Customer Price
	priceStr := getValue("Customer Price")
	if priceStr != "" && priceStr != " " {
		record.CustomerPrice, err = strconv.ParseFloat(priceStr, 64)
		if err != nil {
			// Try with comma as decimal separator
			priceStr = strings.Replace(priceStr, ",", ".", 1)
			record.CustomerPrice, _ = strconv.ParseFloat(priceStr, 64)
		}
	}
	
	// Developer Proceeds
	proceedsStr := getValue("Developer Proceeds")
	if proceedsStr != "" && proceedsStr != " " {
		record.DeveloperProceeds, err = strconv.ParseFloat(proceedsStr, 64)
		if err != nil {
			// Try with comma as decimal separator
			proceedsStr = strings.Replace(proceedsStr, ",", ".", 1)
			record.DeveloperProceeds, _ = strconv.ParseFloat(proceedsStr, 64)
		}
	}
	
	// Validate essential fields
	if record.AppleID == "" {
		return record, fmt.Errorf("missing Apple ID")
	}
	
	return record, nil
}

// ParseDate attempts to parse a date string in various formats
func (p *Parser) ParseDate(dateStr string) time.Time {
	dateStr = strings.TrimSpace(dateStr)
	
	for _, format := range p.dateFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}
	
	// Return zero time if parsing fails
	return time.Time{}
}

// ParseReportDate extracts the report date from the data
func (p *Parser) ParseReportDate(data []byte) (time.Time, error) {
	// Apple reports sometimes include the date in the first few lines
	lines := bytes.Split(data, []byte("\n"))
	
	for i := 0; i < len(lines) && i < 5; i++ {
		line := string(lines[i])
		
		// Look for date patterns
		if strings.Contains(line, "Report Date:") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				dateStr := strings.TrimSpace(parts[1])
				for _, format := range p.dateFormats {
					if t, err := time.Parse(format, dateStr); err == nil {
						return t, nil
					}
				}
			}
		}
	}
	
	return time.Time{}, fmt.Errorf("report date not found")
}