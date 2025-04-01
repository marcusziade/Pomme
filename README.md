# Pomme

A powerful, production-ready App Store Connect CLI tool built with Go.

## Features

- JWT authentication for App Store Connect API
- Sales and financial reporting
- App metadata management
- User reviews monitoring
- Analytics and subscription data
- Multiple output formats (JSON, CSV, Table)

## Installation

```bash
# From binary releases
curl -sfL https://raw.githubusercontent.com/marcusziade/pomme/master/scripts/install.sh | bash
```

## Quick Start

```bash
# Configure your API credentials
pomme config init

# Get your latest sales report
pomme sales report --period=DAILY --date=latest

# Show sales summary for the past month
pomme sales monthly

# List your apps
pomme apps list

# Get recent reviews
pomme reviews list --app=YOUR_APP_ID --limit=10
```

## Documentation

See the [docs](./docs) directory for detailed documentation.

## License

MIT
