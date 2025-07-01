# 🍎 Pomme

A powerful, production-ready App Store Connect CLI tool built with Go.

**[📖 View Documentation](https://marcusziade.github.io/pomme) | [🚀 Quick Start Guide](https://marcusziade.github.io/pomme#installation) | [⭐ Star on GitHub](https://github.com/marcusziade/pomme)**

```
   _____                              
  |  __ \                             
  | |__) |___  _ __ ___  _ __ ___   ___ 
  |  ___/ _ \| '_ ` _ \| '_ ` _ \ / _ \
  | |  | (_) | | | | | | | | | | |  __/
  |_|   \___/|_| |_| |_|_| |_| |_|\___|
                                     
  Your App Store Connect CLI companion
```

## Features

- 🔐 JWT authentication for App Store Connect API
- 📊 Sales and financial reporting
- 📱 App metadata management
- 💬 User reviews monitoring
- 📈 Analytics and subscription data
- 🔄 Multiple output formats (JSON, CSV, Table)

## Installation

```bash
# From binary releases
curl -sfL https://raw.githubusercontent.com/marcusziade/pomme/master/scripts/install.sh | bash

# Using Go
go install github.com/marcusziade/pomme/cmd/pomme@latest
```

## Quick Start

### Configure API Credentials

```bash
# Configure your API credentials
pomme config init
```

This will create a configuration file at `~/.config/pomme/pomme.yaml` where you can add your App Store Connect API credentials.

### Check Your Apps

```bash
# List your apps
pomme apps list
```

Output:
```
type    id              attributes.name    attributes.bundleId   attributes.sku
App     1234567890      Awesome App        com.you.app           APP001
App     9876543210      Amazing Game       com.you.game          GAME002
```

### View Recent Reviews

```bash
# Get recent reviews
pomme reviews list --app=YOUR_APP_ID --limit=10
```

Output:
```
type        id          attributes.rating  attributes.title                  attributes.reviewerNickname  attributes.createdDate
Review      r123456     5                  Love this app!                    HappyUser123                 2023-04-01
Review      r234567     4                  Great but needs one feature       AlmostSatisfied              2023-03-28
Review      r345678     3                  App crashes sometimes             BugReporter                  2023-03-25
```

### View Sales Reports

```bash
# Get your latest sales report
pomme sales report --period=DAILY --date=latest
```

Output:
```
Fetching daily sales report for 2023-04-01...
Report successfully retrieved!
Report size: 123456 bytes
```

```bash
# Show sales summary for the past month
pomme sales monthly
```

Output:
```
Fetching monthly sales report for 2023-03...

AppName         AppID           Units           TotalProceeds   Currency
Awesome App     1234567890      156             109.20          USD  
Premium Pass    5678901234      42              84.00           USD
Pro Upgrade     9876543210      18              89.99           USD

Total: 216 units, 283.19 USD
```

### View Analytics

```bash
# Get app analytics
pomme analytics usage --app=YOUR_APP_ID --metric=activeDevices --start-date=2023-01-01 --end-date=2023-01-31
```

Note: Analytics functionality is planned for future releases.

## Output Formats

You can get output in different formats using the global `-o` flag:

### JSON Format

```bash
pomme apps list -o json
```

Output:
```json
[
  {
    "type": "App",
    "id": "1234567890",
    "attributes": {
      "name": "Awesome App",
      "bundleId": "com.you.app",
      "sku": "APP001"
    }
  }
]
```

### CSV Format

```bash
pomme sales monthly -o csv
```

Output:
```
AppName,AppID,Units,TotalProceeds,Currency
Awesome App,1234567890,156,109.20,USD
Premium Pass,5678901234,42,84.00,USD
Pro Upgrade,9876543210,18,89.99,USD
```

## Documentation

See the [docs](./docs) directory for detailed documentation:

- [Quick Start Guide](./docs/quick-start.md)
- [Authentication Guide](./docs/authentication.md)

## License

MIT