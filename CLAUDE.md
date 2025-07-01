# Pomme Development Notes

## Architecture

- Clean separation of concerns with service layers
- Robust error handling with clear messages
- Concurrent processing using Go routines
- In-memory caching with automatic expiration
- No third-party dependencies except cobra (CLI framework) and jwt

## Key Design Decisions

1. **Removed Viper** - Implemented simple YAML parser to reduce dependencies
2. **Multi-Currency Support** - Sales reports properly handle and display multiple currencies
3. **Apple's 5-day Delay** - Automatically adjusts for monthly report availability
4. **Service Pattern** - Each feature (sales, analytics, reviews) has its own service

## Testing

```bash
go fmt ./...
go vet ./...
go test ./...
```

## Known Issues

1. TestFlight features not yet implemented
2. Analytics API endpoints need real-world testing
3. Watch mode for reviews pending implementation

## Future Enhancements

- Excel/PDF export formats
- Real-time notifications
- Predictive analytics
- Full App Store Connect API coverage