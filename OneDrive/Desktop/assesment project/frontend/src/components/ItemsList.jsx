import React, { useState, useEffect } from 'react'
import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080'

function ItemsList({ token, onLogout }) {
  const [items, setItems] = useState([])
  const [loading, setLoading] = useState(true)
  const [cartId, setCartId] = useState(null)

  useEffect(() => {
    fetchItems()
    fetchUserCart()
  }, [])

  const fetchItems = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/items`)
      setItems(response.data)
    } catch (error) {
      console.error('Error fetching items:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchUserCart = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/carts/me`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      if (response.data.id) {
        setCartId(response.data.id)
      }
    } catch (error) {
      // Cart might not exist yet, that's okay
      console.log('No cart found or error:', error)
    }
  }

  const addToCart = async (itemId) => {
    try {
      const response = await axios.post(
        `${API_BASE_URL}/carts`,
        { item_ids: [itemId] },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      )
      if (response.data.id) {
        setCartId(response.data.id)
        window.alert('Item added to cart!')
      }
    } catch (error) {
      console.error('Error adding to cart:', error)
      window.alert('Failed to add item to cart')
    }
  }

  const handleCheckout = async () => {
    if (!cartId) {
      window.alert('No items in cart')
      return
    }

    try {
      const response = await axios.post(
        `${API_BASE_URL}/orders`,
        { cart_id: cartId },
        {
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      )
      window.alert('Order successful!')
      setCartId(null)
      // Optionally refresh items
      fetchItems()
    } catch (error) {
      console.error('Error creating order:', error)
      window.alert('Failed to create order: ' + (error.response?.data?.error || 'Unknown error'))
    }
  }

  const handleViewCart = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/carts/me`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      
      if (response.data && response.data.cart_items && response.data.cart_items.length > 0) {
        const itemsList = response.data.cart_items
          .map((ci) => ci.item?.name || `Item ID: ${ci.item_id}`)
          .join('\n')
        window.alert(`Cart Items:\n${itemsList}`)
      } else {
        window.alert('Cart is empty')
      }
    } catch (error) {
      console.error('Error fetching cart:', error)
      window.alert('Failed to fetch cart')
    }
  }

  const handleOrderHistory = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/orders`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      
      const orders = response.data || []
      if (orders.length > 0) {
        const orderIds = orders.map((o) => `Order ID: ${o.id}`).join('\n')
        window.alert(`Order History:\n${orderIds}`)
      } else {
        window.alert('No orders found')
      }
    } catch (error) {
      console.error('Error fetching orders:', error)
      window.alert('Failed to fetch order history')
    }
  }

  if (loading) {
    return <div className="container">Loading items...</div>
  }

  return (
    <div className="container">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h1>Shopping Cart - Items</h1>
        <button className="btn btn-secondary" onClick={onLogout}>
          Logout
        </button>
      </div>

      <div className="header-buttons">
        <button className="btn btn-success" onClick={handleCheckout}>
          Checkout
        </button>
        <button className="btn btn-primary" onClick={handleViewCart}>
          Cart
        </button>
        <button className="btn btn-primary" onClick={handleOrderHistory}>
          Order History
        </button>
      </div>

      <div>
        <h2>Available Items</h2>
        {items.length === 0 ? (
          <p>No items available</p>
        ) : (
          items.map((item) => (
            <div key={item.id} className="item-card">
              <div className="item-info">
                <div className="item-name">{item.name}</div>
                <div className="item-status">Status: {item.status}</div>
              </div>
              <button
                className="btn btn-primary"
                onClick={() => addToCart(item.id)}
              >
                Add to Cart
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  )
}

export default ItemsList

