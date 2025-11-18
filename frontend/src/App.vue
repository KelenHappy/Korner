<template>
    <div id="app" class="app-container">
        <!-- Floating Icon (hidden during screenshot) -->
        <FloatingIcon
            v-if="!showScreenshotOverlay"
            :icon="settings.floatingIcon"
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

        <!-- Settings Window -->
        <SettingsWindow
            v-if="showSettingsWindow"
            :current-settings="settings"
            @close="closeSettings"
            @save="saveSettings"
        />
    </div>
</template>

<script>
import { ref, onMounted, onUnmounted, nextTick } from "vue";
import FloatingIcon from "./components/FloatingIcon.vue";
import PieMenu from "./components/PieMenu.vue";
import ScreenshotOverlay from "./components/ScreenshotOverlay.vue";
import QueryWindow from "./components/QueryWindow.vue";
import ResponseWindow from "./components/ResponseWindow.vue";
import SettingsWindow from "./components/SettingsWindow.vue";

import {
    WindowSetAlwaysOnTop,
    WindowFullscreen,
    WindowUnfullscreen,
    WindowSetSize,
    WindowSetPosition,
    WindowGetPosition,
    WindowGetSize,
    WindowCenter,
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
        SettingsWindow,
    },
    setup() {
        const showPieMenu = ref(false);
        const menuCenterX = ref(250);
        const menuCenterY = ref(250);
        const showScreenshotOverlay = ref(false);
        const showQueryWindow = ref(false);
        const showResponseWindow = ref(false);
        const showSettingsWindow = ref(false);
        const currentQuery = ref(null);
        const latestResponse = ref("");
        const isLoadingResponse = ref(false);
        const lastScreenshot = ref(null);
        const beforeScreenshotState = ref(null);
        const beforePieMenuState = ref(null);
        const fixedCenterPos = ref(null);
        const iconScreenPos = ref(null);
        const settings = ref({
            apiProvider: "openai",
            apiKey: "",
            apiEndpoint: "",
            floatingIcon: "🌸",
        });

        // 設置窗口大小和位置，基於浮動圖標的屏幕位置
        const setWindowForIcon = async (windowSize, iconOffset) => {
            try {
                let iconPos = iconScreenPos.value;
                if (!iconPos) {
                    // 如果沒有記錄圖標位置，獲取當前窗口位置作為圖標位置
                    const pos = await WindowGetPosition();
                    iconPos = {
                        x: pos.x + 50,
                        y: pos.y + 50,
                    };
                }

                // 同步內容動畫：先設置位置，再設置大小
                WindowSetPosition(
                    iconPos.x - iconOffset,
                    iconPos.y - iconOffset,
                );
                await new Promise((resolve) => setTimeout(resolve, 100));
                WindowSetSize(windowSize, windowSize);
            } catch (error) {
                console.log("Failed to set window:", error);
            }
        };

        onMounted(async () => {
            try {
                WindowSetAlwaysOnTop(true);
            } catch {}
            // 初始設置小視窗 100x100
            try {
                WindowSetSize(100, 100);
            } catch (error) {
                console.log("Failed to set initial window size:", error);
            }

            // Load settings from backend
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    const savedSettings =
                        await window.go.main.App.GetSettings();
                    if (savedSettings) {
                        settings.value = savedSettings;
                    }
                }
            } catch (error) {
                console.error("Failed to load settings:", error);
            }

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
                } else if (showSettingsWindow.value) {
                    closeSettings();
                }
            }
        };

        const showPieMenuAt = async (x, y) => {
            // Toggle pie menu - if already showing, hide it
            if (showPieMenu.value) {
                await hidePieMenu();
                return;
            }

            // 記錄浮動圖標的屏幕位置
            try {
                const pos = await WindowGetPosition();
                iconScreenPos.value = {
                    x: pos.x + 50,
                    y: pos.y + 50,
                };
            } catch (error) {
                console.log("Failed to save icon position:", error);
            }

            menuCenterX.value = x;
            menuCenterY.value = y;

            // 先顯示菜單，再擴大窗口（避免閃爍）
            showPieMenu.value = true;

            // 延遲擴大窗口
            await new Promise((resolve) => setTimeout(resolve, 16));

            try {
                const newX = iconScreenPos.value.x - 150;
                const newY = iconScreenPos.value.y - 150;
                WindowSetSize(300, 300);
                WindowSetPosition(newX, newY);
            } catch (error) {
                console.log("Failed to resize window:", error);
            }
        };

        const hidePieMenu = async () => {
            showPieMenu.value = false;

            // 等待 CSS 淡出動畫完成（0.2s）
            await new Promise((resolve) => setTimeout(resolve, 220));

            // 縮小窗口回到圖標大小
            if (iconScreenPos.value) {
                try {
                    const newX = iconScreenPos.value.x - 50;
                    const newY = iconScreenPos.value.y - 50;
                    WindowSetSize(100, 100);
                    WindowSetPosition(newX, newY);
                } catch (error) {
                    console.log("Failed to shrink window:", error);
                }
            }
        };

        // 檢查是否需要縮小窗口的函數
        const checkAndShrinkWindow = async () => {
            // 如果沒有主要內容顯示，縮小窗口
            if (
                !showPieMenu.value &&
                !showScreenshotOverlay.value &&
                !showQueryWindow.value &&
                !showResponseWindow.value &&
                !showSettingsWindow.value
            ) {
                try {
                    WindowSetSize(100, 100);
                    if (iconScreenPos.value) {
                        WindowSetPosition(
                            iconScreenPos.value.x - 50,
                            iconScreenPos.value.y - 50,
                        );
                    }
                } catch (error) {
                    console.log("Failed to shrink window:", error);
                }
            }
        };

        const handleScreenshot = async () => {
            // 記錄當前視窗狀態
            try {
                const pos = await WindowGetPosition();
                const size = await WindowGetSize();
                beforeScreenshotState.value = { pos, size };
            } catch {}

            // 隱藏菜單
            showPieMenu.value = false;

            // 等待菜單隱藏動畫
            await new Promise((resolve) => setTimeout(resolve, 100));

            // 開啟全螢幕
            try {
                WindowFullscreen();
                await new Promise((resolve) => setTimeout(resolve, 100));
            } catch {}

            // 顯示截圖覆蓋層
            showScreenshotOverlay.value = true;
        };

        const handleAskQuestion = async () => {
            // 隱藏菜單
            await hidePieMenu();

            // 調整視窗大小並置中
            try {
                WindowSetSize(1200, 800);
                await new Promise((resolve) => setTimeout(resolve, 100));
                WindowCenter();
            } catch (error) {
                console.log("Failed to resize/center window:", error);
            }

            // 設置查詢狀態並顯示
            currentQuery.value = {
                screenshot: null,
                timestamp: new Date(),
            };
            showQueryWindow.value = true;
        };

        const handleSettings = async () => {
            await hidePieMenu();

            // 調整視窗大小並置中
            try {
                WindowSetSize(720, 600);
                await new Promise((resolve) => setTimeout(resolve, 100));
                WindowCenter();
            } catch (error) {
                console.log("Failed to resize/center window:", error);
            }

            showSettingsWindow.value = true;
        };

        const closeSettings = async () => {
            showSettingsWindow.value = false;
            await new Promise((resolve) => setTimeout(resolve, 100));
            await checkAndShrinkWindow();
        };

        const saveSettings = async (newSettings) => {
            settings.value = { ...newSettings };
            // Save to backend
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    await window.go.main.App.SaveSettings(newSettings);
                }
            } catch (error) {
                console.error("Failed to save settings:", error);
            }
            await closeSettings();
        };

        const cancelScreenshot = async () => {
            showScreenshotOverlay.value = false;

            try {
                WindowUnfullscreen();
                await new Promise((resolve) => setTimeout(resolve, 100));
            } catch {}

            // 恢復到截圖前的狀態
            if (beforeScreenshotState.value) {
                try {
                    WindowSetSize(
                        beforeScreenshotState.value.size.w,
                        beforeScreenshotState.value.size.h,
                    );
                    await new Promise((resolve) => setTimeout(resolve, 50));
                    WindowSetPosition(
                        beforeScreenshotState.value.pos.x,
                        beforeScreenshotState.value.pos.y,
                    );
                } catch {}
                beforeScreenshotState.value = null;
            } else {
                await checkAndShrinkWindow();
            }
        };

        const handleScreenshotCaptured = async (screenshotData) => {
            // 保存截圖數據
            currentQuery.value = {
                screenshot: screenshotData,
                timestamp: new Date(),
            };
            lastScreenshot.value = screenshotData;

            // 隱藏覆蓋層
            showScreenshotOverlay.value = false;

            // 退出全螢幕
            try {
                WindowUnfullscreen();
                await new Promise((resolve) => setTimeout(resolve, 150));
            } catch {}

            // 調整視窗大小並置中
            try {
                WindowSetSize(1200, 800);
                await new Promise((resolve) => setTimeout(resolve, 100));
                WindowCenter();
            } catch (error) {
                console.log("Failed to resize/center window:", error);
            }

            beforeScreenshotState.value = null;

            // 顯示查詢窗口
            await new Promise((resolve) => setTimeout(resolve, 50));
            showQueryWindow.value = true;
        };

        const cancelQueryWindow = async () => {
            showQueryWindow.value = false;
            currentQuery.value = null;
            await new Promise((resolve) => setTimeout(resolve, 100));
            await checkAndShrinkWindow();
        };

        const handleQuerySubmit = async (queryText, callback) => {
            if (!currentQuery.value) {
                console.error("[Korner] No current query");
                return;
            }

            let screenshotB64 = currentQuery.value.screenshot || "";
            if (screenshotB64 && screenshotB64.startsWith("data:image")) {
                const comma = screenshotB64.indexOf(",");
                if (comma !== -1) {
                    screenshotB64 = screenshotB64.substring(comma + 1);
                }
            }

            try {
                let response;
                if (window.go && window.go.main && window.go.main.App) {
                    console.log("[Korner] Sending query to backend...");
                    response = await window.go.main.App.QueryLLM(
                        queryText,
                        screenshotB64,
                    );
                    console.log("[Korner] Received response from backend");
                } else {
                    // Dev mode - simulate response
                    console.log("[Korner] Dev mode - simulating response");
                    await new Promise((resolve) => setTimeout(resolve, 1200));
                    response = `This is a simulated response to: "${queryText}"\n\nIn production, this would be the actual AI response from the backend.`;
                }

                // Call the callback with the response
                if (callback && typeof callback === "function") {
                    callback(response);
                } else {
                    console.error("[Korner] Invalid callback function");
                }
            } catch (error) {
                console.error("[Korner] Error in handleQuerySubmit:", error);
                let errorMsg = "Unknown error occurred";

                if (error && typeof error === "string") {
                    errorMsg = error;
                } else if (error && error.message) {
                    errorMsg = error.message;
                } else if (error) {
                    errorMsg = String(error);
                }

                const fullErrorMsg = `Error: ${errorMsg}`;
                console.error("[Korner] Full error message:", fullErrorMsg);

                if (callback && typeof callback === "function") {
                    callback(fullErrorMsg);
                }
            }
        };

        const closeResponseWindow = async () => {
            showResponseWindow.value = false;
            await new Promise((resolve) => setTimeout(resolve, 100));
            await checkAndShrinkWindow();
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
            showSettingsWindow,
            settings,
            closeSettings,
            saveSettings,
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
