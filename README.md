# 🔥 Mikhaz Fire

A lightweight, high-performance cross-platform desktop application for managing **Firebase Auth** and **Firestore** databases. Built with **Go**, **Nuxt 4**, and **Wails**, it allows developers to quickly inspect and control their Firebase projects with zero-configuration setup.

---

## ✨ Key Features

- 🔑 **Multi-Source Authentication**:
  - **1-Click Google Sign In**: Sign in directly with any Google account in your browser with zero manual configuration.
  - **Local CLI Discovery**: Automatically scans and loads credentials from **Firebase CLI** (`firebase-tools.json`) and **gcloud SDK** (`legacy_credentials`, ADC).
  - **Service Account JSON**: Import custom Service Account `.json` key files directly.
- 👥 **Firebase Auth Management**:
  - Browse, filter, and search users in real-time.
  - Create new users, update user profiles, or delete users.
  - View and edit custom user claims (JSON metadata).
  - 📤 **Export Auth Users**: Export users with custom claims, password hashes, salts, and embedded Password Hash Parameters (`--hash-key`, `--salt-separator`, `--rounds`, `--mem-cost`).
  - 📥 **Import Auth Users**: Import user batches with custom claims and password hash parameters so users can sign in without resetting passwords. Supports both Mikhaz Fire and Firebase CLI export formats.
- 🗄️ **Firestore Explorer**:
  - Traverse collections, subcollections, and documents hierarchically.
  - View, modify, create, and overwrite documents.
  - Perform recursive deletes of documents and all their nested subcollections.
- ⚡ **Firestore Query & JS Engine**:
  - Execute standard structured Firestore queries using a visual filter interface.
  - Run custom JavaScript scripts locally against your Firestore database via an embedded JS runner (`goja`).
- 🔄 **Direct Cross-Project Copying / Syncing**:
  - Copy collections or single documents directly from a source project to a destination project (even across different authenticated Google accounts).
  - Built-in conflict checking to prevent accidental overwrites.

---

## 🔐 Authentication Guide

Mikhaz Fire supports 3 seamless ways to connect your Google Cloud / Firebase projects:

### Method 1: 1-Click Google Sign In (Recommended)
1. Click the **Google Sign In** button in the sidebar.
2. Your browser will open the official Google authentication page.
3. Grant access to your account.
4. Your account and all associated Firebase projects will automatically appear in the sidebar!

> **Note on Project Permissions**:  
> No modifications to your GCP/Firebase project settings are required! As long as your Google user account has permission to view/manage Firebase or GCP projects (e.g. Owner, Editor, or Firebase Admin role), Mikhaz Fire will discover and access your projects out of the box.

### Method 2: Service Account JSON Key
1. Go to your **Firebase Console** → **Project Settings** → **Service Accounts** tab.
2. Click **Generate new private key** to download the `.json` file.
3. In Mikhaz Fire, click the **Key** icon (`+`) in the sidebar and upload/paste your Service Account JSON.

### Method 3: Local CLI Session Scan
If you already use the Firebase CLI or Google Cloud SDK:
- Run `firebase login` or `gcloud auth login` in your terminal.
- Click **Rescan** in Mikhaz Fire. The app automatically picks up your active CLI tokens!


---

## 🔒 Security & Privacy

Mikhaz Fire is designed with security and privacy as first-class priorities:
- **No Remote Servers**: Your credentials, session tokens, and database payloads are processed entirely **in-memory** on your local machine.
- **Direct Requests**: The application communicates directly and exclusively with official Google Cloud and Firebase REST APIs via HTTPS.
- **No Key Storage**: It does not persist or cache your service account keys or refresh tokens to disk.
- **Auditable**: Because the code is open-source, you can easily inspect [credentials.go](credentials.go) to see exactly how your credentials are read and used.

---

## 🛠️ Prerequisites

To run or build Mikhaz Fire from source, you will need:
- **Go** (version 1.25 or later recommended)
- **Node.js** & **npm** (for the Nuxt frontend)
- **Wails CLI** (for packaging the desktop app)

To install the Wails CLI, run:
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

---

## 🚀 Getting Started

### Run in Development Mode
To start the app in live-development mode (with hot-reloading for both Go and Nuxt):
```bash
wails dev
```

### Build the Production Binary
To compile the production-ready standalone executable:
```bash
wails build
```
Once the build completes, the executable will be generated at:
- **Windows**: `build/bin/mikhaz-fire.exe`
- **macOS/Linux**: `build/bin/mikhaz-fire`

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
