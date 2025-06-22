# Automated Swagger UI Authentication

This custom Swagger UI implementation provides automated authentication for the Fluently API, eliminating the need to manually copy and paste Bearer tokens.

## Features

🔐 **Easy Login**: Simply enter your email and password
🚀 **Auto Token Management**: Bearer tokens are automatically applied to all API requests
💾 **Session Persistence**: Your login session persists across page refreshes
🔄 **Token Validation**: Automatically checks token expiration and manages logout
📱 **Responsive Design**: Works on desktop and mobile devices

## How to Use

### 1. Access the Authentication Page

Navigate to: `https://your-domain.com/swagger-auth`
Or: `https://your-domain.com/swagger` (redirects to auth page)

### 2. Login

1. Enter your **email** and **password** in the authentication form at the top
2. Click the **"Login"** button
3. Upon successful login:
   - You'll see a success message
   - The login form will be hidden
   - Your user info and token expiration will be displayed
   - The Bearer token is automatically applied to all Swagger requests

### 3. Use the API

- All API endpoints now automatically include the Bearer token
- You can test any protected endpoint without manually setting authorization headers
- The generated curl commands will still need the "Bearer " prefix added manually (this is a Swagger UI limitation)

### 4. Logout

- Click the **"Logout"** button to clear your session
- This will remove the stored token and reset the authentication form

## Technical Details

### Session Management
- Tokens are stored in `sessionStorage` (cleared when browser tab is closed)
- Automatic token expiration checking
- Seamless session restoration on page refresh

### Security Features
- No tokens stored in localStorage (better security)
- Automatic logout on token expiration
- CSRF protection through proper headers

### API Integration
- Uses your existing `/auth/login` endpoint
- Compatible with your current JWT token format
- Works with all existing API endpoints

## Troubleshooting

### Login Issues
- **"Login failed"**: Check your email and password
- **"Network error"**: Ensure the backend server is running
- **Invalid credentials**: Verify your account exists and password is correct

### Token Issues
- **Requests return 401**: Your token may have expired, try logging in again
- **"Token expired"**: The system will automatically prompt you to login again

### Fallback Options
- Original Swagger UI available at: `/swagger/index.html`
- Manual token entry still works in the original UI
- API endpoints can still be tested with tools like curl or Postman

## Development Notes

### File Structure
```
backend/
├── static/
│   └── swagger-auth.html    # Custom authentication UI
├── internal/router/
│   └── router.go           # Routes for serving the auth UI
└── docs/                   # Auto-generated Swagger docs
    ├── swagger.json
    └── swagger.yaml
```

### Customization
You can modify `static/swagger-auth.html` to:
- Change the styling/theme
- Add additional authentication methods (Google OAuth, etc.)
- Customize the user interface
- Add more user information display

### Adding New Features
The JavaScript in the HTML file can be extended to:
- Support refresh token rotation
- Add remember me functionality
- Integrate with SSO providers
- Add user role-based UI changes 