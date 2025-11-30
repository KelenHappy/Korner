<template>
    <transition name="menu">
        <div
            v-if="visible"
            class="hover-menu"
            :style="{ transform: `translate(${position.x}px, ${position.y}px)` }"
        >
            <button @click="$emit('screenshot')" class="menu-item">
                <span class="icon">📸</span>
                <span class="text">Screenshot</span>
            </button>
            <button @click="$emit('toggle-chat')" class="menu-item">
                <span class="icon">💬</span>
                <span class="text">{{ chatVisible ? 'Hide Chat' : 'Show Chat' }}</span>
            </button>
            <button @click="$emit('settings')" class="menu-item">
                <span class="icon">⚙️</span>
                <span class="text">Settings</span>
            </button>
            <button @click="$emit('minimize')" class="menu-item close">
                <span class="icon">✕</span>
                <span class="text">Minimize</span>
            </button>
        </div>
    </transition>
</template>

<script>
export default {
    name: 'PetMenu',
    props: {
        visible: {
            type: Boolean,
            default: false
        },
        position: {
            type: Object,
            default: () => ({ x: 0, y: 0 })
        },
        chatVisible: {
            type: Boolean,
            default: false
        }
    },
    emits: ['screenshot', 'toggle-chat', 'settings', 'minimize']
};
</script>

<style scoped>
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
    background: #000000;
    color: white;
    transform: translateX(4px);
}

.menu-item.close:hover {
    background: #f5576c;
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
</style>
