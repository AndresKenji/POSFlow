const { contextBridge, ipcRenderer } = require('electron');

// Expose protected methods that allow the renderer process to use
// the ipcRenderer without exposing the entire object
contextBridge.exposeInMainWorld('electron', {
  // Window controls
  openKitchenView: () => ipcRenderer.send('open-kitchen-view'),
  openCustomerView: () => ipcRenderer.send('open-customer-view'),

  // API info
  apiUrl: 'http://localhost:8000/api'
});

// Expose API for frontend to use
contextBridge.exposeInMainWorld('api', {
  baseUrl: 'http://localhost:8000/api'
});