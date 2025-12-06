<template>
    <transition name="bubble-fade">
        <div class="action-bubbles">
            <button 
                v-for="btn in buttons"
                :key="btn.name"
                :class="['bubble-btn', btn.name]"
                @click="$emit(btn.event)"
                :disabled="btn.disabled"
                :title="btn.title"
            >
                <span class="bubble-icon">{{ btn.icon }}</span>
            </button>
        </div>
    </transition>
</template>

<script>
export default {
    name: 'ActionButtons',
    props: {
        buttons: {
            type: Array,
            required: true
        }
    },
    emits: ['summary', 'folder', 'selectFile']
};
</script>

<style scoped>
.action-bubbles {
    display: flex;
    gap: 10px;
    animation: bubble-slide-up 0.3s ease-out;
}

@keyframes bubble-slide-up {
    from {
        transform: translateY(20px);
        opacity: 0;
    }
    to {
        transform: translateY(0);
        opacity: 1;
    }
}

.bubble-btn {
    width: 35px;
    height: 35px;
    border-radius: 50%;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    transition: all 0.2s;
    pointer-events: auto;
}

.bubble-btn:hover:not(:disabled) {
    transform: translateY(-3px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.bubble-btn:active:not(:disabled) {
    transform: translateY(-1px);
}

.bubble-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.bubble-btn.back {
    background: #6b7280;
}

.bubble-btn.back:hover:not(:disabled) {
    background: #4b5563;
}

.bubble-btn.summary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.bubble-btn.folder {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.bubble-btn.file-select {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.bubble-icon {
    font-size: 16px;
}
</style>
