// ═══════════════════════════════════════════════
// Страница экспозиции — JavaScript
// ═══════════════════════════════════════════════

let currentExhibition = null;

document.addEventListener('DOMContentLoaded', () => {
    initHeader();
    initBurger();
    initExhibitModal();
    trackVisit();

    const urlParams = new URLSearchParams(window.location.search);
    const exhibitionId = urlParams.get('id');

    if (!exhibitionId) {
        showError('ID экспозиции не указан');
        return;
    }

    loadExhibition(exhibitionId);
});

// ── Header scroll ──
function initHeader() {
    const header = document.getElementById('header');
    if (!header) return;
    const onScroll = () => header.classList.toggle('scrolled', window.scrollY > 50);
    window.addEventListener('scroll', onScroll, { passive: true });
    onScroll();
}

// ── Burger ──
function initBurger() {
    const burger = document.getElementById('burger');
    const nav = document.getElementById('nav');
    if (!burger || !nav) return;
    burger.addEventListener('click', () => {
        burger.classList.toggle('active');
        nav.classList.toggle('open');
    });
    nav.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', () => {
            burger.classList.remove('active');
            nav.classList.remove('open');
        });
    });
}

// ── Load exhibition ──
async function loadExhibition(id) {
    currentExhibition = await api.getExhibitionById(id);

    if (!currentExhibition) {
        showError('Экспозиция не найдена');
        return;
    }

    // Update page title
    document.title = `${currentExhibition.title} — Музей «Страницы истории»`;

    // Hero
    const titleEl = document.getElementById('ex-title');
    const descEl = document.getElementById('ex-description');
    if (titleEl) titleEl.textContent = currentExhibition.title;
    if (descEl) descEl.textContent = currentExhibition.description || '';

    // About block
    const aboutBlock = document.getElementById('ex-about-block');
    const aboutText = document.getElementById('ex-about-text');
    if (currentExhibition.description && aboutBlock && aboutText) {
        aboutText.textContent = currentExhibition.description;
        aboutBlock.style.display = 'block';
    }

    // Exhibits
    renderExhibits(currentExhibition.exhibits || []);
}

// ── Render exhibits ──
function renderExhibits(exhibits) {
    const grid = document.getElementById('exhibits-grid');
    if (!grid) return;

    if (!exhibits || exhibits.length === 0) {
        grid.innerHTML = '<div class="empty-state">Экспонаты пока не добавлены</div>';
        return;
    }

    grid.innerHTML = exhibits.map(exhibit => {
        const imgs = exhibit.image_urls || [];
        const firstImg = imgs.length > 0 ? imgs[0] : '';
        return `
        <div class="exhibit-card" onclick="openExhibitModal('${exhibit.id}')">
            <div class="exhibit-card-image">
                ${firstImg
                    ? `<img src="${firstImg}" alt="${exhibit.title}">`
                    : `<span class="exhibit-card-placeholder">✻</span>`
                }
                ${imgs.length > 1 ? `<span class="exhibit-card-count">${imgs.length} фото</span>` : ''}
            </div>
            <div class="exhibit-card-body">
                <h4 class="exhibit-card-title">${exhibit.title}</h4>
                <p class="exhibit-card-desc">${truncateText(exhibit.description || '', 100)}</p>
            </div>
        </div>
    `}).join('');

    // Scroll animations
    initExhibitAnimations();
}

function initExhibitAnimations() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
                observer.unobserve(entry.target);
            }
        });
    }, { threshold: 0.1, rootMargin: '0px 0px -40px 0px' });

    document.querySelectorAll('.exhibit-card').forEach(el => {
        el.classList.add('fade-in');
        observer.observe(el);
    });
}

// ── Exhibit Modal ──
function initExhibitModal() {
    const modal = document.getElementById('exhibit-modal');
    const closeBtn = document.getElementById('exhibit-modal-close');
    if (!modal || !closeBtn) return;

    closeBtn.addEventListener('click', () => closeExhibitModal());
    modal.addEventListener('click', (e) => {
        if (e.target === modal) closeExhibitModal();
    });
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') closeExhibitModal();
    });
}

function openExhibitModal(exhibitId) {
    if (!currentExhibition) return;
    const exhibits = currentExhibition.exhibits || [];
    const exhibit = exhibits.find(e => e.id === exhibitId);
    if (!exhibit) return;

    const modal = document.getElementById('exhibit-modal');
    const body = document.getElementById('exhibit-modal-body');
    if (!modal || !body) return;

    const imgs = exhibit.image_urls || [];

    body.innerHTML = `
        ${buildImageCarousel(imgs, exhibit.title)}
        <h2 class="modal-title">${exhibit.title}</h2>
        <p class="modal-description">${exhibit.description || 'Описание экспоната отсутствует'}</p>
    `;

    initModalCarousel(body);

    modal.classList.add('active');
    document.body.style.overflow = 'hidden';
}

function closeExhibitModal() {
    const modal = document.getElementById('exhibit-modal');
    if (modal) {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    }
}

// ── Track page visit ──
function trackVisit() {
    const data = {
        page: window.location.pathname + window.location.search,
        referrer: document.referrer || '',
        screen_width: window.screen.width || 0,
        screen_height: window.screen.height || 0,
        language: navigator.language || navigator.userLanguage || ''
    };

    fetch('/museum/visit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    }).catch(() => {});
}

// ── Error ──
function showError(message) {
    const titleEl = document.getElementById('ex-title');
    const descEl = document.getElementById('ex-description');
    const grid = document.getElementById('exhibits-grid');

    if (titleEl) titleEl.textContent = 'Ошибка';
    if (descEl) {
        descEl.innerHTML = `${message} <a href="index.html" style="color:#c9a96e;">← Вернуться на главную</a>`;
    }
    if (grid) grid.innerHTML = '';
}
