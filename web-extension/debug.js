const API_BASE_URL = 'https://fluently-app.ru';
const AUTH_URL = 'https://fluently-app.ru/api/v1/google/auth';

function log(message) {
    const logDiv = document.getElementById('log');
    const timestamp = new Date().toISOString().split('T')[1].split('.')[0];
    logDiv.innerHTML += `[${timestamp}] ${message}\n`;
    logDiv.scrollTop = logDiv.scrollHeight;
    console.log(message);
}

function clearLog() {
    document.getElementById('log').innerHTML = '';
}

function updatePermissionStatus(name, available, className) {
    const container = document.getElementById('permissionStatus');
    const div = document.createElement('div');
    div.className = `permission-check ${className}`;
    div.textContent = `${name}: ${available ? '✅ Available' : '❌ Not Available'}`;
    container.appendChild(div);
}

function checkPermissions() {
    log('Checking permissions...');
    const container = document.getElementById('permissionStatus');
    container.innerHTML = '';
    
    // Test chrome object
    if (!chrome) {
        updatePermissionStatus('Chrome API', false, 'error');
        log('❌ Chrome API not available!');
        return;
    }
    
    updatePermissionStatus('Chrome API', true, 'success');
    log('✅ Chrome API available');
    
    // Test storage permission
    if (chrome.storage) {
        updatePermissionStatus('Storage API', true, 'success');
        log('✅ Storage API available');
    } else {
        updatePermissionStatus('Storage API', false, 'error');
        log('❌ Storage API not available');
    }
    
    // Test tabs permission
    if (chrome.tabs) {
        updatePermissionStatus('Tabs API', true, 'success');
        log('✅ Tabs API available');
    } else {
        updatePermissionStatus('Tabs API', false, 'error');
        log('❌ Tabs API not available');
    }
    
    // Test scripting permission
    if (chrome.scripting) {
        updatePermissionStatus('Scripting API', true, 'success');
        log('✅ Scripting API available');
    } else {
        updatePermissionStatus('Scripting API', false, 'error');
        log('❌ Scripting API not available');
    }
    
    // Test runtime
    if (chrome.runtime) {
        updatePermissionStatus('Runtime API', true, 'success');
        log('✅ Runtime API available');
    } else {
        updatePermissionStatus('Runtime API', false, 'error');
        log('❌ Runtime API not available');
    }
    
    // Test permissions API
    if (chrome.permissions) {
        updatePermissionStatus('Permissions API', true, 'success');
        log('✅ Permissions API available');
    } else {
        updatePermissionStatus('Permissions API', false, 'error');
        log('❌ Permissions API not available');
    }
}

function requestPermissions() {
    log('Requesting permissions...');
    
    if (!chrome || !chrome.permissions) {
        log('❌ Permissions API not available');
        return;
    }
    
    chrome.permissions.request({
        permissions: ['storage', 'tabs', 'scripting'],
        origins: ['https://fluently-app.ru/*']
    }, function(granted) {
        if (granted) {
            log('✅ Permissions granted successfully!');
            checkPermissions();
        } else {
            log('❌ Permissions were denied');
        }
    });
}

function testStorage() {
    log('Testing storage...');
    const statusDiv = document.getElementById('storageStatus');
    
    if (!chrome || !chrome.storage) {
        statusDiv.innerHTML = '<div class="permission-check error">Storage API not available</div>';
        log('❌ Storage API not available');
        return;
    }
    
    // Test write
    chrome.storage.local.set({testKey: 'testValue', timestamp: Date.now()}, function() {
        if (chrome.runtime.lastError) {
            statusDiv.innerHTML = '<div class="permission-check error">Storage write failed: ' + chrome.runtime.lastError.message + '</div>';
            log('❌ Storage write failed: ' + chrome.runtime.lastError.message);
            return;
        }
        
        log('✅ Storage write successful');
        
        // Test read
        chrome.storage.local.get(['testKey', 'timestamp'], function(result) {
            if (chrome.runtime.lastError) {
                statusDiv.innerHTML = '<div class="permission-check error">Storage read failed: ' + chrome.runtime.lastError.message + '</div>';
                log('❌ Storage read failed: ' + chrome.runtime.lastError.message);
                return;
            }
            
            log('✅ Storage read successful: ' + JSON.stringify(result));
            statusDiv.innerHTML = '<div class="permission-check success">Storage test passed. Data: ' + JSON.stringify(result) + '</div>';
        });
    });
}

function clearStorage() {
    log('Clearing storage...');
    
    if (!chrome || !chrome.storage) {
        log('❌ Storage API not available');
        return;
    }
    
    chrome.storage.local.clear(function() {
        if (chrome.runtime.lastError) {
            log('❌ Storage clear failed: ' + chrome.runtime.lastError.message);
        } else {
            log('✅ Storage cleared successfully');
            document.getElementById('storageStatus').innerHTML = '<div class="permission-check success">Storage cleared</div>';
        }
    });
}

function testAuth() {
    log('Testing authentication status...');
    const statusDiv = document.getElementById('authStatus');
    
    if (!chrome || !chrome.storage) {
        statusDiv.innerHTML = '<div class="permission-check error">Cannot check auth - Storage API not available</div>';
        log('❌ Cannot check auth - Storage API not available');
        return;
    }
    
    chrome.storage.local.get(['accessToken', 'userEmail'], function(result) {
        if (chrome.runtime.lastError) {
            statusDiv.innerHTML = '<div class="permission-check error">Auth check failed: ' + chrome.runtime.lastError.message + '</div>';
            log('❌ Auth check failed: ' + chrome.runtime.lastError.message);
            return;
        }
        
        log('Auth data: ' + JSON.stringify(result));
        
        if (result.accessToken) {
            statusDiv.innerHTML = '<div class="permission-check success">Authenticated as: ' + (result.userEmail || 'Unknown') + '</div>';
            log('✅ User is authenticated');
        } else {
            statusDiv.innerHTML = '<div class="permission-check error">Not authenticated</div>';
            log('❌ User is not authenticated');
        }
    });
}

function openAuthPage() {
    log('Opening authentication page...');
    
    if (!chrome || !chrome.tabs) {
        log('❌ Tabs API not available');
        return;
    }
    
    chrome.tabs.create({ 
        url: AUTH_URL,
        active: true
    }, function(tab) {
        if (chrome.runtime.lastError) {
            log('❌ Failed to create tab: ' + chrome.runtime.lastError.message);
        } else {
            log('✅ Auth tab created: ' + tab.id);
        }
    });
}

// Initialize when popup loads
document.addEventListener('DOMContentLoaded', function() {
    log('Debug panel loaded');
    checkPermissions();
    testAuth();
});
