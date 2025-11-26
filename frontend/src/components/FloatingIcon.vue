<template>
    <div class="floating-icon" @mousedown="onMouseDown" @click="onClick">
        <div class="icon-container" :class="{ pulse: showPulse }">
            <div class="icon-emoji">{{ icon }}</div>
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
    props: {
        icon: {
            type: String,
            default: "🌸",
        },
    },
    emits: ["show-menu"],
    setup(props, { emit }) {
        const isDragging = ref(false);
        const dragStart = reactive({ x: 0, y: 0 });
        const windowStart = reactive({ x: 0, y: 0 });
        const showTooltip = ref(false);
        const showPulse = ref(true);
        const hasMoved = ref(false);

        let tooltipTimer = null;
        let rafId = null;
        let pendingPosition = null;

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

            // Store pending position
            pendingPosition = {
                x: windowStart.x + deltaX,
                y: windowStart.y + deltaY,
            };

            // Schedule update using requestAnimationFrame
            if (!rafId) {
                rafId = requestAnimationFrame(updateWindowPosition);
            }
        };

        const updateWindowPosition = () => {
            rafId = null;
            if (pendingPosition && isDragging.value) {
                try {
                    WindowSetPosition(pendingPosition.x, pendingPosition.y);
                } catch (error) {
                    console.log("Failed to set window position:", error);
                }
            }
        };

        const onMouseUp = () => {
            isDragging.value = false;
            document.removeEventListener("mousemove", onMouseMove);
            document.removeEventListener("mouseup", onMouseUp);

            // Cancel any pending animation frame
            if (rafId) {
                cancelAnimationFrame(rafId);
                rafId = null;
            }

            // Apply final position if there's a pending one
            if (pendingPosition) {
                try {
                    WindowSetPosition(pendingPosition.x, pendingPosition.y);
                } catch (error) {
                    console.log("Failed to set window position:", error);
                }
                pendingPosition = null;
            }
        };

        const onClick = async (e) => {
            // Only trigger menu if it wasn't a drag
            if (!hasMoved.value) {
                // Menu appears at center of window (in window coordinates)
                const centerX = window.innerWidth / 2;
                const centerY = window.innerHeight / 2;
                emit("show-menu", centerX, centerY);
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
            if (rafId) {
                cancelAnimationFrame(rafId);
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
    background: #fff;
    border: 3px solid #4a4a4a;
    border-radius: 50%;
    box-shadow: 2px 4px 0px rgba(0,0,0,0.3);
    transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1), box-shadow 0.2s;
}

.icon-container:hover {
    transform: scale(1.2) rotate(5deg);
    box-shadow: 4px 6px 0px rgba(0,0,0,0.4);
}

.icon-container:active {
    transform: scale(0.95);
}

.icon-emoji {
    font-size: 18px;
    position: relative;
    z-index: 2;
    transition: transform 0.2s ease;
}

.icon-glow {
    display: none;
}

/* Pulse animation */
.icon-container.pulse {
    animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
    0%,
    100% {
        transform: scale(1);
        box-shadow: 2px 4px 0px rgba(0,0,0,0.3);
    }
    50% {
        transform: scale(1.05);
        box-shadow: 4px 6px 0px rgba(0,0,0,0.4);
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
