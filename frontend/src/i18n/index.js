import { createI18n } from 'vue-i18n';
import en from './locales/en.json';
import zhTW from './locales/zh-TW.json';

// Get saved language or detect browser language
function getDefaultLocale() {
    try {
        const saved = localStorage.getItem('korner-language');
        if (saved) return saved;
    } catch (e) {
        console.log('[i18n] localStorage not available:', e);
    }
    
    try {
        const browserLang = navigator.language || navigator.userLanguage;
        if (browserLang.startsWith('zh')) {
            return 'zh-TW';
        }
    } catch (e) {
        console.log('[i18n] navigator.language not available:', e);
    }
    
    return 'en';
}

const i18n = createI18n({
    legacy: false,
    locale: getDefaultLocale(),
    fallbackLocale: 'en',
    messages: {
        en,
        'zh-TW': zhTW
    }
});

export default i18n;
