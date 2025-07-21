// Content script to handle authentication token extraction
chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
    if (request.action === 'getAuthToken') {
        // Try to extract token from localStorage
        try {
            const token = localStorage.getItem('accessToken') || 
                         localStorage.getItem('token') ||
                         sessionStorage.getItem('accessToken') ||
                         sessionStorage.getItem('token');
            
            // Try to get user email from various sources
            const email = localStorage.getItem('userEmail') ||
                         localStorage.getItem('email') ||
                         sessionStorage.getItem('userEmail') ||
                         sessionStorage.getItem('email');
            
            if (token) {
                sendResponse({ 
                    token: token, 
                    email: email || 'User' 
                });
            } else {
                // If no token in storage, try to extract from page content
                // This is a fallback method - you might need to adjust based on your auth page structure
                const tokenElement = document.querySelector('[data-token], .token, #token');
                const emailElement = document.querySelector('[data-email], .email, #email');
                
                if (tokenElement) {
                    sendResponse({ 
                        token: tokenElement.textContent || tokenElement.value, 
                        email: emailElement ? (emailElement.textContent || emailElement.value) : 'User' 
                    });
                } else {
                    sendResponse({ error: 'No token found' });
                }
            }
        } catch (error) {
            console.error('Error extracting auth token:', error);
            sendResponse({ error: 'Failed to extract token' });
        }
        return true; // Keep the message channel open for async response
    }
});

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