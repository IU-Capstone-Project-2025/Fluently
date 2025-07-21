// Configuration
const API_BASE_URL = 'https://fluently-app.ru';
const AUTH_URL = `${API_BASE_URL}/auth/google`;

// DOM elements
const authSection = document.getElementById('authSection');
const userInfo = document.getElementById('userInfo');
const wordInputSection = document.getElementById('wordInputSection');
const loginBtn = document.getElementById('loginBtn');
const logoutBtn = document.getElementById('logoutBtn');
const userEmail = document.getElementById('userEmail');
const wordInput = document.getElementById('wordInput');
const addWordBtn = document.getElementById('addWordBtn');
const loading = document.getElementById('loading');
const message = document.getElementById('message');

// Initialize popup
document.addEventListener('DOMContentLoaded', function() {
    checkAuthStatus();
    setupEventListeners();
});

function setupEventListeners() {
    loginBtn.addEventListener('click', handleLogin);
    logoutBtn.addEventListener('click', handleLogout);
    addWordBtn.addEventListener('click', handleAddWord);
    wordInput.addEventListener('keypress', function(e) {
        if (e.key === 'Enter') {
            handleAddWord();
        }
    });
}

// Authentication functions
function checkAuthStatus() {
    chrome.storage.local.get(['accessToken', 'userEmail'], function(result) {
        if (result.accessToken) {
            showAuthenticatedState(result.userEmail);
        } else {
            showUnauthenticatedState();
        }
    });
}

function handleLogin() {
    // Open auth page in a new tab
    chrome.tabs.create({ url: AUTH_URL }, function(tab) {
        // Listen for the auth completion
        chrome.tabs.onUpdated.addListener(function listener(tabId, changeInfo, updatedTab) {
            if (tabId === tab.id && changeInfo.status === 'complete') {
                // Check if we're back from auth (you might need to adjust this logic)
                if (updatedTab.url.includes('fluently-app.ru') && !updatedTab.url.includes('/auth/')) {
                    chrome.tabs.onUpdated.removeListener(listener);
                    // Try to get the token from the page
                    chrome.tabs.sendMessage(tabId, { action: 'getAuthToken' }, function(response) {
                        if (response && response.token) {
                            saveAuthToken(response.token, response.email);
                            chrome.tabs.remove(tabId);
                        }
                    });
                }
            }
        });
    });
}

function handleLogout() {
    chrome.storage.local.remove(['accessToken', 'userEmail'], function() {
        showUnauthenticatedState();
        showMessage('Logged out successfully', 'success');
    });
}

function saveAuthToken(token, email) {
    chrome.storage.local.set({
        accessToken: token,
        userEmail: email
    }, function() {
        showAuthenticatedState(email);
        showMessage('Login successful!', 'success');
    });
}

function showAuthenticatedState(email) {
    authSection.style.display = 'none';
    userInfo.style.display = 'block';
    wordInputSection.style.display = 'block';
    userEmail.textContent = email || 'User';
}

function showUnauthenticatedState() {
    authSection.style.display = 'block';
    userInfo.style.display = 'none';
    wordInputSection.style.display = 'none';
}

// Word management functions
async function handleAddWord() {
    const word = wordInput.value.trim();
    
    if (!word) {
        showMessage('Please enter a word', 'error');
        return;
    }

    // Get auth token
    chrome.storage.local.get(['accessToken'], async function(result) {
        if (!result.accessToken) {
            showMessage('Please login first', 'error');
            return;
        }

        setLoading(true);
        
        try {
            const response = await addWordToNotLearned(word, result.accessToken);
            if (response.success) {
                showMessage(`"${word}" added to your learning list!`, 'success');
                wordInput.value = '';
            } else {
                showMessage(response.error || 'Failed to add word', 'error');
            }
        } catch (error) {
            console.error('Error adding word:', error);
            showMessage('Network error. Please try again.', 'error');
        } finally {
            setLoading(false);
        }
    });
}

async function addWordToNotLearned(word, token) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/not-learned-words`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ word: word })
        });

        const data = await response.json();

        if (response.ok) {
            return { success: true, data: data };
        } else {
            return { 
                success: false, 
                error: data.error || `HTTP ${response.status}: ${response.statusText}` 
            };
        }
    } catch (error) {
        throw error;
    }
}

// UI helper functions
function setLoading(isLoading) {
    if (isLoading) {
        loading.style.display = 'block';
        addWordBtn.disabled = true;
        wordInput.disabled = true;
    } else {
        loading.style.display = 'none';
        addWordBtn.disabled = false;
        wordInput.disabled = false;
    }
}

function showMessage(text, type) {
    message.textContent = text;
    message.className = `message ${type} show`;
    
    // Auto-hide after 3 seconds
    setTimeout(() => {
        message.classList.remove('show');
    }, 3000);
} 