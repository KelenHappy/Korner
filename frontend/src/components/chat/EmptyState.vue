<template>
    <div class="empty-state">
        <div class="empty-icon">💭</div>
        <p class="empty-text">{{ message }}</p>
        <div class="prompts-grid" v-if="prompts.length > 0">
            <button
                v-for="prompt in prompts"
                :key="prompt"
                @click="$emit('select-prompt', prompt)"
                class="prompt-btn"
            >
                <span class="prompt-text">{{ prompt }}</span>
            </button>
        </div>
    </div>
</template>

<script>
export default {
    name: 'EmptyState',
    props: {
        message: {
            type: String,
            required: true
        },
        prompts: {
            type: Array,
            default: () => []
        }
    },
    emits: ['select-prompt']
};
</script>

<style scoped>
.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 16px;
    padding: 40px 20px;
    text-align: center;
}

.empty-icon {
    font-size: 48px;
    opacity: 0.5;
}

.empty-text {
    font-size: 15px;
    color: #64748b;
    margin: 0;
}

.prompts-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
}

.prompt-btn {
    padding: 11px 16px;
    font-size: 13px;
    background: white;
    border: 1.5px solid rgba(0, 0, 0, 0.08);
    color: #334155;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    position: relative;
}

.prompt-btn::before {
    content: "";
    position: absolute;
    inset: 0;
    border-radius: 10px;
    background: #667eea;
    opacity: 0;
    transition: opacity 0.2s ease;
}

.prompt-text {
    position: relative;
    z-index: 1;
}

.prompt-btn:hover {
    border-color: #667eea;
    color: white;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.25), 0 0 0 1px rgba(102, 126, 234, 0.1);
}

.prompt-btn:hover::before {
    opacity: 1;
}

.prompt-btn:active {
    transform: translateY(0);
}

@media (max-width: 600px) {
    .prompts-grid {
        grid-template-columns: 1fr;
    }
}
</style>
