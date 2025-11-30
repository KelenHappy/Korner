<template>
    <div class="screenshot-column">
        <div v-if="hasScreenshot" class="screenshot-wrapper">
            <div class="screenshot-frame">
                <img :src="screenshotSrc" alt="Screenshot" class="screenshot-img" />
            </div>
        </div>
        <div v-else class="screenshot-placeholder">
            <div class="placeholder-icon">
                <svg width="64" height="64" viewBox="0 0 64 64" fill="none">
                    <rect x="8" y="12" width="48" height="40" rx="4" stroke="currentColor" stroke-width="2" />
                    <circle cx="32" cy="28" r="6" stroke="currentColor" stroke-width="2" />
                    <path
                        d="M8 44L18 34L28 44L42 30L56 44V48C56 50.2091 54.2091 52 52 52H12C9.79086 52 8 50.2091 8 48V44Z"
                        fill="currentColor"
                        opacity="0.2"
                    />
                </svg>
            </div>
            <p class="placeholder-text">{{ placeholderText }}</p>
        </div>
    </div>
</template>

<script>
import { computed } from 'vue';

export default {
    name: 'ScreenshotPreview',
    props: {
        screenshot: {
            type: String,
            default: ''
        },
        placeholderText: {
            type: String,
            default: 'No screenshot available'
        }
    },
    setup(props) {
        const hasScreenshot = computed(() => {
            return typeof props.screenshot === 'string' && props.screenshot.trim().length > 0;
        });

        const screenshotSrc = computed(() => {
            if (!hasScreenshot.value) return '';
            return props.screenshot.startsWith('data:image')
                ? props.screenshot
                : `data:image/png;base64,${props.screenshot}`;
        });

        return {
            hasScreenshot,
            screenshotSrc
        };
    }
};
</script>

<style scoped>
.screenshot-column {
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: auto;
    border-radius: 16px;
    background: #ffffff;
    padding: 24px;
}

.screenshot-wrapper {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
}

.screenshot-frame {
    position: relative;
    max-width: 100%;
    max-height: 100%;
    border-radius: 12px;
    padding: 8px;
    background: #f5f5f5;
    box-shadow: 0 0 0 1px rgba(102, 126, 234, 0.1), 0 8px 32px rgba(0, 0, 0, 0.08);
}

.screenshot-img {
    max-width: 100%;
    max-height: 100%;
    width: auto;
    height: auto;
    object-fit: contain;
    border-radius: 8px;
    background: white;
    display: block;
}

.screenshot-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    color: #9ca3af;
}

.placeholder-icon {
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #cbd5e1;
}

.placeholder-text {
    margin: 0;
    font-size: 15px;
    font-weight: 500;
    color: #94a3b8;
}

.screenshot-column::-webkit-scrollbar {
    width: 6px;
}

.screenshot-column::-webkit-scrollbar-track {
    background: transparent;
}

.screenshot-column::-webkit-scrollbar-thumb {
    background: rgba(102, 126, 234, 0.2);
    border-radius: 3px;
}

.screenshot-column::-webkit-scrollbar-thumb:hover {
    background: rgba(102, 126, 234, 0.3);
}
</style>
