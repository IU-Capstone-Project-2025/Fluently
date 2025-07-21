// Configuration
const API_BASE_URL = 'https://fluently-app.ru';
const AUTH_URL = `${API_BASE_URL}/auth/google`;

// DOM elements - will be set after DOM loads
let authSection, userInfo, wordInputSection, loginBtn, logoutBtn, userEmail, wordInput, addWordBtn, loading, message;

// Initialize popup
document.addEventListener('DOMContentLoaded', function() {
    console.log('Popup DOM loaded');
    
    // Test permissions first
    testPermissions();
    
    // Get DOM elements
    authSection = document.getElementById('authSection');
    userInfo = document.getElementById('userInfo');
    wordInputSection = document.getElementById('wordInputSection');
    loginBtn = document.getElementById('loginBtn');
    logoutBtn = document.getElementById('logoutBtn');
    userEmail = document.getElementById('userEmail');
    wordInput = document.getElementById('wordInput');
    addWordBtn = document.getElementById('addWordBtn');
    loading = document.getElementById('loading');
    message = document.getElementById('message');
    
    console.log('DOM elements found:', {
        authSection: !!authSection,
        userInfo: !!userInfo,
        loginBtn: !!loginBtn,
        logoutBtn: !!logoutBtn
    });
    
    checkAuthStatus();
    setupEventListeners();
});

function testPermissions() {
    console.log('Testing extension permissions...');
    
    // Test chrome object
    if (!chrome) {
        console.error('Chrome API not available!');
        return;
    }
    
    // Test storage permission
    if (chrome.storage) {
        console.log('✅ Storage permission available');
        // Test a simple storage operation
        chrome.storage.local.set({test: 'value'}, function() {
            if (chrome.runtime.lastError) {
                console.error('❌ Storage test failed:', chrome.runtime.lastError);
            } else {
                console.log('✅ Storage test successful');
                chrome.storage.local.remove('test');
            }
        });
        
        // Also test checking existing auth data
        chrome.storage.local.get(['accessToken', 'userEmail'], function(result) {
            console.log('💾 Current storage contents:', result);
            console.log('💾 Storage keys:', Object.keys(result));
        });
    } else {
        console.error('❌ Storage permission NOT available');
    }
    
    // Test tabs permission
    if (chrome.tabs) {
        console.log('✅ Tabs permission available');
    } else {
        console.error('❌ Tabs permission NOT available');
    }
    
    // Test scripting permission
    if (chrome.scripting) {
        console.log('✅ Scripting permission available');
    } else {
        console.error('❌ Scripting permission NOT available');
    }
    
    // Test runtime
    if (chrome.runtime) {
        console.log('✅ Runtime API available');
    } else {
        console.error('❌ Runtime API NOT available');
    }
}

// Listen for changes in chrome.storage.local to update UI on auth changes (only if available)
if (chrome && chrome.storage && chrome.storage.onChanged) {
    chrome.storage.onChanged.addListener(function(changes, area) {
        if (area === 'local' && (changes.accessToken || changes.userEmail)) {
            console.log('Storage changed - auth data updated:', changes);
            
            // Re-check auth status when storage changes
            setTimeout(() => {
                checkAuthStatus();
            }, 100);
        }
    });
}

function setupEventListeners() {
    console.log('Setting up event listeners...');
    
    if (loginBtn) {
        loginBtn.addEventListener('click', handleLogin);
        console.log('Login button listener added');
    } else {
        console.error('Login button not found!');
    }
    
    if (logoutBtn) {
        logoutBtn.addEventListener('click', handleLogout);
        console.log('Logout button listener added');
    }
    
    if (addWordBtn) {
        addWordBtn.addEventListener('click', handleAddWord);
    }
    
    if (wordInput) {
        wordInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                handleAddWord();
            }
        });
    }
}

// Authentication functions
function checkAuthStatus() {
    console.log('Checking auth status...');
    
    // Test if chrome.storage is available
    if (!chrome || !chrome.storage) {
        console.error('Chrome storage API not available! Extension may not be properly loaded.');
        showUnauthenticatedState();
        showError('Разрешения расширения недоступны. Попробуйте перезагрузить расширение.');
        return;
    }
    
    chrome.storage.local.get(['accessToken', 'userEmail'], function(result) {
        if (chrome.runtime.lastError) {
            console.error('Storage error:', chrome.runtime.lastError);
            showUnauthenticatedState();
            showError('Не удаётся получить доступ к хранилищу: ' + chrome.runtime.lastError.message);
            return;
        }
        
        console.log('Storage result:', result);
        console.log('Storage keys found:', Object.keys(result));
        console.log('accessToken exists:', !!result.accessToken);
        console.log('userEmail exists:', !!result.userEmail);
        
        if (result.accessToken) {
            console.log('User is authenticated');
            showAuthenticatedState(result.userEmail);
        } else {
            console.log('User is not authenticated - no accessToken found');
            showUnauthenticatedState();
        }
    });
}

function handleLogin() {
    console.log('Starting login process...');
    
    // Prevent multiple clicks
    if (loginBtn && loginBtn.disabled) {
        console.log('Login already in progress');
        return;
    }
    
    // Check if tabs API is available
    if (!chrome || !chrome.tabs) {
        console.error('Chrome tabs API not available!');
        showError('Разрешение на работу с вкладками недоступно. Проверьте разрешения расширения.');
        return;
    }
    
    // Disable button immediately
    if (loginBtn) {
        loginBtn.disabled = true;
        loginBtn.textContent = 'Начинаем аутентификацию...';
    }
    
    // Send message to background script to handle auth
    // Background script persists even when popup closes
    chrome.runtime.sendMessage({ action: 'startAuth' }, function(response) {
        if (chrome.runtime.lastError) {
            console.error('Failed to start auth:', chrome.runtime.lastError);
            showError('Не удалось начать аутентификацию: ' + chrome.runtime.lastError.message);
            
            // Re-enable button on error
            if (loginBtn) {
                loginBtn.disabled = false;
                loginBtn.textContent = 'Войти через Google';
            }
            return;
        }
        
        if (response && !response.success) {
            console.error('Auth start failed:', response.error);
            showError('Ошибка аутентификации: ' + response.error);
            
            // Re-enable button on error
            if (loginBtn) {
                loginBtn.disabled = false;
                loginBtn.textContent = 'Войти через Google';
            }
            return;
        }
        
        console.log('Auth request sent to background script');
        
        // Show a message to user and inform them they can close the popup
        showMessage('Открываем страницу аутентификации... Можете закрыть это окно во время аутентификации.', 'success');
        
        // Update button text
        if (loginBtn) {
            loginBtn.textContent = 'Идёт аутентификация...';
            
            // Re-enable after a few seconds
            setTimeout(() => {
                if (loginBtn) {
                    loginBtn.textContent = 'Войти через Google';
                    loginBtn.disabled = false;
                }
            }, 10000); // Increased timeout to 10 seconds
        }
    });
}

function handleLogout() {
    chrome.storage.local.remove(['accessToken', 'userEmail'], function() {
        showUnauthenticatedState();
        showMessage('Вы успешно вышли из системы', 'success');
    });
}

function saveAuthToken(token, email) {
    console.log('Saving auth token...', { token: !!token, email });
    
    if (!chrome || !chrome.storage) {
        console.error('Chrome storage API not available!');
        showError('Невозможно сохранить аутентификацию - разрешение на хранение недоступно');
        return;
    }
    
    chrome.storage.local.set({
        accessToken: token,
        userEmail: email
    }, function() {
        if (chrome.runtime.lastError) {
            console.error('Error saving to storage:', chrome.runtime.lastError);
            showError('Не удалось сохранить аутентификацию: ' + chrome.runtime.lastError.message);
            return;
        }
        
        console.log('Auth data saved successfully');
        showAuthenticatedState(email);
        showMessage('Вход выполнен успешно!', 'success');
    });
}

function showAuthenticatedState(email) {
    console.log('Showing authenticated state for:', email);
    
    if (authSection) {
        authSection.style.display = 'none';
    }
    if (userInfo) {
        userInfo.style.display = 'block';
    }
    if (wordInputSection) {
        wordInputSection.style.display = 'block';
    }
    if (userEmail) {
        userEmail.textContent = email || 'User';
    }
}

function showUnauthenticatedState() {
    console.log('Showing unauthenticated state');
    
    if (authSection) {
        authSection.style.display = 'block';
    }
    if (userInfo) {
        userInfo.style.display = 'none';
    }
    if (wordInputSection) {
        wordInputSection.style.display = 'none';
    }
}

// Word management functions
async function handleAddWord() {
    const word = wordInput.value.trim();
    
    // Client-side validation
    if (!word) {
        showMessage('Пожалуйста, введите слово', 'error');
        return;
    }
    
    // Check if word contains only letters (basic validation)
    if (!/^[a-zA-Z\s-']+$/.test(word)) {
        showMessage('Пожалуйста, введите корректное английское слово (только буквы)', 'error');
        return;
    }
    
    // Check word length
    if (word.length > 50) {
        showMessage('Слово слишком длинное (максимум 50 символов)', 'error');
        return;
    }
    
    if (word.length < 2) {
        showMessage('Слово слишком короткое (минимум 2 символа)', 'error');
        return;
    }

    // Get auth token
    chrome.storage.local.get(['accessToken'], async function(result) {
        if (!result.accessToken) {
            showMessage('Пожалуйста, сначала войдите в систему', 'error');
            return;
        }

        setLoading(true);
        
        try {
            const response = await addWordToNotLearned(word, result.accessToken);
            if (response.success) {
                showMessage(`"${word}" добавлено в ваш список изучения!`, 'success');
                wordInput.value = '';
            } else {
                // Handle authentication errors
                if (response.authError) {
                    // Clear stored auth and show login screen
                    chrome.storage.local.remove(['accessToken', 'userEmail']);
                    showUnauthenticatedState();
                    showMessage(response.error, 'error');
                } else {
                    showMessage(response.error || 'Не удалось добавить слово', 'error');
                }
            }
        } catch (error) {
            console.error('Error adding word:', error);
            showMessage('Неожиданная ошибка. Пожалуйста, попробуйте снова.', 'error');
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

        // Handle different response status codes
        if (response.ok) {
            const data = await response.json();
            return { success: true, data: data };
        } else {
            // Try to get error message from response body
            let errorMessage = `HTTP ${response.status}: ${response.statusText}`;
            
            try {
                const errorText = await response.text();
                if (errorText) {
                    errorMessage = errorText;
                }
            } catch (e) {
                // If we can't read the error body, use the status text
            }
            
            // Handle specific error cases
            switch (response.status) {
                case 400:
                    if (errorMessage.includes('word cannot be empty')) {
                        return { success: false, error: 'Пожалуйста, введите слово для добавления' };
                    }
                    return { success: false, error: 'Некорректный запрос: ' + errorMessage };
                    
                case 401:
                    return { success: false, error: 'Сессия истекла. Пожалуйста, войдите снова.', authError: true };
                    
                case 404:
                    if (errorMessage.includes('word not found in database')) {
                        return { success: false, error: `"${word}" нет в нашей базе слов. Попробуйте другое слово.` };
                    }
                    return { success: false, error: 'Слово не найдено: ' + errorMessage };
                    
                case 409:
                    if (errorMessage.includes('word already in not learned list')) {
                        return { success: false, error: `"${word}" уже в вашем списке изучения!` };
                    }
                    return { success: false, error: 'Слово уже существует: ' + errorMessage };
                    
                case 500:
                    return { success: false, error: 'Ошибка сервера. Пожалуйста, попробуйте позже.' };
                    
                default:
                    return { success: false, error: errorMessage };
            }
        }
    } catch (error) {
        // Network errors
        if (error.name === 'TypeError' && error.message.includes('fetch')) {
            return { success: false, error: 'Ошибка сети. Проверьте подключение к интернету.' };
        }
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

function showError(text) {
    showMessage(text, 'error');
}

// Add a manual permission test function for debugging
function requestPermissions() {
    if (!chrome || !chrome.permissions) {
        showError('API разрешений недоступно');
        return;
    }
    
    chrome.permissions.request({
        permissions: ['storage', 'tabs', 'scripting'],
        origins: ['https://fluently-app.ru/*']
    }, function(granted) {
        if (granted) {
            showMessage('Разрешения успешно предоставлены!', 'success');
            // Re-test permissions
            testPermissions();
        } else {
            showError('Разрешения были отклонены');
        }
    });
} 