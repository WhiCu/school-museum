// ═══════════════════════════════════════════════
// Главная страница — JavaScript
// ═══════════════════════════════════════════════

document.addEventListener('DOMContentLoaded', () => {
    initHeader();
    initBurger();
    initScrollAnimations();
    initStatCounters();
    loadExhibitions();
    loadNews();
    initNewsModal();
});

// ── Header scroll effect ──
function initHeader() {
    const header = document.getElementById('header');
    if (!header) return;

    const onScroll = () => {
        header.classList.toggle('scrolled', window.scrollY > 50);
    };
    window.addEventListener('scroll', onScroll, { passive: true });
    onScroll();
}

// ── Burger menu ──
function initBurger() {
    const burger = document.getElementById('burger');
    const nav = document.getElementById('nav');
    if (!burger || !nav) return;

    burger.addEventListener('click', () => {
        burger.classList.toggle('active');
        nav.classList.toggle('open');
    });

    // Close on link click
    nav.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', () => {
            burger.classList.remove('active');
            nav.classList.remove('open');
        });
    });
}

// ── Scroll animations (Intersection Observer) ──
function initScrollAnimations() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
                observer.unobserve(entry.target);
            }
        });
    }, { threshold: 0.1, rootMargin: '0px 0px -40px 0px' });

    // Add fade-in to key elements
    const selectors = [
        '.about-content',
        '.stat-card',
        '.exhibition-card',
        '.news-card',
        '.team-card',
        '.visit-card',
        '.contacts-inner'
    ];

    // Delay to let DOM render
    requestAnimationFrame(() => {
        selectors.forEach(sel => {
            document.querySelectorAll(sel).forEach(el => {
                el.classList.add('fade-in');
                observer.observe(el);
            });
        });
    });
}

// ── Stat counters animation ──
function initStatCounters() {
    const counters = document.querySelectorAll('.stat-number[data-target]');
    if (!counters.length) return;

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                animateCounter(entry.target);
                observer.unobserve(entry.target);
            }
        });
    }, { threshold: 0.5 });

    counters.forEach(counter => observer.observe(counter));
}

function animateCounter(el) {
    const target = parseInt(el.dataset.target, 10);
    const duration = 1800;
    const start = performance.now();

    function update(now) {
        const elapsed = now - start;
        const progress = Math.min(elapsed / duration, 1);
        // easeOutCubic
        const eased = 1 - Math.pow(1 - progress, 3);
        el.textContent = Math.round(target * eased);
        if (progress < 1) {
            requestAnimationFrame(update);
        }
    }

    requestAnimationFrame(update);
}

// ── Load Exhibitions ──
async function loadExhibitions() {
    const grid = document.getElementById('exhibitions-grid');
    if (!grid) return;

    const exhibitions = await api.getAllExhibitions();

    // Update stat counter
    const statEl = document.getElementById('stat-exhibitions');
    if (statEl && exhibitions.length > 0) {
        statEl.dataset.target = exhibitions.length;
    }

    if (!exhibitions || exhibitions.length === 0) {
        grid.innerHTML = '<div class="empty-state">Экспозиции пока не добавлены</div>';
        return;
    }

    const icons = ['🏛', '⭐', '📚', '🎨', '🗿', '🔬', '🌍', '🎭'];

    grid.innerHTML = exhibitions.map((ex, i) => {
        const exhibitCount = (ex.exhibits || []).length;
        const icon = icons[i % icons.length];
        return `
            <div class="exhibition-card" onclick="openExhibition('${ex.id}')">
                <div class="exhibition-card-image">
                    <span class="exhibition-card-icon">${icon}</span>
                    ${exhibitCount > 0 ? `<span class="exhibition-card-count">${exhibitCount} экспонатов</span>` : ''}
                </div>
                <div class="exhibition-card-body">
                    <h3 class="exhibition-card-title">${ex.title}</h3>
                    <p class="exhibition-card-desc">${truncateText(ex.description || '', 140)}</p>
                    <span class="exhibition-card-link">Подробнее →</span>
                </div>
            </div>
        `;
    }).join('');

    // Re-run scroll animations for new elements
    requestAnimationFrame(() => initScrollAnimations());
}

function openExhibition(id) {
    window.location.href = `exhibition.html?id=${id}`;
}

// ── Load News ──
async function loadNews() {
    const container = document.getElementById('news-container');
    if (!container) return;

    const newsData = await api.getAllNews();

    if (!newsData || newsData.length === 0) {
        container.innerHTML = '<div class="empty-state">Новости пока отсутствуют</div>';
        return;
    }

    container.innerHTML = newsData.map(news => `
        <div class="news-card" onclick="openNewsModal('${news.id}')">
            ${news.image_url
                ? `<img src="${news.image_url}" alt="${news.title}" class="news-card-image">`
                : `<div class="news-card-image-placeholder"><span>📰</span></div>`
            }
            <div class="news-card-body">
                <div class="news-card-date">${formatDate(news.created_at)}</div>
                <h3 class="news-card-title">${news.title}</h3>
                <p class="news-card-text">${truncateText(news.content, 150)}</p>
            </div>
        </div>
    `).join('');

    // Re-run scroll animations for new elements
    requestAnimationFrame(() => initScrollAnimations());
}

// ── News Modal ──
let newsCache = [];

function initNewsModal() {
    const modal = document.getElementById('news-modal');
    const closeBtn = document.getElementById('news-modal-close');
    if (!modal || !closeBtn) return;

    closeBtn.addEventListener('click', () => closeNewsModal());
    modal.addEventListener('click', (e) => {
        if (e.target === modal) closeNewsModal();
    });
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') closeNewsModal();
    });
}

async function openNewsModal(id) {
    const modal = document.getElementById('news-modal');
    const body = document.getElementById('news-modal-body');
    if (!modal || !body) return;

    const news = await api.getNewsById(id);
    if (!news) return;

    body.innerHTML = `
        ${news.image_url ? `<img src="${news.image_url}" alt="${news.title}" class="modal-image">` : ''}
        <h2 class="modal-title">${news.title}</h2>
        <div class="news-card-date" style="margin-bottom: 16px;">${formatDate(news.created_at)}</div>
        <p class="modal-description">${news.content || 'Содержание новости отсутствует'}</p>
    `;

    modal.classList.add('active');
    document.body.style.overflow = 'hidden';
}

function closeNewsModal() {
    const modal = document.getElementById('news-modal');
    if (modal) {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    }
}
