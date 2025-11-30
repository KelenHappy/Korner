<template>
    <div class="tab-panel">
        <h3 class="section-title">{{ t("settings.icon") }}</h3>

        <div class="form-group">
            <label class="form-label">{{ t("settings.chooseIcon") }}</label>
            <div class="icon-grid">
                <button
                    v-for="icon in availableIcons"
                    :key="icon"
                    @click="selectIcon(icon)"
                    :class="['icon-option', { selected: localIcon === icon }]"
                >
                    {{ icon }}
                </button>
            </div>
        </div>

        <div class="form-group">
            <label class="form-label">{{ t("settings.customIcon") }}</label>
            <input
                v-model="localIcon"
                type="text"
                class="form-input"
                :placeholder="t('settings.customIcon') + '...'"
                maxlength="2"
            />
        </div>

        <div class="preview-section">
            <div class="preview-label">{{ t("settings.preview") }}:</div>
            <div class="icon-preview">{{ localIcon }}</div>
        </div>
    </div>
</template>

<script>
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

export default {
    name: 'IconSettingsTab',
    props: {
        icon: {
            type: String,
            required: true
        }
    },
    emits: ['update:icon'],
    setup(props, { emit }) {
        const { t } = useI18n();
        const localIcon = ref(props.icon);

        const availableIcons = [
            "🌸", "✨", "🔮", "🎯", "🚀", "💡", "🎨", "🌈",
            "⭐", "💎", "🦄", "🐱", "🐶", "🦊", "🐼", "🦋"
        ];

        const selectIcon = (icon) => {
            localIcon.value = icon;
        };

        watch(localIcon, (newVal) => {
            emit('update:icon', newVal);
        });

        watch(() => props.icon, (newVal) => {
            localIcon.value = newVal;
        });

        return {
            t,
            localIcon,
            availableIcons,
            selectIcon
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

.form-input {
    width: 100%;
    padding: 10px 12px;
    border: 2px solid #e2e8f0;
    border-radius: 8px;
    font-size: 14px;
    color: #1e293b;
    transition: all 0.2s ease;
    box-sizing: border-box;
}

.form-input:focus {
    outline: none;
    border-color: #94a3b8;
    box-shadow: 0 0 0 3px rgba(148, 163, 184, 0.1);
}

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
</style>
