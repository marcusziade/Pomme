# Authentication

Pomme uses JWT (JSON Web Token) authentication to securely communicate with the App Store Connect API. Before using Pomme, you'll need to set up your authentication credentials.

## Prerequisites

To authenticate with the App Store Connect API, you'll need:

1. An App Store Connect account with the appropriate role (Admin, Developer, Marketing, etc.)
2. API key access enabled for your account
3. The following credentials:
   - Key ID (from your API key)
   - Issuer ID (from your App Store Connect account)
   - API Private Key (p8 file) - downloaded when you create your API key

## Setting Up Credentials

### Step 1: Create an API Key in App Store Connect

1. Log in to [App Store Connect](https://appstoreconnect.apple.com/)
2. Go to "Users and Access" > "Keys"
3. Click the "+" button to create a new key
4. Enter a name for your key and select the appropriate access level
5. Click "Generate"
6. Download the API key (p8 file) and note the Key ID

### Step 2: Configure Pomme

Initialize the Pomme configuration file:

```bash
pomme config init
```

This will create a configuration file at `~/.config/pomme/pomme.yaml` (on macOS/Linux) or `%APPDATA%\pomme\pomme.yaml` (on Windows).

Edit the configuration file and add your authentication credentials:

```yaml
auth:
  key_id: "YOUR_KEY_ID"
  issuer_id: "YOUR_ISSUER_ID"
  private_key_path: "/path/to/your/private_key.p8"
```

### Step 3: Test Authentication

Verify that your authentication credentials are working correctly:

```bash
pomme auth test
```

If successful, you should see:

```
Authentication successful!
JWT token generated: eyJhbGciOi...
```

## Security Best Practices

1. **Never commit your private key to version control**
2. **Set restrictive permissions on your config file**:
   ```bash
   chmod 600 ~/.config/pomme/pomme.yaml
   ```
3. **Use environment variables in CI/CD environments** instead of config files
4. **Rotate your API keys** regularly, especially in shared environments
5. **Use the minimum required access level** for your needs

## Troubleshooting

### Invalid Authentication (NOT_AUTHORIZED Error)

If you see an error like:

```
API Error [NOT_AUTHORIZED]: Provide a properly configured and signed bearer token, and make sure that it has not expired.
```

Check that:

1. Your Key ID and Issuer ID are correct and match exactly what's shown in App Store Connect
2. The private key file path is correct and the file is readable
3. Your API key hasn't expired or been revoked in App Store Connect
4. The vendor number in your config is correct (find this in App Store Connect > Agreements, Tax, and Banking)
5. The private key file (.p8) is in the correct format and contains the headers:
   ```
   -----BEGIN PRIVATE KEY-----
   ...key content...
   -----END PRIVATE KEY-----
   ```
6. Make sure your system clock is accurate, as time skew can affect JWT token validation
7. If you recently created your API key, there could be a propagation delay - wait a few minutes and try again
8. Try generating a new API key in App Store Connect and update your config with the new credentials

### Permission Issues

If you see an error like:

```
API Error [FORBIDDEN]: Not authorized - You don't have permission to access this resource (Status: 403)
```

Check that:

1. Your API key has the necessary permissions for the requested resource
2. Your App Store Connect user account has the necessary role
