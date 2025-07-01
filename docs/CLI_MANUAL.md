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
<summary>🚀 Quick Start - Interactive Setup Wizard</summary>

The easiest way to configure Pomme is using our interactive setup wizard:

```bash
pomme config init
```

This wizard will:
- ✅ Guide you through creating an App Store Connect API key
- ✅ Help you download and save your private key
- ✅ Prompt for your credentials with examples
- ✅ Automatically validate your configuration
- ✅ Get you ready to use Pomme in 5 minutes!

</details>

<details>
<summary>📋 Configuration Commands</summary>

### Available Commands

```bash
# Interactive setup wizard (recommended for first-time users)
pomme config init

# View current configuration (credentials are masked)
pomme config show

# Validate your configuration
pomme config validate

# Show detailed setup instructions
pomme config help
```

### Command Details

#### `pomme config init`
Interactive wizard that walks you through the entire setup process:
- Opens relevant App Store Connect pages
- Provides step-by-step instructions
- Validates inputs as you go
- Tests your connection before finishing

#### `pomme config show`
Displays your current configuration with sensitive values masked:
```
🔧 Current Configuration
────────────────────────────────────────────────────────────

Config file: /Users/john/.config/pomme/pomme.yaml

API Settings:
  Base URL: https://api.appstoreconnect.apple.com/v1
  Timeout: 30 seconds

Authentication:
  Key ID: 73****5R
  Issuer ID: a5****7b
  Private Key: /Users/john/.config/pomme/AuthKey_73TT63DP5R.p8
               ✓ File exists

Defaults:
  Output Format: table
  Vendor Number: 93****63
```

#### `pomme config validate`
Tests your configuration by:
- Checking all required fields are set
- Verifying your private key file exists
- Testing API connection with your credentials
- Providing clear error messages if something's wrong

#### `pomme config help`
Comprehensive guide showing:
- How to create API keys in App Store Connect
- Where to find each required value
- Security best practices
- Troubleshooting tips

</details>

<details>
<summary>🔐 Manual Configuration</summary>

### Config File Location

Pomme looks for configuration in these locations (in order):
1. `./pomme.yaml` (current directory)
2. `~/.config/pomme/pomme.yaml` (recommended)
3. `~/.pomme.yaml`
4. Environment variables (with `POMME_` prefix)

### Config File Format

Create `~/.config/pomme/pomme.yaml`:

```yaml
api:
  base_url: https://api.appstoreconnect.apple.com/v1
  timeout: 30
defaults:
  output_format: table
  vendor_number: YOUR_VENDOR_NUMBER  # Optional
auth:
  key_id: YOUR_KEY_ID
  issuer_id: YOUR_ISSUER_ID
  private_key_path: /path/to/your/AuthKey.p8
```

### Environment Variables

You can also use environment variables (useful for CI/CD):

```bash
export POMME_AUTH_KEY_ID=YOUR_KEY_ID
export POMME_AUTH_ISSUER_ID=YOUR_ISSUER_ID
export POMME_AUTH_PRIVATE_KEY_PATH=/path/to/key.p8
export POMME_DEFAULTS_VENDOR_NUMBER=93036463
```

Environment variables override values from the config file.

</details>

<details>
<summary>🔑 Getting Your API Credentials</summary>

### Step 1: Create an API Key

1. Sign in to [App Store Connect](https://appstoreconnect.apple.com)
2. Navigate to **Users and Access**
3. Click **Keys** tab under "App Store Connect API"
4. Click the **+** button to generate a new key
5. Enter a name (e.g., "Pomme CLI")
6. Choose access level:
   - **Admin** - Full access (recommended)
   - **Finance** - Sales reports only
   - **Sales** - Limited sales access

### Step 2: Download Your Private Key

⚠️ **Important**: You can only download the private key once!

1. Click **Generate**
2. Click **Download API Key**
3. Save the `.p8` file securely, for example:
   - `~/.config/pomme/AuthKey_XXXXXXXXXX.p8`
   - `~/Documents/AppStoreConnect/AuthKey_XXXXXXXXXX.p8`

### Step 3: Note Your Credentials

From the Keys page, you'll need:
- **Key ID**: 10 characters (e.g., `73TT63DP5R`)
- **Issuer ID**: UUID format (e.g., `a5ebdab5-0ceb-463c-8151-195b902f117b`)

### Step 4: Find Your Vendor Number (Optional)

1. Go to **Payments and Financial Reports**
2. Your vendor number is displayed at the top
3. Format: 8 digits (e.g., `93036463`)

</details>

<details>
<summary>🛡️ Security Best Practices</summary>

### Private Key Security

- **Never commit** your `.p8` file to version control
- **Store securely** in your home directory with restricted permissions:
  ```bash
  chmod 600 ~/.config/pomme/AuthKey_*.p8
  ```
- **Use environment variables** for CI/CD instead of files
- **Rotate keys** periodically through App Store Connect

### Access Control

- Create keys with **minimum required permissions**
- Use **Finance** role for sales-only access
- **Revoke unused keys** in App Store Connect
- Monitor key usage in your account

### Configuration Security

- Config file permissions:
  ```bash
  chmod 600 ~/.config/pomme/pomme.yaml
  ```
- Use `.gitignore` for local config files:
  ```
  pomme.yaml
  .pomme.yaml
  *.p8
  ```

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