// Background service worker
chrome.runtime.onMessage.addListener(function(request, sender, sendResponse) {
    if (request.action === 'authCompleted') {
        // Store the auth data
        chrome.storage.local.set({
            accessToken: request.token,
            userEmail: request.email
        }, function() {
            // Close the auth tab if it's still open
            chrome.tabs.query({ url: 'https://fluently-app.ru/auth/*' }, function(tabs) {
                tabs.forEach(tab => {
                    chrome.tabs.remove(tab.id);
                });
            });
        });
    }
});

// Handle extension installation
chrome.runtime.onInstalled.addListener(function(details) {
    if (details.reason === 'install') {
        console.log('Fluently Word Learner extension installed');
    }
}); 