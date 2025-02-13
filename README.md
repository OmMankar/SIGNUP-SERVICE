# User Authentication Service

This project is a user authentication service built using Golang. It supports traditional email/password authentication with OTP-based email verification and Google OAuth-based authentication.

## Features
- User registration with email verification (OTP-based)
- Secure password hashing
- Google OAuth authentication
- JWT token-based authentication
- Caching OTPs using ttlcache
- MySQL database for storing user credentials

## Tech Stack
- **Golang**: Backend implementation
- **MySQL**: Database management
- **Gorilla Mux**: HTTP routing
- **TTLCache**: OTP caching mechanism
- **Google OAuth2**: Third-party authentication
- **JWT**: Token-based authentication

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/your-repo/auth-service.git
   cd auth-service
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Set up environment variables by creating a `.env` file:
   ```env
   DBUSER=your_db_user
   DBPASS=your_db_password
   DBNAME=your_db_name
   JWT_SECRET=your_jwt_secret
   CLIENT_ID=your_google_client_id
   CLIENT_SECRET=your_google_client_secret
   SENDER_EMAIL=your_smtp_email
   SENDER_PASSWORD=your_smtp_password
   ```

4. Run the application:
   ```sh
   go run main.go
   ```

## API Endpoints

### User Authentication

- **`POST /signup`** - Register a new user (triggers OTP email verification)
- **`POST /verify-email`** - Verify OTP and complete registration

### Google OAuth Authentication

- **`GET /google/login`** - Redirects user to Google login page
- **`GET /google/callback`** - Handles Google OAuth response, registers/logs in the user

### Sample JSON Requests

#### Register User
```json
{
  "name": "John Doe",
  "emailid": "john.doe@example.com",
  "password": "securepassword"
}
```

#### Verify Email (OTP)
```json
{
  "emailid": "john.doe@example.com",
  "otp": 123456
}
```

## Project Structure
```
.
├── config/           # Database and OAuth configuration
├── controller/       # Handlers for authentication
├── route/            # API route definitions
├── models/           # Database models
├── utiles/           # Helper functions (JWT, DB interactions, etc.)
├── main.go           # Application entry point
└── .env              # Environment variables (ignored in Git)
```

## Security Considerations
- Passwords are hashed before storage using **bcrypt**.
- JWT tokens are signed using a secret key stored in `.env`.
- OTPs are cached and expire after 10 minutes.
