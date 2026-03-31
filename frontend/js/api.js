// ═══════════════════════════════════════════════
// API — взаимодействие с бэкендом
// ═══════════════════════════════════════════════

const API_BASE_URL = '/museum';

function escapeHtml(value) {
    return String(value ?? '')
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

function sanitizeUrl(value) {
    if (!value || typeof value !== 'string') return '';
    const trimmed = value.trim();
    if (!trimmed) return '';

    try {
        const parsed = new URL(trimmed, window.location.origin);
        if (parsed.protocol === 'http:' || parsed.protocol === 'https:') {
            return parsed.href;
        }
    } catch (_) {
        return '';
    }

    return '';
}

function isVideoUrl(value) {
    const safe = sanitizeUrl(value);
    if (!safe) return false;

    try {
        const pathname = new URL(safe).pathname.toLowerCase();
        return /\.(mp4|webm|ogg|ogv|mov|m4v)$/i.test(pathname);
    } catch (_) {
        return false;
    }
}

function parseExternalEmbed(value) {
    const safe = sanitizeUrl(value);
    if (!safe) return null;

    let url;
    try {
        url = new URL(safe);
    } catch (_) {
        return null;
    }

    const host = url.hostname.replace(/^www\./, '').toLowerCase();
    const path = url.pathname;

    if (host === 'imgur.com') {
        const albumMatch = path.match(/^\/(a|gallery)\/([a-zA-Z0-9]+)\/?$/);
        if (albumMatch) {
            return {
                provider: 'imgur',
                source: safe
            };
        }

        const imageMatch = path.match(/^\/([a-zA-Z0-9]{5,})\/?$/);
        if (imageMatch) {
            return {
                provider: 'imgur',
                source: safe
            };
        }
    }

    if (host === 'vk.com') {
        const pathVideoMatch = path.match(/^\/video(-?\d+)_(\d+)\/?$/);
        if (pathVideoMatch) {
            const oid = pathVideoMatch[1];
            const id = pathVideoMatch[2];
            return {
                provider: 'vk',
                base: `https://vk.com/video_ext.php?oid=${oid}&id=${id}&hd=2`
            };
        }

        if (path === '/video_ext.php') {
            const oid = url.searchParams.get('oid');
            const id = url.searchParams.get('id');
            if (oid && id) {
                return {
                    provider: 'vk',
                    base: `https://vk.com/video_ext.php?oid=${oid}&id=${id}&hd=2`
                };
            }
        }
    }

    return null;
}

function buildExternalEmbedSrc(base, provider, autoplay = false) {
    let parsed;
    try {
        parsed = new URL(base);
    } catch (_) {
        return '';
    }

    if (provider === 'vk') {
        parsed.searchParams.set('autoplay', autoplay ? '1' : '0');
        parsed.searchParams.set('js_api', '0');
    }

    return parsed.toString();
}

function normalizeMediaUrls(items) {
    if (!Array.isArray(items)) return [];
    return items
        .map(sanitizeUrl)
        .filter(Boolean);
}

function buildCardMedia(url, altText, imageClass, videoClass, embedClass = imageClass) {
    const safeUrl = sanitizeUrl(url);
    if (!safeUrl) return '';

    if (isVideoUrl(safeUrl)) {
        return `<video class="${videoClass}" src="${safeUrl}" muted loop playsinline preload="metadata"></video>`;
    }

    const external = parseExternalEmbed(safeUrl);
    if (external) {
        if (external.provider === 'imgur') {
            return `<div class="${embedClass} external-media-loading" data-imgur-source="${escapeHtml(external.source || safeUrl)}" data-imgur-image-class="${escapeHtml(imageClass)}" data-imgur-video-class="${escapeHtml(videoClass)}" data-imgur-embed-class="${escapeHtml(embedClass)}"><span class="external-media-loading__label">Загрузка медиа...</span></div>`;
        }

        const initialSrc = buildExternalEmbedSrc(external.base, external.provider, false);
        return `<iframe class="${embedClass}" src="${escapeHtml(initialSrc)}" data-embed-provider="${external.provider}" data-embed-base="${escapeHtml(external.base)}" loading="lazy" referrerpolicy="strict-origin-when-cross-origin" allow="autoplay; fullscreen; picture-in-picture; encrypted-media" allowfullscreen title="${escapeHtml(altText || 'Медиа')}" tabindex="-1"></iframe>`;
    }

    return `<img src="${safeUrl}" alt="${escapeHtml(altText)}" class="${imageClass}">`;
}

async function resolveImgurMedia(source) {
    const safe = sanitizeUrl(source);
    if (!safe) return null;

    try {
        const response = await fetch(`${API_BASE_URL}/media/resolve?url=${encodeURIComponent(safe)}`);
        if (!response.ok) return null;
        const data = await response.json();
        if (!data || !data.url) return null;
        return {
            url: sanitizeUrl(data.url),
            type: String(data.type || '').toLowerCase()
        };
    } catch (_) {
        return null;
    }
}

async function hydrateImgurMedia(scope = document) {
    const placeholders = Array.from(scope.querySelectorAll('[data-imgur-source]'));
    if (placeholders.length === 0) return;

    await Promise.all(placeholders.map(async (node) => {
        if (node.dataset.imgurHydrated === '1') return;
        node.dataset.imgurHydrated = '1';

        const source = node.dataset.imgurSource || '';
        const resolved = await resolveImgurMedia(source);
        if (!resolved || !resolved.url) {
            node.classList.remove('external-media-loading');
            node.classList.add('external-media-error');
            node.innerHTML = '<span class="external-media-loading__label">Медиа недоступно</span>';
            return;
        }

        const imageClass = node.dataset.imgurImageClass || '';
        const videoClass = node.dataset.imgurVideoClass || '';
        const embedClass = node.dataset.imgurEmbedClass || '';
        const mediaType = resolved.type || (isVideoUrl(resolved.url) ? 'video' : 'image');

        if (mediaType === 'video') {
            node.outerHTML = `<video class="${videoClass || embedClass}" src="${resolved.url}" muted loop playsinline preload="metadata"></video>`;
        } else {
            node.outerHTML = `<img src="${resolved.url}" alt="Media" class="${imageClass || embedClass}">`;
        }
    }));
}

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
    const media = normalizeMediaUrls(imgs);
    if (media.length === 0) return '';

    if (media.length === 1) {
        return `<div class="modal-carousel">
            <div class="modal-carousel-track">
                <div class="modal-carousel-slide">
                    ${buildCardMedia(media[0], altText || '', 'modal-carousel-media', 'modal-carousel-media', 'modal-carousel-embed')}
                </div>
            </div>
        </div>`;
    }

    const slides = media.map(url =>
        `<div class="modal-carousel-slide">${buildCardMedia(url, altText || '', 'modal-carousel-media', 'modal-carousel-media', 'modal-carousel-embed')}</div>`
    ).join('');

    const dots = media.map((_, i) =>
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
    if (slides.length === 0) return;

    let current = 0;
    const total = slides.length;
    let isHovered = false;

    const prevBtn = carousel.querySelector('.modal-carousel-arrow--prev');
    const nextBtn = carousel.querySelector('.modal-carousel-arrow--next');
    const dots    = carousel.querySelectorAll('.modal-carousel-dot');

    function mediaInSlide(index) {
        const slide = slides[index];
        if (!slide) return null;
        return slide.querySelector('video, iframe[data-embed-base][data-embed-provider]');
    }

    function startMedia(node) {
        if (!node) return;

        if (node.tagName === 'VIDEO') {
            node.currentTime = 0;
            node.play().catch(() => {});
            return;
        }

        if (node.tagName === 'IFRAME') {
            const provider = node.dataset.embedProvider;
            const base = node.dataset.embedBase;
            const src = buildExternalEmbedSrc(base, provider, true);
            if (src && node.src !== src) {
                node.src = src;
            }
        }
    }

    function stopMedia(node) {
        if (!node) return;

        if (node.tagName === 'VIDEO') {
            node.pause();
            node.currentTime = 0;
            return;
        }

        if (node.tagName === 'IFRAME') {
            const provider = node.dataset.embedProvider;
            const base = node.dataset.embedBase;
            const src = buildExternalEmbedSrc(base, provider, false);
            if (src && node.src !== src) {
                node.src = src;
            }
        }
    }

    function stopAllMedia() {
        slides.forEach((slide) => {
            const node = slide.querySelector('video, iframe[data-embed-base][data-embed-provider]');
            stopMedia(node);
        });
    }

    function syncMedia() {
        stopAllMedia();
        if (!isHovered) return;
        startMedia(mediaInSlide(current));
    }

    function goTo(index) {
        if (index < 0) index = total - 1;
        if (index >= total) index = 0;
        current = index;
        track.style.transform = `translateX(-${current * 100}%)`;

        dots.forEach((d, i) => d.classList.toggle('active', i === current));
        syncMedia();
    }

    if (prevBtn) prevBtn.addEventListener('click', () => goTo(current - 1));
    if (nextBtn) nextBtn.addEventListener('click', () => goTo(current + 1));

    dots.forEach(dot => {
        dot.addEventListener('click', () => goTo(Number(dot.dataset.index)));
    });

    carousel.addEventListener('mouseenter', () => {
        isHovered = true;
        syncMedia();
    });

    carousel.addEventListener('mouseleave', () => {
        isHovered = false;
        syncMedia();
    });
}


