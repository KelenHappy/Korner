<template>
    <div
        class="screenshot-overlay"
        @mousedown="startSelection"
        @mousemove="updateSelection"
        @mouseup="endSelection"
        @keydown.esc="cancel"
        tabindex="0"
        ref="overlayRef"
    >
        <!-- Semi-transparent overlay -->
        <div class="overlay-backdrop"></div>

        <!-- Selection Box -->
        <SelectionBox
            :visible="isSelecting || (selectionRect.width > 0 && selectionRect.height > 0)"
            :rect="selectionRect"
        />

        <!-- Instructions -->
        <InstructionBanner
            :instructionText="t('screenshot.selectArea')"
            :escKey="t('common.esc')"
            :cancelText="t('screenshot.cancel')"
        />
    </div>
</template>

<script>
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import SelectionBox from './screenshot/SelectionBox.vue';
import InstructionBanner from './screenshot/InstructionBanner.vue';

export default {
    name: 'ScreenshotOverlay',
    components: {
        SelectionBox,
        InstructionBanner
    },
    emits: ['screenshot-captured', 'cancel'],
    setup(props, { emit }) {
        const { t } = useI18n();
        const overlayRef = ref(null);
        const isSelecting = ref(false);
        const startPoint = reactive({ x: 0, y: 0 });
        const endPoint = reactive({ x: 0, y: 0 });
        const selectionRect = reactive({ x: 0, y: 0, width: 0, height: 0 });
        const screenRect = reactive({ x: 0, y: 0, width: 0, height: 0 });

        let rafId = null;
        let pendingUpdate = false;

        onMounted(() => {
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

            selectionRect.x = x;
            selectionRect.y = y;
            selectionRect.width = width;
            selectionRect.height = height;

            screenRect.x = x;
            screenRect.y = y;
            screenRect.width = width;
            screenRect.height = height;
        };

        const endSelection = async () => {
            if (!isSelecting.value) return;
            isSelecting.value = false;

            if (rafId) {
                cancelAnimationFrame(rafId);
                rafId = null;
                pendingUpdate = false;
            }

            updateSelectionRect();

            if (selectionRect.width < 10 || selectionRect.height < 10) {
                cancel();
                return;
            }

            await captureScreenshot();
        };

        const captureScreenshot = async () => {
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    try {
                        const overlayElement = overlayRef.value;
                        if (overlayElement) {
                            overlayElement.style.display = 'none';
                        }

                        await new Promise((resolve) => setTimeout(resolve, 50));

                        const dpr = window.devicePixelRatio || 1;
                        const physicalX = Math.round(screenRect.x * dpr);
                        const physicalY = Math.round(screenRect.y * dpr);
                        const physicalWidth = Math.round(screenRect.width * dpr);
                        const physicalHeight = Math.round(screenRect.height * dpr);

                        console.log(
                            `[Korner][ScreenshotOverlay] Scaling coordinates by devicePixelRatio=${dpr}: viewport (${screenRect.x}, ${screenRect.y}, ${screenRect.width}, ${screenRect.height}) -> physical (${physicalX}, ${physicalY}, ${physicalWidth}, ${physicalHeight})`
                        );

                        const dataUrl = await window.go.main.App.CaptureScreenshot(
                            physicalX,
                            physicalY,
                            physicalWidth,
                            physicalHeight
                        );

                        if (overlayElement) {
                            overlayElement.style.display = '';
                        }

                        emit('screenshot-captured', dataUrl);
                        return;
                    } catch (backendError) {
                        console.warn('Backend screenshot failed, using placeholder:', backendError);
                        const overlayElement = overlayRef.value;
                        if (overlayElement) {
                            overlayElement.style.display = '';
                        }
                    }
                }

                // Fallback: Create placeholder
                const canvas = document.createElement('canvas');
                canvas.width = screenRect.width;
                canvas.height = screenRect.height;
                const ctx = canvas.getContext('2d');

                ctx.fillStyle = '#777777';
                ctx.fillRect(0, 0, canvas.width, canvas.height);

                ctx.fillStyle = 'white';
                ctx.font = '16px sans-serif';
                ctx.textAlign = 'center';
                ctx.fillText('Screenshot Placeholder', canvas.width / 2, canvas.height / 2);
                ctx.font = '12px sans-serif';
                ctx.fillText(
                    `${screenRect.width} × ${screenRect.height}`,
                    canvas.width / 2,
                    canvas.height / 2 + 25
                );

                const dataUrl = canvas.toDataURL('image/png');
                emit('screenshot-captured', dataUrl);
            } catch (error) {
                console.error('Screenshot capture failed:', error);
                cancel();
            }
        };

        const cancel = () => {
            if (rafId) {
                cancelAnimationFrame(rafId);
                rafId = null;
                pendingUpdate = false;
            }
            emit('cancel');
        };

        onUnmounted(() => {
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
            cancel
        };
    }
};
</script>

<style scoped>
.screenshot-overlay {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    z-index: 50;
    cursor: crosshair;
}

.overlay-backdrop {
    position: absolute;
    top: 0;
    right: 0;
    bottom: 0;
    left: 0;
    background-color: rgba(0, 0, 0, 0.5);
}
</style>
