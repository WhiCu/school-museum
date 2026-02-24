// Главная страница - скрипт

let exhibitions = [];
let currentExhibitionIndex = 0;

// Инициализация страницы
document.addEventListener('DOMContentLoaded', async () => {
    await loadNews();
    await loadExhibitions();
});

// Загрузка новостей
async function loadNews() {
    const newsContainer = document.getElementById('news-container');
    newsContainer.innerHTML = '<div class="loading">Загрузка новостей</div>';

    const newsData = await api.getAllNews();
    
    if (newsData && newsData.length > 0) {
        newsContainer.innerHTML = newsData.map(news => `
            <div class="news-card" onclick="openNews('${news.id}')">
                ${news.image_url ? `<img src="${news.image_url}" alt="${news.title}" class="news-image">` : ''}
                <h3>${news.title}</h3>
                <p>${truncateText(news.content, 150)}</p>
                <div class="news-date">${formatDate(news.created_at)}</div>
            </div>
        `).join('');
    } else {
        newsContainer.innerHTML = '<p>Новости пока отсутствуют</p>';
    }
}

// Загрузка экспозиций
async function loadExhibitions() {
    const exhibitionCard = document.getElementById('exhibition-card');
    const dotsContainer = document.getElementById('carousel-dots');
    
    exhibitions = await api.getAllExhibitions();
    
    if (exhibitions && exhibitions.length > 0) {
        // Создаем точки навигации
        dotsContainer.innerHTML = exhibitions.map((_, index) => `
            <span class="dot ${index === 0 ? 'active' : ''}" onclick="goToExhibition(${index})"></span>
        `).join('');
        
        // Показываем первую экспозицию
        showExhibition(0);
    } else {
        exhibitionCard.innerHTML = `
            <div class="exhibition-info">
                <p class="exhibition-note">Экспозиции пока не добавлены</p>
            </div>
        `;
    }
}

// Показать экспозицию по индексу
function showExhibition(index) {
    if (!exhibitions || exhibitions.length === 0) return;
    
    currentExhibitionIndex = index;
    const exhibition = exhibitions[index];
    const exhibitionCard = document.getElementById('exhibition-card');
    
    exhibitionCard.innerHTML = `
        <div class="exhibition-info" onclick="openExhibition('${exhibition.id}')">
            <h3 class="exhibition-title">${exhibition.title}</h3>
            <p class="exhibition-text">${exhibition.description || 'Нажмите, чтобы узнать больше'}</p>
            <p class="exhibition-description">Нажмите для перехода к экспозиции →</p>
        </div>
    `;
    
    // Обновляем точки
    updateDots();
}

// Обновить активную точку
function updateDots() {
    const dots = document.querySelectorAll('.dot');
    dots.forEach((dot, index) => {
        dot.classList.toggle('active', index === currentExhibitionIndex);
    });
}

// Предыдущая экспозиция
function prevExhibition() {
    if (exhibitions.length === 0) return;
    currentExhibitionIndex = (currentExhibitionIndex - 1 + exhibitions.length) % exhibitions.length;
    showExhibition(currentExhibitionIndex);
}

// Следующая экспозиция
function nextExhibition() {
    if (exhibitions.length === 0) return;
    currentExhibitionIndex = (currentExhibitionIndex + 1) % exhibitions.length;
    showExhibition(currentExhibitionIndex);
}

// Перейти к конкретной экспозиции
function goToExhibition(index) {
    showExhibition(index);
}

// Открыть страницу экспозиции
function openExhibition(id) {
    window.location.href = `exhibition.html?id=${id}`;
}

// Открыть новость (можно реализовать модальное окно)
function openNews(id) {
    console.log('Open news:', id);
    // TODO: Можно добавить модальное окно с полным текстом новости
}

// Обрезать текст
function truncateText(text, maxLength) {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

// Автопрокрутка карусели
let autoSlideInterval;

function startAutoSlide() {
    autoSlideInterval = setInterval(() => {
        nextExhibition();
    }, 5000);
}

function stopAutoSlide() {
    clearInterval(autoSlideInterval);
}

// Запуск автопрокрутки после загрузки
document.addEventListener('DOMContentLoaded', () => {
    startAutoSlide();
    
    // Останавливаем при наведении на карусель
    const carousel = document.querySelector('.exhibitions-carousel');
    if (carousel) {
        carousel.addEventListener('mouseenter', stopAutoSlide);
        carousel.addEventListener('mouseleave', startAutoSlide);
    }
});
