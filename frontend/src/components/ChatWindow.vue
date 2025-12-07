<template>
    <div class="chat-overlay">
        <div class="chat-modal">
            <!-- Header -->
            <ChatHeader
                :title="t('query.title')"
                :showClearButton="messages.length > 0"
                :clearTitle="t('query.clearChat')"
                :closeTitle="t('query.cancel')"
                @clear="clearChat"
                @close="cancel"
            />

            <!-- Body with two-column layout -->
            <div class="modal-body">
                <!-- Left: Screenshot Preview -->
                <ScreenshotPreview :screenshot="screenshot" :placeholderText="t('query.noScreenshot')" />

                <!-- Right: Chat Area -->
                <div class="chat-column">
                    <!-- Messages -->
                    <div class="messages-container" ref="messagesContainer">
                        <EmptyState
                            v-if="messages.length === 0"
                            :message="t('query.emptyChat')"
                            :prompts="quickPrompts"
                            @select-prompt="queryText = $event"
                        />

                        <ChatMessage
                            v-for="(message, index) in messages"
                            :key="index"
                            :content="message.content"
                            :role="message.role"
                            :timestamp="message.timestamp"
                        />

                        <LoadingIndicator v-if="isLoading" />
                    </div>

                    <!-- Input Area -->
                    <ChatInput
                        v-model="queryText"
                        :placeholder="t('query.placeholder')"
                        :disabled="isLoading"
                        :charCountLabel="t('query.charCount')"
                        @submit="submit"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, computed, nextTick } from 'vue';
import { useI18n } from 'vue-i18n';
import ChatHeader from './chat/ChatHeader.vue';
import ChatMessage from './chat/ChatMessage.vue';
import ChatInput from './chat/ChatInput.vue';
import ScreenshotPreview from './chat/ScreenshotPreview.vue';
import EmptyState from './chat/EmptyState.vue';
import LoadingIndicator from './chat/LoadingIndicator.vue';

export default {
    name: 'ChatWindow',
    components: {
        ChatHeader,
        ChatMessage,
        ChatInput,
        ScreenshotPreview,
        EmptyState,
        LoadingIndicator
    },
    props: {
        screenshot: {
            type: String,
            required: true
        }
    },
    emits: ['submit', 'cancel'],
    setup(props, { emit }) {
        const { t } = useI18n();
        const queryText = ref('');
        const messages = ref([]);
        const isLoading = ref(false);
        const messagesContainer = ref(null);

        const quickPrompts = computed(() => [
            t('query.promptExplain'),
            t('query.promptWrong'),
            t('query.promptSummarize'),
            t('query.promptTranslate'),
            t('query.promptImprove'),
            t('query.promptBugs')
        ]);

        const scrollToBottom = async () => {
            await nextTick();
            if (messagesContainer.value) {
                messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
            }
        };

        const submit = async (submitData) => {
            // 支持舊格式（純文字）和新格式（對象）
            const text = typeof submitData === 'string' ? submitData : submitData.text;
            const webSearch = typeof submitData === 'object' ? submitData.webSearch : false;
            
            if (!text || isLoading.value) return;

            messages.value.push({
                role: 'user',
                content: text + (webSearch ? ' 🌐' : ''),
                timestamp: new Date()
            });

            isLoading.value = true;
            scrollToBottom();

            emit('submit', { text, webSearch }, (response) => {
                messages.value.push({
                    role: 'assistant',
                    content: response,
                    timestamp: new Date()
                });
                isLoading.value = false;
                scrollToBottom();
            });
        };

        const cancel = () => {
            emit('cancel');
        };

        const clearChat = () => {
            if (confirm(t('query.clearConfirm'))) {
                messages.value = [];
            }
        };

        return {
            t,
            queryText,
            messages,
            isLoading,
            messagesContainer,
            quickPrompts,
            submit,
            cancel,
            clearChat
        };
    }
};
</script>

<style scoped>
.chat-overlay {
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

.chat-modal {
    background: #ffffff;
    border-radius: 24px;
    box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.05), 0 10px 60px rgba(0, 0, 0, 0.15), 0 30px 100px rgba(102, 126, 234, 0.2);
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

.modal-body {
    flex: 1;
    display: grid;
    grid-template-columns: 1fr 400px;
    gap: 24px;
    padding: 32px;
    overflow: hidden;
    min-height: 0;
}

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

.messages-container::-webkit-scrollbar {
    width: 6px;
}

.messages-container::-webkit-scrollbar-track {
    background: transparent;
}

.messages-container::-webkit-scrollbar-thumb {
    background: rgba(102, 126, 234, 0.2);
    border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb:hover {
    background: rgba(102, 126, 234, 0.3);
}

@media (max-width: 900px) {
    .chat-overlay {
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
    .modal-body {
        padding: 16px;
        grid-template-rows: 200px 1fr;
    }
}
</style>
