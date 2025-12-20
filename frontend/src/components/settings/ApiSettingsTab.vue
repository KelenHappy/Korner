<template>
    <div class="tab-panel">
        <h3 class="section-title">{{ t("settings.api") }}</h3>

        <div class="form-group">
            <label class="form-label">{{ t("settings.provider") }}</label>
            <select v-model="localSettings.apiProvider" class="form-select">
                <option value="ollama">🦙 Ollama (推薦)</option>
                <option value="gptoss">🚀 AMD GPT-OSS-120B</option>
                <option value="openai" disabled>OpenAI (已停用)</option>
                <option value="anthropic" disabled>Anthropic (已停用)</option>
                <option value="gemini" disabled>Google Gemini (已停用)</option>
            </select>
        </div>

        <div class="form-group" v-if="localSettings.apiProvider === 'ollama'">
            <label class="form-label">Ollama 端點</label>
            <input
                v-model="localSettings.ollamaEndpoint"
                type="text"
                class="form-input"
                placeholder="http://127.0.0.1:11434"
            />
            <p class="form-hint">本地 Ollama 服務端點，使用 qwen3-vl:4b 模型</p>
        </div>

        <div class="form-group" v-if="localSettings.apiProvider === 'gptoss'">
            <label class="form-label">{{ t("settings.endpoint") }}</label>
            <select v-model="localSettings.apiEndpoint" class="form-select">
                <option value="http://210.61.209.139:45014/v1/">{{ t("settings.endpoint1") }}</option>
                <option value="http://210.61.209.139:45005/v1/">{{ t("settings.endpoint2") }}</option>
            </select>
            <p class="form-hint">{{ t("settings.endpointHint") }}</p>
        </div>

        <div class="form-group" v-if="localSettings.apiProvider !== 'gptoss' && localSettings.apiProvider !== 'ollama'">
            <label class="form-label">{{ t("settings.apiKey") }}</label>
            <div class="input-with-icon">
                <input
                    v-model="localSettings.apiKey"
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
</template>

<script>
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

export default {
    name: 'ApiSettingsTab',
    props: {
        settings: {
            type: Object,
            required: true
        }
    },
    emits: ['update:settings'],
    setup(props, { emit }) {
        const { t } = useI18n();
        const showApiKey = ref(false);
        const localSettings = ref({ ...props.settings });

        watch(localSettings, (newVal) => {
            emit('update:settings', newVal);
        }, { deep: true });

        watch(() => props.settings, (newVal) => {
            localSettings.value = { ...newVal };
        }, { deep: true });

        return {
            t,
            showApiKey,
            localSettings
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
</style>
