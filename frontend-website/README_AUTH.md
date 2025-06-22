# Google OAuth Authentication - Frontend Implementation

## Overview

This implementation provides Google OAuth authentication for the Fluently web application. Users can sign in with their Google accounts and access their profile information.

## Features

- 🔐 Google OAuth 2.0 authentication
- 👤 User profile page with Google account information
- 🔄 Persistent login state with localStorage
- 📱 Responsive design for mobile and desktop
- 🎨 Modern UI with smooth animations
- 🚪 Secure logout functionality

## Files

### Frontend Files
- `index.html` - Main landing page with Google login buttons
- `profile.html` - User profile page displaying Google account info
- `script.js` - Main page JavaScript with OAuth initiation
- `profile.js` - Profile page JavaScript for user data handling
- `style.css` - Enhanced CSS with profile page styles

### Backend Changes
- `backend/internal/api/v1/handlers/auth_handler.go` - Modified to redirect to frontend with user data

## How It Works

### Authentication Flow

1. **User clicks "Login with Google"** on the main page
2. **Frontend redirects** to `/auth/google` endpoint
3. **Backend redirects** to Google OAuth consent screen
4. **User authorizes** the application on Google
5. **Google redirects back** to `/auth/google/callback`
6. **Backend processes** the OAuth response and creates user/preferences
7. **Backend redirects** to `/profile.html` with user data in URL parameters
8. **Frontend displays** user profile with Google account information

### Data Flow

```
Frontend (index.html) → Backend (/auth/google) → Google OAuth → 
Backend (/auth/google/callback) → Frontend (profile.html) → User Profile
```

## Testing

### Prerequisites

1. Make sure your backend is running with Google OAuth configured
2. Ensure you have valid Google OAuth credentials in your environment variables:
   - `WEB_GOOGLE_CLIENT_ID`
   - `WEB_GOOGLE_CLIENT_SECRET`

### Steps to Test

1. **Start the backend server**:
   ```bash
   cd backend
   go run cmd/main.go
   ```

2. **Serve the frontend** (you can use any static file server):
   ```bash
   # Using Python
   cd frontend-website
   python -m http.server 8080
   
   # Or using Node.js (if you have http-server installed)
   npx http-server -p 8080
   
   # Or using Go
   go run -m http.fileserver -addr :8080 .
   ```

3. **Open your browser** and navigate to `http://localhost:8080`

4. **Click "Login with Google"** button

5. **Complete Google OAuth flow**:
   - You'll be redirected to Google's consent screen
   - Sign in with your Google account
   - Grant permissions to the application

6. **Verify profile page**:
   - You should be redirected to the profile page
   - Your Google profile picture, name, and email should be displayed
   - The access token should be stored in localStorage

7. **Test persistence**:
   - Refresh the page - you should remain logged in
   - Close and reopen the browser - login state should persist
   - Click logout to clear the session

## Features Explained

### Profile Page Components

- **User Avatar**: Displays Google profile picture
- **User Info**: Shows name and email from Google account
- **Welcome Message**: Personalized greeting for new users
- **Learning Dashboard**: Placeholder for future learning statistics
- **Quick Actions**: Buttons for future app features

### Security Features

- Access tokens are stored securely in localStorage
- URL parameters are cleaned after processing
- Logout clears all stored authentication data
- HTTPS recommended for production

### Mobile Responsive

- Responsive design works on all device sizes
- Touch-friendly buttons and navigation
- Mobile menu for smaller screens

## Customization

### Styling
You can customize the appearance by modifying the CSS variables in `style.css`:

```css
:root {
    --primary-accent: #A043DB;  /* Main brand color */
    --secondary-accent: #CC8EF1; /* Secondary brand color */
    --cta-blue: #119CE2;        /* Button color */
    /* ... other variables */
}
```

### Functionality
Add new features by:

1. Modifying the profile page HTML
2. Adding new JavaScript functions in `profile.js`
3. Making authenticated API calls using the `makeAuthenticatedRequest()` function

## Production Deployment

### Security Considerations

1. **Use HTTPS** in production
2. **Configure CORS** properly in your backend
3. **Set secure cookie flags** for state management
4. **Validate redirect URLs** to prevent open redirect attacks
5. **Implement proper error handling** for failed authentication

### Environment Variables

Make sure these are set in your production environment:
- `WEB_GOOGLE_CLIENT_ID`
- `WEB_GOOGLE_CLIENT_SECRET`
- `JWT_SECRET`
- Proper database connection settings

## Troubleshooting

### Common Issues

1. **"Invalid client" error**: Check your Google OAuth client ID
2. **Redirect URI mismatch**: Ensure your OAuth redirect URI matches the configured one
3. **CORS errors**: Configure CORS in your backend for your frontend domain
4. **Profile page not loading**: Check browser console for JavaScript errors

### Debug Mode

Enable debug logging by opening browser console and checking for:
- Authentication flow logs
- API request/response logs
- Error messages

## Next Steps

This implementation provides the foundation for:
- User preferences management
- Learning progress tracking
- API authenticated requests
- Additional social login providers 