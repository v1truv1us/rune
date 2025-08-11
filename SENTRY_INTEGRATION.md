# Sentry Integration Guide

Rune CLI includes comprehensive Sentry integration for error tracking, performance monitoring, and crash reporting. This guide explains how to set up and use Sentry with Rune.

## Features

### Error Tracking
- Automatic capture of command failures and exceptions
- Detailed error context including command, arguments, and system information
- Stack traces for debugging
- Error grouping and deduplication

### Performance Monitoring
- Command execution time tracking
- Performance transactions for key operations
- Slow command detection and alerting
- Resource usage monitoring

### User Context
- Anonymous user identification for session tracking
- System information (OS, version, architecture)
- Application version and environment context
- Command usage patterns

### Breadcrumbs
- Command execution history
- User actions and events
- System state changes
- Debug information trail

## Setup

### 1. Get Your Sentry DSN

1. Create a Sentry account at [sentry.io](https://sentry.io)
2. Create a new project for your Rune CLI
3. Copy the DSN from your project settings
4. The DSN format is: `https://public_key@sentry.io/project_id`

### 2. Configuration Methods

#### Method 1: Environment Variables (Recommended)
```bash
export RUNE_SENTRY_DSN="https://your_public_key@sentry.io/your_project_id"
# Optional: configure OTLP logs endpoint for OpenTelemetry
export RUNE_OTLP_ENDPOINT="http://localhost:4318/v1/logs"
```

#### Method 2: Configuration File
Add to your `~/.rune/config.yaml`:
```yaml
integrations:
  telemetry:
    enabled: true
    sentry_dsn: "https://your_public_key@sentry.io/your_project_id"
    # Optional OpenTelemetry logs endpoint
    otlp_endpoint: "http://localhost:4318/v1/logs"
```

#### Method 3: Example Configuration
Use the provided example:
```bash
cp examples/config-telemetry.yaml ~/.rune/config.yaml
# Edit the file to add your actual DSN
```

### 3. Environment Configuration

Set the environment for better error categorization:
```bash
export RUNE_ENV="production"  # or "development", "staging"
```

## Privacy and Data Collection

### What Data is Collected

#### Error Data
- Error messages and stack traces
- Command that caused the error
- System information (OS, architecture)
- Application version
- Anonymous user ID

#### Performance Data
- Command execution times
- Success/failure rates
- System resource usage
- Performance bottlenecks

#### User Context
- Anonymous user identifier (hostname-based)
- Operating system and version
- Application version and environment
- Command usage patterns

### What Data is NOT Collected
- Personal information or credentials
- File contents or sensitive data
- Network requests or API keys
- User-specific paths or filenames
- Any personally identifiable information

### Data Anonymization
- User IDs are generated anonymously using hostname and timestamp
- No personal information is transmitted
- All data is aggregated and anonymized
- Sensitive information is filtered out before transmission

## Disabling Telemetry

### Temporary Disable
```bash
export RUNE_TELEMETRY_DISABLED=true
rune start  # Telemetry disabled for this session
```

### Permanent Disable
Add to your configuration:
```yaml
integrations:
  telemetry:
    enabled: false
```

Or set environment variable permanently:
```bash
echo 'export RUNE_TELEMETRY_DISABLED=true' >> ~/.bashrc
```

## Monitoring and Alerts

### Error Alerts
Set up alerts in Sentry for:
- Command failures
- Performance degradation
- High error rates
- New error types

### Performance Monitoring
Monitor:
- Average command execution time
- Slow commands (>5 seconds)
- Memory usage patterns
- System resource consumption

### Custom Dashboards
Create dashboards for:
- Command usage statistics
- Error trends over time
- Performance metrics
- User engagement patterns

## Development and Testing

### Local Development
```bash
export RUNE_ENV="development"
export RUNE_SENTRY_DSN="https://dev_key@sentry.io/dev_project_id"
```

### Testing Error Reporting
Force an error to test Sentry integration:
```bash
# This will trigger error reporting
rune start nonexistent_project_with_invalid_config
```

### Debug Mode
Enable verbose logging:
```bash
rune --verbose status
```

## Advanced Configuration

### Custom Error Handling
The telemetry system provides additional functions for custom error handling:

```go
// In your custom commands
telemetry.CaptureException(err, map[string]string{
    "command": "custom_command",
    "context": "additional_context",
}, map[string]interface{}{
    "extra_data": "value",
})

telemetry.CaptureMessage("Custom message", sentry.LevelWarning, map[string]string{
    "category": "custom",
})
```

### Performance Monitoring
Track custom operations:

```go
transaction := telemetry.StartTransaction("custom_operation", "operation")
if transaction != nil {
    defer transaction.Finish()
    // Your operation here
}
```

## Troubleshooting

### Common Issues

#### Sentry Not Receiving Data
1. Verify DSN format is correct
2. Check network connectivity
3. Ensure telemetry is enabled
4. Verify Sentry project settings

#### Performance Impact
- Telemetry is designed to be lightweight
- Errors are sent asynchronously
- Minimal performance overhead (<1ms per command)
- Automatic rate limiting prevents spam

#### Privacy Concerns
- All data is anonymized
- No sensitive information is collected
- Easy to disable completely
- Transparent data collection practices

### Debug Commands
```bash
# Check telemetry status
rune config show

# Test with verbose output
rune --verbose start

# Disable temporarily
RUNE_TELEMETRY_DISABLED=true rune start
```

## Support

For issues with Sentry integration:
1. Check this documentation
2. Verify your Sentry project configuration
3. Test with verbose logging enabled
4. Create an issue in the Rune repository

## Security

- All data transmission uses HTTPS
- No credentials or secrets are transmitted
- Anonymous user identification only
- Compliant with privacy regulations
- Regular security audits of telemetry code

---

**Note**: Telemetry helps improve Rune by identifying common issues and usage patterns. All data is collected anonymously and used solely for product improvement.