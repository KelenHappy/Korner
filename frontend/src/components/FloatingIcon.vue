<template>
    <div class="floating-icon" @mousedown="onMouseDown" @click="onClick">
        <div class="icon-container" :class="{ pulse: showPulse }">
            <div class="icon-emoji">🌸</div>
            <div class="icon-glow"></div>
        </div>
        <div class="icon-tooltip" v-if="showTooltip">Korner AI</div>
    </div>
</template>

<script>
import { ref, reactive, onMounted, onUnmounted } from "vue";
import {
    WindowSetPosition,
    WindowGetPosition,
} from "../../wailsjs/runtime/runtime";

export default {
    name: "FloatingIcon",
    emits: ["show-menu"],
    setup(props, { emit }) {
        const isDragging = ref(false);
        const dragStart = reactive({ x: 0, y: 0 });
        const windowStart = reactive({ x: 0, y: 0 });
        const showTooltip = ref(false);
        const showPulse = ref(true);
        const hasMoved = ref(false);

        let tooltipTimer = null;

        const onMouseDown = async (e) => {
            if (e.button !== 0) return; // Only left click

            isDragging.value = true;
            hasMoved.value = false;
            dragStart.x = e.screenX; // Use screenX/Y for window positioning
            dragStart.y = e.screenY;

            // Get current window position
            try {
                const winPos = await WindowGetPosition();
                windowStart.x = winPos.x;
                windowStart.y = winPos.y;
            } catch (error) {
                console.log("Failed to get window position:", error);
            }

            document.addEventListener("mousemove", onMouseMove);
            document.addEventListener("mouseup", onMouseUp);

            e.preventDefault();
            e.stopPropagation();
        };

        const onMouseMove = (e) => {
            if (!isDragging.value) return;

            const deltaX = e.screenX - dragStart.x;
            const deltaY = e.screenY - dragStart.y;

            // If moved more than 5px, consider it a drag
            if (Math.abs(deltaX) > 5 || Math.abs(deltaY) > 5) {
                hasMoved.value = true;
            }

            // Move the entire window
            try {
                WindowSetPosition(
                    windowStart.x + deltaX,
                    windowStart.y + deltaY,
                );
            } catch (error) {
                console.log("Failed to set window position:", error);
            }
        };

        const onMouseUp = () => {
            isDragging.value = false;
            document.removeEventListener("mousemove", onMouseMove);
            document.removeEventListener("mouseup", onMouseUp);
        };

        const onClick = (e) => {
            // Only trigger menu if it wasn't a drag
            if (!hasMoved.value) {
                // Menu appears at center of window
                emit(
                    "show-menu",
                    window.innerWidth / 2,
                    window.innerHeight / 2,
                );
            }
            hasMoved.value = false;
        };

        const onMouseEnter = () => {
            showPulse.value = false;
            tooltipTimer = setTimeout(() => {
                showTooltip.value = true;
            }, 500);
        };

        const onMouseLeave = () => {
            showTooltip.value = false;
            if (tooltipTimer) {
                clearTimeout(tooltipTimer);
            }
        };

        onMounted(() => {
            // Stop pulse after a few seconds
            setTimeout(() => {
                showPulse.value = false;
            }, 3000);
        });

        onUnmounted(() => {
            document.removeEventListener("mousemove", onMouseMove);
            document.removeEventListener("mouseup", onMouseUp);
            if (tooltipTimer) {
                clearTimeout(tooltipTimer);
            }
        });

        return {
            showTooltip,
            showPulse,
            onMouseDown,
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
}

.floating-icon:active {
    cursor: grabbing;
}

.icon-container {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 50%;
    box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
    transition: all 0.3s ease;
}

.icon-container:hover {
    transform: scale(1.1);
    box-shadow: 0 12px 32px rgba(102, 126, 234, 0.6);
}

.icon-container:active {
    transform: scale(0.95);
}

.icon-emoji {
    font-size: 18px;
    position: relative;
    z-index: 2;
    transition: transform 0.3s ease;
}

.icon-container:hover .icon-emoji {
    transform: rotate(10deg);
}

.icon-glow {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(255, 255, 255, 0.3), transparent);
    opacity: 0;
    transition: opacity 0.3s ease;
}

.icon-container:hover .icon-glow {
    opacity: 1;
}

/* Pulse animation */
.icon-container.pulse {
    animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
    0%,
    100% {
        transform: scale(1);
        box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
    }
    50% {
        transform: scale(1.05);
        box-shadow: 0 12px 32px rgba(102, 126, 234, 0.6);
    }
}

/* Tooltip */
.icon-tooltip {
    position: absolute;
    top: -35px;
    left: 50%;
    transform: translateX(-50%);
    padding: 6px 12px;
    background: rgba(0, 0, 0, 0.8);
    color: white;
    font-size: 12px;
    font-weight: 500;
    border-radius: 6px;
    white-space: nowrap;
    pointer-events: none;
    animation: tooltipFadeIn 0.2s ease;
}

.icon-tooltip::after {
    content: "";
    position: absolute;
    top: 100%;
    left: 50%;
    transform: translateX(-50%);
    border: 5px solid transparent;
    border-top-color: rgba(0, 0, 0, 0.8);
}

@keyframes tooltipFadeIn {
    from {
        opacity: 0;
        transform: translateX(-50%) translateY(-5px);
    }
    to {
        opacity: 1;
        transform: translateX(-50%) translateY(0);
    }
}
</style>
