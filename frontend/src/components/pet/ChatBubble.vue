<template>
    <transition name="bubble">
        <div
            v-if="visible"
            class="chat-bubble"
            :class="{ dragging: isDragging }"
            :style="{ transform: `translate(${position.x}px, ${position.y}px)` }"
        >
            <div class="bubble-header" @pointerdown="startDrag">
                <span class="bubble-title">💬 Response</span>
                <button @click="$emit('close')" class="close-bubble">✕</button>
            </div>
            <div class="bubble-content">
                <!-- Screenshot Preview -->
                <div v-if="screenshot" class="screenshot-display">
                    <img :src="screenshot" alt="Last screenshot" class="response-screenshot" />
                </div>

                <!-- Loading State -->
                <div v-if="loading" class="loading">
                    <div class="loading-dots">
                        <span></span>
                        <span></span>
                        <span></span>
                    </div>
                    <p class="loading-text">Thinking...</p>
                </div>

                <!-- Response Display (Plain text) -->
                <div v-else-if="message" class="chat-text">{{ message }}</div>

                <!-- Empty State -->
                <div v-else class="chat-text empty">Response will appear here</div>

                <!-- Screenshot Action Button -->
                <div v-if="!loading && !message" class="action-area">
                    <button @click="$emit('screenshot')" class="screenshot-action-btn">
                        <span class="btn-icon">📸</span>
                        <span class="btn-label">Take Screenshot</span>
                    </button>
                </div>
            </div>
        </div>
    </transition>
</template>

<script>
import { ref, onUnmounted } from 'vue';

export default {
    name: 'ChatBubble',
    props: {
        visible: {
            type: Boolean,
            default: false
        },
        message: {
            type: String,
            default: ''
        },
        loading: {
            type: Boolean,
            default: false
        },
        screenshot: {
            type: String,
            default: null
        }
    },
    emits: ['close', 'screenshot'],
    setup() {
        const isDragging = ref(false);
        const position = ref({ x: 180, y: 20 });
        
        let dragStartX = 0;
        let dragStartY = 0;
        let dragStartPosX = 0;
        let dragStartPosY = 0;

        const startDrag = (e) => {
            isDragging.value = true;
            dragStartX = e.clientX;
            dragStartY = e.clientY;
            dragStartPosX = position.value.x;
            dragStartPosY = position.value.y;

            document.addEventListener('pointermove', onDrag, { passive: true });
            document.addEventListener('pointerup', stopDrag);
        };

        const onDrag = (e) => {
            if (!isDragging.value) return;

            const deltaX = e.clientX - dragStartX;
            const deltaY = e.clientY - dragStartY;

            let newX = dragStartPosX + deltaX;
            let newY = dragStartPosY + deltaY;

            const chatWidth = 400;
            const chatHeight = 400;
            newX = Math.max(0, Math.min(newX, window.innerWidth - chatWidth));
            newY = Math.max(0, Math.min(newY, window.innerHeight - chatHeight));

            position.value = { x: newX, y: newY };
        };

        const stopDrag = () => {
            isDragging.value = false;
            document.removeEventListener('pointermove', onDrag);
            document.removeEventListener('pointerup', stopDrag);
        };

        onUnmounted(() => {
            document.removeEventListener('pointermove', onDrag);
            document.removeEventListener('pointerup', stopDrag);
        });

        return {
            isDragging,
            position,
            startDrag
        };
    }
};
</script>

<style scoped>
.chat-bubble {
    position: fixed;
    width: 380px;
    background: #ffffff;
    border-radius: 20px;
    box-shadow: 0 12px 40px rgba(102, 126, 234, 0.4);
    pointer-events: auto;
    overflow: hidden;
    border: 2px solid rgba(102, 126, 234, 0.2);
    z-index: 999;
    will-change: transform;
}

.bubble-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 15px;
    background: #000000;
    color: white;
    font-weight: 600;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);
    cursor: move;
}

.bubble-title {
    font-weight: 600;
    font-size: 14px;
}

.close-bubble {
    width: 24px;
    height: 24px;
    border: none;
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border-radius: 50%;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    transition: all 0.2s ease;
    pointer-events: auto;
}

.close-bubble:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.1);
}

.bubble-content {
    padding: 16px;
    min-height: 120px;
    max-height: 400px;
    overflow-y: auto;
}

.screenshot-display {
    width: 100%;
    border-radius: 8px;
    overflow: hidden;
    border: 1px solid rgba(102, 126, 234, 0.1);
    background: rgba(102, 126, 234, 0.05);
    margin-bottom: 12px;
}

.response-screenshot {
    width: 100%;
    max-height: 150px;
    object-fit: contain;
    display: block;
    background: white;
}

.chat-text {
    font-size: 14px;
    line-height: 1.7;
    color: #333;
    word-wrap: break-word;
    white-space: pre-wrap;
}

.chat-text.empty {
    color: #999;
    font-style: italic;
    text-align: center;
}

.loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    height: 100px;
}

.loading-text {
    font-size: 13px;
    color: #222222;
    font-weight: 500;
}

.loading-dots {
    display: flex;
    gap: 6px;
}

.loading-dots span {
    width: 8px;
    height: 8px;
    background: #667eea;
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
    0%, 80%, 100% {
        transform: scale(0);
    }
    40% {
        transform: scale(1);
    }
}

.action-area {
    display: flex;
    justify-content: center;
    padding-top: 12px;
    border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.screenshot-action-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    background: #000000;
    color: #FFFFFF;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-size: 13px;
    font-weight: 600;
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(102, 126, 234, 0.3);
    pointer-events: auto;
}

.screenshot-action-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.screenshot-action-btn:active {
    transform: translateY(0);
}

.btn-icon {
    font-size: 16px;
}

.btn-label {
    font-weight: 600;
}

.bubble-enter-active,
.bubble-leave-active {
    transition: all 0.3s ease;
}

.bubble-enter-from {
    opacity: 0;
    transform: translateX(-20px) scale(0.9);
}

.bubble-leave-to {
    opacity: 0;
    transform: translateX(-20px) scale(0.9);
}

.bubble-content::-webkit-scrollbar {
    width: 6px;
}

.bubble-content::-webkit-scrollbar-track {
    background: rgba(102, 126, 234, 0.1);
    border-radius: 10px;
}

.bubble-content::-webkit-scrollbar-thumb {
    background: #667eea;
    border-radius: 10px;
}

.bubble-content::-webkit-scrollbar-thumb:hover {
    background: #667eea;
}
</style>
