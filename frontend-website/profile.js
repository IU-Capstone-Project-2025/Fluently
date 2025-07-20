// Profile page functionality
document.addEventListener('DOMContentLoaded', () => {
    // Parse URL parameters to get user data
    const urlParams = new URLSearchParams(window.location.search);
    const name = urlParams.get('name');
    const email = urlParams.get('email');
    const picture = urlParams.get('picture');
    const accessToken = urlParams.get('access_token');

    // If no user data in URL, check localStorage
    if (!name || !email) {
        const storedUser = localStorage.getItem('fluently_user');
        if (storedUser) {
            const userData = JSON.parse(storedUser);
            displayUserProfile(userData);
        } else {
            // Redirect to home if no authentication data
            window.location.href = '/index.html';
        }
    } else {
        // Store user data and access token
        const userData = {
            name,
            email,
            picture,
            accessToken,
            loginTime: new Date().toISOString()
        };
        
        localStorage.setItem('fluently_user', JSON.stringify(userData));
        localStorage.setItem('fluently_access_token', accessToken);
        
        // Display user profile
        displayUserProfile(userData);
        
        // Clean URL parameters
        window.history.replaceState({}, document.title, '/profile.html');
    }

    // Set up logout functionality
    document.getElementById('logoutBtn').addEventListener('click', logout);
});

function displayUserProfile(userData) {
    // Update profile elements
    document.getElementById('userName').textContent = userData.name || 'Anonymous User';
    document.getElementById('userEmail').textContent = userData.email || 'No email';
    
    if (userData.picture) {
        document.getElementById('userAvatar').src = userData.picture;
        document.getElementById('userAvatar').style.display = 'block';
    } else {
        // Show default avatar
        document.getElementById('userAvatar').src = './logo.jpg';
    }

    // Add loading state removal
    document.querySelector('.profile-container').classList.add('loaded');
}

function logout() {
    // Clear local storage
    localStorage.removeItem('fluently_user');
    localStorage.removeItem('fluently_access_token');
    
    // Show logout message
    if (confirm('Are you sure you want to logout?')) {
        // Optional: Call backend logout endpoint
        // fetch('/auth/logout', { method: 'POST' });
        
        // Redirect to home page
        window.location.href = '/index.html';
    }
}

// Utility function to make authenticated API calls
function makeAuthenticatedRequest(url, options = {}) {
    const token = localStorage.getItem('fluently_access_token');
    
    const defaultOptions = {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
            ...options.headers
        }
    };
    
    return fetch(url, { ...options, ...defaultOptions });
}

// Example usage for future API calls
function loadUserProgress() {
    makeAuthenticatedRequest('/api/v1/users/progress')
        .then(response => response.json())
        .then(data => {
            // Update dashboard with real data
            console.log('User progress:', data);
        })
        .catch(error => {
            console.error('Error loading progress:', error);
        });
}

// Add smooth animations
document.addEventListener('DOMContentLoaded', () => {
    // Animate stats cards
    const statCards = document.querySelectorAll('.stat-card');
    statCards.forEach((card, index) => {
        setTimeout(() => {
            card.style.opacity = '1';
            card.style.transform = 'translateY(0)';
        }, index * 100);
    });
    
    // Set up action button handlers
    document.querySelectorAll('.action-btn').forEach(btn => {
        btn.addEventListener('click', (e) => {
            const icon = e.target.closest('.action-btn').querySelector('.icon').textContent;
            
            switch(icon) {
                case '📚':
                    alert('Learning feature coming soon!');
                    break;
                case '⚙️':
                    alert('Preferences feature coming soon!');
                    break;
                case '📊':
                    alert('Progress tracking feature coming soon!');
                    break;
            }
        });
    });
}); 