# SMTP Email Configuration Guide

This application now uses SMTP for sending emails instead of Google Apps Script. Follow this guide to configure email notifications.

## Supported Email Providers

### Gmail Configuration
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM=your-email@gmail.com
```

**Important for Gmail:**
- You need to use an "App Password" instead of your regular Gmail password
- Enable 2-Factor Authentication on your Gmail account
- Generate an App Password: Google Account → Security → 2-Step Verification → App passwords

### Outlook/Hotmail Configuration
```env
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USERNAME=your-email@outlook.com
SMTP_PASSWORD=your-password
SMTP_FROM=your-email@outlook.com
```

### Other SMTP Providers
- **SendGrid**: smtp.sendgrid.net:587
- **Mailgun**: smtp.mailgun.org:587
- **Amazon SES**: email-smtp.region.amazonaws.com:587

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `SMTP_HOST` | SMTP server hostname | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USERNAME` | Email username | `your-email@gmail.com` |
| `SMTP_PASSWORD` | Email password/app password | `your-app-password` |
| `SMTP_FROM` | From email address | `your-email@gmail.com` |

## Security Features

- **STARTTLS Encryption**: Automatically upgrades to TLS encryption
- **Authentication**: SMTP authentication is required
- **Error Handling**: Comprehensive error handling and logging

## Email Types Sent

The application sends the following types of notifications:

1. **Employee Creation** - Welcome emails with login credentials
2. **Leave Applications** - Notifications to managers and admins
3. **Leave Approvals/Rejections** - Status updates to employees and admins
4. **Leave Withdrawals** - Withdrawal notifications
5. **Password Updates** - Security notifications
6. **Payslip Withdrawals** - Payroll notifications

## Troubleshooting

### Common Issues

1. **Authentication Failed (535 5.7.8)**
   - Verify username and password
   - For Gmail, ensure you're using an App Password
   - Check if 2FA is enabled

2. **Connection Timeout**
   - Verify SMTP host and port
   - Check firewall settings
   - Ensure internet connectivity

3. **TLS Handshake Errors**
   - Use port 587 for STARTTLS (recommended)
   - Use port 465 for SSL/TLS (if supported)
   - Avoid port 25 (often blocked)

### Fixed Issues

- **TLS Handshake Error**: Fixed by using `smtp.SendMail` with STARTTLS instead of direct TLS connection
- **Connection Issues**: Simplified SMTP implementation for better compatibility

### Testing SMTP Configuration

You can test your SMTP configuration by checking the application logs when an email is sent. Look for:

```
Attempting to send email to: recipient@example.com with subject: Test Subject
Email sent successfully to: recipient@example.com
```

## Migration from Google Apps Script

The application has been migrated from Google Apps Script to SMTP. Key changes:

- **Removed**: `GOOGLE_SCRIPT_URL` environment variable
- **Added**: SMTP configuration variables
- **Improved**: Better error handling and logging
- **Enhanced**: Support for multiple email providers
- **Fixed**: TLS connection issues with proper STARTTLS implementation

## Production Deployment

For production deployment, ensure:

1. Use environment variables or secrets management
2. Never commit SMTP credentials to version control
3. Use strong, unique passwords
4. Monitor email delivery logs
5. Set up proper DNS records (SPF, DKIM) for better deliverability