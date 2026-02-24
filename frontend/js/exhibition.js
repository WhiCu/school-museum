// Страница экспозиции - скрипт

let currentExhibition = null;

// Инициализация страницы
document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const exhibitionId = urlParams.get('id');
    
    if (!exhibitionId) {
        showError('ID экспозиции не указан');
        return;
    }
    
    await loadExhibition(exhibitionId);
});

// Загрузка экспозиции
async function loadExhibition(id) {
    const descriptionBlock = document.getElementById('exhibition-description');
    const exhibitsGrid = document.getElementById('exhibits-grid');
    
    descriptionBlock.innerHTML = '<p class="loading">Загрузка экспозиции</p>';
    exhibitsGrid.innerHTML = '<p class="loading">Загрузка экспонатов</p>';
    
    currentExhibition = await api.getExhibitionById(id);
    
    if (!currentExhibition) {
        showError('Экспозиция не найдена');
        return;
    }
    
    // Обновляем заголовок страницы
    document.title = `${currentExhibition.title} - Музей лицея №76`;
    
    // Показываем описание экспозиции
    descriptionBlock.innerHTML = `
        <h2>${currentExhibition.title}</h2>
        <p>${currentExhibition.description || 'Описание экспозиции'}</p>
    `;
    
    // Показываем экспонаты
    renderExhibits(currentExhibition.exhibits || []);
}

// Отображение экспонатов
function renderExhibits(exhibits) {
    const exhibitsGrid = document.getElementById('exhibits-grid');
    
    if (!exhibits || exhibits.length === 0) {
        exhibitsGrid.innerHTML = '<p>Экспонаты пока не добавлены</p>';
        return;
    }
    
    exhibitsGrid.innerHTML = exhibits.map(exhibit => `
        <div class="exhibit-card" onclick="openExhibitModal('${exhibit.id}')">
            <div class="exhibit-image">
                ${exhibit.image_url 
                    ? `<img src="${exhibit.image_url}" alt="${exhibit.title}">`
                    : '<span class="exhibit-placeholder">✻</span>'
                }
            </div>
            <div class="exhibit-info">
                <h4 class="exhibit-title">${exhibit.title}</h4>
                <div class="exhibit-description-box">
                    ${truncateText(exhibit.description || 'описание экспоната', 100)}
                </div>
            </div>
        </div>
    `).join('');
}

// Открыть модальное окно экспоната
function openExhibitModal(exhibitId) {
    const exhibits = currentExhibition.exhibits || [];
    const exhibit = exhibits.find(e => e.id === exhibitId);
    if (!exhibit) return;
    
    // Создаем модальное окно, если его нет
    let modal = document.querySelector('.modal');
    if (!modal) {
        modal = document.createElement('div');
        modal.className = 'modal';
        modal.innerHTML = `
            <div class="modal-content">
                <span class="modal-close" onclick="closeModal()">&times;</span>
                <div class="modal-body"></div>
            </div>
        `;
        document.body.appendChild(modal);
        
        // Закрытие по клику вне модального окна
        modal.addEventListener('click', (e) => {
            if (e.target === modal) closeModal();
        });
    }
    
    // Заполняем модальное окно
    const modalBody = modal.querySelector('.modal-body');
    modalBody.innerHTML = `
        ${exhibit.image_url 
            ? `<img src="${exhibit.image_url}" alt="${exhibit.title}" class="modal-image">`
            : ''
        }
        <h2 class="modal-title">${exhibit.title}</h2>
        <p class="modal-description">${exhibit.description || 'Описание экспоната отсутствует'}</p>
    `;
    
    modal.classList.add('active');
    document.body.style.overflow = 'hidden';
}

// Закрыть модальное окно
function closeModal() {
    const modal = document.querySelector('.modal');
    if (modal) {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    }
}

// Закрытие по Escape
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') closeModal();
});

// Показать ошибку
function showError(message) {
    const descriptionBlock = document.getElementById('exhibition-description');
    const exhibitsGrid = document.getElementById('exhibits-grid');
    
    descriptionBlock.innerHTML = `
        <h2>Ошибка</h2>
        <p>${message}</p>
        <p><a href="index.html">← Вернуться на главную</a></p>
    `;
    
    exhibitsGrid.innerHTML = '';
}

// Обрезать текст
function truncateText(text, maxLength) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}
