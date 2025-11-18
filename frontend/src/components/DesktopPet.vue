<template>
    <div class="desktop-pet-wrapper">
        <!-- Draggable pet avatar -->
        <div
            class="pet-container"
            :class="{ dragging: isDraggingPet }"
            :style="{
                transform: `translate(${petPosition.x}px, ${petPosition.y}px)`,
            }"
            @mousedown="startDragPet"
        >
            <!-- Pet character with logo -->
            <div class="pet-character">
                <img src="/icon.png" alt="Korner Pet" class="pet-logo" />

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

        <!-- Chat bubble (appears next to pet, independently draggable) -->
        <transition name="bubble">
            <div
                class="chat-bubble"
                v-if="showChatBubble"
                :class="{ dragging: isDraggingChat }"
                :style="{
                    transform: `translate(${chatPosition.x}px, ${chatPosition.y}px)`,
                }"
            >
                <div class="bubble-header" @mousedown="startDragChat">
                    <span class="bubble-title">💬 Response</span>
                    <button @click="closeChatBubble" class="close-bubble">
                        ✕
                    </button>
                </div>
                <div class="bubble-content">
                    <!-- Screenshot Preview -->
                    <div v-if="currentScreenshot" class="screenshot-display">
                        <img
                            :src="currentScreenshot"
                            alt="Last screenshot"
                            class="response-screenshot"
                        />
                    </div>

                    <!-- Loading State -->
                    <div v-if="chatLoading" class="loading">
                        <div class="loading-dots">
                            <span></span>
                            <span></span>
                            <span></span>
                        </div>
                        <p class="loading-text">Thinking...</p>
                    </div>

                    <!-- Response Display -->
                    <div v-else-if="chatMessage" class="chat-text">
                        {{ chatMessage }}
                    </div>

                    <!-- Empty State -->
                    <div v-else class="chat-text empty">
                        Response will appear here
                    </div>

                    <!-- Screenshot Action Button -->
                    <div
                        v-if="!chatLoading && !chatMessage"
                        class="action-area"
                    >
                        <button
                            @click="handleScreenshot"
                            class="screenshot-action-btn"
                        >
                            <span class="btn-icon">📸</span>
                            <span class="btn-label">Take Screenshot</span>
                        </button>
                    </div>
                </div>
            </div>
        </transition>

        <!-- Hover Menu -->
        <transition name="menu">
            <div
                class="hover-menu"
                v-if="showMenu"
                :style="{
                    transform: `translate(${petPosition.x}px, ${petPosition.y}px)`,
                }"
            >
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
                <button @click="handleMinimize" class="menu-item close">
                    <span class="icon">✕</span>
                    <span class="text">Minimize</span>
                </button>
            </div>
        </transition>

        <!-- Status indicator -->
        <transition name="fade">
            <div
                class="status-indicator"
                v-if="status"
                :style="{
                    transform: `translate(${petPosition.x}px, ${petPosition.y - 50}px)`,
                }"
            >
                {{ status }}
            </div>
        </transition>
    </div>
</template>

<script>
import { ref, watch, onUnmounted } from "vue";

export default {
    name: "DesktopPet",
    emits: ["screenshot", "settings", "minimize"],
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
        lastScreenshot: {
            type: String,
            default: null,
        },
    },
    setup(props, { emit }) {
        const showMenu = ref(false);
        const showChatBubble = ref(false);
        const currentScreenshot = ref(null);

        // Position tracking
        const petPosition = ref({ x: 20, y: 20 });
        const chatPosition = ref({ x: 180, y: 20 });

        // Drag state
        const isDraggingPet = ref(false);
        const isDraggingChat = ref(false);
        let dragStartX = 0;
        let dragStartY = 0;
        let dragStartPetX = 0;
        let dragStartPetY = 0;
        let dragStartChatX = 0;
        let dragStartChatY = 0;

        const toggleMenu = () => {
            showMenu.value = !showMenu.value;
        };

        const toggleChatBubble = () => {
            showChatBubble.value = !showChatBubble.value;
            showMenu.value = false;
        };

        const closeChatBubble = () => {
            showChatBubble.value = false;
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

        // Watch for chat message or loading changes to auto-show bubble
        watch(
            () => [props.chatMessage, props.chatLoading],
            ([newMessage, newLoading]) => {
                // Auto-show when loading starts or message arrives
                if (newLoading || newMessage) {
                    if (!showChatBubble.value) {
                        showChatBubble.value = true;
                    }
                }
            },
            { immediate: true },
        );

        // Watch for lastScreenshot prop changes
        watch(
            () => props.lastScreenshot,
            (newScreenshot) => {
                if (newScreenshot) {
                    currentScreenshot.value = newScreenshot;
                }
            },
        );

        // Dragging functions for pet
        const startDragPet = (e) => {
            // Don't drag if clicking on buttons
            if (e.target.closest("button")) return;

            isDraggingPet.value = true;
            dragStartX = e.clientX;
            dragStartY = e.clientY;
            dragStartPetX = petPosition.value.x;
            dragStartPetY = petPosition.value.y;

            document.addEventListener("mousemove", onDragPet);
            document.addEventListener("mouseup", stopDragPet);
        };

        const onDragPet = (e) => {
            if (!isDraggingPet.value) return;

            const deltaX = e.clientX - dragStartX;
            const deltaY = e.clientY - dragStartY;

            petPosition.value.x = Math.max(
                0,
                Math.min(dragStartPetX + deltaX, window.innerWidth - 150),
            );
            petPosition.value.y = Math.max(
                0,
                Math.min(dragStartPetY + deltaY, window.innerHeight - 150),
            );
        };

        const stopDragPet = () => {
            isDraggingPet.value = false;
            document.removeEventListener("mousemove", onDragPet);
            document.removeEventListener("mouseup", stopDragPet);
        };

        // Dragging functions for chat
        const startDragChat = (e) => {
            isDraggingChat.value = true;
            dragStartX = e.clientX;
            dragStartY = e.clientY;
            dragStartChatX = chatPosition.value.x;
            dragStartChatY = chatPosition.value.y;

            document.addEventListener("mousemove", onDragChat);
            document.addEventListener("mouseup", stopDragChat);
        };

        const onDragChat = (e) => {
            if (!isDraggingChat.value) return;

            const deltaX = e.clientX - dragStartX;
            const deltaY = e.clientY - dragStartY;

            chatPosition.value.x = Math.max(
                0,
                Math.min(dragStartChatX + deltaX, window.innerWidth - 400),
            );
            chatPosition.value.y = Math.max(
                0,
                Math.min(dragStartChatY + deltaY, window.innerHeight - 400),
            );
        };

        const stopDragChat = () => {
            isDraggingChat.value = false;
            document.removeEventListener("mousemove", onDragChat);
            document.removeEventListener("mouseup", stopDragChat);
        };

        onUnmounted(() => {
            document.removeEventListener("mousemove", onDragPet);
            document.removeEventListener("mouseup", stopDragPet);
            document.removeEventListener("mousemove", onDragChat);
            document.removeEventListener("mouseup", stopDragChat);
        });

        return {
            showMenu,
            showChatBubble,
            currentScreenshot,
            petPosition,
            chatPosition,
            toggleMenu,
            toggleChatBubble,
            closeChatBubble,
            handleScreenshot,
            handleSettings,
            handleMinimize,
            startDragPet,
            startDragChat,
            isDraggingPet,
            isDraggingChat,
        };
    },
};
</script>

<style scoped>
.desktop-pet-wrapper {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    overflow: hidden;
    z-index: 1000;
}

.pet-container {
    position: fixed;
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: grab;
    user-select: none;
    pointer-events: auto;
    transition: transform 0.1s ease-out;
}

.pet-container.dragging {
    transition: none;
}

.pet-container:active {
    cursor: grabbing;
}

/* Pet Character */
.pet-character {
    position: relative;
    width: 60px;
    height: 60px;
    background: transparent;
    border-radius: 50%;
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    transition:
        transform 0.3s ease,
        box-shadow 0.3s ease;
    animation: float 3s ease-in-out infinite;
    overflow: hidden;
}

.pet-character:hover {
    transform: scale(1.1);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.pet-logo {
    width: 100%;
    height: 100%;
    object-fit: cover;
    pointer-events: none;
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

/* Action Button */
.action-btn {
    position: absolute;
    bottom: -3px;
    right: -3px;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: 2px solid white;
    box-shadow: 0 3px 10px rgba(102, 126, 234, 0.4);
    font-size: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    pointer-events: auto;
    z-index: 10;
}

.action-btn:hover {
    transform: scale(1.15) rotate(15deg);
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.6);
}

.action-btn:active {
    transform: scale(0.95);
}

/* Menu Toggle */
.menu-toggle {
    position: absolute;
    bottom: -3px;
    left: -3px;
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: 2px solid rgba(102, 126, 234, 0.3);
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    pointer-events: auto;
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
    position: fixed;
    width: 380px;
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(20px);
    border-radius: 20px;
    box-shadow: 0 12px 40px rgba(102, 126, 234, 0.4);
    pointer-events: auto;
    overflow: hidden;
    border: 2px solid rgba(102, 126, 234, 0.2);
    z-index: 999;
}

.chat-bubble.dragging {
    transition: none;
}

.chat-bubble::before {
    content: "";
    position: absolute;
    left: -10px;
    top: 50px;
    width: 0;
    height: 0;
    border-style: solid;
    border-width: 10px 10px 10px 0;
    border-color: transparent rgba(255, 255, 255, 0.98) transparent transparent;
    z-index: -1;
}

.bubble-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 15px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

/* Screenshot Display */
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

/* Chat Text */
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

/* Loading */
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

/* Action Area */
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
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
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

/* Hover Menu */
.hover-menu {
    position: fixed;
    top: 0;
    left: 0;
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
    pointer-events: auto;
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
    pointer-events: auto;
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
    position: fixed;
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
    pointer-events: none;
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
    transform: translateY(-10px);
}

.menu-leave-to {
    opacity: 0;
    transform: translateY(-10px);
}

.fade-enter-active,
.fade-leave-active {
    transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: translateY(-5px);
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
