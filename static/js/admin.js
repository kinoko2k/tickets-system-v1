let categoriesConfig = [];
let currentState = { currentNumbers: [], waitingNumbers: [] };
let orderData = { orders: {} };

document.addEventListener('DOMContentLoaded', async () => {
    try {
        const res = await fetch('/api/categories');
        categoriesConfig = await res.json();
    } catch (e) {
        console.error(e);
    }

    document.getElementById('add-calling-btn').addEventListener('click', () => {
        const input = document.getElementById('calling-number-input');
        const num = parseInt(input.value);
        if (num > 0) addCallingNumber(num);
        input.value = '';
    });

    document.getElementById('add-waiting-btn').addEventListener('click', () => {
        const input = document.getElementById('waiting-number-input');
        const num = parseInt(input.value);
        if (num > 0) addWaitingNumber(num);
        input.value = '';
    });

    document.getElementById('call-next-btn').addEventListener('click', callNext);
    
    document.getElementById('clear-current-btn').addEventListener('click', () => {
        if(confirm("現在呼び出し中の番号をすべてクリアしますか？")) {
            updateServerCurrent([]);
        }
    });

    setupSSE();
});

function setupSSE() {
    const eventSource = new EventSource('/api/events');

    eventSource.onmessage = function(event) {
        const data = JSON.parse(event.data);
        currentState = data;
        updateUI();
    };

    eventSource.onerror = function(error) {
        console.error(error);
    };
}

function getCategoryForNumber(number) {
    for (const cat of categoriesConfig) {
        if (number >= cat.rangeStart && number <= cat.rangeEnd) {
            return cat;
        }
    }
    return null;
}

function updateUI() {
    const currentContainer = document.getElementById('admin-current-numbers');
    if (currentState.currentNumbers && currentState.currentNumbers.length > 0) {
        let currentHTML = '';
        currentState.currentNumbers.forEach(number => {
            const cat = getCategoryForNumber(number);
            const productName = cat ? cat.name : '';
            const productText = productName ? `(${productName})` : '';
            
            currentHTML += `
                <div class="current-item" data-number="${number}">
                    <span class="current-number">${number} ${productText}</span>
                    <div class="current-actions">
                        <button class="mui-btn action-btn handover-btn" onclick="removeCurrent(${number})">
                            受け渡し
                        </button>
                    </div>
                </div>
            `;
        });
        currentContainer.innerHTML = currentHTML;
    } else {
        currentContainer.innerHTML = '<div class="no-number-message"></div>';
    }

    const waitingContainer = document.getElementById('admin-waiting-numbers');
    document.getElementById('waiting-count').textContent = currentState.waitingNumbers ? currentState.waitingNumbers.length : 0;

    if (currentState.waitingNumbers && currentState.waitingNumbers.length > 0) {
        let waitingHTML = '';
        currentState.waitingNumbers.forEach(number => {
            const cat = getCategoryForNumber(number);
            const productName = cat ? cat.name : '';
            const productText = productName ? `(${productName})` : '';
            
            let categoryButtons = '';
            categoriesConfig.forEach(c => {
                const isActive = (cat && cat.id === c.id) ? 'active' : '';
                const btnStyle = isActive 
                    ? `background-color: ${c.color}; color: white; border-color: ${c.color}; box-shadow: 0 1px 3px rgba(0,0,0,0.2);` 
                    : `background-color: rgba(0, 0, 0, 0.05); color: var(--mui-text-secondary); border: 1px solid var(--mui-grey-300);`;
                categoryButtons += `
                    <button class="mui-btn action-btn ${isActive}" style="${btnStyle}" onclick="dummyAssign(${number}, '${c.id}')">
                        ${c.name}
                    </button>
                `;
            });
            
            waitingHTML += `
                <div class="waiting-item" data-number="${number}">
                    <span class="waiting-number">${number} ${productText}</span>
                    <div class="waiting-actions">
                        ${categoryButtons}
                        <button class="mui-btn action-btn call-btn" onclick="callSpecificNumber(${number})">
                            呼び出し
                        </button>
                        <button class="mui-btn action-btn" style="background-color: #f44336; color: white;" onclick="removeWaiting(${number})">
                            削除
                        </button>
                    </div>
                </div>
            `;
        });
        waitingContainer.innerHTML = waitingHTML;
    } else {
        waitingContainer.innerHTML = '<div class="no-number-message"></div>';
    }
}

window.dummyAssign = function(num, catId) {
}

function addCallingNumber(num) {
    if (currentState.currentNumbers.includes(num)) {
        return;
    }
    const newNumbers = [...currentState.currentNumbers, num];
    
    if (currentState.waitingNumbers.includes(num)) {
        removeWaiting(num, true);
    }
    updateServerCurrent(newNumbers);
}

function addWaitingNumber(num) {
    if (currentState.waitingNumbers.includes(num) || currentState.currentNumbers.includes(num)) {
        return;
    }
    const newNumbers = [...currentState.waitingNumbers, num];
    updateServerWaiting(newNumbers);
}

function callNext() {
    if (!currentState.waitingNumbers || currentState.waitingNumbers.length === 0) {
        return;
    }
    const nextNum = currentState.waitingNumbers[0];
    
    const newWaiting = currentState.waitingNumbers.slice(1);
    const newCurrent = [...currentState.currentNumbers, nextNum];
    
    updateServerWaiting(newWaiting).then(() => {
        updateServerCurrent(newCurrent);
    });
}

window.callSpecificNumber = function(num) {
    const newWaiting = currentState.waitingNumbers.filter(n => n !== num);
    const newCurrent = [...currentState.currentNumbers, num];
    
    updateServerWaiting(newWaiting).then(() => {
        updateServerCurrent(newCurrent);
        showNotification(`${num}`);
    });
}

window.removeCurrent = function(num) {
    const newNumbers = currentState.currentNumbers.filter(n => n !== num);
    updateServerCurrent(newNumbers);
}

window.removeWaiting = function(num, skipUpdate = false) {
    const newNumbers = currentState.waitingNumbers.filter(n => n !== num);
    if (!skipUpdate) {
        updateServerWaiting(newNumbers);
    }
}

function showNotification(msg) {
    const notification = document.createElement('div');
    notification.className = 'call-notification';
    notification.textContent = msg;
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.classList.add('show');
        setTimeout(() => {
            notification.classList.remove('show');
            setTimeout(() => {
                document.body.removeChild(notification);
            }, 300);
        }, 2000);
    }, 10);
}

async function updateServerCurrent(numbers) {
    try {
        await fetch('/api/numbers/current', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(numbers)
        });
    } catch (e) {
        console.error(e);
    }
}

async function updateServerWaiting(numbers) {
    try {
        await fetch('/api/numbers/waiting', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(numbers)
        });
    } catch (e) {
        console.error(e);
    }
}

