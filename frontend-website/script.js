// Basic interactive functionality
document.addEventListener('DOMContentLoaded', () => {
    // Mobile menu toggle
    const menuToggle = document.createElement('div');
    menuToggle.className = 'mobile-menu-toggle';
    menuToggle.innerHTML = '☰';
    
    document.querySelector('nav').appendChild(menuToggle);
    
    menuToggle.addEventListener('click', () => {
        document.querySelector('.nav-links').classList.toggle('active');
    });

    // Check if user is already logged in
    checkAuthStatus();

    // Google OAuth login handlers
    const loginBtn = document.getElementById('loginBtn');
    const mainLoginBtn = document.getElementById('mainLoginBtn');
    
    if (loginBtn) {
        loginBtn.addEventListener('click', initiateGoogleLogin);
    }
    
    if (mainLoginBtn) {
        mainLoginBtn.addEventListener('click', initiateGoogleLogin);
    }
});

function checkAuthStatus() {
    const userData = localStorage.getItem('fluently_user');
    if (userData) {
        // User is logged in, update UI
        const user = JSON.parse(userData);
        updateUIForLoggedInUser(user);
    }
}

function updateUIForLoggedInUser(user) {
    // Update login buttons to show user info
    const loginBtn = document.getElementById('loginBtn');
    const mainLoginBtn = document.getElementById('mainLoginBtn');
    
    if (loginBtn) {
        loginBtn.innerHTML = `
            <img src="${user.picture || './logo.jpg'}" alt="User Avatar" 
                 style="width: 20px; height: 20px; border-radius: 50%; margin-right: 8px;">
            ${user.name}
        `;
        loginBtn.onclick = () => window.location.href = '/auth-success.html';
    }
    
    if (mainLoginBtn) {
        mainLoginBtn.innerHTML = `
            <svg width="24" height="24" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
            </svg>
            Continue to Extension
        `;
        mainLoginBtn.onclick = () => window.location.href = '/auth-success.html';
    }
}

function initiateGoogleLogin() {
    // Show loading state
    const button = event.target.closest('button');
    const originalText = button.innerHTML;
    button.innerHTML = '🔄 Signing in...';
    button.disabled = true;
    
    // Get current domain for redirect URL
    const currentDomain = window.location.origin;
    
    // Redirect to backend Google OAuth endpoint
    window.location.href = `/auth/google`;
}

// Utility function to handle authentication errors
function handleAuthError(error) {
    console.error('Authentication error:', error);
    alert('Authentication failed. Please try again.');
    
    // Reset login buttons
    const loginBtn = document.getElementById('loginBtn');
    const mainLoginBtn = document.getElementById('mainLoginBtn');
    
    if (loginBtn) {
        loginBtn.innerHTML = `
            <svg width="20" height="20" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Login with Google
        `;
        loginBtn.disabled = false;
    }
    
    if (mainLoginBtn) {
        mainLoginBtn.innerHTML = `
            <svg width="24" height="24" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Get Started with Google
        `;
        mainLoginBtn.disabled = false;
    }
}