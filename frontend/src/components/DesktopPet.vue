<template>
    <div class="desktop-pet-container" :class="{ 'with-chat': showChatBubble }">
        <!-- Draggable pet avatar -->
        <div
            class="pet-avatar"
            @mousedown="startDrag"
            :style="{ cursor: isDragging ? 'grabbing' : 'grab' }"
        >
            <!-- Pet character -->
            <div class="pet-character">
                <div class="pet-face">
                    <div class="eyes">
                        <div class="eye left">
                            <div class="pupil"></div>
                        </div>
                        <div class="eye right">
                            <div class="pupil"></div>
                        </div>
                    </div>
                    <div class="mouth"></div>
                </div>

                <!-- Floating action button -->
                <button
                    @click="handleScreenshot"
                    class="action-btn screenshot-btn"
                    title="Take Screenshot"
                >
                    📸
                </button>
            </div>

            <!-- Menu toggle button -->
            <button
                @click="toggleMenu"
                class="menu-toggle"
                :class="{ active: showMenu }"
            >
                {{ showMenu ? "✕" : "☰" }}
            </button>
        </div>

        <!-- Chat bubble (appears next to pet) -->
        <transition name="bubble">
            <div class="chat-bubble" v-if="showChatBubble">
                <div class="bubble-header">
                    <span class="bubble-title">💬 Chat</span>
                    <button @click="closeChatBubble" class="close-bubble">
                        ✕
                    </button>
                </div>
                <div class="bubble-content">
                    <div v-if="chatLoading" class="loading">
                        <div class="loading-dots">
                            <span></span>
                            <span></span>
                            <span></span>
                        </div>
                        <p class="loading-text">Thinking...</p>
                    </div>
                    <div v-else-if="chatMessage" class="chat-text">
                        {{ chatMessage }}
                    </div>
                    <div v-else class="chat-text empty">
                        Hi! Click 📸 to take a screenshot!
                    </div>
                </div>
            </div>
        </transition>

        <!-- Compact menu (appears when menu button clicked) -->
        <transition name="menu">
            <div class="hover-menu" v-if="showMenu">
                <button @click="handleScreenshot" class="menu-item">
                    <span class="icon">📸</span>
                    <span class="text">Screenshot</span>
                </button>
                <button @click="toggleChatBubble" class="menu-item">
                    <span class="icon">💬</span>
                    <span class="text">{{
                        showChatBubble ? "Hide Chat" : "Show Chat"
                    }}</span>
                </button>
                <button @click="handleSettings" class="menu-item">
                    <span class="icon">⚙️</span>
                    <span class="text">Settings</span>
                </button>
                <button @click="handleMinimize" class="menu-item">
                    <span class="icon">➖</span>
                    <span class="text">Minimize</span>
                </button>
                <button @click="handleClose" class="menu-item close">
                    <span class="icon">✕</span>
                    <span class="text">Close</span>
                </button>
            </div>
        </transition>

        <!-- Status indicator -->
        <transition name="fade">
            <div class="status-indicator" v-if="status">
                {{ status }}
            </div>
        </transition>
    </div>
</template>

<script>
import { ref, watch, onMounted } from "vue";
import { WindowSetSize } from "../../wailsjs/runtime/runtime";

export default {
    name: "DesktopPet",
    emits: ["screenshot", "settings", "minimize", "close"],
    props: {
        status: {
            type: String,
            default: "",
        },
        chatMessage: {
            type: String,
            default: "",
        },
        chatLoading: {
            type: Boolean,
            default: false,
        },
    },
    setup(props, { emit }) {
        const showMenu = ref(false);
        const showChatBubble = ref(false);
        const isDragging = ref(false);

        const toggleMenu = () => {
            showMenu.value = !showMenu.value;
        };

        const toggleChatBubble = () => {
            showChatBubble.value = !showChatBubble.value;
            showMenu.value = false;
            updateWindowSize();
        };

        const closeChatBubble = () => {
            showChatBubble.value = false;
            updateWindowSize();
        };

        const updateWindowSize = () => {
            try {
                if (showChatBubble.value) {
                    // Pet (150px) + Chat bubble (320px) + spacing (30px)
                    WindowSetSize(500, 240);
                } else {
                    // Just pet
                    WindowSetSize(220, 220);
                }
            } catch (e) {
                // Ignore if runtime not available
            }
        };

        const startDrag = (e) => {
            isDragging.value = true;
            if (window.wails && window.wails.Window) {
                window.wails.Window.StartDrag();
            }
        };

        const handleScreenshot = () => {
            showMenu.value = false;
            emit("screenshot");
        };

        const handleSettings = () => {
            showMenu.value = false;
            emit("settings");
        };

        const handleMinimize = () => {
            showMenu.value = false;
            emit("minimize");
        };

        const handleClose = () => {
            showMenu.value = false;
            emit("close");
        };

        // Watch for chat message or loading changes to auto-show bubble
        watch(
            () => [props.chatMessage, props.chatLoading],
            ([newMessage, newLoading]) => {
                // Auto-show when loading starts or message arrives
                if (newLoading || newMessage) {
                    if (!showChatBubble.value) {
                        showChatBubble.value = true;
                        updateWindowSize();
                    }
                } else {
                    // Hide chatbox when no message and not loading
                    if (showChatBubble.value) {
                        showChatBubble.value = false;
                        updateWindowSize();
                    }
                }
            },
            { immediate: true },
        );

        onMounted(() => {
            updateWindowSize();
        });

        return {
            showMenu,
            showChatBubble,
            isDragging,
            toggleMenu,
            toggleChatBubble,
            closeChatBubble,
            startDrag,
            handleScreenshot,
            handleSettings,
            handleMinimize,
            handleClose,
        };
    },
};
</script>

<style scoped>
.desktop-pet-container {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: flex-start;
    gap: 30px;
    padding: 20px;
}

/* Pet Avatar */
.pet-avatar {
    position: relative;
    width: 150px;
    height: 150px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: grab;
    user-select: none;
    -webkit-app-region: drag;
}

.pet-avatar:active {
    cursor: grabbing;
}

/* Pet Character */
.pet-character {
    position: relative;
    width: 120px;
    height: 120px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 50%;
    box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    transition:
        transform 0.3s ease,
        box-shadow 0.3s ease;
    animation: float 3s ease-in-out infinite;
}

.pet-character:hover {
    transform: scale(1.05);
    box-shadow: 0 15px 40px rgba(102, 126, 234, 0.6);
}

@keyframes float {
    0%,
    100% {
        transform: translateY(0px);
    }
    50% {
        transform: translateY(-10px);
    }
}

/* Face */
.pet-face {
    position: relative;
    width: 80%;
    height: 80%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

/* Eyes */
.eyes {
    display: flex;
    gap: 20px;
    margin-bottom: 10px;
}

.eye {
    width: 20px;
    height: 20px;
    background: white;
    border-radius: 50%;
    position: relative;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
    animation: blink 4s infinite;
}

@keyframes blink {
    0%,
    48%,
    52%,
    100% {
        transform: scaleY(1);
    }
    50% {
        transform: scaleY(0.1);
    }
}

.pupil {
    width: 10px;
    height: 10px;
    background: #333;
    border-radius: 50%;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    transition: all 0.1s ease;
}

/* Mouth */
.mouth {
    width: 30px;
    height: 15px;
    border: 3px solid white;
    border-top: none;
    border-radius: 0 0 30px 30px;
    margin-top: 5px;
}

/* Action Button */
.action-btn {
    position: absolute;
    bottom: -5px;
    right: -5px;
    width: 45px;
    height: 45px;
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: 3px solid white;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
    font-size: 22px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    -webkit-app-region: no-drag;
    z-index: 10;
}

.action-btn:hover {
    transform: scale(1.15) rotate(10deg);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.action-btn:active {
    transform: scale(0.95);
}

/* Menu Toggle */
.menu-toggle {
    position: absolute;
    bottom: -5px;
    left: -5px;
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: 2px solid rgba(102, 126, 234, 0.3);
    font-size: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    -webkit-app-region: no-drag;
    z-index: 10;
}

.menu-toggle:hover {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    transform: scale(1.1) rotate(90deg);
    border-color: transparent;
}

.menu-toggle.active {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    transform: rotate(180deg);
}

/* Chat Bubble */
.chat-bubble {
    position: relative;
    flex: 1;
    min-width: 280px;
    max-width: 320px;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    border-radius: 20px;
    box-shadow: 0 12px 40px rgba(102, 126, 234, 0.4);
    padding: 0;
    -webkit-app-region: no-drag;
    overflow: hidden;
    border: 2px solid rgba(102, 126, 234, 0.2);
}

.chat-bubble::before {
    content: "";
    position: absolute;
    left: -10px;
    top: 50%;
    transform: translateY(-50%);
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 10px 10px 10px 0;
    border-color: transparent rgba(255, 255, 255, 0.98) transparent transparent;
}

.bubble-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
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
}

.close-bubble:hover {
    background: rgba(255, 255, 255, 0.3);
    transform: scale(1.1);
}

.bubble-content {
    padding: 16px;
    min-height: 100px;
    max-height: 180px;
    overflow-y: auto;
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
    color: #667eea;
    font-weight: 500;
}

.loading-dots {
    display: flex;
    gap: 6px;
}

.loading-dots span {
    width: 8px;
    height: 8px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

/* Hover Menu */
.hover-menu {
    position: absolute;
    top: 50%;
    left: calc(100% + 10px);
    transform: translateY(-50%);
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    padding: 8px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 160px;
    z-index: 100;
    -webkit-app-region: no-drag;
}

.menu-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 14px;
    border: none;
    background: transparent;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 14px;
    color: #333;
    text-align: left;
}

.menu-item:hover {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    transform: translateX(4px);
}

.menu-item.close:hover {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.menu-item .icon {
    font-size: 18px;
    width: 24px;
    text-align: center;
}

.menu-item .text {
    flex: 1;
    font-weight: 500;
}

/* Status Indicator */
.status-indicator {
    position: absolute;
    top: -35px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    padding: 8px 16px;
    border-radius: 20px;
    font-size: 12px;
    color: #667eea;
    font-weight: 600;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    white-space: nowrap;
    z-index: 10;
}

/* Transitions */
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

.menu-enter-active,
.menu-leave-active {
    transition: all 0.3s ease;
}

.menu-enter-from {
    opacity: 0;
    transform: translateY(-50%) translateX(-10px);
}

.menu-leave-to {
    opacity: 0;
    transform: translateY(-50%) translateX(-10px);
}

.fade-enter-active,
.fade-leave-active {
    transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: translateX(-50%) translateY(-5px);
}

/* Scrollbar styling */
.bubble-content::-webkit-scrollbar {
    width: 6px;
}

.bubble-content::-webkit-scrollbar-track {
    background: rgba(102, 126, 234, 0.1);
    border-radius: 10px;
}

.bubble-content::-webkit-scrollbar-thumb {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 10px;
}

.bubble-content::-webkit-scrollbar-thumb:hover {
    background: linear-gradient(135deg, #764ba2 0%, #667eea 100%);
}
</style>
