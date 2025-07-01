# Pomme CLI Manual

Pomme is a powerful App Store Connect CLI tool that provides comprehensive access to sales reports, analytics, and review management.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Sales Commands](#sales-commands)
- [Analytics Commands](#analytics-commands)
- [Reviews Commands](#reviews-commands)
- [Tips & Tricks](#tips--tricks)

## Installation

<details>
<summary>📦 Installation Options</summary>

### Install with Go
```bash
go install github.com/marcusziade/pomme/cmd/pomme@latest
```

### Download Binary
```bash
curl -L https://github.com/marcusziade/pomme/releases/latest/download/pomme_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv pomme /usr/local/bin/
```

### Build from Source
```bash
git clone https://github.com/marcusziade/pomme.git
cd pomme
go build -o pomme ./cmd/pomme/main.go
```

</details>

## Configuration

<details>
<summary>⚙️ Setting up Pomme</summary>

### Config File

Create `~/.config/pomme/pomme.yaml`:

```yaml
api:
  base_url: https://api.appstoreconnect.apple.com/v1
  timeout: 30
defaults:
  output_format: table
  vendor_number: YOUR_VENDOR_NUMBER
auth:
  key_id: YOUR_KEY_ID
  issuer_id: YOUR_ISSUER_ID
  private_key_path: /path/to/your/AuthKey.p8
```

### Environment Variables

Override config with environment variables:

```bash
export POMME_AUTH_KEY_ID=YOUR_KEY_ID
export POMME_AUTH_ISSUER_ID=YOUR_ISSUER_ID
export POMME_AUTH_PRIVATE_KEY_PATH=/path/to/key.p8
export POMME_DEFAULTS_VENDOR_NUMBER=93036463
```

### Getting API Keys

1. Log in to [App Store Connect](https://appstoreconnect.apple.com)
2. Go to Users and Access → Keys
3. Create a new key with appropriate permissions
4. Download the .p8 file and note the Key ID and Issuer ID

</details>

## Sales Commands

<details>
<summary>💰 Sales Reports</summary>

### Basic Commands

```bash
# Latest available month
pomme sales

# Specific month
pomme sales monthly 2025-03

# With details
pomme sales monthly --details
```

### Filtering Options

```bash
# Group by app
pomme sales monthly --by-app

# Group by country
pomme sales monthly --by-country

# Combined
pomme sales monthly --by-app --by-country
```

### Output Formats

```bash
# JSON output
pomme sales monthly --output json

# CSV export
pomme sales monthly --output csv > sales.csv

# Disable cache
pomme sales monthly --no-cache
```

</details>

<details>
<summary>📊 Sales Comparison</summary>

### Compare Months

```bash
# Compare two months
pomme sales compare --current 2025-03 --previous 2025-02

# Year-over-year
pomme sales compare --current 2025-03 --previous 2024-03
```

### Output

The comparison shows:
- Units sold change
- Revenue change by currency
- Percentage differences
- Visual indicators (↑↓)

</details>

<details>
<summary>📈 Sales Trends</summary>

### Analyze Trends

```bash
# 6-month trend
pomme sales trends --months 6

# Full year
pomme sales trends --months 12

# Export data
pomme sales trends --months 6 --output json
```

### Features

- Automatic trend detection
- Growth rate calculations
- Seasonal pattern identification
- Top performing apps/regions

</details>

## Analytics Commands

<details>
<summary>🚀 Performance Metrics</summary>

### View Metrics

```bash
# All metrics for an app
pomme analytics show APP_ID

# Specific metric
pomme analytics show APP_ID --metric launch
pomme analytics show APP_ID --metric memory
pomme analytics show APP_ID --metric battery

# Filter by device
pomme analytics show APP_ID --device iPhone
pomme analytics show APP_ID --device iPad

# Show performance goals
pomme analytics show APP_ID --goals
```

### Available Metrics

- **Launch Time** - Time to first frame (ms)
- **Memory Usage** - Peak memory consumption (MB)
- **Battery Usage** - Battery drain (%/hr)
- **Hang Rate** - App hangs per hour
- **Disk Writes** - Disk I/O (MB)

</details>

<details>
<summary>🔄 Version Comparison</summary>

### Compare Versions

```bash
# Compare two versions
pomme analytics compare APP_ID --version1 1.2.0 --version2 1.3.0

# Show improvements/regressions
pomme analytics compare APP_ID --version1 1.2.0 --version2 1.3.0 --details
```

</details>

<details>
<summary>📉 Performance Trends</summary>

### View Trends

```bash
# Overall trends
pomme analytics trends APP_ID

# Specific metric trend
pomme analytics trends APP_ID --metric launch

# Device-specific
pomme analytics trends APP_ID --device iPhone
```

</details>

## Reviews Commands

<details>
<summary>⭐ List Reviews</summary>

### Basic Listing

```bash
# List recent reviews
pomme reviews list APP_ID

# Limit results
pomme reviews list APP_ID --limit 50

# Show full content
pomme reviews list APP_ID --verbose
```

### Filtering

```bash
# By rating
pomme reviews list APP_ID --rating 1  # 1-star reviews
pomme reviews list APP_ID --rating 5  # 5-star reviews

# By territory
pomme reviews list APP_ID --territory US
pomme reviews list APP_ID --territory GB

# Sort options
pomme reviews list APP_ID --sort recent    # Default
pomme reviews list APP_ID --sort critical  # Low ratings first
pomme reviews list APP_ID --sort helpful   # Most helpful
```

</details>

<details>
<summary>📊 Review Analytics</summary>

### Summary Statistics

```bash
# Overall summary
pomme reviews summary APP_ID

# Territory-specific
pomme reviews summary APP_ID --territory US
```

### Summary Includes

- Average rating
- Rating distribution
- Territory breakdown
- Recent review preview
- Total review count

</details>

<details>
<summary>🔍 Search & Respond</summary>

### Search Reviews

```bash
# Search by keyword
pomme reviews search APP_ID "crash"
pomme reviews search APP_ID "love"
pomme reviews search APP_ID "bug"

# Case-insensitive search
pomme reviews search APP_ID "CRASH"
```

### Respond to Reviews

```bash
# Simple response
pomme reviews respond REVIEW_ID "Thank you for your feedback!"

# Multi-line response
pomme reviews respond REVIEW_ID "Thank you for the 5-star review! 
We're glad you're enjoying the app.
Stay tuned for more updates!"
```

</details>

<details>
<summary>👀 Watch Mode</summary>

### Monitor Reviews

```bash
# Watch for new reviews
pomme reviews watch APP_ID

# Filter by minimum rating
pomme reviews watch APP_ID --min-rating 4

# Custom interval (seconds)
pomme reviews watch APP_ID --interval 300
```

</details>

## Tips & Tricks

<details>
<summary>🛠️ Advanced Usage</summary>

### Shell Aliases

Add to your `.bashrc` or `.zshrc`:

```bash
alias sales='pomme sales'
alias reviews='pomme reviews list'
alias metrics='pomme analytics show'
alias review-summary='pomme reviews summary'
```

### JSON Processing with jq

```bash
# Total revenue across currencies
pomme sales monthly --output json | jq '[.revenue[]] | add'

# Count reviews by rating
pomme reviews list APP_ID --output json | \
  jq '.data | group_by(.attributes.rating) | 
      map({rating: .[0].attributes.rating, count: length})'

# Extract app names and units
pomme sales monthly --output json | \
  jq '.apps[] | {name: .name, units: .units}'
```

### Automation Scripts

```bash
#!/bin/bash
# daily-report.sh - Email daily sales summary

REPORT=$(pomme sales --output json | jq -r '
  "Sales Report for " + .month + "\n" +
  "Total Units: " + (.totalUnits|tostring) + "\n" +
  "Revenue: " + (.revenue | to_entries | 
    map(.key + " " + (.value|tostring)) | join(", "))
')

echo "$REPORT" | mail -s "Daily Sales Report" you@example.com
```

### Monitoring Critical Reviews

```bash
#!/bin/bash
# monitor-reviews.sh - Alert on critical reviews

CRITICAL=$(pomme reviews list APP_ID --rating 1 --output json | \
  jq -r '.data[] | "New 1-star review from " + 
         .attributes.reviewerNickname + ": " + 
         .attributes.title')

if [ ! -z "$CRITICAL" ]; then
  echo "$CRITICAL" | mail -s "Critical Review Alert" you@example.com
fi
```

</details>

<details>
<summary>🔧 Troubleshooting</summary>

### Common Issues

#### Config File Errors
```
Error: yaml: control characters are not allowed
```
**Solution**: Recreate the config file ensuring no hidden characters

#### Authentication Failed
```
Error: API error: UNAUTHORIZED
```
**Solution**: 
- Verify API key permissions
- Check key hasn't expired
- Ensure correct key ID and issuer ID

#### Rate Limiting
```
Error: API error: RATE_LIMIT_EXCEEDED
```
**Solution**:
- Use built-in caching (default)
- Add delays in automation
- Batch operations when possible

### Debug Mode

```bash
# Enable verbose logging
POMME_DEBUG=1 pomme sales

# Check config
pomme config show

# Validate auth
pomme auth validate
```

</details>

<details>
<summary>📝 Examples</summary>

### Daily Automation

```bash
# Crontab entries
# Daily sales report at 9 AM
0 9 * * * /usr/local/bin/pomme sales >> ~/pomme-sales.log

# Weekly review summary on Mondays
0 10 * * 1 /usr/local/bin/pomme reviews summary APP_ID

# Monitor critical reviews every hour
0 * * * * /home/user/scripts/check-critical-reviews.sh
```

### Reporting Dashboard

```bash
#!/bin/bash
# dashboard.sh - Generate HTML dashboard

cat > dashboard.html << EOF
<html>
<head><title>App Dashboard</title></head>
<body>
<h1>App Performance Dashboard</h1>
<h2>Sales This Month</h2>
<pre>$(pomme sales)</pre>
<h2>Recent Reviews</h2>
<pre>$(pomme reviews list APP_ID --limit 5)</pre>
<h2>Performance Metrics</h2>
<pre>$(pomme analytics show APP_ID)</pre>
</body>
</html>
EOF
```

### CSV Export for Excel

```bash
# Export sales to CSV
pomme sales monthly --output csv > sales_$(date +%Y%m).csv

# Export reviews to CSV
pomme reviews list APP_ID --output csv > reviews_$(date +%Y%m%d).csv

# Combine multiple months
for month in 2025-01 2025-02 2025-03; do
  pomme sales monthly $month --output csv
done > q1_sales.csv
```

</details>

## Support

- **Issues**: [GitHub Issues](https://github.com/marcusziade/pomme/issues)
- **Documentation**: [GitHub Wiki](https://github.com/marcusziade/pomme/wiki)
- **Updates**: Watch the repository for new features