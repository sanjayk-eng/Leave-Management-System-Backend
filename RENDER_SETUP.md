# Quick Render Setup Guide

## Why Deployment Isn't Working

The GitHub Action is failing because **you don't need the API call at all!** Render can auto-deploy directly from GitHub.

## Fix: Use Render's Built-in Auto-Deploy

### Step 1: Connect to Render (One-time setup)

1. Go to https://dashboard.render.com/
2. Click **"New +"** → **"Web Service"**
3. Click **"Connect GitHub"** and authorize Render
4. Select your repository: `UserMenagmentSystem_Backend`
5. Click **"Connect"**

### Step 2: Configure Service

Render will detect your `render.yaml` file automatically. If not, use these settings:

- **Name:** `zenithive-backend` (or your preferred name)
- **Region:** Oregon (or closest to you)
- **Branch:** `main`
- **Build Command:** `go build -o ums-backend main.go`
- **Start Command:** `./ums-backend`
- **Plan:** Free

### Step 3: Add Environment Variables

In Render Dashboard → Your Service → Environment tab, add:

```
DB_URL=your_postgresql_connection_string
JWT_SECRET=your_jwt_secret_key
PORT=8082
EMAIL_FROM=your_email@gmail.com
EMAIL_PASSWORD=your_app_password
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
```

### Step 4: Deploy

Click **"Create Web Service"** or **"Manual Deploy"**

That's it! From now on:
- Every push to `main` branch will auto-deploy
- No GitHub secrets needed
- No API calls needed

## How It Works

1. You push code to GitHub `main` branch
2. GitHub Actions runs tests and builds (for verification)
3. Render detects the push automatically
4. Render builds and deploys your app
5. Your app is live!

## Verify Deployment

After deployment:
1. Check Render Dashboard → Logs
2. Look for: "✅ Database connection established successfully"
3. Look for: "Migrations ran successfully!"
4. Test your API: `https://your-service.onrender.com/api/auth/login`

## Troubleshooting

### Build Fails
- Check Render logs for errors
- Verify Go version (should be 1.22)
- Ensure all dependencies are in `go.mod`

### App Crashes
- Check environment variables are set correctly
- Verify database connection string
- Check logs for error messages

### Migrations Fail
- Ensure `pkg/migration` folder is included
- Check migration SQL files are valid
- Verify database permissions

## Current Status

Your GitHub Action shows:
```
% Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
```

This means the API call is failing because:
1. Secrets are not configured, OR
2. You don't need the API call at all!

**Solution:** Just use Render's auto-deploy feature (Steps above). Remove the deploy job from GitHub Actions or add the deploy hook secret.

## Recommended Approach

**Use render.yaml auto-deploy** (already configured in your repo):
- ✅ No secrets needed
- ✅ Automatic deployment
- ✅ Simpler setup
- ✅ Works out of the box

The GitHub Action will just verify the build works, and Render handles the actual deployment.
