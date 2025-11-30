<template>
    <div
        class="fixed inset-0 z-50 cursor-crosshair"
        @mousedown="startSelection"
        @mousemove="updateSelection"
        @mouseup="endSelection"
        @keydown.esc="cancel"
        tabindex="0"
        ref="overlayRef"
    >
        <!-- Semi-transparent overlay -->
        <div class="absolute inset-0 bg-black bg-opacity-50"></div>

        <!-- Selection rectangle -->
        <div
            v-if="
                isSelecting ||
                (selectionRect.width > 0 && selectionRect.height > 0)
            "
            class="absolute border-2 border-blue-500 bg-blue-500 bg-opacity-10"
            :style="{
                left: selectionRect.x + 'px',
                top: selectionRect.y + 'px',
                width: selectionRect.width + 'px',
                height: selectionRect.height + 'px',
            }"
        >
            <!-- Corner indicators -->
            <div
                class="absolute -top-1 -left-1 w-2 h-2 bg-blue-500 rounded-full"
            ></div>
            <div
                class="absolute -top-1 -right-1 w-2 h-2 bg-blue-500 rounded-full"
            ></div>
            <div
                class="absolute -bottom-1 -left-1 w-2 h-2 bg-blue-500 rounded-full"
            ></div>
            <div
                class="absolute -bottom-1 -right-1 w-2 h-2 bg-blue-500 rounded-full"
            ></div>

            <!-- Dimension display -->
            <div
                class="absolute -top-8 left-0 bg-blue-500 text-white text-xs px-2 py-1 rounded whitespace-nowrap"
            >
                {{ Math.abs(selectionRect.width) }} ×
                {{ Math.abs(selectionRect.height) }}
            </div>
        </div>

        <!-- Instructions -->
        <div class="absolute top-8 left-1/2 transform -translate-x-1/2 z-10">
            <div
                class="bg-white rounded-lg shadow-lg px-6 py-3 flex items-center space-x-3"
            >
                <svg
                    class="w-5 h-5 text-blue-500"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"
                    />
                </svg>
                <span class="text-slate-700 font-medium"
                    >{{ t("screenshot.selectArea") }}</span
                >
                <kbd
                    class="ml-4 px-2 py-1 bg-slate-100 text-slate-600 text-xs rounded border border-slate-300"
                    >{{ t("common.esc") }}</kbd
                >
                <span class="text-slate-500 text-sm">{{ t("screenshot.cancel") }}</span>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, reactive, onMounted, onUnmounted } from "vue";
import { useI18n } from "vue-i18n";

export default {
    name: "ScreenshotOverlay",
    emits: ["screenshot-captured", "cancel"],
    setup(props, { emit }) {
        const { t } = useI18n();
        const overlayRef = ref(null);
        const isSelecting = ref(false);
        const startPoint = reactive({ x: 0, y: 0 });
        const endPoint = reactive({ x: 0, y: 0 });
        const selectionRect = reactive({
            x: 0,
            y: 0,
            width: 0,
            height: 0,
        });
        const screenRect = reactive({
            x: 0,
            y: 0,
            width: 0,
            height: 0,
        });

        let rafId = null;
        let pendingUpdate = false;

        onMounted(async () => {
            // Focus the overlay for keyboard events
            overlayRef.value?.focus();
        });

        const startSelection = (event) => {
            isSelecting.value = true;
            startPoint.x = event.clientX;
            startPoint.y = event.clientY;
            endPoint.x = event.clientX;
            endPoint.y = event.clientY;
            updateSelectionRect();
        };

        const updateSelection = (event) => {
            if (!isSelecting.value) return;
            endPoint.x = event.clientX;
            endPoint.y = event.clientY;

            // Use requestAnimationFrame to optimize updates
            if (!pendingUpdate) {
                pendingUpdate = true;
                rafId = requestAnimationFrame(() => {
                    updateSelectionRect();
                    pendingUpdate = false;
                });
            }
        };

        const updateSelectionRect = () => {
            const x = Math.min(startPoint.x, endPoint.x);
            const y = Math.min(startPoint.y, endPoint.y);
            const width = Math.abs(endPoint.x - startPoint.x);
            const height = Math.abs(endPoint.y - startPoint.y);

            // Viewport coordinates for display (what user sees)
            selectionRect.x = x;
            selectionRect.y = y;
            selectionRect.width = width;
            selectionRect.height = height;

            // Store viewport coordinates - backend will handle DPI conversion
            screenRect.x = x;
            screenRect.y = y;
            screenRect.width = width;
            screenRect.height = height;
        };

        const endSelection = async () => {
            if (!isSelecting.value) return;
            isSelecting.value = false;

            // Cancel any pending animation frame
            if (rafId) {
                cancelAnimationFrame(rafId);
                rafId = null;
                pendingUpdate = false;
            }

            // Apply final update if needed
            updateSelectionRect();

            // Minimum selection size
            if (selectionRect.width < 10 || selectionRect.height < 10) {
                cancel();
                return;
            }

            // Capture the screenshot
            await captureScreenshot();
        };

        const captureScreenshot = async () => {
            try {
                // Try to call Wails backend first
                if (window.go && window.go.main && window.go.main.App) {
                    try {
                        // Hide the overlay temporarily so it doesn't appear in screenshot
                        const overlayElement = overlayRef.value;
                        if (overlayElement) {
                            overlayElement.style.display = "none";
                        }

                        // Wait a moment for overlay to hide
                        await new Promise((resolve) => setTimeout(resolve, 50));

                        // Scale coordinates by devicePixelRatio to get physical pixels
                        const dpr = window.devicePixelRatio || 1;
                        const physicalX = Math.round(screenRect.x * dpr);
                        const physicalY = Math.round(screenRect.y * dpr);
                        const physicalWidth = Math.round(
                            screenRect.width * dpr,
                        );
                        const physicalHeight = Math.round(
                            screenRect.height * dpr,
                        );

                        console.log(
                            `[Korner][ScreenshotOverlay] Scaling coordinates by devicePixelRatio=${dpr}: viewport (${screenRect.x}, ${screenRect.y}, ${screenRect.width}, ${screenRect.height}) -> physical (${physicalX}, ${physicalY}, ${physicalWidth}, ${physicalHeight})`,
                        );

                        // Pass physical pixel coordinates to backend
                        const dataUrl =
                            await window.go.main.App.CaptureScreenshot(
                                physicalX,
                                physicalY,
                                physicalWidth,
                                physicalHeight,
                            );

                        // Restore overlay (will be closed by parent anyway)
                        if (overlayElement) {
                            overlayElement.style.display = "";
                        }

                        emit("screenshot-captured", dataUrl);
                        return;
                    } catch (backendError) {
                        console.warn(
                            "Backend screenshot failed, using placeholder:",
                            backendError,
                        );
                        // Restore overlay on error
                        const overlayElement = overlayRef.value;
                        if (overlayElement) {
                            overlayElement.style.display = "";
                        }
                        // Fall through to placeholder
                    }
                }

                // Fallback: Create a simulated screenshot (canvas) for dev mode
                const canvas = document.createElement("canvas");
                canvas.width = screenRect.width;
                canvas.height = screenRect.height;
                const ctx = canvas.getContext("2d");

                // Create a solid color as placeholder
                ctx.fillStyle = "#777777";
                ctx.fillRect(0, 0, canvas.width, canvas.height);

                // Add text
                ctx.fillStyle = "white";
                ctx.font = "16px sans-serif";
                ctx.textAlign = "center";
                ctx.fillText(
                    "Screenshot Placeholder",
                    canvas.width / 2,
                    canvas.height / 2,
                );
                ctx.font = "12px sans-serif";
                ctx.fillText(
                    `${screenRect.width} × ${screenRect.height}`,
                    canvas.width / 2,
                    canvas.height / 2 + 25,
                );

                const dataUrl = canvas.toDataURL("image/png");
                emit("screenshot-captured", dataUrl);
            } catch (error) {
                console.error("Screenshot capture failed:", error);
                cancel();
            }
        };

        const cancel = () => {
            // Cancel any pending animation frame
            if (rafId) {
                cancelAnimationFrame(rafId);
                rafId = null;
                pendingUpdate = false;
            }
            emit("cancel");
        };

        onUnmounted(() => {
            // Clean up animation frame on unmount
            if (rafId) {
                cancelAnimationFrame(rafId);
                rafId = null;
            }
        });

        return {
            t,
            overlayRef,
            isSelecting,
            selectionRect,
            startSelection,
            updateSelection,
            endSelection,
            cancel,
        };
    },
};
</script>

<style scoped>
.fixed {
    position: fixed;
}

.inset-0 {
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
}

.z-50 {
    z-index: 50;
}

.cursor-crosshair {
    cursor: crosshair;
}

.absolute {
    position: absolute;
}

.bg-black {
    background-color: rgb(0, 0, 0);
}

.bg-opacity-50 {
    background-color: rgba(0, 0, 0, 0.5);
}

.border-2 {
    border-width: 2px;
}

.border-blue-500 {
    border-color: rgb(59, 130, 246);
}

.bg-blue-500 {
    background-color: rgb(59, 130, 246);
}

.bg-opacity-10 {
    background-color: rgba(0, 0, 0, 0.1);
}

.w-2 {
    width: 0.5rem;
}

.h-2 {
    height: 0.5rem;
}

.rounded-full {
    border-radius: 9999px;
}

.-top-1 {
    top: -0.25rem;
}

.-left-1 {
    left: -0.25rem;
}

.-right-1 {
    right: -0.25rem;
}

.-bottom-1 {
    bottom: -0.25rem;
}

.-top-8 {
    top: -2rem;
}

.left-0 {
    left: 0;
}

.text-white {
    color: rgb(255, 255, 255);
}

.text-xs {
    font-size: 0.75rem;
    line-height: 1rem;
}

.px-2 {
    padding-left: 0.5rem;
    padding-right: 0.5rem;
}

.py-1 {
    padding-top: 0.25rem;
    padding-bottom: 0.25rem;
}

.rounded {
    border-radius: 0.25rem;
}

.whitespace-nowrap {
    white-space: nowrap;
}

.top-8 {
    top: 2rem;
}

.left-1-2 {
    left: 50%;
}

.transform {
    transform: translateX(-50%);
}

.bg-white {
    background-color: rgb(255, 255, 255);
}

.rounded-lg {
    border-radius: 0.5rem;
}

.shadow-lg {
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
}

.px-6 {
    padding-left: 1.5rem;
    padding-right: 1.5rem;
}

.py-3 {
    padding-top: 0.75rem;
    padding-bottom: 0.75rem;
}

.flex {
    display: flex;
}

.items-center {
    align-items: center;
}

.space-x-3 > * + * {
    margin-left: 0.75rem;
}

.w-5 {
    width: 1.25rem;
}

.h-5 {
    height: 1.25rem;
}

.text-blue-500 {
    color: rgb(59, 130, 246);
}

.fill-none {
    fill: none;
}

.stroke-current {
    stroke: currentColor;
}

.stroke-linecap-round {
    stroke-linecap: round;
}

.stroke-linejoin-round {
    stroke-linejoin: round;
}

.stroke-width-2 {
    stroke-width: 2;
}

.text-slate-700 {
    color: rgb(51, 65, 85);
}

.font-medium {
    font-weight: 500;
}

.ml-4 {
    margin-left: 1rem;
}

.bg-slate-100 {
    background-color: rgb(241, 245, 249);
}

.text-slate-600 {
    color: rgb(71, 85, 105);
}

.border {
    border-width: 1px;
}

.border-slate-300 {
    border-color: rgb(203, 213, 225);
}

.text-slate-500 {
    color: rgb(100, 116, 139);
}

.text-sm {
    font-size: 0.875rem;
    line-height: 1.25rem;
}

.bg-yellow-100 {
    background-color: rgb(254, 249, 195);
}

.text-black {
    color: rgb(0, 0, 0);
}

.p-2 {
    padding: 0.5rem;
}

.font-mono {
    font-family: ui-monospace, monospace;
}

.z-20 {
    z-index: 20;
}

.top-2 {
    top: 0.5rem;
}

.left-2 {
    left: 0.5rem;
}
</style>
