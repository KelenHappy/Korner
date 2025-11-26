<template>
    <div class="floating-icon" @dblclick="onClick">
        <div class="icon-container" :class="{ pulse: showPulse }">
            <div class="icon-emoji">{{ icon }}</div>
        </div>
    </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from "vue";

export default {
    name: "FloatingIcon",
    props: {
        icon: {
            type: String,
            default: "🌸",
        },
    },
    emits: ["show-menu"],
    setup(props, { emit }) {
        const showPulse = ref(true);

        const onClick = () => {
            emit("show-menu", 50, 50);
        };

        onMounted(() => {
            setTimeout(() => {
                showPulse.value = false;
            }, 3000);
        });

        return {
            showPulse,
            onClick,
        };
    },
};
</script>

<style scoped>
.floating-icon {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 40px;
    height: 40px;
    cursor: grab;
    user-select: none;
    z-index: 9999;
    --wails-draggable: drag;
}

.floating-icon:active {
    cursor: grabbing;
}

.icon-container {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #fff;
    border: 3px solid #4a4a4a;
    border-radius: 50%;
    box-shadow: 2px 4px 0px rgba(0, 0, 0, 0.3);
}

.icon-emoji {
    font-size: 18px;
}

.icon-container.pulse {
    animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
    0%, 100% { box-shadow: 2px 4px 0px rgba(0, 0, 0, 0.3); }
    50% { box-shadow: 4px 6px 0px rgba(0, 0, 0, 0.4); }
}
</style>
