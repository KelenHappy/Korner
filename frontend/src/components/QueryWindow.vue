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

            <!-- Body with two-column layout -->
            <div class="modal-body">
                <!-- Left: Screenshot Preview -->
                <div class="screenshot-column">
                    <div v-if="hasScreenshot" class="screenshot-wrapper">
                        <img
                            :src="screenshotPreview"
                            alt="Screenshot"
                            class="screenshot-img"
                        />
                    </div>
                    <div v-else class="screenshot-placeholder">
                        <p>📷</p>
                        <p>Screenshot unavailable</p>
                    </div>
                </div>

                <!-- Right: Query Input -->
                <div class="input-column">
                    <div class="input-section">
                        <label class="input-label">Your Question</label>
                        <textarea
                            v-model="queryText"
                            @keydown.ctrl.enter="submit"
                            @keydown.meta.enter="submit"
                            class="query-textarea"
                            rows="8"
                            placeholder="Ask anything about this screenshot..."
                            autofocus
                        ></textarea>
                        <div class="input-info">
                            <span class="char-count"
                                >{{ queryText.length }} / 1000</span
                            >
                            <span class="shortcut-hint">
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
import { ref, computed, watch, onMounted } from "vue";

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

        const hasScreenshot = computed(() => {
            return (
                typeof props.screenshot === "string" &&
                props.screenshot.trim().length > 0
            );
        });

        const screenshotPreview = computed(() => {
            if (!hasScreenshot.value) {
                return "";
            }
            return props.screenshot.startsWith("data:image")
                ? props.screenshot
                : `data:image/png;base64,${props.screenshot}`;
        });

        onMounted(() => {
            console.log("[Korner][QueryWindow] mounted", {
                hasScreenshot: hasScreenshot.value,
                length: props.screenshot ? props.screenshot.length : 0,
            });
        });

        watch(
            () => props.screenshot,
            (newValue, oldValue) => {
                console.log("[Korner][QueryWindow] screenshot prop changed", {
                    hasScreenshot: !!newValue,
                    newLength: newValue ? newValue.length : 0,
                    oldLength: oldValue ? oldValue.length : 0,
                });
            },
            { immediate: true },
        );

        const isMac = computed(() => {
            return navigator.userAgent.toLowerCase().includes("mac");
        });

        const submitKey = computed(() => {
            return isMac.value ? "⌘ + Enter" : "Ctrl + Enter";
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
            hasScreenshot,
            screenshotPreview,
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
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10000;
    padding: 40px;
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
    border-radius: 20px;
    box-shadow: 0 25px 80px rgba(0, 0, 0, 0.3);
    width: 100%;
    max-width: 1200px;
    max-height: 85vh;
    display: flex;
    flex-direction: column;
    animation: slideUp 0.3s ease-out;
    overflow: hidden;
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(40px) scale(0.95);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

/* Header */
.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 24px 32px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-bottom: none;
    flex-shrink: 0;
}

.modal-title {
    font-size: 20px;
    font-weight: 700;
    color: white;
    margin: 0;
}

.close-btn {
    width: 36px;
    height: 36px;
    border: none;
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border-radius: 10px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 22px;
    transition: all 0.2s ease;
    padding: 0;
}

.close-btn:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.05);
}

/* Body - Two Column Layout */
.modal-body {
    flex: 1;
    display: grid;
    grid-template-columns: 1fr 400px;
    gap: 24px;
    padding: 32px;
    overflow: hidden;
    min-height: 0;
}

/* Left Column: Screenshot */
.screenshot-column {
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: auto;
    border-radius: 12px;
    background: rgba(102, 126, 234, 0.03);
}

.screenshot-wrapper {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 16px;
}

.screenshot-img {
    max-width: 100%;
    max-height: 100%;
    width: auto;
    height: auto;
    object-fit: contain;
    border-radius: 8px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    background: white;
}

.screenshot-placeholder {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #9ca3af;
    gap: 8px;
}

.screenshot-placeholder p {
    margin: 0;
    font-size: 14px;
}

.screenshot-placeholder p:first-child {
    font-size: 48px;
}

/* Right Column: Input */
.input-column {
    display: flex;
    flex-direction: column;
    gap: 20px;
    overflow-y: auto;
}

.input-section {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.input-label {
    font-size: 13px;
    font-weight: 700;
    color: #4b5563;
    text-transform: uppercase;
    letter-spacing: 0.8px;
}

.query-textarea {
    width: 100%;
    padding: 14px;
    border: 2px solid rgba(102, 126, 234, 0.2);
    border-radius: 12px;
    font-size: 15px;
    line-height: 1.6;
    color: #1f2937;
    resize: vertical;
    min-height: 180px;
    font-family:
        -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    transition: all 0.2s ease;
    box-sizing: border-box;
}

.query-textarea:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 4px rgba(102, 126, 234, 0.1);
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
    font-weight: 500;
}

.shortcut-hint {
    color: #6b7280;
    display: flex;
    align-items: center;
    gap: 4px;
}

.kbd {
    display: inline-block;
    padding: 3px 8px;
    background: rgba(102, 126, 234, 0.1);
    border: 1px solid rgba(102, 126, 234, 0.2);
    border-radius: 6px;
    font-size: 11px;
    font-family: monospace;
    color: #667eea;
    font-weight: 600;
}

/* Prompts */
.prompts-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.prompts-label {
    font-size: 12px;
    font-weight: 700;
    color: #6b7280;
    margin: 0;
    text-transform: uppercase;
    letter-spacing: 0.8px;
}

.prompts-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
}

.prompt-btn {
    padding: 10px 14px;
    font-size: 13px;
    background: linear-gradient(
        135deg,
        rgba(102, 126, 234, 0.08) 0%,
        rgba(118, 75, 162, 0.08) 100%
    );
    border: 1px solid rgba(102, 126, 234, 0.15);
    color: #4b5563;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.2s ease;
    font-weight: 600;
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
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

/* Footer */
.modal-footer {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 12px;
    padding: 20px 32px;
    background: rgba(102, 126, 234, 0.03);
    border-top: 1px solid rgba(102, 126, 234, 0.1);
    flex-shrink: 0;
}

.btn-cancel {
    padding: 12px 24px;
    font-size: 15px;
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
    padding: 12px 28px;
    font-size: 15px;
    font-weight: 600;
    color: white;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
}

.btn-submit:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 6px 24px rgba(102, 126, 234, 0.4);
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
    font-size: 18px;
}

.btn-text {
    font-weight: 700;
}

/* Scrollbar styling */
.screenshot-column::-webkit-scrollbar,
.input-column::-webkit-scrollbar {
    width: 8px;
    height: 8px;
}

.screenshot-column::-webkit-scrollbar-track,
.input-column::-webkit-scrollbar-track {
    background: rgba(102, 126, 234, 0.05);
    border-radius: 4px;
}

.screenshot-column::-webkit-scrollbar-thumb,
.input-column::-webkit-scrollbar-thumb {
    background: rgba(102, 126, 234, 0.2);
    border-radius: 4px;
}

.screenshot-column::-webkit-scrollbar-thumb:hover,
.input-column::-webkit-scrollbar-thumb:hover {
    background: rgba(102, 126, 234, 0.3);
}

/* Responsive: Stack on smaller screens */
@media (max-width: 900px) {
    .query-overlay {
        padding: 20px;
    }

    .modal-body {
        grid-template-columns: 1fr;
        grid-template-rows: 300px 1fr;
        gap: 20px;
        padding: 24px;
    }

    .input-column {
        overflow-y: visible;
    }
}

@media (max-width: 600px) {
    .modal-header {
        padding: 16px 20px;
    }

    .modal-title {
        font-size: 16px;
    }

    .modal-body {
        padding: 16px;
        grid-template-rows: 200px 1fr;
    }

    .prompts-grid {
        grid-template-columns: 1fr;
    }

    .modal-footer {
        padding: 12px 16px;
    }
}
</style>
