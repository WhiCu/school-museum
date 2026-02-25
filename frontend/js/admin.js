// ═══════════════════════════════════════════════
// Админ-панель — JavaScript
// ═══════════════════════════════════════════════

const ADMIN_API = '/admin';
const MUSEUM_API = '/museum';

// ==================== АВТОРИЗАЦИЯ ====================

function getAuthHeader() {
    const creds = sessionStorage.getItem('admin_auth');
    return creds ? 'Basic ' + creds : '';
}

function isLoggedIn() {
    return !!sessionStorage.getItem('admin_auth');
}

function showAdmin() {
    document.getElementById('login-screen').style.display = 'none';
    document.getElementById('admin-main').classList.remove('admin-hidden');
    document.getElementById('admin-layout').classList.remove('admin-hidden');
}

function showLogin() {
    document.getElementById('login-screen').style.display = '';
    document.getElementById('admin-main').classList.add('admin-hidden');
    document.getElementById('admin-layout').classList.add('admin-hidden');
    sessionStorage.removeItem('admin_auth');
}

// Login form
document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const login = document.getElementById('login-username').value.trim();
    const password = document.getElementById('login-password').value;
    const errorEl = document.getElementById('login-error');

    errorEl.style.display = 'none';

    const creds = btoa(login + ':' + password);

    try {
        const resp = await fetch(`${ADMIN_API}/ping?message=ping`, {
            headers: { 'Authorization': 'Basic ' + creds }
        });

        if (resp.ok) {
            sessionStorage.setItem('admin_auth', creds);
            showAdmin();
            initAdminData();
        } else {
            errorEl.textContent = 'Неверный логин или пароль';
            errorEl.style.display = 'block';
        }
    } catch (err) {
        errorEl.textContent = 'Ошибка подключения к серверу';
        errorEl.style.display = 'block';
    }
});

// Logout
document.getElementById('btn-logout').addEventListener('click', () => {
    showLogin();
});

// Check if already authenticated
if (isLoggedIn()) {
    // Verify session is still valid
    fetch(`${ADMIN_API}/ping?message=ping`, {
        headers: { 'Authorization': getAuthHeader() }
    }).then(resp => {
        if (resp.ok) {
            showAdmin();
            initAdminData();
        } else {
            showLogin();
        }
    }).catch(() => showLogin());
}

function initAdminData() {
    loadExhibitions();
    loadAllExhibits();
    loadNews();
    loadAdminStats();
}

// ==================== Навигация ====================

document.querySelectorAll('.sidebar-item').forEach(item => {
    item.addEventListener('click', () => {
        document.querySelectorAll('.sidebar-item').forEach(i => i.classList.remove('active'));
        document.querySelectorAll('.admin-section').forEach(s => s.classList.remove('active'));
        item.classList.add('active');
        const section = item.getAttribute('data-section');
        document.getElementById('section-' + section).classList.add('active');
    });
});

// ==================== Общие утилиты ====================

async function apiRequest(url, method = 'GET', body = null) {
    const opts = { method, headers: {} };

    // Add auth header for admin API calls
    if (url.startsWith(ADMIN_API)) {
        opts.headers['Authorization'] = getAuthHeader();
    }

    if (body) {
        opts.headers['Content-Type'] = 'application/json';
        opts.body = JSON.stringify(body);
    }
    const resp = await fetch(url, opts);

    if (resp.status === 401) {
        showLogin();
        throw new Error('Сессия истекла, войдите снова');
    }

    if (method === 'DELETE') {
        if (!resp.ok) throw new Error('Ошибка удаления');
        return null;
    }
    if (!resp.ok) {
        const text = await resp.text();
        throw new Error(text || 'Ошибка запроса');
    }
    return resp.json();
}

function formatDate(dateStr) {
    if (!dateStr) return '';
    const d = new Date(dateStr);
    return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' });
}

function truncate(text, len = 100) {
    if (!text || text.length <= len) return text || '';
    return text.substring(0, len) + '...';
}

// ==================== ЭКСПОЗИЦИИ ====================

let exhibitionsCache = [];

async function loadExhibitions() {
    const container = document.getElementById('exhibitions-list');
    try {
        const data = await apiRequest(`${MUSEUM_API}/exhibitions`);
        exhibitionsCache = Array.isArray(data) ? data : (data && data.exhibitions ? data.exhibitions : []);

        if (exhibitionsCache.length === 0) {
            container.innerHTML = '<div class="empty-state">Экспозиций пока нет</div>';
            return;
        }

        container.innerHTML = exhibitionsCache.map(ex => `
            <div class="item-card">
                <div class="item-info">
                    <h3 class="item-title">${ex.title}</h3>
                    <p class="item-desc">${truncate(ex.description, 120)}</p>
                    <span class="item-meta">${(ex.exhibits || []).length} экспонатов</span>
                </div>
                <div class="item-actions">
                    <button class="btn btn-small btn-edit" onclick="showExhibitionForm('${ex.id}')">✏️</button>
                    <button class="btn btn-small btn-delete" onclick="deleteExhibition('${ex.id}')">🗑️</button>
                </div>
            </div>
        `).join('');
    } catch (e) {
        container.innerHTML = '<div class="error-state">Ошибка загрузки: ' + e.message + '</div>';
    }
}

function showExhibitionForm(id = null) {
    const ex = id ? exhibitionsCache.find(e => e.id === id) : null;
    const isEdit = !!ex;

    document.getElementById('modal-title').textContent = isEdit ? 'Редактировать экспозицию' : 'Новая экспозиция';
    document.getElementById('modal-body').innerHTML = `
        <form id="exhibition-form" onsubmit="saveExhibition(event, '${id || ''}')">
            <div class="form-group">
                <label>Название *</label>
                <input type="text" id="ex-title" value="${isEdit ? ex.title : ''}" required>
            </div>
            <div class="form-group">
                <label>Описание</label>
                <textarea id="ex-description" rows="4">${isEdit ? (ex.description || '') : ''}</textarea>
            </div>
            <div class="form-actions">
                <button type="button" class="btn btn-secondary" onclick="closeModal()">Отмена</button>
                <button type="submit" class="btn btn-primary">${isEdit ? 'Сохранить' : 'Создать'}</button>
            </div>
        </form>
    `;
    openModal();
}

async function saveExhibition(event, id) {
    event.preventDefault();
    const body = {
        title: document.getElementById('ex-title').value.trim(),
        description: document.getElementById('ex-description').value.trim()
    };
    try {
        if (id) {
            await apiRequest(`${ADMIN_API}/exhibitions/${id}`, 'PUT', body);
        } else {
            await apiRequest(`${ADMIN_API}/exhibitions`, 'POST', body);
        }
        closeModal();
        await loadExhibitions();
        await loadAllExhibits(); // обновить выпадающий список
    } catch (e) {
        alert('Ошибка: ' + e.message);
    }
}

async function deleteExhibition(id) {
    if (!confirm('Удалить экспозицию и все её экспонаты?')) return;
    try {
        await apiRequest(`${ADMIN_API}/exhibitions/${id}`, 'DELETE');
        await loadExhibitions();
        await loadAllExhibits();
    } catch (e) {
        alert('Ошибка удаления: ' + e.message);
    }
}

// ==================== ЭКСПОНАТЫ ====================

let allExhibits = [];

async function loadAllExhibits() {
    const container = document.getElementById('exhibits-list');
    try {
        // Загружаем все экспозиции с экспонатами
        const data = await apiRequest(`${MUSEUM_API}/exhibitions`);
        const exhibitions = Array.isArray(data) ? data : (data && data.exhibitions ? data.exhibitions : []);

        allExhibits = [];
        exhibitions.forEach(ex => {
            (ex.exhibits || []).forEach(exhibit => {
                allExhibits.push({ ...exhibit, exhibition_id: ex.id, exhibition_title: ex.title });
            });
        });

        if (allExhibits.length === 0 && exhibitions.length === 0) {
            container.innerHTML = '<div class="empty-state">Экспонатов пока нет</div>';
            return;
        }

        container.innerHTML = exhibitions.map(ex => {
            const exhibits = allExhibits.filter(e => e.exhibition_id === ex.id);
            return `
                <div class="exhibit-group">
                    <div class="exhibit-group-header" onclick="toggleGroup(this)">
                        <span class="group-toggle">▾</span>
                        <h3 class="group-title">🏛️ ${ex.title}</h3>
                        <span class="group-count">${exhibits.length} экспонатов</span>
                    </div>
                    <div class="exhibit-group-items">
                        ${exhibits.length === 0
                            ? '<div class="empty-state" style="padding:16px;font-size:14px;">Нет экспонатов</div>'
                            : exhibits.map(item => `
                                <div class="item-card">
                                    <div class="item-info">
                                        ${item.image_url ? `<img src="${item.image_url}" class="item-thumb" alt="">` : ''}
                                        <div>
                                            <h3 class="item-title">${item.title}</h3>
                                            <p class="item-desc">${truncate(item.description, 100)}</p>
                                        </div>
                                    </div>
                                    <div class="item-actions">
                                        <button class="btn btn-small btn-edit" onclick="showExhibitForm('${item.id}')">✏️</button>
                                        <button class="btn btn-small btn-delete" onclick="deleteExhibit('${item.id}')">🗑️</button>
                                    </div>
                                </div>
                            `).join('')
                        }
                    </div>
                </div>
            `;
        }).join('');
    } catch (e) {
        container.innerHTML = '<div class="error-state">Ошибка загрузки: ' + e.message + '</div>';
    }
}

function showExhibitForm(id = null) {
    const item = id ? allExhibits.find(e => e.id === id) : null;
    const isEdit = !!item;

    const exhibitionOptions = exhibitionsCache.map(ex =>
        `<option value="${ex.id}" ${item && item.exhibition_id === ex.id ? 'selected' : ''}>${ex.title}</option>`
    ).join('');

    document.getElementById('modal-title').textContent = isEdit ? 'Редактировать экспонат' : 'Новый экспонат';
    document.getElementById('modal-body').innerHTML = `
        <form id="exhibit-form" onsubmit="saveExhibit(event, '${id || ''}')">
            ${!isEdit ? `
            <div class="form-group">
                <label>Экспозиция *</label>
                <select id="exhibit-exhibition" required>
                    <option value="">Выберите экспозицию</option>
                    ${exhibitionOptions}
                </select>
            </div>
            ` : ''}
            <div class="form-group">
                <label>Название *</label>
                <input type="text" id="exhibit-title" value="${isEdit ? item.title : ''}" required>
            </div>
            <div class="form-group">
                <label>Описание</label>
                <textarea id="exhibit-description" rows="4">${isEdit ? (item.description || '') : ''}</textarea>
            </div>
            <div class="form-group">
                <label>URL изображения</label>
                <input type="url" id="exhibit-image" value="${isEdit ? (item.image_url || '') : ''}" placeholder="https://...">
            </div>
            <div class="form-actions">
                <button type="button" class="btn btn-secondary" onclick="closeModal()">Отмена</button>
                <button type="submit" class="btn btn-primary">${isEdit ? 'Сохранить' : 'Создать'}</button>
            </div>
        </form>
    `;
    openModal();
}

async function saveExhibit(event, id) {
    event.preventDefault();
    try {
        if (id) {
            const body = {
                title: document.getElementById('exhibit-title').value.trim(),
                description: document.getElementById('exhibit-description').value.trim(),
                image_url: document.getElementById('exhibit-image').value.trim()
            };
            await apiRequest(`${ADMIN_API}/exhibits/${id}`, 'PUT', body);
        } else {
            const body = {
                exhibition_id: document.getElementById('exhibit-exhibition').value,
                title: document.getElementById('exhibit-title').value.trim(),
                description: document.getElementById('exhibit-description').value.trim(),
                image_url: document.getElementById('exhibit-image').value.trim()
            };
            await apiRequest(`${ADMIN_API}/exhibits`, 'POST', body);
        }
        closeModal();
        await loadAllExhibits();
        await loadExhibitions();
    } catch (e) {
        alert('Ошибка: ' + e.message);
    }
}

async function deleteExhibit(id) {
    if (!confirm('Удалить экспонат?')) return;
    try {
        await apiRequest(`${ADMIN_API}/exhibits/${id}`, 'DELETE');
        await loadAllExhibits();
        await loadExhibitions();
    } catch (e) {
        alert('Ошибка удаления: ' + e.message);
    }
}

// ==================== НОВОСТИ ====================

let newsCache = [];

async function loadNews() {
    const container = document.getElementById('news-list');
    try {
        const data = await apiRequest(`${MUSEUM_API}/news`);
        newsCache = Array.isArray(data) ? data : (data && data.news ? data.news : []);

        if (newsCache.length === 0) {
            container.innerHTML = '<div class="empty-state">Новостей пока нет</div>';
            return;
        }

        container.innerHTML = newsCache.map(n => `
            <div class="item-card">
                <div class="item-info">
                    ${n.image_url ? `<img src="${n.image_url}" class="item-thumb" alt="">` : ''}
                    <div>
                        <h3 class="item-title">${n.title}</h3>
                        <p class="item-desc">${truncate(n.content, 120)}</p>
                        <span class="item-meta">${formatDate(n.created_at)}</span>
                    </div>
                </div>
                <div class="item-actions">
                    <button class="btn btn-small btn-edit" onclick="showNewsForm('${n.id}')">✏️</button>
                    <button class="btn btn-small btn-delete" onclick="deleteNewsItem('${n.id}')">🗑️</button>
                </div>
            </div>
        `).join('');
    } catch (e) {
        container.innerHTML = '<div class="error-state">Ошибка загрузки: ' + e.message + '</div>';
    }
}

function showNewsForm(id) {
    const n = id ? newsCache.find(x => x.id === id) : null;
    const isEdit = !!n;

    document.getElementById('modal-title').textContent = isEdit ? 'Редактировать новость' : 'Новая новость';
    document.getElementById('modal-body').innerHTML = `
        <form id="news-form" onsubmit="saveNews(event, '${id || ''}')">
            <div class="form-group">
                <label>Заголовок *</label>
                <input type="text" id="news-title" value="${isEdit ? n.title : ''}" required>
            </div>
            <div class="form-group">
                <label>Содержание</label>
                <textarea id="news-content" rows="6">${isEdit ? (n.content || '') : ''}</textarea>
            </div>
            <div class="form-group">
                <label>URL изображения</label>
                <input type="url" id="news-image" value="${isEdit ? (n.image_url || '') : ''}" placeholder="https://...">
            </div>
            <div class="form-actions">
                <button type="button" class="btn btn-secondary" onclick="closeModal()">Отмена</button>
                <button type="submit" class="btn btn-primary">${isEdit ? 'Сохранить' : 'Создать'}</button>
            </div>
        </form>
    `;
    openModal();
}

async function saveNews(event, id) {
    event.preventDefault();
    const body = {
        title: document.getElementById('news-title').value.trim(),
        content: document.getElementById('news-content').value.trim(),
        image_url: document.getElementById('news-image').value.trim()
    };
    try {
        if (id) {
            await apiRequest(`${ADMIN_API}/news/${id}`, 'PUT', body);
        } else {
            await apiRequest(`${ADMIN_API}/news`, 'POST', body);
        }
        closeModal();
        await loadNews();
    } catch (e) {
        alert('Ошибка: ' + e.message);
    }
}

async function deleteNewsItem(id) {
    if (!confirm('Удалить новость?')) return;
    try {
        await apiRequest(`${ADMIN_API}/news/${id}`, 'DELETE');
        await loadNews();
    } catch (e) {
        alert('Ошибка удаления: ' + e.message);
    }
}

// ==================== СТАТИСТИКА ====================

async function loadAdminStats() {
    try {
        const data = await apiRequest(`${ADMIN_API}/stats`);
        if (!data) return;

        const set = (id, val) => {
            const el = document.getElementById(id);
            if (el) el.textContent = val != null ? val : 0;
        };

        // Visitors by last visit
        set('admin-stat-total', data.total_visits);
        set('admin-stat-today', data.today_visits);
        set('admin-stat-week', data.week_visits);
        set('admin-stat-month', data.month_visits);

        // New visitors by first visit
        set('admin-stat-new-today', data.new_today);
        set('admin-stat-new-week', data.new_week);
        set('admin-stat-new-month', data.new_month);

        // Engagement
        set('admin-stat-returning', data.returning_visitors);
        set('admin-stat-pageviews', data.total_page_views);
        set('admin-stat-avg', data.avg_visits_per_user != null
            ? data.avg_visits_per_user.toFixed(1)
            : '0');

        // Entity counts
        set('admin-stat-exhibitions', data.exhibition_count);
        set('admin-stat-exhibits', data.exhibit_count);
        set('admin-stat-news', data.news_count);

        // Daily chart
        renderDailyChart(data.daily_visits || []);
    } catch (err) {
        console.error('Failed to load stats', err);
    }
}

function renderDailyChart(dailyVisits) {
    const container = document.getElementById('stats-chart');
    if (!container) return;

    if (!dailyVisits.length) {
        container.innerHTML = '<div class="chart-empty">Нет данных за последние 7 дней</div>';
        return;
    }

    const maxCount = Math.max(...dailyVisits.map(d => d.count), 1);
    const dayNames = ['Вс', 'Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб'];

    const bars = dailyVisits.map(d => {
        const pct = Math.round((d.count / maxCount) * 100);
        const dt = new Date(d.date + 'T00:00:00');
        const dayLabel = dayNames[dt.getDay()];
        const dateLabel = dt.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
        return `
            <div class="chart-col">
                <div class="chart-value">${d.count}</div>
                <div class="chart-bar-wrap">
                    <div class="chart-bar" style="height:${Math.max(pct, 4)}%"></div>
                </div>
                <div class="chart-day">${dayLabel}</div>
                <div class="chart-date">${dateLabel}</div>
            </div>`;
    }).join('');

    container.innerHTML = bars;
}

// ==================== ГРУППИРОВКА ====================

function toggleGroup(header) {
    const group = header.closest('.exhibit-group');
    group.classList.toggle('collapsed');
    const toggle = header.querySelector('.group-toggle');
    toggle.textContent = group.classList.contains('collapsed') ? '▸' : '▾';
}

// ==================== МОДАЛЬНОЕ ОКНО ====================

function openModal() {
    document.getElementById('modal-overlay').classList.add('active');
}

function closeModal(event) {
    if (event && event.target !== event.currentTarget) return;
    document.getElementById('modal-overlay').classList.remove('active');
}

// Закрытие по Escape
document.addEventListener('keydown', e => {
    if (e.key === 'Escape') closeModal();
});
