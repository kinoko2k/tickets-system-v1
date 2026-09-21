let categoriesConfig = [];

document.addEventListener('DOMContentLoaded', async () => {
    try {
        const res = await fetch('/api/categories');
        categoriesConfig = await res.json();
        renderCategoryContainers();
    } catch (e) {
        console.error(e);
    }

    setupSSE();
    
    setupAboutPopup();
});

function renderCategoryContainers() {
    const container = document.getElementById('waiting-categories');
    if (!categoriesConfig || categoriesConfig.length === 0) {
        container.innerHTML = '<div class="no-number-message"></div>';
        return;
    }

    let html = '';
    categoriesConfig.forEach(cat => {
        const bgColor = `${cat.color}1a`;
        html += `
            <div class="category-tab">
                <h3 class="category-title" style="background-color: ${bgColor}; color: ${cat.color}; border-left: 4px solid ${cat.color};">
                    ${cat.rangeStart}番台 <span class="category-label">(${cat.name})</span>
                </h3>
                <div class="waiting-container" id="waiting-numbers-${cat.id}">
                    <div class="no-number-message"></div>
                </div>
            </div>
        `;
    });
    container.innerHTML = html;
    
    const popupDesc = document.getElementById('popup-categories-desc');
    if(popupDesc) {
        let descHtml = '<strong></strong><br>';
        categoriesConfig.forEach(cat => {
            descHtml += `<span style="color: ${cat.color};">${cat.rangeStart}</span> - ${cat.name}<br>`;
        });
        popupDesc.innerHTML = descHtml;
    }
}

function setupSSE() {
    const eventSource = new EventSource('/api/events');

    eventSource.onmessage = function(event) {
        const data = JSON.parse(event.data);
        updateUI(data);
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

function updateUI(data) {
    const currentContainer = document.getElementById('current-numbers');
    if (data.currentNumbers && data.currentNumbers.length > 0) {
        let currentHTML = '';
        data.currentNumbers.forEach(number => {
            currentHTML += `<div class="number-card">${number}</div>`;
        });
        currentContainer.innerHTML = currentHTML;
    } else {
        currentContainer.innerHTML = '<div class="no-number-message"></div>';
    }

    if (!categoriesConfig || categoriesConfig.length === 0) return;

    const categorized = {};
    categoriesConfig.forEach(cat => categorized[cat.id] = []);

    if (data.waitingNumbers) {
        data.waitingNumbers.forEach(number => {
            const cat = getCategoryForNumber(number);
            if (cat) {
                categorized[cat.id].push(number);
            }
        });
    }

    categoriesConfig.forEach(cat => {
        const container = document.getElementById(`waiting-numbers-${cat.id}`);
        if (!container) return;
        
        const numbers = categorized[cat.id];
        if (numbers.length > 0) {
            const bgColor = `${cat.color}1a`;
            let html = '<ul class="number-list">';
            numbers.forEach(number => {
                html += `<li class="number-list-item" style="background-color: ${bgColor}; border-color: ${cat.color}; color: ${cat.color};">${number}</li>`;
            });
            html += '</ul>';
            container.innerHTML = html;
        } else {
            container.innerHTML = '<div class="no-number-message"></div>';
        }
    });
}

function setupAboutPopup() {
    const aboutBtn = document.getElementById('about-btn');
    const aboutPopup = document.getElementById('about-popup');
    const closeButtons = document.querySelectorAll('.popup-close');

    if(aboutBtn) {
        aboutBtn.addEventListener('click', function () {
            aboutPopup.classList.add('active');
            document.body.style.overflow = 'hidden';
        });
    }

    closeButtons.forEach(button => {
        button.addEventListener('click', function () {
            const popup = this.closest('.popup-overlay');
            popup.classList.remove('active');
            document.body.style.overflow = '';
        });
    });

    document.querySelectorAll('.popup-overlay').forEach(overlay => {
        overlay.addEventListener('click', function (e) {
            if (e.target === this) {
                this.classList.remove('active');
                document.body.style.overflow = '';
            }
        });
    });
}

