# Fluently Chrome Extension - Authentication Debug Guide

## Recent Fixes Applied

### Issues Identified:
1. **localStorage key mismatch**: Auth page stored `fluently_access_token` but extension looked for `accessToken`
2. **Content script communication**: Message passing between auth page and extension was broken
3. **Manifest permissions**: Missing required permissions for proper operation
4. **Authentication flow**: Extension didn't properly detect auth success

### Fixes Applied:

#### 1. Updated `content.js`:
- Added proper localStorage key detection for both `fluently_user` and `fluently_access_token`
- Added message listener for `AUTH_SUCCESS` from auth page
- Added automatic detection when on auth-success page
- Added comprehensive logging for debugging

#### 2. Updated `manifest.json`:
- Added `scripting` and `tabs` permissions
- Fixed content script injection to only run on auth-success pages
- Set `run_at: "document_end"` for proper timing

#### 3. Updated `background.js`:
- Added proper logging and error handling
- Improved tab closing logic to handle both auth and auth-success tabs
- Added async response handling

#### 4. Updated `popup.js`:
- Improved auth tab detection to specifically look for auth-success pages
- Added longer timeout for content script communication
- Added proper cleanup of event listeners

#### 5. Updated `auth-success.html`:
- Added automatic notification to extension on page load
- Improved window messaging with actual token data
- Added better logging for debugging

## Testing Steps:

### 1. Load the Extension:
1. Open Chrome and go to `chrome://extensions/`
2. Enable "Developer mode"
3. Click "Load unpacked" and select the `web-extension` folder
4. Note the extension ID for debugging

### 2. Test Authentication:
1. Click the extension icon to open popup
2. Click "Login with Google"
3. Complete Google authentication
4. Check if extension popup shows authenticated state

### 3. Debug with Console:
1. Open extension popup
2. Right-click and select "Inspect" to open DevTools
3. Check Console tab for logs
4. For background script logs: go to `chrome://extensions/`, click "service worker" link

### 4. Check Storage:
1. In extension DevTools, go to Application tab
2. Look at Local Storage for the extension
3. Should see `accessToken` and `userEmail` keys

## Common Issues and Solutions:

### Extension Doesn't Detect Authentication:
- Check if content script is injected: DevTools → Sources → Content Scripts
- Verify localStorage contains the correct keys
- Check console for error messages

### Authentication Page Doesn't Close:
- Verify extension has `tabs` permission
- Check background script console for tab closing errors

### Token Not Saved:
- Check if `chrome.storage.local` is working
- Verify extension has `storage` permission
- Check for JavaScript errors in popup

## Debug Commands:

### Check Extension Storage:
```javascript
// In extension popup console
chrome.storage.local.get(['accessToken', 'userEmail'], console.log);
```

### Check Page LocalStorage:
```javascript
// In auth-success page console
console.log('User data:', localStorage.getItem('fluently_user'));
console.log('Token:', localStorage.getItem('fluently_access_token'));
```

### Test Content Script Communication:
```javascript
// In auth-success page console
window.postMessage({ type: 'AUTH_SUCCESS', token: 'test', email: 'test@example.com' }, '*');
```

## Expected Logs:

### Successful Flow:
1. Popup: "Starting login process..."
2. Content: "Fluently extension content script loaded on: https://fluently-app.ru/auth-success.html"
3. Content: "Received AUTH_SUCCESS message"
4. Background: "Background received message: {action: 'authCompleted', ...}"
5. Background: "Auth data saved to storage"
6. Popup: Shows authenticated state

If any of these logs are missing, that indicates where the issue is occurring.
