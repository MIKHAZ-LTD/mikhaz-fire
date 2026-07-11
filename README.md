# 🔥 Mikhaz Fire

A lightweight, high-performance cross-platform desktop application for managing **Firebase Auth** and **Firestore** databases. Built with **Go**, **Nuxt 4**, and **Wails**, it allows developers to quickly inspect and control their Firebase projects with zero-configuration setup.

---

## ✨ Key Features

- 🔑 **Zero-Config Authenticated Session Discovery**:
  - Automatically scans and loads credentials from local **Firebase CLI** configurations (`firebase-tools.json`).
  - Automatically scans and loads credentials from **gcloud SDK** configurations (legacy credentials and Application Default Credentials).
  - Supports adding custom **Service Account JSON** keys directly.
- 👥 **Firebase Auth Management**:
  - Browse, filter, and search users in real-time.
  - Create new users, update user profiles, or delete users.
  - Update custom user claims (JSON metadata) and toggle account status.
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
