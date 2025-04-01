# Quick Start Guide

This guide will help you get up and running with Pomme quickly.

## Installation

Install Pomme using Go:

```bash
go install github.com/marcus/pomme/cmd/pomme@latest
```

Or download a pre-built binary from the [releases page](https://github.com/marcus/pomme/releases).

## Setup

### 1. Initialize Configuration

Create a default configuration file:

```bash
pomme config init
```

### 2. Add Authentication Credentials

Edit the configuration file located at `~/.config/pomme/pomme.yaml` (on macOS/Linux) or `%APPDATA%\pomme\pomme.yaml` (on Windows) to add your App Store Connect API credentials:

```yaml
auth:
  key_id: "YOUR_KEY_ID"
  issuer_id: "YOUR_ISSUER_ID"
  private_key_path: "/path/to/your/private_key.p8"
```

### 3. Test Authentication

Verify your authentication is working:

```bash
pomme auth test
```

## Common Tasks

### List Your Apps

```bash
pomme apps list
```

### Get Sales Reports

Get the latest daily sales report:

```bash
pomme sales report --period=DAILY --date=latest
```

Get a specific sales report:

```bash
pomme sales report --period=MONTHLY --date=2023-01 --type=SUBSCRIPTION
```

Get a simple summary of the past month's sales:

```bash
pomme sales monthly
```

### View Recent Reviews

```bash
pomme reviews list --app=YOUR_APP_ID --limit=10
```

### Get App Analytics

```bash
pomme analytics usage --app=YOUR_APP_ID --metric=activeDevices --start-date=2023-01-01 --end-date=2023-01-31
```

## Output Formats

You can get output in different formats using the global `-o` flag:

### JSON

```bash
pomme apps list -o json
```

### CSV

```bash
pomme sales report --period=MONTHLY --date=2023-01 -o csv > sales.csv
```

### Table (Default)

```bash
pomme reviews list --app=YOUR_APP_ID
```

## Next Steps

For more detailed information, check out:

- [Full Command Reference](./commands.md)
- [Configuration Options](./configuration.md)
- [Authentication Details](./authentication.md)
- [Sales Reporting Guide](./sales.md)
