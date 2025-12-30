// Wait for DOM to load
document.addEventListener('DOMContentLoaded', () => {
  initializeAdmin();
});

function initializeAdmin() {
  // Button event listeners
  document.getElementById('openKitchen').addEventListener('click', () => {
    window.electron.openKitchenView();
  });

  document.getElementById('openCustomer').addEventListener('click', () => {
    window.electron.openCustomerView();
  });

  // Load initial data
  loadDashboardData();
}

async function loadDashboardData() {
  try {
    const sales = await apiClient.getDailySales();

    // Update stats (mock data for now)
    document.getElementById('todaySales').textContent = '$0.00';
    document.getElementById('todayOrders').textContent = '0';
    document.getElementById('totalProducts').textContent = '0';

    console.log('Dashboard data loaded');
  } catch (error) {
    console.error('Failed to load dashboard data:', error);
    // Show error message to user
    alert('Failed to connect to backend. Make sure the Python API is running.');
  }
}