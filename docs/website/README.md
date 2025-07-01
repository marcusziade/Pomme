# Pomme Landing Page

This is the landing page for Pomme, hosted on GitHub Pages.

## 🚀 Local Development

### Quick Start

```bash
# Python 3 (recommended)
./serve.py

# Or use the shell script
./serve.sh

# Or manually with Python
python3 -m http.server 8000
```

The server will:
- Start on http://localhost:8000
- Automatically open your browser
- Disable caching for development
- Show request logs

### Alternative Servers

If you prefer other tools:

```bash
# Node.js
npx http-server -p 8000 -c-1

# Go
go run -mod=mod github.com/shurcooL/goexec@latest 'http.ListenAndServe(":8000", http.FileServer(http.Dir(".")))'

# Ruby
ruby -run -e httpd . -p 8000
```

## 📁 Files

- `index.html` - Main landing page
- `styles.css` - All styling
- `script.js` - Interactivity
- `serve.py` - Local preview server
- `.nojekyll` - Disables Jekyll processing on GitHub Pages

## 🚢 Deployment

The site automatically deploys to GitHub Pages when changes are pushed to the `master` branch.

Workflow: `.github/workflows/deploy-website.yml`

## 🎨 Features

- Dark theme with App Store Connect branding
- Interactive terminal demos
- Smooth animations and transitions
- Mobile responsive
- Fast loading (no build step)

## 🛠 Making Changes

1. Edit the HTML/CSS/JS files
2. Preview locally with `./serve.py`
3. Commit and push to `master`
4. GitHub Actions deploys automatically

The site will be available at: https://marcusziade.github.io/pomme