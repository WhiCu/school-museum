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


