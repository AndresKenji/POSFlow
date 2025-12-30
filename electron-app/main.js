const { app, BrowserWindow, ipcMain } = require('electron');
const path = require('path');

// Keep references to windows to prevent garbage collection
let mainWindow;
let kitchenWindow;
let customerWindow;

function createMainWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: false,
      contextIsolation: true
    }
  });

  mainWindow.loadFile('src/views/admin.html');

  // Open DevTools in development
  if (process.env.NODE_ENV === 'development') {
    mainWindow.webContents.openDevTools();
  }

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
}

function createKitchenWindow() {
  kitchenWindow = new BrowserWindow({
    width: 1024,
    height: 768,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: false,
      contextIsolation: true
    }
  });

  kitchenWindow.loadFile('src/views/kitchen.html');

  kitchenWindow.on('closed', () => {
    kitchenWindow = null;
  });
}

function createCustomerWindow() {
  customerWindow = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      nodeIntegration: false,
      contextIsolation: true
    }
  });

  customerWindow.loadFile('src/views/customer.html');

  customerWindow.on('closed', () => {
    customerWindow = null;
  });
}

// App lifecycle
app.whenReady().then(() => {
  createMainWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createMainWindow();
    }
  });
});

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

// IPC handlers for opening different windows
ipcMain.on('open-kitchen-view', () => {
  if (!kitchenWindow) {
    createKitchenWindow();
  } else {
    kitchenWindow.focus();
  }
});

ipcMain.on('open-customer-view', () => {
  if (!customerWindow) {
    createCustomerWindow();
  } else {
    customerWindow.focus();
  }
});