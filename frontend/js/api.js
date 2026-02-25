// ═══════════════════════════════════════════════
// API — взаимодействие с бэкендом
// ═══════════════════════════════════════════════

const API_BASE_URL = '/museum';

const api = {
    async getAllExhibitions() {
        try {
            const response = await fetch(`${API_BASE_URL}/exhibitions`);
            if (!response.ok) throw new Error('Ошибка загрузки экспозиций');
            const data = await response.json();
            return Array.isArray(data) ? data : (data.exhibitions || []);
        } catch (error) {
            console.error('API Error (exhibitions):', error);
            return [];
        }
    },

    async getExhibitionById(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/exhibitions/${id}`);
            if (!response.ok) throw new Error('Экспозиция не найдена');
            const data = await response.json();
            return data.exhibition || data || null;
        } catch (error) {
            console.error('API Error (exhibition):', error);
            return null;
        }
    },

    async getAllNews() {
        try {
            const response = await fetch(`${API_BASE_URL}/news`);
            if (!response.ok) throw new Error('Ошибка загрузки новостей');
            const data = await response.json();
            return data.news || data || [];
        } catch (error) {
            console.error('API Error (news):', error);
            return [];
        }
    },

    async getNewsById(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/news/${id}`);
            if (!response.ok) throw new Error('Новость не найдена');
            const data = await response.json();
            return data.news || data || null;
        } catch (error) {
            console.error('API Error (news item):', error);
            return null;
        }
    }
};

function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

function truncateText(text, maxLength) {
    if (!text) return '';
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
}

// ═══════════════════════════════════════════════
// Карусель изображений (общие функции)
// ═══════════════════════════════════════════════

/**
 * Возвращает HTML-строку карусели для массива URL-ов изображений.
 * Если изображений нет — пустая строка. Если одно — просто <img>.
 */
function buildImageCarousel(imgs, altText) {
    if (!imgs || imgs.length === 0) return '';

    if (imgs.length === 1) {
        return `<div class="modal-carousel">
            <div class="modal-carousel-track">
                <div class="modal-carousel-slide">
                    <img src="${imgs[0]}" alt="${altText || ''}">
                </div>
            </div>
        </div>`;
    }

    const slides = imgs.map(url =>
        `<div class="modal-carousel-slide"><img src="${url}" alt="${altText || ''}"></div>`
    ).join('');

    const dots = imgs.map((_, i) =>
        `<span class="modal-carousel-dot${i === 0 ? ' active' : ''}" data-index="${i}"></span>`
    ).join('');

    return `<div class="modal-carousel">
        <div class="modal-carousel-track">${slides}</div>
        <button class="modal-carousel-arrow modal-carousel-arrow--prev" aria-label="Назад">&#10094;</button>
        <button class="modal-carousel-arrow modal-carousel-arrow--next" aria-label="Вперёд">&#10095;</button>
        <div class="modal-carousel-dots">${dots}</div>
    </div>`;
}

/**
 * Инициализирует логику карусели внутри переданного контейнера.
 * Навешивает обработчики на стрелки и точки.
 */
function initModalCarousel(container) {
    const carousel = container.querySelector('.modal-carousel');
    if (!carousel) return;

    const track = carousel.querySelector('.modal-carousel-track');
    const slides = carousel.querySelectorAll('.modal-carousel-slide');
    if (slides.length <= 1) return;          // стрелки/точки не нужны

    let current = 0;
    const total = slides.length;

    const prevBtn = carousel.querySelector('.modal-carousel-arrow--prev');
    const nextBtn = carousel.querySelector('.modal-carousel-arrow--next');
    const dots    = carousel.querySelectorAll('.modal-carousel-dot');

    function goTo(index) {
        if (index < 0) index = total - 1;
        if (index >= total) index = 0;
        current = index;
        track.style.transform = `translateX(-${current * 100}%)`;

        dots.forEach((d, i) => d.classList.toggle('active', i === current));
    }

    if (prevBtn) prevBtn.addEventListener('click', () => goTo(current - 1));
    if (nextBtn) nextBtn.addEventListener('click', () => goTo(current + 1));

    dots.forEach(dot => {
        dot.addEventListener('click', () => goTo(Number(dot.dataset.index)));
    });
}


