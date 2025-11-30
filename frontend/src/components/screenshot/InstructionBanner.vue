<template>
    <div class="instruction-banner">
        <div class="banner-content">
            <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"
                />
            </svg>
            <span class="instruction-text">{{ instructionText }}</span>
            <kbd class="key-hint">{{ escKey }}</kbd>
            <span class="cancel-text">{{ cancelText }}</span>
        </div>
    </div>
</template>

<script>
export default {
    name: 'InstructionBanner',
    props: {
        instructionText: {
            type: String,
            required: true
        },
        escKey: {
            type: String,
            default: 'ESC'
        },
        cancelText: {
            type: String,
            required: true
        }
    }
};
</script>

<style scoped>
.instruction-banner {
    position: absolute;
    top: 32px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 10;
}

.banner-content {
    background: white;
    border-radius: 12px;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    padding: 12px 24px;
    display: flex;
    align-items: center;
    gap: 12px;
}

.icon {
    width: 20px;
    height: 20px;
    color: rgb(59, 130, 246);
}

.instruction-text {
    color: rgb(51, 65, 85);
    font-weight: 500;
    font-size: 14px;
}

.key-hint {
    margin-left: 8px;
    padding: 4px 8px;
    background-color: rgb(241, 245, 249);
    color: rgb(71, 85, 105);
    font-size: 12px;
    border-radius: 4px;
    border: 1px solid rgb(203, 213, 225);
    font-family: monospace;
}

.cancel-text {
    color: rgb(100, 116, 139);
    font-size: 14px;
}
</style>
