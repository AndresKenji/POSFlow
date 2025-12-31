const { app, BrowserWindow, ipcMain } = require('electron');
const path = require('path');
const { spawn } = require('child_process');

// Keep references to windows to prevent garbage collection
let mainWindow;
let kitchenWindow;
let customerWindow;
let backendProcess;

// Backend server configuration
const BACKEND_PORT = 8080;
const BACKEND_STARTUP_DELAY = 2000; // 2 seconds

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

/**
 * Starts the Go backend server
 * In development, assumes the backend is started manually
 * In production, spawns the backend executable
 */
function startBackend() {
  // Skip starting backend in development (assume it's running separately)
  if (process.env.NODE_ENV === 'development') {
    console.log('Development mode: Assuming backend is running on port', BACKEND_PORT);
    return;
  }

  // Determine backend executable path and name
  const isWin = process.platform === 'win32';
  const backendName = isWin ? 'posflow-server.exe' : 'posflow-server';

  // In production, the backend is in the resources folder
  const backendPath = path.join(process.resourcesPath, backendName);

  console.log('Starting backend from:', backendPath);

  // Spawn the backend process
  backendProcess = spawn(backendPath, [], {
    env: {
      ...process.env,
      PORT: BACKEND_PORT.toString(),
    },
  });

  // Log backend output
  backendProcess.stdout.on('data', (data) => {
    console.log('[Backend]', data.toString().trim());
  });

  backendProcess.stderr.on('data', (data) => {
    console.error('[Backend Error]', data.toString().trim());
  });

  backendProcess.on('close', (code) => {
    console.log(`Backend process exited with code ${code}`);
    backendProcess = null;
  });

  backendProcess.on('error', (err) => {
    console.error('Failed to start backend:', err);
  });
}

/**
 * Stops the backend server
 */
function stopBackend() {
  if (backendProcess) {
    console.log('Stopping backend server...');
    backendProcess.kill();
    backendProcess = null;
  }
}

/**
 * Checks if backend is ready by pinging health endpoint
 */
async function waitForBackend(maxAttempts = 10) {
  for (let i = 0; i < maxAttempts; i++) {
    try {
      const response = await fetch(`http://localhost:${BACKEND_PORT}/api/health`);
      if (response.ok) {
        console.log('Backend is ready!');
        return true;
      }
    } catch (err) {
      console.log(`Waiting for backend... (attempt ${i + 1}/${maxAttempts})`);
    }
    await new Promise(resolve => setTimeout(resolve, 1000));
  }
  console.error('Backend failed to start in time');
  return false;
}

// App lifecycle
app.whenReady().then(async () => {
  // Start backend server
  startBackend();

  // Wait for backend to be ready
  await waitForBackend();

  // Create main window
  createMainWindow();

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createMainWindow();
    }
  });
});

app.on('window-all-closed', () => {
  stopBackend();
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', () => {
  stopBackend();
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