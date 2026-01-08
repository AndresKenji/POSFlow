/**
 * API Client for POSFlow backend (Go server)
 * Handles all communication between Electron frontend and Go backend
 */
class APIClient {
  constructor() {
    // Get API URL from preload script or use default
    this.baseURL = window.api?.baseUrl || 'http://localhost:8000/api/v1';
    console.log('🔧 APIClient initialized with baseURL:', this.baseURL);
    this.timeout = 10000; // 10 seconds timeout
    this.retryAttempts = 3;
    this.retryDelay = 1000; // 1 second
  }

  /**
   * Makes an HTTP request to the backend
   * @param {string} endpoint - API endpoint (e.g., '/orders')
   * @param {Object} options - Fetch options
   * @param {number} attempt - Current retry attempt
   * @returns {Promise<Object>} Response data
   */
  async request(endpoint, options = {}, attempt = 1) {
    const url = `${this.baseURL}${endpoint}`;
    console.log('Making request to:', url);
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers
        },
        signal: controller.signal,
        ...options
      });

      clearTimeout(timeoutId);

      // Parse response
      const data = await response.json();

      // Handle standard response format from Go backend
      if (!response.ok) {
        throw new Error(data.message || `HTTP error! status: ${response.status}`);
      }

      return data;
    } catch (error) {
      clearTimeout(timeoutId);

      // Retry on network errors
      if (attempt < this.retryAttempts && this.shouldRetry(error)) {
        console.warn(`Request failed (attempt ${attempt}/${this.retryAttempts}), retrying...`);
        await this.delay(this.retryDelay);
        return this.request(endpoint, options, attempt + 1);
      }

      console.error('API request failed:', error);
      throw this.handleError(error);
    }
  }

  /**
   * Determines if a request should be retried
   * @param {Error} error - The error that occurred
   * @returns {boolean}
   */
  shouldRetry(error) {
    return (
      error.name === 'AbortError' ||
      error.message.includes('Failed to fetch') ||
      error.message.includes('NetworkError')
    );
  }

  /**
   * Delays execution for the specified time
   * @param {number} ms - Milliseconds to delay
   * @returns {Promise<void>}
   */
  delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * Handles and formats errors
   * @param {Error} error - The error to handle
   * @returns {Error}
   */
  handleError(error) {
    if (error.name === 'AbortError') {
      return new Error('Request timeout. Please check your connection and try again.');
    }
    if (error.message.includes('Failed to fetch')) {
      return new Error('Cannot connect to server. Please ensure the backend is running.');
    }
    return error;
  }

  // ==================== Health Check ====================

  /**
   * Checks if the backend server is running
   * @returns {Promise<boolean>}
   */
  async healthCheck() {
    try {
      const response = await fetch('http://localhost:8000/health');
      return response.success === true;
    } catch (error) {
      return false;
    }
  }

  // ==================== Orders ====================

  /**
   * Gets all orders
   * @param {string} status - Optional status filter (pending, preparing, ready, completed)
   * @returns {Promise<Object>}
   */
  async getOrders(status = null) {
    const query = status ? `?status=${status}` : '';
    return this.request(`/orders${query}`);
  }

  /**
   * Gets a specific order by ID
   * @param {number} orderId - Order ID
   * @returns {Promise<Object>}
   */
  async getOrder(orderId) {
    return this.request(`/orders/${orderId}`);
  }

  /**
   * Creates a new order
   * @param {Object} orderData - Order information
   * @returns {Promise<Object>}
   */
  async createOrder(orderData) {
    return this.request('/orders', {
      method: 'POST',
      body: JSON.stringify(orderData)
    });
  }

  /**
   * Updates order status
   * @param {number} orderId - Order ID
   * @param {string} status - New status (pending, preparing, ready, completed)
   * @returns {Promise<Object>}
   */
  async updateOrderStatus(orderId, status) {
    return this.request(`/orders/${orderId}/status`, {
      method: 'PUT',
      body: JSON.stringify({ status })
    });
  }

  /**
   * Deletes an order
   * @param {number} orderId - Order ID
   * @returns {Promise<Object>}
   */
  async deleteOrder(orderId) {
    return this.request(`/orders/${orderId}`, {
      method: 'DELETE'
    });
  }

  // ==================== Products ====================

  /**
   * Gets all products
   * @param {boolean} available - Filter by availability
   * @returns {Promise<Object>}
   */
  async getProducts(available = null) {
    const query = available !== null ? `?available=${available}` : '';
    return this.request(`/products${query}`);
  }

  /**
   * Gets a specific product by ID
   * @param {number} productId - Product ID
   * @returns {Promise<Object>}
   */
  async getProduct(productId) {
    return this.request(`/products/${productId}`);
  }

  /**
   * Creates a new product
   * @param {Object} productData - Product information
   * @returns {Promise<Object>}
   */
  async createProduct(productData) {
    return this.request('/products', {
      method: 'POST',
      body: JSON.stringify(productData)
    });
  }

  /**
   * Updates a product
   * @param {number} productId - Product ID
   * @param {Object} productData - Updated product information
   * @returns {Promise<Object>}
   */
  async updateProduct(productId, productData) {
    return this.request(`/products/${productId}`, {
      method: 'PUT',
      body: JSON.stringify(productData)
    });
  }

  /**
   * Deletes a product
   * @param {number} productId - Product ID
   * @returns {Promise<Object>}
   */
  async deleteProduct(productId) {
    return this.request(`/products/${productId}`, {
      method: 'DELETE'
    });
  }

  // ==================== Inventory ====================

  /**
   * Gets current inventory
   * @returns {Promise<Object>}
   */
  async getInventory() {
    return this.request('/inventory');
  }

  /**
   * Updates inventory stock
   * @param {number} productId - Product ID
   * @param {number} quantity - Quantity change (positive or negative)
   * @param {string} reason - Reason for change
   * @returns {Promise<Object>}
   */
  async updateInventory(productId, quantity, reason) {
    return this.request('/inventory', {
      method: 'POST',
      body: JSON.stringify({ productId, quantity, reason })
    });
  }

  // ==================== Sales ====================

  /**
   * Gets daily sales report
   * @param {string} date - Date in YYYY-MM-DD format (optional, defaults to today)
   * @returns {Promise<Object>}
   */
  async getDailySales(date = null) {
    const query = date ? `?date=${date}` : '';
    return this.request(`/sales/daily${query}`);
  }

  /**
   * Gets sales report for a date range
   * @param {string} startDate - Start date in YYYY-MM-DD format
   * @param {string} endDate - End date in YYYY-MM-DD format
   * @returns {Promise<Object>}
   */
  async getSalesReport(startDate, endDate) {
    return this.request(`/sales/report?start=${startDate}&end=${endDate}`);
  }

  /**
   * Closes the day and generates report
   * @returns {Promise<Object>}
   */
  async closeDay() {
    return this.request('/sales/close-day', {
      method: 'POST'
    });
  }

  // ==================== Users/Authentication ====================

  /**
   * Logs in a user
   * @param {string} username - Username
   * @param {string} password - Password
   * @returns {Promise<Object>}
   */
  async login(username, password) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password })
    });
  }

  /**
   * Logs out current user
   * @returns {Promise<Object>}
   */
  async logout() {
    return this.request('/auth/logout', {
      method: 'POST'
    });
  }

  /**
   * Gets current user info
   * @returns {Promise<Object>}
   */
  async getCurrentUser() {
    return this.request('/auth/me');
  }

  // ==================== Menu ====================

  /**
   * Gets today's menu
   * @returns {Promise<Object>}
   */
  async getMenu() {
    return this.request('/menu');
  }

  /**
   * Updates menu item availability
   * @param {number} itemId - Menu item ID
   * @param {boolean} available - Availability status
   * @returns {Promise<Object>}
   */
  async updateMenuAvailability(itemId, available) {
    return this.request(`/menu/${itemId}/availability`, {
      method: 'PUT',
      body: JSON.stringify({ available })
    });
  }
}

// Create global instance
const apiClient = new APIClient();