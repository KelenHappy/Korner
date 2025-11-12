<template>
    <div class="query-window-container">
        <div
            ref="queryDialog"
            class="query-dialog"
            :style="{ top: position.y + 'px', left: position.x + 'px' }"
        >
            <!-- Draggable header -->
            <div
                class="query-header"
                @mousedown="startDrag"
                :style="{ cursor: isDragging ? 'grabbing' : 'grab' }"
            >
                <h3 class="query-title">💭 Ask about this screenshot</h3>
                <button @click="cancel" class="close-btn" title="Close">
                    ✕
                </button>
            </div>

            <!-- Screenshot Preview -->
            <div class="screenshot-preview">
                <img
                    :src="screenshot"
                    alt="Screenshot"
                    class="screenshot-img"
                />
            </div>

            <!-- Query Input -->
            <div class="query-input-section">
                <label class="input-label">Your Question</label>
                <textarea
                    v-model="queryText"
                    @keydown.ctrl.enter="submit"
                    @keydown.meta.enter="submit"
                    class="query-textarea"
                    rows="3"
                    placeholder="e.g., 'Explain this code', 'What's wrong with this error?', 'Summarize this chart'..."
                    autofocus
                ></textarea>
                <div class="input-footer">
                    <span class="shortcut-hint">
                        Press
                        <kbd class="kbd">{{ submitKey }}</kbd>
                        to submit
                    </span>
                    <span class="char-count"
                        >{{ queryText.length }} / 1000</span
                    >
                </div>
            </div>

            <!-- Quick Prompts -->
            <div class="quick-prompts">
                <p class="prompts-label">Quick prompts:</p>
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

            <!-- Actions -->
            <div class="actions">
                <button @click="cancel" class="btn-cancel">Cancel</button>
                <button
                    @click="submit"
                    :disabled="!queryText.trim()"
                    class="btn-submit"
                >
                    <svg
                        class="icon"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M13 10V3L4 14h7v7l9-11h-7z"
                        />
                    </svg>
                    <span>Ask AI</span>
                </button>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from "vue";

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
        const queryDialog = ref(null);
        const position = ref({ x: 0, y: 0 });
        const isDragging = ref(false);
        const dragOffset = ref({ x: 0, y: 0 });

        const quickPrompts = [
            "Explain this code",
            "What does this mean?",
            "Find the error",
            "Summarize this",
            "Translate this",
            "Improve this code",
        ];

        const isMac = computed(() => {
            return navigator.userAgent.toLowerCase().includes("mac");
        });

        const submitKey = computed(() => {
            return isMac.value ? "Cmd+Enter" : "Ctrl+Enter";
        });

        // Position dialog in top-right corner on mount
        onMounted(() => {
            const margin = 20;
            const dialogWidth = 420;
            position.value = {
                x: window.innerWidth - dialogWidth - margin,
                y: margin,
            };
        });

        // Dragging functionality
        const startDrag = (e) => {
            isDragging.value = true;
            dragOffset.value = {
                x: e.clientX - position.value.x,
                y: e.clientY - position.value.y,
            };

            document.addEventListener("mousemove", onDrag);
            document.addEventListener("mouseup", stopDrag);
        };

        const onDrag = (e) => {
            if (!isDragging.value) return;

            const newX = e.clientX - dragOffset.value.x;
            const newY = e.clientY - dragOffset.value.y;

            // Constrain to window bounds
            const dialogWidth = queryDialog.value?.offsetWidth || 420;
            const dialogHeight = queryDialog.value?.offsetHeight || 600;

            position.value = {
                x: Math.max(0, Math.min(newX, window.innerWidth - dialogWidth)),
                y: Math.max(
                    0,
                    Math.min(newY, window.innerHeight - dialogHeight),
                ),
            };
        };

        const stopDrag = () => {
            isDragging.value = false;
            document.removeEventListener("mousemove", onDrag);
            document.removeEventListener("mouseup", stopDrag);
        };

        onUnmounted(() => {
            document.removeEventListener("mousemove", onDrag);
            document.removeEventListener("mouseup", stopDrag);
        });

        const submit = () => {
            const trimmed = queryText.value.trim();
            if (trimmed) {
                const clamped = trimmed.slice(0, 1000);
                emit("submit", clamped);
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
            queryDialog,
            position,
            isDragging,
            startDrag,
        };
    },
};
</script>

<style scoped>
.query-window-container {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    z-index: 9998;
}

.query-dialog {
    position: fixed;
    width: 420px;
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(20px);
    border-radius: 20px;
    box-shadow: 0 20px 60px rgba(102, 126, 234, 0.3);
    border: 2px solid rgba(102, 126, 234, 0.2);
    overflow: hidden;
    pointer-events: auto;
    animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
    from {
        opacity: 0;
        transform: translateY(-20px) scale(0.95);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

/* Header */
.query-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    user-select: none;
    -webkit-app-region: no-drag;
}

.query-title {
    font-size: 16px;
    font-weight: 700;
    color: white;
    margin: 0;
}

.close-btn {
    width: 28px;
    height: 28px;
    border: none;
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border-radius: 50%;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    transition: all 0.2s ease;
}

.close-btn:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.1);
}

/* Screenshot Preview */
.screenshot-preview {
    padding: 16px 20px 12px;
    background: rgba(102, 126, 234, 0.03);
}

.screenshot-img {
    width: 100%;
    max-height: 180px;
    object-fit: contain;
    border-radius: 12px;
    border: 2px solid rgba(102, 126, 234, 0.15);
    background: white;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

/* Query Input */
.query-input-section {
    padding: 12px 20px;
}

.input-label {
    display: block;
    font-size: 13px;
    font-weight: 600;
    color: #4b5563;
    margin-bottom: 8px;
}

.query-textarea {
    width: 100%;
    padding: 12px;
    border: 2px solid rgba(102, 126, 234, 0.2);
    border-radius: 10px;
    font-size: 14px;
    line-height: 1.5;
    color: #1f2937;
    resize: none;
    transition: all 0.2s ease;
    font-family:
        -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
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

.input-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 8px;
}

.shortcut-hint {
    font-size: 11px;
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

.char-count {
    font-size: 11px;
    color: #9ca3af;
}

/* Quick Prompts */
.quick-prompts {
    padding: 0 20px 12px;
}

.prompts-label {
    font-size: 11px;
    font-weight: 600;
    color: #6b7280;
    margin: 0 0 8px 0;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.prompts-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 6px;
}

.prompt-btn {
    padding: 8px 12px;
    font-size: 12px;
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

/* Actions */
.actions {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 10px;
    padding: 16px 20px;
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
    border-radius: 10px;
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
    border-radius: 10px;
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

.btn-submit .icon {
    width: 18px;
    height: 18px;
}
</style>
