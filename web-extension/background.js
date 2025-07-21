// Background service worker
console.log('Background script loaded');

const AUTH_URL = 'https://fluently-app.ru/auth/google';
let authInProgress = false;

chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
    console.log('Background received message:', request);
    
    if (request.action === 'startAuth') {
        console.log('Starting authentication flow...');
        
        // Prevent multiple auth attempts
        if (authInProgress) {
            console.log('Auth already in progress, ignoring duplicate request');
            sendResponse({ success: false, error: 'Authentication already in progress' });
            return;
        }
        
        authInProgress = true;
        
        // Check if tabs API is available
        if (!chrome || !chrome.tabs) {
            console.error('Chrome tabs API not available in background!');
            authInProgress = false;
            sendResponse({ success: false, error: 'Tabs API not available' });
            return;
        }
        
        // Open auth page in a new tab
        chrome.tabs.create({ 
            url: AUTH_URL,
            active: true
        }, function(tab) {
            if (chrome.runtime.lastError) {
                console.error('Failed to create auth tab:', chrome.runtime.lastError);
                authInProgress = false;
                sendResponse({ success: false, error: chrome.runtime.lastError.message });
                return;
            }
            
            console.log('Auth tab created:', tab.id);
            
            // Set up a listener for when the tab is updated
            const listener = function(tabId, changeInfo, updatedTab) {
                if (tabId === tab.id && changeInfo.status === 'complete') {
                    console.log('Auth tab updated:', updatedTab.url);
                    
                    // Check if we're on the auth-success page
                    if (updatedTab.url && updatedTab.url.includes('auth-success')) {
                        console.log('Detected auth success page, waiting for token...');
                        
                        // Wait a bit for the content script to process
                        setTimeout(() => {
                            chrome.tabs.sendMessage(tabId, { action: 'getAuthToken' }, function(response) {
                                if (chrome.runtime.lastError) {
                                    console.log('Error getting auth token:', chrome.runtime.lastError.message);
                                    return;
                                }
                                
                                if (response && response.token) {
                                    console.log('Got token from content script:', !!response.token);
                                    
                                    // Save the auth token
                                    saveAuthTokenInBackground(response.token, response.email);
                                    
                                    // Close the auth tab
                                    chrome.tabs.remove(tabId);
                                    chrome.tabs.onUpdated.removeListener(listener);
                                    authInProgress = false;
                                } else {
                                    console.log('No token received from content script');
                                    // Don't reset authInProgress here, let it timeout
                                }
                            });
                        }, 2000);
                    }
                }
            };
            
            if (chrome.tabs.onUpdated) {
                chrome.tabs.onUpdated.addListener(listener);
                
                // Clean up listener after 5 minutes to prevent memory leaks
                setTimeout(() => {
                    chrome.tabs.onUpdated.removeListener(listener);
                    authInProgress = false; // Reset auth flag on timeout
                }, 300000);
            } else {
                console.error('tabs.onUpdated not available in background');
                authInProgress = false;
                sendResponse({ success: false, error: 'Cannot listen for tab updates' });
                return;
            }
            
            sendResponse({ success: true, tabId: tab.id });
        });
        
        return true; // Keep message channel open for async response
    }
    
    if (request.action === 'authCompleted') {
        console.log('Processing auth completion with token:', !!request.token);
        saveAuthTokenInBackground(request.token, request.email);
        sendResponse({ success: true });
    }
    
    return true; // Keep message channel open for async response
});

// Function to save auth token in background script
function saveAuthTokenInBackground(token, email) {
    console.log('Saving auth token in background...', { token: !!token, email });
    
    // Test if chrome.storage is available
    if (!chrome || !chrome.storage) {
        console.error('Chrome storage API not available in background!');
        return;
    }
    
    // Store the auth data
    chrome.storage.local.set({
        accessToken: token,
        userEmail: email
    }, function() {
        if (chrome.runtime.lastError) {
            console.error('Error saving to storage in background:', chrome.runtime.lastError);
            return;
        }
        
        console.log('Auth data saved to storage successfully in background');
        
        // Test if chrome.tabs is available
        if (!chrome.tabs) {
            console.error('Chrome tabs API not available!');
            return;
        }
        
        // Close auth-related tabs
        chrome.tabs.query({ url: 'https://fluently-app.ru/auth*' }, function(tabs) {
            if (chrome.runtime.lastError) {
                console.error('Error querying tabs:', chrome.runtime.lastError);
            } else {
                tabs.forEach(tab => {
                    console.log('Closing auth tab:', tab.id);
                    chrome.tabs.remove(tab.id);
                });
            }
        });
        
        // Also close auth-success tabs
        chrome.tabs.query({ url: 'https://fluently-app.ru/auth-success*' }, function(tabs) {
            if (chrome.runtime.lastError) {
                console.error('Error querying auth-success tabs:', chrome.runtime.lastError);
            } else {
                tabs.forEach(tab => {
                    console.log('Closing auth-success tab:', tab.id);
                    chrome.tabs.remove(tab.id);
                });
            }
        });
    });
}

// Handle extension installation
chrome.runtime.onInstalled.addListener(function(details) {
    if (details.reason === 'install') {
        console.log('Fluently Word Learner extension installed');
    }
});

// Handle extension installation
chrome.runtime.onInstalled.addListener(function(details) {
    if (details.reason === 'install') {
        console.log('Fluently Word Learner extension installed');
    }
}); 