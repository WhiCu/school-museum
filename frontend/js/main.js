// ═══════════════════════════════════════════════
// Главная страница — JavaScript
// ═══════════════════════════════════════════════

document.addEventListener('DOMContentLoaded', () => {
    initHeader();
    initBurger();
    initScrollAnimations();
    loadNewsHighlight();
    loadExhibitions();
    initNewsModal();
    initYandexMap();
    trackVisit();
});

// ── Track page visit ──
function trackVisit() {
    const data = {
        page: window.location.pathname,
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

    const selectors = [
        '.about-content',
        '.team-card',
        '.visit-card',
        '.contacts-inner',
        '.carousel'
    ];

    requestAnimationFrame(() => {
        selectors.forEach(sel => {
            document.querySelectorAll(sel).forEach(el => {
                el.classList.add('fade-in');
                observer.observe(el);
            });
        });
    });
}

// ── Load news highlight (dark strip) ──
async function loadNewsHighlight() {
    const track = document.getElementById('news-hl-track');
    if (!track) return;

    const newsData = await api.getAllNews();
    if (!newsData || newsData.length === 0) {
        track.innerHTML = '<div class="empty-state" style="color:rgba(255,255,255,.5)">Новости пока отсутствуют</div>';
        return;
    }

    track.innerHTML = newsData.map(n => {
        const imgs = n.image_urls || [];
        const firstImg = imgs.length > 0 ? imgs[0] : '';
        return `
        <div class="news-hl-card" onclick="openNewsModal('${n.id}')">
            ${firstImg
                ? `<img src="${firstImg}" alt="${n.title}" class="news-hl-image">`
                : `<div class="news-hl-image-placeholder"><span>📰</span></div>`
            }
            <div class="news-hl-body">
                <div class="news-card-date">${formatDate(n.created_at)}</div>
                <h3 class="news-hl-title">${n.title}</h3>
            </div>
        </div>
    `}).join('');

    const ctrl = initCarousel('news-hl-carousel', 'news-hl-dots');
    if (ctrl) ctrl.refresh();
}

// ═══════════════════════════════════════════════
// Carousel controller (generic)
// ═══════════════════════════════════════════════

function initCarousel(carouselId, dotsId, options = {}) {
    const carousel = document.getElementById(carouselId);
    if (!carousel) return null;

    const track = carousel.querySelector('.carousel-track');
    const dotsContainer = document.getElementById(dotsId);
    const prevBtn = carousel.querySelector('.carousel-arrow--prev');
    const nextBtn = carousel.querySelector('.carousel-arrow--next');

    const fixedPerView = options.perView || null; // null = responsive

    let currentIndex = 0;
    let perView = fixedPerView || getPerView();
    let totalSlides = 0;
    let maxIndex = 0;
    let autoTimer = null;

    function getPerView() {
        if (fixedPerView) return fixedPerView;
        const w = window.innerWidth;
        if (w <= 600) return 1;
        if (w <= 960) return 2;
        return 3;
    }

    function recalc() {
        const cards = track.children;
        totalSlides = cards.length;
        perView = getPerView();
        maxIndex = Math.max(0, totalSlides - perView);
        if (currentIndex > maxIndex) currentIndex = maxIndex;
    }

    function slide(animate = true) {
        if (!track.children.length) return;
        const card = track.children[0];
        const gap = 28;
        const cardWidth = card.offsetWidth + gap;
        const offset = currentIndex * cardWidth;
        track.style.transition = animate ? 'transform 0.5s cubic-bezier(0.4, 0, 0.2, 1)' : 'none';
        track.style.transform = `translateX(-${offset}px)`;
        updateArrows();
        updateDots();
    }

    function updateArrows() {
        if (prevBtn) prevBtn.disabled = currentIndex <= 0;
        if (nextBtn) nextBtn.disabled = currentIndex >= maxIndex;
    }

    function updateDots() {
        if (!dotsContainer) return;
        const pages = maxIndex + 1;
        // Rebuild dots if count changed
        if (dotsContainer.children.length !== pages) {
            dotsContainer.innerHTML = '';
            for (let i = 0; i < pages; i++) {
                const dot = document.createElement('button');
                dot.className = 'carousel-dot';
                dot.setAttribute('aria-label', `Страница ${i + 1}`);
                dot.addEventListener('click', () => goTo(i));
                dotsContainer.appendChild(dot);
            }
        }
        Array.from(dotsContainer.children).forEach((d, i) => {
            d.classList.toggle('active', i === currentIndex);
        });
    }

    function goTo(index) {
        currentIndex = Math.max(0, Math.min(index, maxIndex));
        slide();
        resetAuto();
    }

    function next() { goTo(currentIndex + 1); }
    function prev() { goTo(currentIndex - 1); }

    // Arrows
    if (prevBtn) prevBtn.addEventListener('click', prev);
    if (nextBtn) nextBtn.addEventListener('click', next);

    // Touch / swipe
    let startX = 0, startY = 0, dx = 0, swiping = false;

    track.addEventListener('touchstart', (e) => {
        startX = e.touches[0].clientX;
        startY = e.touches[0].clientY;
        dx = 0;
        swiping = true;
        track.style.transition = 'none';
    }, { passive: true });

    track.addEventListener('touchmove', (e) => {
        if (!swiping) return;
        dx = e.touches[0].clientX - startX;
        const dy = Math.abs(e.touches[0].clientY - startY);
        if (dy > Math.abs(dx)) { swiping = false; return; }
        const card = track.children[0];
        const gap = 28;
        const cardWidth = card.offsetWidth + gap;
        const base = currentIndex * cardWidth;
        track.style.transform = `translateX(-${base - dx}px)`;
    }, { passive: true });

    track.addEventListener('touchend', () => {
        if (!swiping) { slide(); return; }
        swiping = false;
        const threshold = 50;
        if (dx < -threshold) next();
        else if (dx > threshold) prev();
        else slide();
    });

    // Auto-play
    function startAuto() {
        stopAuto();
        autoTimer = setInterval(() => {
            if (currentIndex >= maxIndex) goTo(0);
            else next();
        }, 5000);
    }

    function stopAuto() {
        if (autoTimer) { clearInterval(autoTimer); autoTimer = null; }
    }

    function resetAuto() {
        stopAuto();
        startAuto();
    }

    carousel.addEventListener('mouseenter', stopAuto);
    carousel.addEventListener('mouseleave', startAuto);

    // Window resize
    let resizeTimer;
    window.addEventListener('resize', () => {
        clearTimeout(resizeTimer);
        resizeTimer = setTimeout(() => {
            recalc();
            slide(false);
        }, 150);
    });

    return {
        refresh() {
            recalc();
            slide(false);
            startAuto();
        }
    };
}

// ── Load Exhibitions ──
async function loadExhibitions() {
    const grid = document.getElementById('exhibitions-grid');
    if (!grid) return;

    const exhibitions = await api.getAllExhibitions();

    if (!exhibitions || exhibitions.length === 0) {
        grid.innerHTML = '<div class="empty-state">Экспозиции пока не добавлены</div>';
        return;
    }

    const icons = ['🏛', '⭐', '📚', '🎨', '🗿', '🔬', '🌍', '🎭'];

    grid.innerHTML = exhibitions.map((ex, i) => {
        const exhibitCount = (ex.exhibits || []).length;
        const icon = icons[i % icons.length];
        // Find preview image from the selected preview exhibit
        let previewImg = '';
        if (ex.preview_exhibit_id && ex.exhibits) {
            const previewExhibit = ex.exhibits.find(e => e.id === ex.preview_exhibit_id);
            if (previewExhibit) {
                const imgs = previewExhibit.image_urls || [];
                if (imgs.length > 0) previewImg = imgs[0];
            }
        }
        return `
            <div class="exhibition-card" onclick="openExhibition('${ex.id}')">
                <div class="exhibition-card-image">
                    ${previewImg
                        ? `<img src="${previewImg}" alt="${ex.title}" class="exhibition-card-preview-img">`
                        : `<span class="exhibition-card-icon">${icon}</span>`
                    }
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

    const ctrl = initCarousel('exhibitions-carousel', 'exhibitions-dots', { perView: 1 });
    if (ctrl) ctrl.refresh();
}

function openExhibition(id) {
    window.location.href = `exhibition.html?id=${id}`;
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

    const imgs = news.image_urls || [];
    const imagesHtml = buildImageCarousel(imgs, news.title);

    body.innerHTML = `
        ${imagesHtml}
        <h2 class="modal-title">${news.title}</h2>
        <div class="news-card-date" style="margin-bottom: 16px;">${formatDate(news.created_at)}</div>
        <p class="modal-description">${news.content || 'Содержание новости отсутствует'}</p>
    `;

    initModalCarousel(body);

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

// ── Yandex Map ──
function initYandexMap() {
    if (typeof ymaps === 'undefined') return;

    ymaps.ready(() => {
        const map = new ymaps.Map('yandex-map', {
            center: [57.581944, 39.839645], // Ярославль, ул. Зелинского, 6
            zoom: 16,
            controls: ['zoomControl', 'geolocationControl']
        });

        const placemark = new ymaps.Placemark([57.581944, 39.839645], {
            hintContent: 'Лицей №86',
            balloonContentHeader: 'Музей «Страницы истории»',
            balloonContentBody: 'г. Ярославль, ул. Зелинского, 6<br>Лицей №86',
            balloonContentFooter: '+7 (905) 646-41-27'
        }, {
            preset: 'islands#redEducationIcon'
        });

        map.geoObjects.add(placemark);
        map.behaviors.disable('scrollZoom');
    });
}
