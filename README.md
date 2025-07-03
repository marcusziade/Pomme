# 🍎 Pomme

Beautiful App Store Connect CLI for sales reports, analytics, and reviews. Built with Go.

[![Go Report Card](https://goreportcard.com/badge/github.com/marcusziade/pomme)](https://goreportcard.com/report/github.com/marcusziade/pomme)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**[📖 Documentation](https://marcusziade.github.io/pomme) | [🚀 Get Started](https://marcusziade.github.io/pomme#getting-started) | [⭐ Star on GitHub](https://github.com/marcusziade/pomme)**

## ✨ Features

- **📊 Sales Reports** - View monthly sales with multi-currency support and beautiful formatting
- **⭐ Review Management** - Monitor, analyze, and respond to customer reviews
- **🎯 Smart CLI** - Interactive setup wizard, automatic validation, and intuitive commands
- **🚀 Fast & Secure** - Built with Go for speed, uses official App Store Connect API

## 🚀 Quick Start

### 1. Install

#### Arch Linux (AUR) - Recommended for Arch users

```bash
# Using yay
yay -S pomme

# Using paru
paru -S pomme

# For development version
yay -S pomme-git
```

#### Other Installation Methods

```bash
# With Homebrew (macOS)
brew tap marcusziade/tap
brew install pomme

# With Go
go install github.com/marcusziade/pomme/cmd/pomme@latest

# Or download the latest binary
curl -L https://github.com/marcusziade/pomme/releases/latest/download/pomme_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv pomme /usr/local/bin/
```

### 2. Configure (5 minutes)

```bash
pomme config init
```

Our interactive wizard will guide you through:
- Creating an App Store Connect API key
- Downloading your credentials
- Validating everything works

### 3. Start Using

```bash
# View your latest sales
pomme sales

# See customer reviews
pomme reviews list <app-id>

# Get detailed help
pomme --help
```

## 📸 Examples

### Sales Reports
```bash
$ pomme sales

📊 Sales Report for December 2024
────────────────────────────────────────────────────────────

  📦 Total Units
  125,847

  💰 Revenue
  USD 486,392.45 (73.2% of total)
  EUR 124,836.20 (18.8%)
  JPY ¥693,450 (5.2%)
  GBP £15,234.89 (2.8%)

  🌍 Countries
  142 markets
```

### Review Summary
```bash
$ pomme reviews summary <app-id>

📊 Review Summary
════════════════════════════════════════════════════════════

Overall Statistics
────────────────────────────────────────
  Average Rating: 4.6 ⭐⭐⭐⭐⭐
  Total Reviews: 8,743

Rating Distribution
────────────────────────────────────────
  5⭐ ████████████████████████ 5,892 (67.4%)
  4⭐ ████████████              1,834 (21.0%)
  3⭐ ████                        623 ( 7.1%)
  2⭐ ██                          287 ( 3.3%)
  1⭐ █                           107 ( 1.2%)
```

## 🛠 Commands

### Configuration
- `pomme config init` - Interactive setup wizard
- `pomme config validate` - Test your credentials
- `pomme config show` - View current config

### Sales
- `pomme sales` - Latest monthly report
- `pomme sales monthly 2024-03` - Specific month
- `pomme sales compare --current 2024-03 --previous 2024-02` - Compare periods

### Reviews
- `pomme reviews list <app-id>` - List reviews
- `pomme reviews summary <app-id>` - Statistics
- `pomme reviews respond <review-id> "message"` - Respond

### Apps
- `pomme apps list` - List all apps
- `pomme apps info <app-id>` - App details

## 📚 Documentation

- [CLI Manual](docs/CLI_MANUAL.md) - Comprehensive command reference
- [Configuration Guide](https://marcusziade.github.io/pomme#getting-started) - Detailed setup instructions
- [Development Notes](CLAUDE.md) - Architecture and contributing

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines and submit pull requests to our repository.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.