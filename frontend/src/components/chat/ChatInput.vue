<template>
    <div class="input-area">
        <textarea
            v-model="inputText"
            @keydown.ctrl.enter="handleSubmit"
            @keydown.meta.enter="handleSubmit"
            class="chat-input"
            rows="3"
            :placeholder="placeholder"
            :disabled="disabled"
        ></textarea>
        <div class="input-actions">
            <span class="char-count">{{ inputText.length }} / 1000 {{ charCountLabel }}</span>
            <button
                @click="handleSubmit"
                :disabled="!inputText.trim() || disabled"
                class="send-btn"
            >
                <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                    <path d="M2 10L18 2L10 18L8 11L2 10Z" fill="currentColor" />
                </svg>
            </button>
        </div>
    </div>
</template>

<script>
import { ref, watch } from 'vue';

export default {
    name: 'ChatInput',
    props: {
        modelValue: {
            type: String,
            default: ''
        },
        placeholder: {
            type: String,
            default: 'Type your message...'
        },
        disabled: {
            type: Boolean,
            default: false
        },
        charCountLabel: {
            type: String,
            default: 'characters'
        }
    },
    emits: ['update:modelValue', 'submit'],
    setup(props, { emit }) {
        const inputText = ref(props.modelValue);

        watch(() => props.modelValue, (newValue) => {
            inputText.value = newValue;
        });

        watch(inputText, (newValue) => {
            emit('update:modelValue', newValue);
        });

        const handleSubmit = () => {
            const trimmed = inputText.value.trim();
            if (!trimmed || props.disabled) return;
            
            emit('submit', trimmed.slice(0, 1000));
            inputText.value = '';
        };

        return {
            inputText,
            handleSubmit
        };
    }
};
</script>

<style scoped>
.input-area {
    border-top: 1px solid #cccccc;
    padding: 16px 24px;
    background: #ffffff;
}

.chat-input {
    width: 100%;
    padding: 12px 16px;
    border: 1.5px solid rgba(0, 0, 0, 0.08);
    border-radius: 12px;
    font-size: 14px;
    line-height: 1.6;
    color: #1a1a1a;
    resize: none;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    box-sizing: border-box;
    background: white;
    margin-bottom: 12px;
}

.chat-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 1px #667eea, 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.chat-input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.chat-input::placeholder {
    color: #9ca3af;
}

.input-actions {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.char-count {
    font-size: 12px;
    color: #9ca3af;
    font-weight: 500;
}

.send-btn {
    width: 40px;
    height: 40px;
    border: none;
    background: #667eea;
    color: white;
    border-radius: 50%;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 0 0 1px rgba(102, 126, 234, 0.1), 0 4px 12px rgba(102, 126, 234, 0.3);
}

.send-btn:hover:not(:disabled) {
    transform: scale(1.1);
    box-shadow: 0 0 0 1px rgba(102, 126, 234, 0.2), 0 6px 16px rgba(102, 126, 234, 0.4);
}

.send-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: scale(1);
}
</style>
