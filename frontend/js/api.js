// API configuration
const API_BASE_URL = '/museum';

// API функции для работы с бэкендом
const api = {
    // Получить все экспозиции
    async getAllExhibitions() {
        try {
            const response = await fetch(`${API_BASE_URL}/exhibitions`);
            if (!response.ok) throw new Error('Ошибка загрузки экспозиций');
            const data = await response.json();
            return Array.isArray(data) ? data : (data.exhibitions || []);
        } catch (error) {
            console.error('API Error:', error);
            return [];
        }
    },

    // Получить экспозицию по ID
    async getExhibitionById(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/exhibitions/${id}`);
            if (!response.ok) throw new Error('Экспозиция не найдена');
            const data = await response.json();
            return data.exhibition || data || null;
        } catch (error) {
            console.error('API Error:', error);
            return null;
        }
    },

    // Получить все новости
    async getAllNews() {
        try {
            const response = await fetch(`${API_BASE_URL}/news`);
            if (!response.ok) throw new Error('Ошибка загрузки новостей');
            const data = await response.json();
            return data.news || data || [];
        } catch (error) {
            console.error('API Error:', error);
            return [];
        }
    },

    // Получить новость по ID
    async getNewsById(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/news/${id}`);
            if (!response.ok) throw new Error('Новость не найдена');
            const data = await response.json();
            return data.news || data || null;
        } catch (error) {
            console.error('API Error:', error);
            return null;
        }
    }
};

// Форматирование даты
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}
