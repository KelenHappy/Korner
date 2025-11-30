<template>
    <div class="query-overlay">
        <div class="query-modal">
            <!-- Header -->
            <div class="modal-header">
                <div class="header-content">
                    <div class="icon-wrapper">
                        <span class="header-icon">✨</span>
                    </div>
                    <h2 class="modal-title">{{ t("query.title") }}</h2>
                </div>
                <div class="header-actions">
                    <button
                        v-if="messages.length > 0"
                        @click="clearChat"
                        class="clear-btn"
                        :title="t('query.clearChat')"
                    >
                        <svg
                            width="18"
                            height="18"
                            viewBox="0 0 18 18"
                            fill="none"
                        >
                            <path
                                d="M3 5H15M7 8V13M11 8V13M4 5L5 15C5 15.5304 5.21071 16.0391 5.58579 16.4142C5.96086 16.7893 6.46957 17 7 17H11C11.5304 17 12.0391 16.7893 12.4142 16.4142C12.7893 16.0391 13 15.5304 13 15L14 5M6 5V3C6 2.73478 6.10536 2.48043 6.29289 2.29289C6.48043 2.10536 6.73478 2 7 2H11C11.2652 2 11.5196 2.10536 11.7071 2.29289C11.8946 2.48043 12 2.73478 12 3V5"
                                stroke="currentColor"
                                stroke-width="1.5"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                        </svg>
                    </button>
                    <button @click="cancel" class="close-btn" :title="t('query.cancel')">
                        <svg
                            width="20"
                            height="20"
                            viewBox="0 0 20 20"
                            fill="none"
                        >
                            <path
                                d="M15 5L5 15M5 5L15 15"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                            />
                        </svg>
                    </button>
                </div>
            </div>

            <!-- Body with two-column layout -->
            <div class="modal-body">
                <!-- Left: Screenshot Preview -->
                <div class="screenshot-column">
                    <div v-if="hasScreenshot" class="screenshot-wrapper">
                        <div class="screenshot-frame">
                            <img
                                :src="screenshotPreview"
                                alt="Screenshot"
                                class="screenshot-img"
                            />
                        </div>
                    </div>
                    <div v-else class="screenshot-placeholder">
                        <div class="placeholder-icon">
                            <svg
                                width="64"
                                height="64"
                                viewBox="0 0 64 64"
                                fill="none"
                            >
                                <rect
                                    x="8"
                                    y="12"
                                    width="48"
                                    height="40"
                                    rx="4"
                                    stroke="currentColor"
                                    stroke-width="2"
                                />
                                <circle
                                    cx="32"
                                    cy="28"
                                    r="6"
                                    stroke="currentColor"
                                    stroke-width="2"
                                />
                                <path
                                    d="M8 44L18 34L28 44L42 30L56 44V48C56 50.2091 54.2091 52 52 52H12C9.79086 52 8 50.2091 8 48V44Z"
                                    fill="currentColor"
                                    opacity="0.2"
                                />
                            </svg>
                        </div>
                        <p class="placeholder-text">{{ t("query.noScreenshot") }}</p>
                    </div>
                </div>

                <!-- Right: Chat Area -->
                <div class="chat-column">
                    <!-- Messages -->
                    <div class="messages-container" ref="messagesContainer">
                        <div v-if="messages.length === 0" class="empty-state">
                            <div class="empty-icon">💭</div>
                            <p class="empty-text">
                                {{ t("query.emptyChat") }}
                            </p>
                            <div class="prompts-grid">
                                <button
                                    v-for="prompt in quickPrompts"
                                    :key="prompt"
                                    @click="queryText = prompt"
                                    class="prompt-btn"
                                >
                                    <span class="prompt-text">{{
                                        prompt
                                    }}</span>
                                </button>
                            </div>
                        </div>

                        <div
                            v-for="(message, index) in messages"
                            :key="index"
                            :class="[
                                'message',
                                message.role === 'user'
                                    ? 'message-user'
                                    : 'message-ai',
                            ]"
                        >
                            <div class="message-avatar">
                                <span v-if="message.role === 'user'">👤</span>
                                <span v-else>✨</span>
                            </div>
                            <div class="message-content">
                                <div class="message-text">
                                    {{ message.content }}
                                </div>
                                <div class="message-time">
                                    {{ formatTime(message.timestamp) }}
                                </div>
                            </div>
                        </div>

                        <div v-if="isLoading" class="message message-ai">
                            <div class="message-avatar">
                                <span>✨</span>
                            </div>
                            <div class="message-content">
                                <div class="loading-dots">
                                    <span></span>
                                    <span></span>
                                    <span></span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Input Area -->
                    <div class="input-area">
                        <textarea
                            v-model="queryText"
                            @keydown.ctrl.enter="submit"
                            @keydown.meta.enter="submit"
                            class="chat-input"
                            rows="3"
                            :placeholder="t('query.placeholder')"
                            :disabled="isLoading"
                        ></textarea>
                        <div class="input-actions">
                            <span class="char-count"
                                >{{ queryText.length }} / 1000 {{ t("query.charCount") }}</span
                            >
                            <button
                                @click="submit"
                                :disabled="!queryText.trim() || isLoading"
                                class="send-btn"
                            >
                                <svg
                                    width="20"
                                    height="20"
                                    viewBox="0 0 20 20"
                                    fill="none"
                                >
                                    <path
                                        d="M2 10L18 2L10 18L8 11L2 10Z"
                                        fill="currentColor"
                                    />
                                </svg>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, watch, onMounted, nextTick } from "vue";
import { useI18n } from "vue-i18n";

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
        const { t } = useI18n();
        const queryText = ref("");
        const messages = ref([]);
        const isLoading = ref(false);
        const messagesContainer = ref(null);

        const quickPrompts = computed(() => [
            t("query.promptExplain"),
            t("query.promptWrong"),
            t("query.promptSummarize"),
            t("query.promptTranslate"),
            t("query.promptImprove"),
            t("query.promptBugs"),
        ]);

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

        const scrollToBottom = async () => {
            await nextTick();
            if (messagesContainer.value) {
                messagesContainer.value.scrollTop =
                    messagesContainer.value.scrollHeight;
            }
        };

        const submit = async () => {
            const trimmed = queryText.value.trim();
            if (!trimmed || isLoading.value) return;

            // Add user message
            messages.value.push({
                role: "user",
                content: trimmed.slice(0, 1000),
                timestamp: new Date(),
            });

            queryText.value = "";
            isLoading.value = true;
            scrollToBottom();

            // Emit for backend processing
            emit("submit", trimmed.slice(0, 1000), (response) => {
                // Add AI response
                messages.value.push({
                    role: "assistant",
                    content: response,
                    timestamp: new Date(),
                });
                isLoading.value = false;
                scrollToBottom();
            });
        };

        const cancel = () => {
            emit("cancel");
        };

        const clearChat = () => {
            if (confirm(t("query.clearConfirm"))) {
                messages.value = [];
            }
        };

        const formatTime = (date) => {
            return date.toLocaleTimeString([], {
                hour: "2-digit",
                minute: "2-digit",
            });
        };

        return {
            t,
            queryText,
            messages,
            isLoading,
            messagesContainer,
            quickPrompts,
            submitKey,
            submit,
            cancel,
            clearChat,
            formatTime,
            hasScreenshot,
            screenshotPreview,
        };
    },
};
</script>

<style scoped>
.query-overlay {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 10000;
    animation: fadeIn 0.2s ease-out;
    pointer-events: auto;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        backdrop-filter: blur(0px);
    }
    to {
        opacity: 1;
        backdrop-filter: blur(8px);
    }
}

.query-modal {
    background: #ffffff;
    border-radius: 24px;
    box-shadow:
        0 0 0 1px rgba(0, 0, 0, 0.05),
        0 10px 60px rgba(0, 0, 0, 0.15),
        0 30px 100px rgba(102, 126, 234, 0.2);
    width: 1200px;
    max-width: 90vw;
    max-height: 85vh;
    display: flex;
    flex-direction: column;
    animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
    overflow: hidden;
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(50px) scale(0.96);
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
    background: #f5f5f5;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    flex-shrink: 0;
}

.header-actions {
    display: flex;
    align-items: center;
    gap: 8px;
}

.clear-btn {
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

.close-btn:hover {
    background: rgba(100, 116, 139, 0.1);
    color: #1a1a1a;
    transform: rotate(90deg);
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
    box-shadow:
        0 0 0 1px rgba(102, 126, 234, 0.1),
        0 8px 32px rgba(0, 0, 0, 0.08);
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

/* Right Column: Chat */
.chat-column {
    display: flex;
    flex-direction: column;
    gap: 0;
    overflow: hidden;
}

.messages-container {
    flex: 1;
    overflow-y: auto;
    padding: 24px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

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

.message {
    display: flex;
    gap: 12px;
    animation: messageSlideIn 0.3s ease;
}

@keyframes messageSlideIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.message-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    flex-shrink: 0;
}

.message-user .message-avatar {
    background: #ffffff;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
}

.message-ai .message-avatar {
    background: #cccccc;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
}

.message-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.message-text {
    padding: 12px 16px;
    border-radius: 12px;
    font-size: 14px;
    line-height: 1.6;
    white-space: pre-wrap;
    word-wrap: break-word;
}

.message-user .message-text {
    background: #ffffff;
    color: #000000;
    border-top-left-radius: 4px;
}

.message-ai .message-text {
    background: #ffffff;
    color: #6b7280;
    border: 1px solid rgba(0, 0, 0, 0.06);
    border-top-left-radius: 4px;
}

.message-time {
    font-size: 11px;
    color: #94a3b8;
    padding-left: 4px;
}

.loading-dots {
    display: flex;
    gap: 6px;
    padding: 12px 16px;
}

.loading-dots span {
    width: 8px;
    height: 8px;
    background: #f5576c;
    border-radius: 50%;
    animation: bounce 1.4s infinite ease-in-out both;
}

.loading-dots span:nth-child(1) {
    animation-delay: -0.32s;
}

.loading-dots span:nth-child(2) {
    animation-delay: -0.16s;
}

@keyframes bounce {
    0%,
    80%,
    100% {
        transform: scale(0);
    }
    40% {
        transform: scale(1);
    }
}

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
    font-family:
        -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    box-sizing: border-box;
    background: white;
    margin-bottom: 12px;
}

.chat-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow:
        0 0 0 1px #667eea,
        0 0 0 4px rgba(102, 126, 234, 0.1);
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
    box-shadow:
        0 0 0 1px rgba(102, 126, 234, 0.1),
        0 4px 12px rgba(102, 126, 234, 0.3);
}

.send-btn:hover:not(:disabled) {
    transform: scale(1.1);
    box-shadow:
        0 0 0 1px rgba(102, 126, 234, 0.2),
        0 6px 16px rgba(102, 126, 234, 0.4);
}

.send-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: scale(1);
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
    box-shadow:
        0 4px 12px rgba(102, 126, 234, 0.25),
        0 0 0 1px rgba(102, 126, 234, 0.1);
}

.prompt-btn:hover::before {
    opacity: 1;
}

.prompt-btn:active {
    transform: translateY(0);
}

/* Scrollbar styling */
.screenshot-column::-webkit-scrollbar,
.messages-container::-webkit-scrollbar {
    width: 6px;
}

.screenshot-column::-webkit-scrollbar-track,
.messages-container::-webkit-scrollbar-track {
    background: transparent;
}

.screenshot-column::-webkit-scrollbar-thumb,
.messages-container::-webkit-scrollbar-thumb {
    background: rgba(102, 126, 234, 0.2);
    border-radius: 3px;
}

.screenshot-column::-webkit-scrollbar-thumb:hover,
.messages-container::-webkit-scrollbar-thumb:hover {
    background: rgba(102, 126, 234, 0.3);
}

/* Responsive: Stack on smaller screens */
@media (max-width: 900px) {
    .query-overlay {
        padding: 20px;
    }

    .modal-body {
        grid-template-columns: 1fr;
        grid-template-rows: 200px 1fr;
        gap: 20px;
        padding: 24px;
    }

    .chat-column {
        min-height: 400px;
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
