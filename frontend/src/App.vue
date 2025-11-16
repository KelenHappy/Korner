<template>
    <div id="app" class="app-container">
        <!-- Floating Icon (hidden during screenshot) -->
        <FloatingIcon
            v-if="!showScreenshotOverlay"
            @show-menu="showPieMenuAt"
        />

        <!-- Pie Menu (shown when icon is clicked) -->
        <PieMenu
            :visible="showPieMenu"
            :center-x="menuCenterX"
            :center-y="menuCenterY"
            @screenshot="handleScreenshot"
            @ask-question="handleAskQuestion"
            @settings="handleSettings"
            @hide="hidePieMenu"
        />

        <!-- Screenshot Overlay -->
        <ScreenshotOverlay
            v-if="showScreenshotOverlay"
            @screenshot-captured="handleScreenshotCaptured"
            @cancel="cancelScreenshot"
        />

        <!-- Query Window Modal -->
        <QueryWindow
            v-if="showQueryWindow && currentQuery"
            :screenshot="currentQuery.screenshot"
            @submit="handleQuerySubmit"
            @cancel="cancelQueryWindow"
        />

        <!-- Response Window -->
        <ResponseWindow
            v-if="showResponseWindow"
            :response="latestResponse"
            :loading="isLoadingResponse"
            :screenshot="lastScreenshot"
            @close="closeResponseWindow"
        />
    </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from "vue";
import FloatingIcon from "./components/FloatingIcon.vue";
import PieMenu from "./components/PieMenu.vue";
import ScreenshotOverlay from "./components/ScreenshotOverlay.vue";
import QueryWindow from "./components/QueryWindow.vue";
import ResponseWindow from "./components/ResponseWindow.vue";

import {
    WindowSetAlwaysOnTop,
    WindowFullscreen,
    WindowUnfullscreen,
    WindowSetSize,
    WindowSetPosition,
    WindowGetPosition,
    WindowGetSize,
    EventsOn,
    EventsOff,
} from "../wailsjs/runtime/runtime";

export default {
    name: "App",
    components: {
        FloatingIcon,
        PieMenu,
        ScreenshotOverlay,
        QueryWindow,
        ResponseWindow,
    },
    setup() {
        const showPieMenu = ref(false);
        const menuCenterX = ref(250);
        const menuCenterY = ref(250);
        const showScreenshotOverlay = ref(false);
        const showQueryWindow = ref(false);
        const showResponseWindow = ref(false);
        const currentQuery = ref(null);
        const latestResponse = ref("");
        const isLoadingResponse = ref(false);
        const lastScreenshot = ref(null);
        const beforeScreenshotState = ref(null);

        // 統一的視窗大小調整函數，保持中心點不變
        const resizeWindowKeepCenter = async (newWidth, newHeight) => {
            try {
                const pos = await WindowGetPosition();
                const oldSize = await WindowGetSize();
                const newSize = { w: newWidth, h: newHeight };
                // 計算新位置，保持視窗中心不變
                const newX = pos.x + (oldSize.w - newSize.w) / 2;
                const newY = pos.y + (oldSize.h - newSize.h) / 2;
                WindowSetSize(newSize.w, newSize.h);
                WindowSetPosition(newX, newY);
            } catch (error) {
                console.log("Failed to resize window:", error);
            }
        };

        onMounted(async () => {
            try {
                WindowSetAlwaysOnTop(true);
            } catch {}
            // 初始設置小視窗 - 只顯示浮動圖標
            await resizeWindowKeepCenter(100, 100);

            // Listen for system tray trigger
            try {
                EventsOn("trigger-pie-menu", () => {
                    showPieMenuAt(250, 250);
                });
            } catch {}

            // ESC key to hide menu
            document.addEventListener("keydown", onKeyDown);
        });

        onUnmounted(() => {
            try {
                EventsOff("trigger-pie-menu");
            } catch {}
            document.removeEventListener("keydown", onKeyDown);
        });

        const onKeyDown = (e) => {
            if (e.key === "Escape") {
                if (showPieMenu.value) {
                    hidePieMenu();
                } else if (showScreenshotOverlay.value) {
                    cancelScreenshot();
                } else if (showQueryWindow.value) {
                    cancelQueryWindow();
                } else if (showResponseWindow.value) {
                    closeResponseWindow();
                }
            }
        };

        const showPieMenuAt = async (x, y) => {
            // Toggle pie menu - if already showing, hide it
            if (showPieMenu.value) {
                await hidePieMenu();
            } else {
                menuCenterX.value = x;
                menuCenterY.value = y;
                showPieMenu.value = true;
                // 放大視窗以顯示 PieMenu
                await resizeWindowKeepCenter(300, 300);
            }
        };

        const hidePieMenu = async () => {
            showPieMenu.value = false;
            // 縮回小視窗
            await resizeWindowKeepCenter(100, 100);
        };

        const handleScreenshot = async () => {
            await hidePieMenu();
            // 記錄當前視窗狀態
            try {
                const pos = await WindowGetPosition();
                const size = await WindowGetSize();
                beforeScreenshotState.value = { pos, size };
            } catch {}

            showScreenshotOverlay.value = true;
            try {
                // 截圖時需要全螢幕
                WindowFullscreen();
            } catch {}
        };

        const handleAskQuestion = async () => {
            await hidePieMenu();
            currentQuery.value = {
                screenshot: null,
                timestamp: new Date(),
            };
            showQueryWindow.value = true;
            // 調整視窗大小以顯示 QueryWindow
            await resizeWindowKeepCenter(1200, 800);
        };

        const handleSettings = () => {
            hidePieMenu();
            console.log("Settings clicked");
        };

        const cancelScreenshot = async () => {
            showScreenshotOverlay.value = false;
            try {
                WindowUnfullscreen();
            } catch {}

            // 恢復到截圖前的狀態
            if (beforeScreenshotState.value) {
                try {
                    WindowSetSize(
                        beforeScreenshotState.value.size.w,
                        beforeScreenshotState.value.size.h,
                    );
                    WindowSetPosition(
                        beforeScreenshotState.value.pos.x,
                        beforeScreenshotState.value.pos.y,
                    );
                } catch {}
                beforeScreenshotState.value = null;
            } else {
                // 如果沒有記錄，恢復小視窗
                await resizeWindowKeepCenter(100, 100);
            }
        };

        const handleScreenshotCaptured = async (screenshotData) => {
            currentQuery.value = {
                screenshot: screenshotData,
                timestamp: new Date(),
            };
            lastScreenshot.value = screenshotData;
            showScreenshotOverlay.value = false;

            // 截圖完成後取消全螢幕
            try {
                WindowUnfullscreen();
            } catch {}

            // 恢復視窗位置並調整到 QueryWindow 大小
            if (beforeScreenshotState.value) {
                try {
                    WindowSetPosition(
                        beforeScreenshotState.value.pos.x,
                        beforeScreenshotState.value.pos.y,
                    );
                } catch {}
                beforeScreenshotState.value = null;
            }

            await resizeWindowKeepCenter(1200, 800);

            setTimeout(() => {
                showQueryWindow.value = true;
            }, 200);
        };

        const cancelQueryWindow = async () => {
            showQueryWindow.value = false;
            currentQuery.value = null;
            // 恢復小視窗
            await resizeWindowKeepCenter(100, 100);
        };

        const handleQuerySubmit = async (queryText) => {
            if (!currentQuery.value) return;

            let screenshotB64 = currentQuery.value.screenshot || "";
            if (screenshotB64 && screenshotB64.startsWith("data:image")) {
                const comma = screenshotB64.indexOf(",");
                if (comma !== -1) {
                    screenshotB64 = screenshotB64.substring(comma + 1);
                }
            }

            showQueryWindow.value = false;
            isLoadingResponse.value = true;
            latestResponse.value = "";

            try {
                if (window.go && window.go.main && window.go.main.App) {
                    const response = await window.go.main.App.QueryLLM(
                        queryText,
                        screenshotB64,
                    );
                    latestResponse.value = response;
                    isLoadingResponse.value = false;
                } else {
                    setTimeout(() => {
                        latestResponse.value = `Response to: "${queryText}"\n\n(Dev mode)`;
                        isLoadingResponse.value = false;
                    }, 1200);
                }
            } catch (error) {
                latestResponse.value = `Error: ${error?.message || "Unknown"}`;
                isLoadingResponse.value = false;
            }

            // 顯示 ResponseWindow，調整視窗到合適大小
            showResponseWindow.value = true;
            currentQuery.value = null;
            await resizeWindowKeepCenter(450, 500);
        };

        const closeResponseWindow = async () => {
            showResponseWindow.value = false;
            // 恢復小視窗
            await resizeWindowKeepCenter(100, 100);
        };

        return {
            showPieMenu,
            menuCenterX,
            menuCenterY,
            showScreenshotOverlay,
            showQueryWindow,
            showResponseWindow,
            currentQuery,
            latestResponse,
            isLoadingResponse,
            lastScreenshot,
            showPieMenuAt,
            hidePieMenu,
            handleScreenshot,
            handleAskQuestion,
            handleSettings,
            cancelScreenshot,
            handleScreenshotCaptured,
            cancelQueryWindow,
            handleQuerySubmit,
            closeResponseWindow,
        };
    },
};
</script>

<style>
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

html,
body {
    width: 100%;
    height: 100%;
    background: transparent;
    overflow: hidden;
}

#app {
    width: 100%;
    height: 100%;
}

.app-container {
    width: 100%;
    height: 100%;
    background: transparent;
    pointer-events: none;
}

.app-container > * {
    pointer-events: auto;
}
</style>
