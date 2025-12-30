let cart = [];

document.addEventListener('DOMContentLoaded', () => {
  initializeCustomer();
});

function initializeCustomer() {
  loadMenu();

  document.getElementById('placeOrder').addEventListener('click', placeOrder);
}

async function loadMenu() {
  try {
    const response = await apiClient.getProducts();
    displayMenu(response.data || []);
  } catch (error) {
    console.error('Failed to load menu:', error);
    document.getElementById('menuItems').innerHTML =
      '<p>Failed to load menu. Make sure the backend is running.</p>';
  }
}

function displayMenu(products) {
  const menuItems = document.getElementById('menuItems');

  if (products.length === 0) {
    menuItems.innerHTML = '<p>No products available</p>';
    return;
  }

  menuItems.innerHTML = products.map(product => `
    <div class="menu-item" onclick="addToCart(${product.id}, '${product.name}', ${product.price})">
      <h4>${product.name}</h4>
      <p>$${product.price}</p>
    </div>
  `).join('');
}

function addToCart(id, name, price) {
  const existingItem = cart.find(item => item.id === id);

  if (existingItem) {
    existingItem.quantity++;
  } else {
    cart.push({ id, name, price, quantity: 1 });
  }

  updateCartDisplay();
}

function updateCartDisplay() {
  const cartItems = document.getElementById('cartItems');
  const cartTotal = document.getElementById('cartTotal');

  if (cart.length === 0) {
    cartItems.innerHTML = '<p>Cart is empty</p>';
    cartTotal.textContent = '0.00';
    return;
  }

  cartItems.innerHTML = cart.map((item, index) => `
    <div style="display: flex; justify-content: space-between; margin-bottom: 10px;">
      <span>${item.name} x${item.quantity}</span>
      <span>$${(item.price * item.quantity).toFixed(2)}</span>
      <button onclick="removeFromCart(${index})" class="btn-danger" style="padding: 5px 10px;">Remove</button>
    </div>
  `).join('');

  const total = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  cartTotal.textContent = total.toFixed(2);
}

function removeFromCart(index) {
  cart.splice(index, 1);
  updateCartDisplay();
}

async function placeOrder() {
  if (cart.length === 0) {
    alert('Cart is empty!');
    return;
  }

  const orderData = {
    table_number: 'Table 1', // You can add input for this
    items: cart,
    total: cart.reduce((sum, item) => sum + (item.price * item.quantity), 0)
  };

  try {
    await apiClient.createOrder(orderData);
    alert('Order placed successfully!');
    cart = [];
    updateCartDisplay();
  } catch (error) {
    console.error('Failed to place order:', error);
    alert('Failed to place order. Make sure the backend is running.');
  }
}