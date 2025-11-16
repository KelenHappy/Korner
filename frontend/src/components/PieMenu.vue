<template>
    <transition name="pie-fade">
        <div v-if="visible" class="pie-menu-container">
            <!-- Background overlay -->
            <div class="pie-menu-background" @click="$emit('hide')"></div>

            <!-- Pie Menu centered at window center -->
            <div class="pie-menu-content">
                <!-- Center circle -->
                <div class="pie-center" @click="$emit('hide')">
                    <span class="center-icon">⊕</span>
                </div>

                <!-- Menu items -->
                <div class="pie-items">
                    <!-- Screenshot (top) -->
                    <div
                        class="pie-item screenshot-item"
                        :style="getPieItemStyle(0)"
                        @click="$emit('screenshot')"
                        @mouseenter="activeItem = 0"
                        @mouseleave="activeItem = null"
                        :class="{ active: activeItem === 0 }"
                    >
                        <div class="pie-item-icon">📸</div>
                        <div class="pie-item-label">Screenshot</div>
                    </div>

                    <!-- Ask (right) -->
                    <div
                        class="pie-item question-item"
                        :style="getPieItemStyle(1)"
                        @click="$emit('ask-question')"
                        @mouseenter="activeItem = 1"
                        @mouseleave="activeItem = null"
                        :class="{ active: activeItem === 1 }"
                    >
                        <div class="pie-item-icon">💬</div>
                        <div class="pie-item-label">Ask</div>
                    </div>

                    <!-- Settings (bottom) -->
                    <div
                        class="pie-item settings-item"
                        :style="getPieItemStyle(2)"
                        @click="$emit('settings')"
                        @mouseenter="activeItem = 2"
                        @mouseleave="activeItem = null"
                        :class="{ active: activeItem === 2 }"
                    >
                        <div class="pie-item-icon">⚙️</div>
                        <div class="pie-item-label">Settings</div>
                    </div>

                    <!-- Hide (left) -->
                    <div
                        class="pie-item hide-item"
                        :style="getPieItemStyle(3)"
                        @click="$emit('hide')"
                        @mouseenter="activeItem = 3"
                        @mouseleave="activeItem = null"
                        :class="{ active: activeItem === 3 }"
                    >
                        <div class="pie-item-icon">✕</div>
                        <div class="pie-item-label">Hide</div>
                    </div>
                </div>
            </div>
        </div>
    </transition>
</template>

<script>
import { ref } from "vue";

export default {
    name: "PieMenu",
    props: {
        visible: {
            type: Boolean,
            default: false,
        },
        centerX: {
            type: Number,
            default: 250,
        },
        centerY: {
            type: Number,
            default: 250,
        },
    },
    emits: ["screenshot", "ask-question", "settings", "hide"],
    setup() {
        const activeItem = ref(null);
        const radius = 50; // 縮小一半

        const getPieItemStyle = (index) => {
            // 4 items in cardinal directions
            const angle = (index * 360) / 4 - 90; // Start from top
            const radians = (angle * Math.PI) / 180;
            const x = Math.cos(radians) * radius;
            const y = Math.sin(radians) * radius;

            return {
                transform: `translate(${x}px, ${y}px)`,
            };
        };

        return {
            activeItem,
            getPieItemStyle,
        };
    },
};
</script>

<style scoped>
.pie-menu-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 10000;
    pointer-events: auto;
}

.pie-menu-background {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.4);
    backdrop-filter: blur(3px);
    cursor: pointer;
}

.pie-menu-content {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    pointer-events: none;
}

.pie-center {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 30px;
    height: 30px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.5);
    z-index: 10;
    pointer-events: auto;
    cursor: pointer;
    transition: all 0.2s ease;
}

.pie-center:hover {
    transform: translate(-50%, -50%) scale(1.1);
}

.center-icon {
    font-size: 14px;
    color: white;
    font-weight: bold;
}

.pie-items {
    position: absolute;
    left: 50%;
    top: 50%;
    pointer-events: none;
}

.pie-item {
    position: absolute;
    left: -20px;
    top: -20px;
    width: 40px;
    height: 40px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    cursor: pointer;
    pointer-events: auto;
    transition: all 0.2s ease;
    user-select: none;
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.2);
}

.pie-item:hover,
.pie-item.active {
    transform: translate(var(--hover-x, 0), var(--hover-y, 0)) scale(1.2);
    box-shadow: 0 10px 24px rgba(0, 0, 0, 0.3);
}

.pie-item-icon {
    font-size: 16px;
    margin-bottom: 1px;
    transition: transform 0.2s ease;
}

.pie-item:hover .pie-item-icon,
.pie-item.active .pie-item-icon {
    transform: scale(1.1);
}

.pie-item-label {
    font-size: 6px;
    font-weight: 600;
    text-align: center;
    white-space: nowrap;
}

/* Item colors */
.screenshot-item {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    --hover-x: 0px;
    --hover-y: -10px;
}

.screenshot-item .pie-item-label {
    color: white;
}

.question-item {
    background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    --hover-x: 10px;
    --hover-y: 0px;
}

.question-item .pie-item-label {
    color: white;
}

.settings-item {
    background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    --hover-x: 0px;
    --hover-y: 10px;
}

.settings-item .pie-item-label {
    color: white;
}

.hide-item {
    background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
    --hover-x: -10px;
    --hover-y: 0px;
}

.hide-item .pie-item-label {
    color: #333;
}

/* Animations */
.pie-fade-enter-active {
    animation: pieMenuIn 0.3s ease-out;
}

.pie-fade-leave-active {
    animation: pieMenuOut 0.2s ease-in;
}

@keyframes pieMenuIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}

@keyframes pieMenuOut {
    from {
        opacity: 1;
    }
    to {
        opacity: 0;
    }
}

.pie-item {
    animation: itemSlideIn 0.3s ease-out backwards;
}

.pie-item:nth-child(1) {
    animation-delay: 0.05s;
}

.pie-item:nth-child(2) {
    animation-delay: 0.1s;
}

.pie-item:nth-child(3) {
    animation-delay: 0.15s;
}

.pie-item:nth-child(4) {
    animation-delay: 0.2s;
}

@keyframes itemSlideIn {
    from {
        opacity: 0;
        transform: translate(0, 0) scale(0.5);
    }
    to {
        opacity: 1;
    }
}
</style>
