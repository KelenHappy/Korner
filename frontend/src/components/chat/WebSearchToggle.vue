<template>
    <button
        @click="toggle"
        :disabled="disabled"
        class="web-search-btn"
        :class="{ 'active': modelValue }"
        :title="modelValue ? '關閉旅遊搜尋' : '開啟旅遊搜尋'"
    >
        🚆
    </button>
</template>

<script>
export default {
    name: 'WebSearchToggle',
    props: {
        modelValue: {
            type: Boolean,
            default: false
        },
        disabled: {
            type: Boolean,
            default: false
        }
    },
    emits: ['update:modelValue'],
    setup(props, { emit }) {
        const toggle = () => {
            if (!props.disabled) {
                emit('update:modelValue', !props.modelValue);
            }
        };

        return { toggle };
    }
};
</script>

<style scoped>
.web-search-btn {
    width: 36px;
    height: 36px;
    border: 1.5px solid rgba(0, 0, 0, 0.08);
    background: white;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    transition: all 0.2s;
}

.web-search-btn:hover:not(:disabled) {
    border-color: #3b82f6;
    background: rgba(59, 130, 246, 0.05);
}

.web-search-btn.active {
    border-color: #3b82f6;
    background: rgba(59, 130, 246, 0.15);
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.web-search-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}
</style>
