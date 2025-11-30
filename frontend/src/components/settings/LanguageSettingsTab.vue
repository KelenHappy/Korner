<template>
    <div class="tab-panel">
        <h3 class="section-title">Language / 語言</h3>

        <div class="form-group">
            <label class="form-label">{{ t("settings.language") }}</label>
            <div class="language-options">
                <button
                    @click="selectLanguage('en')"
                    :class="['lang-btn', { active: currentLocale === 'en' }]"
                >
                    🇺🇸 English
                </button>
                <button
                    @click="selectLanguage('zh-TW')"
                    :class="['lang-btn', { active: currentLocale === 'zh-TW' }]"
                >
                    🇹🇼 繁體中文
                </button>
            </div>
        </div>
    </div>
</template>

<script>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

export default {
    name: 'LanguageSettingsTab',
    emits: ['change-language'],
    setup(props, { emit }) {
        const { t, locale } = useI18n();
        const currentLocale = computed(() => locale.value);

        const selectLanguage = (lang) => {
            emit('change-language', lang);
        };

        return {
            t,
            currentLocale,
            selectLanguage
        };
    }
};
</script>

<style scoped>
.tab-panel {
    animation: tabFadeIn 0.3s ease;
}

@keyframes tabFadeIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.section-title {
    font-size: 16px;
    font-weight: 700;
    color: #1e293b;
    margin: 0 0 20px 0;
}

.form-group {
    margin-bottom: 20px;
}

.form-label {
    display: block;
    font-size: 13px;
    font-weight: 600;
    color: #475569;
    margin-bottom: 8px;
}

.language-options {
    display: flex;
    gap: 12px;
    margin-top: 8px;
}

.lang-btn {
    flex: 1;
    padding: 16px 24px;
    border: 2px solid #e2e8f0;
    background: white;
    border-radius: 12px;
    font-size: 15px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.lang-btn:hover {
    border-color: #cbd5e1;
    background: #f8fafc;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.lang-btn.active {
    border-color: #000;
    background: #000;
    color: white;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.lang-btn.active:hover {
    background: #1a1a1a;
    border-color: #1a1a1a;
}
</style>
