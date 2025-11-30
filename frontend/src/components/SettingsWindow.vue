<template>
    <div class="settings-overlay">
        <div class="settings-modal">
            <!-- Header -->
            <div class="modal-header">
                <h2 class="modal-title">⚙️ {{ t("settings.title") }}</h2>
                <button @click="close" class="close-btn" :title="t('settings.cancel')">
                    ✕
                </button>
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
                    <!-- API Settings Tab -->
                    <div v-if="activeTab === 'api'" class="tab-panel">
                        <h3 class="section-title">{{ t("settings.api") }}</h3>

                        <div class="form-group">
                            <label class="form-label">{{ t("settings.provider") }}</label>
                            <select
                                v-model="settings.apiProvider"
                                class="form-select"
                            >
                                <option value="gptoss">🚀 AMD GPT-OSS-120B (推薦)</option>
                                <option value="openai" disabled>OpenAI (已停用)</option>
                                <option value="anthropic" disabled>Anthropic (已停用)</option>
                                <option value="gemini" disabled>Google Gemini (已停用)</option>
                                <option value="custom" disabled>Custom (已停用)</option>
                            </select>
                        </div>

                        <div class="form-group" v-if="settings.apiProvider === 'gptoss'">
                            <label class="form-label">{{ t("settings.endpoint") }}</label>
                            <select
                                v-model="settings.apiEndpoint"
                                class="form-select"
                            >
                                <option value="http://210.61.209.139:45014/v1/">{{ t("settings.endpoint1") }}</option>
                                <option value="http://210.61.209.139:45005/v1/">{{ t("settings.endpoint2") }}</option>
                            </select>
                            <p class="form-hint">{{ t("settings.endpointHint") }}</p>
                        </div>

                        <div class="form-group" v-if="settings.apiProvider !== 'gptoss'">
                            <label class="form-label">{{ t("settings.apiKey") }}</label>
                            <div class="input-with-icon">
                                <input
                                    v-model="settings.apiKey"
                                    :type="showApiKey ? 'text' : 'password'"
                                    class="form-input"
                                    placeholder="sk-..."
                                    disabled
                                />
                                <button
                                    @click="showApiKey = !showApiKey"
                                    class="icon-btn"
                                    type="button"
                                    disabled
                                >
                                    {{ showApiKey ? "🙈" : "👁️" }}
                                </button>
                            </div>
                            <p class="form-hint">{{ t("settings.providerDisabled") }}</p>
                        </div>
                    </div>

                    <!-- Icon Settings Tab -->
                    <div v-if="activeTab === 'icon'" class="tab-panel">
                        <h3 class="section-title">{{ t("settings.icon") }}</h3>

                        <div class="form-group">
                            <label class="form-label">{{ t("settings.chooseIcon") }}</label>
                            <div class="icon-grid">
                                <button
                                    v-for="icon in availableIcons"
                                    :key="icon"
                                    @click="settings.floatingIcon = icon"
                                    :class="[
                                        'icon-option',
                                        {
                                            selected:
                                                settings.floatingIcon === icon,
                                        },
                                    ]"
                                >
                                    {{ icon }}
                                </button>
                            </div>
                        </div>

                        <div class="form-group">
                            <label class="form-label"
                                >{{ t("settings.customIcon") }}</label
                            >
                            <input
                                v-model="settings.floatingIcon"
                                type="text"
                                class="form-input"
                                :placeholder="t('settings.customIcon') + '...'"
                                maxlength="2"
                            />
                        </div>

                        <div class="preview-section">
                            <div class="preview-label">{{ t("settings.preview") }}:</div>
                            <div class="icon-preview">
                                {{ settings.floatingIcon }}
                            </div>
                        </div>
                    </div>

                    <!-- Language Settings Tab -->
                    <div v-if="activeTab === 'language'" class="tab-panel">
                        <h3 class="section-title">Language / 語言</h3>

                        <div class="form-group">
                            <label class="form-label">{{ t("settings.language") }}</label>
                            <div class="language-options">
                                <button
                                    @click="changeLanguage('en')"
                                    :class="['lang-btn', { active: currentLocale === 'en' }]"
                                >
                                    🇺🇸 English
                                </button>
                                <button
                                    @click="changeLanguage('zh-TW')"
                                    :class="['lang-btn', { active: currentLocale === 'zh-TW' }]"
                                >
                                    🇹🇼 繁體中文
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Footer -->
            <div class="modal-footer">
                <button @click="resetToDefaults" class="btn-secondary">
                    {{ t("settings.reset") }}
                </button>
                <div class="footer-right">
                    <button @click="close" class="btn-cancel">{{ t("settings.cancel") }}</button>
                    <button @click="save" class="btn-save">
                        💾 {{ t("settings.save") }}
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from "vue";
import { useI18n } from "vue-i18n";

export default {
    name: "SettingsWindow",
    props: {
        currentSettings: {
            type: Object,
            default: () => ({}),
        },
    },
    emits: ["close", "save"],
    setup(props, { emit }) {
        const { t, locale } = useI18n();
        const currentLocale = computed(() => locale.value);
        
        const activeTab = ref("api");
        const showApiKey = ref(false);

        const tabs = computed(() => [
            { id: "api", name: t("tabs.api"), icon: "🤖" },
            { id: "icon", name: t("tabs.icon"), icon: "🎨" },
            { id: "language", name: t("tabs.language"), icon: "🌐" },
        ]);

        const availableIcons = [
            "🌸",
            "✨",
            "🔮",
            "🎯",
            "🚀",
            "💡",
            "🎨",
            "🌈",
            "⭐",
            "💎",
            "🦄",
            "🐱",
            "🐶",
            "🦊",
            "🐼",
            "🦋",
        ];

        const defaultSettings = {
            apiProvider: "openai",
            apiKey: "",
            apiEndpoint: "",
            floatingIcon: "🌸",
        };

        const settings = reactive({
            ...defaultSettings,
            ...props.currentSettings,
        });

        const resetToDefaults = () => {
            if (confirm(t("settings.resetConfirm"))) {
                Object.assign(settings, defaultSettings);
            }
        };

        const save = () => {
            emit("save", { ...settings });
        };

        const close = () => {
            emit("close");
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
            showApiKey,
            tabs,
            availableIcons,
            settings,
            currentLocale,
            changeLanguage,
            resetToDefaults,
            save,
            close,
        };
    },
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

/* Header */
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

/* Body */
.modal-body {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

/* Tabs */
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

/* Tab Content */
.tab-content {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
}

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

/* Form Elements */
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

.form-input,
.form-select {
    width: 100%;
    padding: 10px 12px;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    font-size: 14px;
    color: #1e293b;
    transition: all 0.2s ease;
    box-sizing: border-box;
}

.form-input:focus,
.form-select:focus {
    outline: none;
    border-color: #94a3b8;
    box-shadow: 0 0 0 3px rgba(148, 163, 184, 0.1);
}

.input-with-icon {
    position: relative;
    display: flex;
    align-items: center;
}

.input-with-icon .form-input {
    padding-right: 45px;
}

.icon-btn {
    position: absolute;
    right: 8px;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 18px;
    padding: 4px;
    transition: transform 0.2s ease;
}

.icon-btn:hover {
    transform: scale(1.1);
}

.icon-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.form-hint {
    margin-top: 6px;
    font-size: 12px;
    color: #64748b;
    line-height: 1.4;
}

.form-range {
    width: 100%;
    height: 6px;
    border-radius: 3px;
    background: #e2e8f0;
    outline: none;
    -webkit-appearance: none;
}

.form-range::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: #94a3b8;
    cursor: pointer;
    box-shadow: 0 2px 8px rgba(148, 163, 184, 0.3);
}

.form-range::-moz-range-thumb {
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: #94a3b8;
    cursor: pointer;
    box-shadow: 0 2px 8px rgba(148, 163, 184, 0.3);
    border: none;
}

.range-labels {
    display: flex;
    justify-content: space-between;
    margin-top: 8px;
    font-size: 11px;
    color: #94a3b8;
}

/* Icon Grid */
.icon-grid {
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    gap: 8px;
}

.icon-option {
    aspect-ratio: 1;
    border: 2px solid #e2e8f0;
    background: white;
    border-radius: 8px;
    font-size: 24px;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
}

.icon-option:hover {
    border-color: #94a3b8;
    transform: scale(1.1);
}

.icon-option.selected {
    border-color: #64748b;
    background: #f1f5f9;
    box-shadow: 0 0 0 3px rgba(100, 116, 139, 0.1);
}

.preview-section {
    margin-top: 24px;
    padding: 20px;
    background: #f5f5f5;
    border-radius: 12px;
    display: flex;
    align-items: center;
    gap: 16px;
}

.preview-label {
    font-size: 14px;
    font-weight: 600;
    color: #475569;
}

.icon-preview {
    width: 60px;
    height: 60px;
    background: #e2e8f0;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 30px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
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

/* Footer */
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

/* Scrollbar */
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
