# 🏪 POSFlow

> Modern and free POS system for restaurants and retail stores

POSFlow is a complete point of sale system that allows you to manage your
business efficiently. Designed to work offline with dedicated interfaces
for different roles.

[Demo](#) • [Documentation](#) • [Downloads](#)

## 🚀 Features

- 📦 **Inventory Control**: Real-time product tracking
- 📋 **Order Management**: Queue system by arrival order
- 🍽️ **Daily Menu**: Easily update your menu
- 👥 **Multiple Views**: Customer, Kitchen, and Administrator
- 💰 **Day Closing**: Daily sales reports
- 🔌 **Offline First**: No internet connection required
- 🖥️ **Cross-Platform**: Works on Windows, macOS, and Linux

## 🛠️ Tech Stack

- **Frontend**: Electron + HTML/CSS/JavaScript
- **Backend**: Go + Gin framework
- **Database**: SQLite
- **Architecture**: Offline-first, local-first, single executable

## 📸 Screenshots

[Coming soon]

## 🚀 Quick Start
```bash
# Clone the repository
git clone https://github.com/AndresKenji/POSFlow.git

# Install dependencies
cd POSFlow
npm install
pip install -r requirements.txt

# Install frontend dependencies
cd electron-app
npm install

# Build and run backend (in another terminal)
cd ../backend
go build -o posflow-server
./posflow-server

# Run the Electron app
cd ../electron-app
npm start
```

## 📋 Requirements

- Node.js 16+
- Go 1.21+
- SQLite 3

## 📦 Building for Production

```bash
# Build backend executable
cd backend
go build -ldflags="-s -w" -o posflow-server

# Optional: Compress with UPX
upx --best posflow-server

# Build Electron app with electron-builder
cd ../electron-app
npm run build
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with ❤️ for small and medium businesses

---

**Note**: POSFlow is in active development. Features and documentation are being added regularly.
```

## GitHub About Section (Short Description)
```
Modern offline-first POS system for restaurants and retail. Inventory, orders, and sales management with multi-role interfaces.
```

## Alternative Tagline Options
```
"Streamline your business with offline-first POS"
"Point of Sale reimagined for modern businesses"
"Free, fast, and offline POS for restaurants & retail"
"Your business, your data, offline and in control"