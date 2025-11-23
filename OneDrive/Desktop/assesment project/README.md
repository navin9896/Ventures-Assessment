# Shopping Cart Application

A full-stack shopping cart application built with Go (Gin) backend and React frontend.

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Web Framework**: Gin (github.com/gin-gonic/gin)
- **ORM**: GORM (github.com/jinzhu/gorm)
- **Database**: PostgreSQL
- **Testing**: Ginkgo (github.com/onsi/ginkgo)

### Frontend
- **Framework**: React 18
- **Build Tool**: Vite
- **HTTP Client**: Axios
- **Routing**: React Router DOM

## Project Structure

```
.
├── backend/
│   ├── models/          # Database models (User, Item, Cart, CartItem, Order)
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # Authentication middleware
│   ├── database/        # Database connection and seeding
│   ├── main.go          # Application entry point
│   ├── cart_test.go     # Ginkgo tests
│   ├── go.mod           # Go dependencies
│   └── .env.example     # Environment variables template
├── frontend/
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── App.jsx      # Main app component
│   │   └── main.jsx     # React entry point
│   ├── package.json     # Node dependencies
│   └── vite.config.js   # Vite configuration
├── README.md            # This file
└── postman_collection.json  # Postman API collection
```

## Prerequisites

- Go 1.21 or higher
- Node.js 16+ and npm/yarn
- PostgreSQL 12+ (or Docker for running PostgreSQL)

## Backend Setup

### 1. Install Dependencies

```bash
cd backend
go mod download
```

### 2. Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE shopping_cart;
```

Or using Docker:

```bash
docker run --name shopping-cart-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=shopping_cart -p 5432:5432 -d postgres
```

### 3. Environment Variables

Create a `.env` file in the `backend` directory (or copy from `.env.example`):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shopping_cart
DB_SSLMODE=disable
```

### 4. Run the Backend

```bash
cd backend
go run main.go
```

The server will start on `http://localhost:8080`

### 5. Run Tests

```bash
cd backend
go test ./...
```

## Frontend Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Run the Frontend

```bash
npm run dev
```

The frontend will start on `http://localhost:3000`

### 3. Build for Production

```bash
npm run build
```

## API Endpoints

### Users

- `POST /users` - Create a new user
  ```json
  {
    "username": "john_doe",
    "password": "password123"
  }
  ```

- `GET /users` - List all users

- `POST /users/login` - Login user
  ```json
  {
    "username": "john_doe",
    "password": "password123"
  }
  ```
  Returns: `{ "token": "...", "user": {...} }`

### Items

- `POST /items` - Create a new item (requires authentication for production)
  ```json
  {
    "name": "Laptop",
    "status": "active"
  }
  ```

- `GET /items` - List all items

### Carts (Requires Authentication)

All cart endpoints require `Authorization: Bearer <token>` header.

- `POST /carts` - Create or update cart with items
  ```json
  {
    "item_ids": [1, 2, 3]
  }
  ```

- `GET /carts` - List all carts (optional query: `?user_id=1`)

- `GET /carts/me` - Get current user's cart

### Orders (Requires Authentication)

All order endpoints require `Authorization: Bearer <token>` header.

- `POST /orders` - Create an order from a cart
  ```json
  {
    "cart_id": 1
  }
  ```

- `GET /orders` - List all orders (optional query: `?user_id=1`)

## Postman Collection

A Postman collection is provided in `postman_collection.json`. Import it into Postman to test all API endpoints.

### Using the Postman Collection

1. Open Postman
2. Click "Import" button
3. Select `postman_collection.json` file
4. The collection will include all endpoints with example requests

**Important**: After logging in, copy the token from the response and set it as a collection variable:
1. Go to the collection variables
2. Set `token` variable with the value from login response
3. All authenticated requests will automatically use this token

## Frontend Features

### Login Screen
- Username and password login
- Redirects to items list on successful login
- Shows alert on invalid credentials

**Note**: To create a new user account, use the `POST /users` endpoint via Postman or API client before logging in through the frontend.

### Items List Screen
- Displays all available items
- "Add to Cart" button for each item
- Header buttons:
  - **Checkout**: Converts current cart to an order
  - **Cart**: Shows current cart items in an alert
  - **Order History**: Shows user's order IDs in an alert

## Database Schema

### Users
- `id` (primary key)
- `username` (unique)
- `password` (hashed)
- `token` (nullable, for session management)
- `cart_id` (nullable, FK to carts)
- `created_at`

### Items
- `id` (primary key)
- `name`
- `status`
- `created_at`

### Carts
- `id` (primary key)
- `user_id` (FK to users)
- `name`
- `status`
- `created_at`

### Cart Items
- `id` (primary key)
- `cart_id` (FK to carts)
- `item_id` (FK to items)

### Orders
- `id` (primary key)
- `cart_id` (FK to carts)
- `user_id` (FK to users)
- `created_at`

## Authentication

- Users log in with username and password
- On successful login, a token is generated and stored in the user's `token` field
- Only one active token per user (logging in again invalidates previous token)
- Protected endpoints require `Authorization: Bearer <token>` header
- Token is validated via middleware that injects user info into request context

## Sample Data

The database is automatically seeded with sample items on first run:
- Laptop
- Mouse
- Keyboard
- Monitor
- Headphones

## Development Notes

- CORS is enabled for all origins (adjust for production)
- Passwords are hashed using bcrypt
- Tokens are randomly generated hex strings
- Cart status is set to "checked_out" when converted to an order
- User's `cart_id` is cleared after checkout

## Troubleshooting

### Backend Issues

1. **Database connection error**: Ensure PostgreSQL is running and credentials in `.env` are correct
2. **Port already in use**: Change port in `main.go` or kill the process using port 8080
3. **Migration errors**: Drop and recreate the database

### Frontend Issues

1. **Cannot connect to backend**: Ensure backend is running on port 8080
2. **CORS errors**: Check that backend CORS middleware is enabled
3. **Token not working**: Ensure you're sending `Authorization: Bearer <token>` header

## License

This project is for assessment purposes.

