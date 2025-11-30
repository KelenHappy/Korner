<template>
    <div class="modal-header">
        <div class="header-content">
            <div class="icon-wrapper">
                <span class="header-icon">✨</span>
            </div>
            <h2 class="modal-title">{{ title }}</h2>
        </div>
        <div class="header-actions">
            <button
                v-if="showClearButton"
                @click="$emit('clear')"
                class="clear-btn"
                :title="clearTitle"
            >
                <svg width="18" height="18" viewBox="0 0 18 18" fill="none">
                    <path
                        d="M3 5H15M7 8V13M11 8V13M4 5L5 15C5 15.5304 5.21071 16.0391 5.58579 16.4142C5.96086 16.7893 6.46957 17 7 17H11C11.5304 17 12.0391 16.7893 12.4142 16.4142C12.7893 16.0391 13 15.5304 13 15L14 5M6 5V3C6 2.73478 6.10536 2.48043 6.29289 2.29289C6.48043 2.10536 6.73478 2 7 2H11C11.2652 2 11.5196 2.10536 11.7071 2.29289C11.8946 2.48043 12 2.73478 12 3V5"
                        stroke="currentColor"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
            </button>
            <button @click="$emit('close')" class="close-btn" :title="closeTitle">
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                    <path d="M15 5L5 15M5 5L15 15" stroke="currentColor" stroke-width="2" stroke-linecap="round" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script>
export default {
    name: 'ChatHeader',
    props: {
        title: {
            type: String,
            required: true
        },
        showClearButton: {
            type: Boolean,
            default: false
        },
        clearTitle: {
            type: String,
            default: 'Clear'
        },
        closeTitle: {
            type: String,
            default: 'Close'
        }
    },
    emits: ['clear', 'close']
};
</script>

<style scoped>
.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 24px 32px;
    background: #f5f5f5;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    flex-shrink: 0;
}

.header-actions {
    display: flex;
    align-items: center;
    gap: 8px;
}

.clear-btn,
.close-btn {
    width: 36px;
    height: 36px;
    border: none;
    background: transparent;
    color: #64748b;
    border-radius: 10px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;
    padding: 0;
}

.clear-btn:hover {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
}

.close-btn:hover {
    background: rgba(100, 116, 139, 0.1);
    color: #1a1a1a;
    transform: rotate(90deg);
}

.header-content {
    display: flex;
    align-items: center;
    gap: 14px;
}

.icon-wrapper {
    width: 44px;
    height: 44px;
    border-radius: 12px;
    background: #ffffff;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
}

.header-icon {
    font-size: 24px;
    filter: brightness(1.2);
}

.modal-title {
    font-size: 22px;
    font-weight: 600;
    color: #1a1a1a;
    margin: 0;
    letter-spacing: -0.02em;
}
</style>
