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
                    >Click and drag to select an area</span
                >
                <kbd
                    class="ml-4 px-2 py-1 bg-slate-100 text-slate-600 text-xs rounded border border-slate-300"
                    >ESC</kbd
                >
                <span class="text-slate-500 text-sm">to cancel</span>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, reactive, onMounted, onUnmounted } from "vue";

export default {
    name: "ScreenshotOverlay",
    emits: ["screenshot-captured", "cancel"],
    setup(props, { emit }) {
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

        onMounted(() => {
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
            updateSelectionRect();
        };

        const updateSelectionRect = () => {
            const x = Math.min(startPoint.x, endPoint.x);
            const y = Math.min(startPoint.y, endPoint.y);
            const width = Math.abs(endPoint.x - startPoint.x);
            const height = Math.abs(endPoint.y - startPoint.y);

            selectionRect.x = x;
            selectionRect.y = y;
            selectionRect.width = width;
            selectionRect.height = height;
        };

        const endSelection = async () => {
            if (!isSelecting.value) return;
            isSelecting.value = false;

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
                        const dataUrl =
                            await window.go.main.App.CaptureScreenshot(
                                selectionRect.x,
                                selectionRect.y,
                                selectionRect.width,
                                selectionRect.height,
                            );
                        emit("screenshot-captured", dataUrl);
                        return;
                    } catch (backendError) {
                        console.warn(
                            "Backend screenshot failed, using placeholder:",
                            backendError,
                        );
                        // Fall through to placeholder
                    }
                }

                // Fallback: Create a simulated screenshot (canvas) for dev mode
                const canvas = document.createElement("canvas");
                canvas.width = selectionRect.width;
                canvas.height = selectionRect.height;
                const ctx = canvas.getContext("2d");

                // Create a gradient as placeholder
                const gradient = ctx.createLinearGradient(
                    0,
                    0,
                    canvas.width,
                    canvas.height,
                );
                gradient.addColorStop(0, "#667eea");
                gradient.addColorStop(1, "#764ba2");
                ctx.fillStyle = gradient;
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
                    `${canvas.width} × ${canvas.height}`,
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
            emit("cancel");
        };

        return {
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
