<template>
    <div class="query-overlay">
        <div class="query-modal">
            <!-- Header -->
            <div class="modal-header">
                <h2 class="modal-title">💭 Ask about this screenshot</h2>
                <button @click="cancel" class="close-btn" title="Close">
                    ✕
                </button>
            </div>

            <!-- Screenshot Preview -->
            <div class="modal-body">
                <div class="screenshot-section">
                    <img
                        :src="screenshot"
                        alt="Screenshot"
                        class="screenshot-img"
                    />
                </div>

                <!-- Query Input -->
                <div class="input-section">
                    <label class="input-label">Your Question</label>
                    <textarea
                        v-model="queryText"
                        @keydown.ctrl.enter="submit"
                        @keydown.meta.enter="submit"
                        class="query-textarea"
                        rows="4"
                        placeholder="Ask anything about this screenshot..."
                        autofocus
                    ></textarea>
                    <div class="input-info">
                        <span class="char-count"
                            >{{ queryText.length }} / 1000</span
                        >
                        <span class="shortcut-hint">
                            Press
                            <kbd class="kbd">{{ submitKey }}</kbd>
                            to submit
                        </span>
                    </div>
                </div>

                <!-- Quick Prompts -->
                <div class="prompts-section">
                    <p class="prompts-label">Quick suggestions:</p>
                    <div class="prompts-grid">
                        <button
                            v-for="prompt in quickPrompts"
                            :key="prompt"
                            @click="queryText = prompt"
                            class="prompt-btn"
                        >
                            {{ prompt }}
                        </button>
                    </div>
                </div>
            </div>

            <!-- Footer Actions -->
            <div class="modal-footer">
                <button @click="cancel" class="btn-cancel">Cancel</button>
                <button
                    @click="submit"
                    :disabled="!queryText.trim()"
                    class="btn-submit"
                >
                    <span class="btn-icon">✨</span>
                    <span class="btn-text">Ask AI</span>
                </button>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed } from "vue";

export default {
    name: "QueryWindow",
    props: {
        screenshot: {
            type: String,
            required: true,
        },
    },
    emits: ["submit", "cancel"],
    setup(props, { emit }) {
        const queryText = ref("");

        const quickPrompts = [
            "Explain this",
            "What's wrong?",
            "Summarize",
            "Translate",
            "Improve this",
            "Find bugs",
        ];

        const isMac = computed(() => {
            return navigator.userAgent.toLowerCase().includes("mac");
        });

        const submitKey = computed(() => {
            return isMac.value ? "Cmd+Enter" : "Ctrl+Enter";
        });

        const submit = () => {
            const trimmed = queryText.value.trim();
            if (trimmed) {
                emit("submit", trimmed.slice(0, 1000));
                queryText.value = "";
            }
        };

        const cancel = () => {
            emit("cancel");
        };

        return {
            queryText,
            quickPrompts,
            submitKey,
            submit,
            cancel,
        };
    },
};
</script>

<style scoped>
.query-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10000;
    padding: 20px;
    animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}

.query-modal {
    background: white;
    border-radius: 16px;
    box-shadow: 0 25px 80px rgba(0, 0, 0, 0.2);
    width: 100%;
    max-width: 500px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    animation: slideUp 0.3s ease-out;
    overflow: hidden;
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(30px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* Header */
.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-bottom: none;
}

.modal-title {
    font-size: 18px;
    font-weight: 700;
    color: white;
    margin: 0;
}

.close-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: rgba(255, 255, 255, 0.2);
    color: white;
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
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.05);
}

/* Body */
.modal-body {
    flex: 1;
    padding: 24px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

/* Screenshot */
.screenshot-section {
    width: 100%;
    border-radius: 12px;
    overflow: hidden;
    border: 2px solid rgba(102, 126, 234, 0.15);
    background: rgba(102, 126, 234, 0.05);
}

.screenshot-img {
    width: 100%;
    max-height: 200px;
    object-fit: contain;
    display: block;
    background: white;
}

/* Input Section */
.input-section {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.input-label {
    font-size: 13px;
    font-weight: 600;
    color: #4b5563;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.query-textarea {
    width: 100%;
    padding: 12px;
    border: 2px solid rgba(102, 126, 234, 0.2);
    border-radius: 10px;
    font-size: 14px;
    line-height: 1.6;
    color: #1f2937;
    resize: vertical;
    font-family:
        -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    transition: all 0.2s ease;
    box-sizing: border-box;
}

.query-textarea:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.query-textarea::placeholder {
    color: #9ca3af;
}

.input-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 12px;
}

.char-count {
    color: #9ca3af;
}

.shortcut-hint {
    color: #6b7280;
}

.kbd {
    display: inline-block;
    padding: 2px 6px;
    background: rgba(102, 126, 234, 0.1);
    border: 1px solid rgba(102, 126, 234, 0.2);
    border-radius: 4px;
    font-size: 11px;
    font-family: monospace;
    color: #667eea;
}

/* Prompts */
.prompts-section {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.prompts-label {
    font-size: 12px;
    font-weight: 600;
    color: #6b7280;
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.prompts-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
}

.prompt-btn {
    padding: 9px 12px;
    font-size: 13px;
    background: linear-gradient(
        135deg,
        rgba(102, 126, 234, 0.08) 0%,
        rgba(118, 75, 162, 0.08) 100%
    );
    border: 1px solid rgba(102, 126, 234, 0.15);
    color: #4b5563;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.prompt-btn:hover {
    background: linear-gradient(
        135deg,
        rgba(102, 126, 234, 0.15) 0%,
        rgba(118, 75, 162, 0.15) 100%
    );
    border-color: rgba(102, 126, 234, 0.3);
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.15);
}

/* Footer */
.modal-footer {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 24px;
    background: rgba(102, 126, 234, 0.03);
    border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.btn-cancel {
    padding: 10px 20px;
    font-size: 14px;
    font-weight: 600;
    color: #6b7280;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.btn-cancel:hover {
    color: #374151;
    background: rgba(0, 0, 0, 0.05);
}

.btn-submit {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    font-size: 14px;
    font-weight: 600;
    color: white;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-submit:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.btn-submit:active:not(:disabled) {
    transform: translateY(0);
}

.btn-submit:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    box-shadow: none;
}

.btn-icon {
    font-size: 16px;
}

.btn-text {
    font-weight: 600;
}

/* Scrollbar styling */
.modal-body::-webkit-scrollbar {
    width: 6px;
}

.modal-body::-webkit-scrollbar-track {
    background: rgba(102, 126, 234, 0.05);
}

.modal-body::-webkit-scrollbar-thumb {
    background: rgba(102, 126, 234, 0.2);
    border-radius: 3px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
    background: rgba(102, 126, 234, 0.3);
}

/* Mobile responsiveness */
@media (max-width: 600px) {
    .query-overlay {
        padding: 10px;
    }

    .query-modal {
        max-height: 95vh;
    }

    .modal-header {
        padding: 16px 20px;
    }

    .modal-title {
        font-size: 16px;
    }

    .modal-body {
        padding: 16px;
        gap: 16px;
    }

    .prompts-grid {
        grid-template-columns: repeat(2, 1fr);
    }

    .modal-footer {
        padding: 12px 16px;
    }
}
</style>
