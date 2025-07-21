// content.js
(function() {
    chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
        if (request.action === 'getAuthToken') {
            // Parse URL parameters from the current page
            const params = new URLSearchParams(window.location.search);
            const token = params.get('access_token');
            const email = params.get('email');
            sendResponse({ token, email });
        }
    });
})();

// Listen for auth completion by monitoring URL changes
let lastUrl = location.href;
new MutationObserver(() => {
    const url = location.href;
    if (url !== lastUrl) {
        lastUrl = url;
        
        // If we're on the main site after auth, try to extract token
        if (url.includes('fluently-app.ru') && !url.includes('/auth/')) {
            setTimeout(() => {
                // Give the page time to load and set localStorage
                const token = localStorage.getItem('accessToken') || 
                             localStorage.getItem('token');
                
                if (token) {
                    // Notify the popup that we have a token
                    chrome.runtime.sendMessage({
                        action: 'authCompleted',
                        token: token,
                        email: localStorage.getItem('userEmail') || 
                               localStorage.getItem('email') || 'User'
                    });
                }
            }, 1000);
        }
    }
}).observe(document, { subtree: true, childList: true }); 