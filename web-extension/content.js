// content.js
(function() {
    console.log('Fluently extension content script loaded on:', window.location.href);
    
    // Listen for messages from the popup
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
        if (request.action === 'getAuthToken') {
            console.log('Getting auth token from page...');
            
            // First try URL parameters
            const params = new URLSearchParams(window.location.search);
            let token = params.get('access_token');
            let email = params.get('email');
            
            // If not in URL, try localStorage
            if (!token) {
                const fluentlyUser = localStorage.getItem('fluently_user');
                const fluentlyToken = localStorage.getItem('fluently_access_token');
                
                if (fluentlyUser) {
                    try {
                        const userData = JSON.parse(fluentlyUser);
                        token = userData.accessToken;
                        email = userData.email;
                    } catch (e) {
                        console.error('Error parsing user data:', e);
                    }
                }
                
                if (!token && fluentlyToken) {
                    token = fluentlyToken;
                }
            }
            
            console.log('Found token:', !!token, 'email:', email);
            sendResponse({ token, email });
        }
    });
    
    // Listen for messages from the auth-success page
    window.addEventListener('message', function(event) {
        if (event.data && event.data.type === 'AUTH_SUCCESS') {
            console.log('Received AUTH_SUCCESS message');
            
            // Wait a bit for localStorage to be set
            setTimeout(() => {
                const fluentlyUser = localStorage.getItem('fluently_user');
                const fluentlyToken = localStorage.getItem('fluently_access_token');
                
                let token = null;
                let email = null;
                
                if (fluentlyUser) {
                    try {
                        const userData = JSON.parse(fluentlyUser);
                        token = userData.accessToken;
                        email = userData.email;
                    } catch (e) {
                        console.error('Error parsing user data:', e);
                    }
                }
                
                if (!token && fluentlyToken) {
                    token = fluentlyToken;
                }
                
                if (token) {
                    console.log('Sending auth completion to background');
                    chrome.runtime.sendMessage({
                        action: 'authCompleted',
                        token: token,
                        email: email || 'User'
                    });
                }
            }, 500);
        }
    });
})();

// Monitor for auth success page load
if (window.location.href.includes('auth-success.html')) {
    console.log('On auth success page, checking for auth data...');
    
    // Check immediately
    setTimeout(() => {
        const fluentlyUser = localStorage.getItem('fluently_user');
        const fluentlyToken = localStorage.getItem('fluently_access_token');
        
        if (fluentlyUser || fluentlyToken) {
            let token = null;
            let email = null;
            
            if (fluentlyUser) {
                try {
                    const userData = JSON.parse(fluentlyUser);
                    token = userData.accessToken;
                    email = userData.email;
                } catch (e) {
                    console.error('Error parsing user data:', e);
                }
            }
            
            if (!token && fluentlyToken) {
                token = fluentlyToken;
            }
            
            if (token) {
                console.log('Auto-sending auth completion to background');
                chrome.runtime.sendMessage({
                    action: 'authCompleted',
                    token: token,
                    email: email || 'User'
                });
            }
        }
    }, 1000);
} 