<template>
    <div class="desktop-pet-wrapper">
        <!-- Pet Avatar -->
        <PetAvatar
            :menuActive="showMenu"
            @screenshot="handleScreenshot"
            @toggle-menu="toggleMenu"
            @position-change="updatePetPosition"
        />

        <!-- Hover Menu -->
        <PetMenu
            :visible="showMenu"
            :position="petPosition"
            :chatVisible="showChatBubble"
            @screenshot="handleScreenshot"
            @toggle-chat="toggleChatBubble"
            @settings="handleSettings"
            @minimize="handleMinimize"
        />

        <!-- Chat Bubble -->
        <ChatBubble
            :visible="showChatBubble"
            :message="chatMessage"
            :loading="chatLoading"
            :screenshot="currentScreenshot"
            @close="closeChatBubble"
            @screenshot="handleScreenshot"
        />

        <!-- Status indicator -->
        <transition name="fade">
            <div
                class="status-indicator"
                v-if="status"
                :style="{ transform: `translate(${petPosition.x}px, ${petPosition.y - 50}px)` }"
            >
                {{ status }}
            </div>
        </transition>
    </div>
</template>

<script>
import { ref, watch } from 'vue';
import PetAvatar from './pet/PetAvatar.vue';
import PetMenu from './pet/PetMenu.vue';
import ChatBubble from './pet/ChatBubble.vue';

export default {
    name: 'DesktopPet',
    components: {
        PetAvatar,
        PetMenu,
        ChatBubble
    },
    props: {
        status: {
            type: String,
            default: ''
        },
        chatMessage: {
            type: String,
            default: ''
        },
        chatLoading: {
            type: Boolean,
            default: false
        },
        lastScreenshot: {
            type: String,
            default: null
        }
    },
    emits: ['screenshot', 'settings', 'minimize'],
    setup(props, { emit }) {
        const showMenu = ref(false);
        const showChatBubble = ref(false);
        const currentScreenshot = ref(null);
        const petPosition = ref({ x: 20, y: 20 });

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
            emit('screenshot');
        };

        const handleSettings = () => {
            showMenu.value = false;
            emit('settings');
        };

        const handleMinimize = () => {
            showMenu.value = false;
            emit('minimize');
        };

        const updatePetPosition = (pos) => {
            petPosition.value = pos;
        };

        // Auto-show chat bubble when loading or message arrives
        watch(
            () => [props.chatMessage, props.chatLoading],
            ([newMessage, newLoading]) => {
                if (newLoading || newMessage) {
                    if (!showChatBubble.value) {
                        showChatBubble.value = true;
                    }
                }
            },
            { immediate: true }
        );

        // Update screenshot when prop changes
        watch(
            () => props.lastScreenshot,
            (newScreenshot) => {
                if (newScreenshot) {
                    currentScreenshot.value = newScreenshot;
                }
            }
        );

        return {
            showMenu,
            showChatBubble,
            currentScreenshot,
            petPosition,
            toggleMenu,
            toggleChatBubble,
            closeChatBubble,
            handleScreenshot,
            handleSettings,
            handleMinimize,
            updatePetPosition
        };
    }
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

.fade-enter-active,
.fade-leave-active {
    transition: all 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
    transform: translateY(-5px);
}
</style>
