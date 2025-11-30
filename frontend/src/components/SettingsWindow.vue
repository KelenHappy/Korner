<template>
    <div class="settings-overlay">
        <div class="settings-modal">
            <!-- Header -->
            <div class="modal-header">
                <h2 class="modal-title">⚙️ {{ t("settings.title") }}</h2>
                <button @click="close" class="close-btn" :title="t('settings.cancel')">✕</button>
            </div>

            <!-- Body -->
            <div class="modal-body">
                <!-- Tabs -->
                <div class="tabs">
                    <button
                        v-for="tab in tabs"
                        :key="tab.id"
                        @click="activeTab = tab.id"
                        :class="['tab-btn', { active: activeTab === tab.id }]"
                    >
                        {{ tab.icon }} {{ tab.name }}
                    </button>
                </div>

                <!-- Tab Content -->
                <div class="tab-content">
                    <ApiSettingsTab
                        v-if="activeTab === 'api'"
                        :settings="settings"
                        @update:settings="updateSettings"
                    />
                    <IconSettingsTab
                        v-if="activeTab === 'icon'"
                        :icon="settings.floatingIcon"
                        @update:icon="settings.floatingIcon = $event"
                    />
                    <LanguageSettingsTab
                        v-if="activeTab === 'language'"
                        @change-language="changeLanguage"
                    />
                </div>
            </div>

            <!-- Footer -->
            <div class="modal-footer">
                <button @click="resetToDefaults" class="btn-secondary">
                    {{ t("settings.reset") }}
                </button>
                <div class="footer-right">
                    <button @click="close" class="btn-cancel">{{ t("settings.cancel") }}</button>
                    <button @click="save" class="btn-save">💾 {{ t("settings.save") }}</button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, reactive, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import ApiSettingsTab from './settings/ApiSettingsTab.vue';
import IconSettingsTab from './settings/IconSettingsTab.vue';
import LanguageSettingsTab from './settings/LanguageSettingsTab.vue';

export default {
    name: 'SettingsWindow',
    components: {
        ApiSettingsTab,
        IconSettingsTab,
        LanguageSettingsTab
    },
    props: {
        currentSettings: {
            type: Object,
            default: () => ({})
        }
    },
    emits: ['close', 'save'],
    setup(props, { emit }) {
        const { t, locale } = useI18n();
        const activeTab = ref('api');

        const tabs = computed(() => [
            { id: 'api', name: t('tabs.api'), icon: '🤖' },
            { id: 'icon', name: t('tabs.icon'), icon: '🎨' },
            { id: 'language', name: t('tabs.language'), icon: '🌐' }
        ]);

        const defaultSettings = {
            apiProvider: 'openai',
            apiKey: '',
            apiEndpoint: '',
            floatingIcon: '🌸'
        };

        const settings = reactive({
            ...defaultSettings,
            ...props.currentSettings
        });

        const updateSettings = (newSettings) => {
            Object.assign(settings, newSettings);
        };

        const resetToDefaults = () => {
            if (confirm(t('settings.resetConfirm'))) {
                Object.assign(settings, defaultSettings);
            }
        };

        const save = () => {
            emit('save', { ...settings });
        };

        const close = () => {
            emit('close');
        };

        const changeLanguage = (lang) => {
            locale.value = lang;
            try {
                localStorage.setItem('korner-language', lang);
            } catch (e) {
                console.log('[Settings] Failed to save language:', e);
            }
        };

        return {
            t,
            activeTab,
            tabs,
            settings,
            updateSettings,
            resetToDefaults,
            save,
            close,
            changeLanguage
        };
    }
};
</script>

<style scoped>
.settings-overlay {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 10000;
    animation: fadeIn 0.2s ease-out;
    pointer-events: auto;
}

@keyframes fadeIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}

.settings-modal {
    background: #ffffff;
    border-radius: 16px;
    box-shadow: 0 25px 80px rgba(0, 0, 0, 0.5);
    width: 700px;
    max-width: 90vw;
    max-height: 85vh;
    display: flex;
    flex-direction: column;
    animation: slideUp 0.3s ease-out;
    overflow: hidden;
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(40px) scale(0.95);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 24px;
    background: #f5f5f5;
    border-bottom: 1px solid #e2e8f0;
    flex-shrink: 0;
}

.modal-title {
    font-size: 18px;
    font-weight: 700;
    color: #1e293b;
    margin: 0;
}

.close-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: #e2e8f0;
    color: #64748b;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    transition: all 0.2s ease;
    padding: 0;
}

.close-btn:hover {
    background: #cbd5e1;
    transform: scale(1.05);
}

.modal-body {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.tabs {
    display: flex;
    gap: 4px;
    padding: 16px 16px 0 16px;
    background: rgba(248, 249, 250, 0.05);
    border-bottom: 1px solid #e2e8f0;
}

.tab-btn {
    padding: 10px 20px;
    border: none;
    background: transparent;
    color: #64748b;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    border-radius: 8px 8px 0 0;
    transition: all 0.2s ease;
}

.tab-btn:hover {
    background: #e2e8f0;
    color: #1e293b;
}

.tab-btn.active {
    background: #ffffff;
    color: #1e293b;
}

.tab-content {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
}

.modal-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 24px;
    background: #f5f5f5;
    border-top: 1px solid #e2e8f0;
    flex-shrink: 0;
}

.footer-right {
    display: flex;
    gap: 12px;
}

.btn-secondary,
.btn-cancel,
.btn-save {
    padding: 10px 20px;
    font-size: 14px;
    font-weight: 600;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    border: none;
}

.btn-secondary {
    background: transparent;
    color: #64748b;
    border: 1px solid #cbd5e1;
}

.btn-secondary:hover {
    background: #f1f5f9;
}

.btn-cancel {
    background: transparent;
    color: #64748b;
}

.btn-cancel:hover {
    background: #f1f5f9;
}

.btn-save {
    background: #1e293b;
    color: white;
    box-shadow: 0 4px 12px rgba(30, 41, 59, 0.2);
}

.btn-save:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(30, 41, 59, 0.3);
    background: #334155;
}

.tab-content::-webkit-scrollbar {
    width: 8px;
}

.tab-content::-webkit-scrollbar-track {
    background: #f1f5f9;
}

.tab-content::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 4px;
}

.tab-content::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
}
</style>
