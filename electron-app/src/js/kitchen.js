document.addEventListener('DOMContentLoaded', () => {
  initializeKitchen();
});

function initializeKitchen() {
  loadOrders();

  // Refresh orders every 5 seconds
  setInterval(loadOrders, 5000);
}

async function loadOrders() {
  try {
    const response = await apiClient.getOrders();
    displayOrders(response.data || []);
  } catch (error) {
    console.error('Failed to load orders:', error);
  }
}

function displayOrders(orders) {
  const ordersList = document.getElementById('ordersList');
  const orderCount = document.getElementById('orderCount');

  const pendingOrders = orders.filter(o => o.status === 'pending' || o.status === 'preparing');
  orderCount.textContent = `${pendingOrders.length} pending orders`;

  if (orders.length === 0) {
    ordersList.innerHTML = '<p>No pending orders</p>';
    return;
  }

  ordersList.innerHTML = orders.map(order => `
    <div class="order-card ${order.status}">
      <h3>Order #${order.id}</h3>
      <p><strong>Table:</strong> ${order.table_number}</p>
      <p><strong>Status:</strong> ${order.status}</p>
      <p><strong>Total:</strong> $${order.total}</p>
      <button onclick="updateStatus(${order.id}, 'preparing')">Start Preparing</button>
      <button onclick="updateStatus(${order.id}, 'ready')">Mark Ready</button>
    </div>
  `).join('');
}

async function updateStatus(orderId, status) {
  try {
    await apiClient.updateOrderStatus(orderId, status);
    loadOrders();
  } catch (error) {
    console.error('Failed to update order:', error);
    alert('Failed to update order status');
  }
}