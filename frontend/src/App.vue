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

            // Load settings from localStorage
            try {
                const savedSettings = localStorage.getItem("korner-settings");
                if (savedSettings) {
                    settings.value = JSON.parse(savedSettings);
                }
            } catch {}

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
            } else {
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
                showPieMenu.value = true;

                // 打開後立即調整窗口
                try {
                    const newX = iconScreenPos.value.x - 150;
                    const newY = iconScreenPos.value.y - 150;
                    WindowSetSize(300, 300);
                    WindowSetPosition(newX, newY);
                } catch (error) {
                    console.log("Failed to resize window:", error);
                }
            }
        };

        const hidePieMenu = async () => {
            showPieMenu.value = false;

            // 等待 CSS 淡出動畫完成（0.2s）後再縮小窗口
            await new Promise((resolve) => setTimeout(resolve, 200));

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
            // 記錄當前視窗狀態（在隱藏菜單之前）
            try {
                const pos = await WindowGetPosition();
                const size = await WindowGetSize();
                beforeScreenshotState.value = { pos, size };
            } catch {}

            // 先開全螢幕
            try {
                WindowFullscreen();
            } catch {}

            // 顯示截圖覆蓋層
            showScreenshotOverlay.value = true;

            // 隱藏菜單狀態（全螢幕後菜單會被覆蓋，這裡只是更新狀態）
            showPieMenu.value = false;
        };

        const handleAskQuestion = async () => {
            await hidePieMenu();
            currentQuery.value = {
                screenshot: null,
                timestamp: new Date(),
            };

            // 調整視窗大小並置中
            try {
                WindowSetSize(1200, 800);
                // 等待大小調整完成
                await new Promise((resolve) => setTimeout(resolve, 50));
                // 置中視窗
                WindowCenter();
            } catch (error) {
                console.log("Failed to resize/center window:", error);
            }

            showQueryWindow.value = true;
        };

        const handleSettings = async () => {
            await hidePieMenu();

            // 調整視窗大小並置中
            try {
                WindowSetSize(720, 600);
                // 等待大小調整完成
                await new Promise((resolve) => setTimeout(resolve, 50));
                // 置中視窗
                WindowCenter();
            } catch (error) {
                console.log("Failed to resize/center window:", error);
            }

            showSettingsWindow.value = true;
        };

        const closeSettings = async () => {
            showSettingsWindow.value = false;
            // 檢查是否需要縮小窗口
            await checkAndShrinkWindow();
        };

        const saveSettings = async (newSettings) => {
            settings.value = { ...newSettings };
            // Save to localStorage
            try {
                localStorage.setItem(
                    "korner-settings",
                    JSON.stringify(settings.value),
                );
            } catch (error) {
                console.error("Failed to save settings:", error);
            }
            await closeSettings();
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
                // 檢查是否需要縮小窗口
                await checkAndShrinkWindow();
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

            // 等待取消全螢幕完成
            await new Promise((resolve) => setTimeout(resolve, 100));

            // 調整視窗大小並置中
            try {
                WindowSetSize(1200, 800);
                // 等待大小調整完成
                await new Promise((resolve) => setTimeout(resolve, 50));
                // 置中視窗
                WindowCenter();
            } catch (error) {
                console.log("Failed to resize/center window:", error);
            }

            beforeScreenshotState.value = null;

            setTimeout(() => {
                showQueryWindow.value = true;
            }, 100);
        };

        const cancelQueryWindow = async () => {
            showQueryWindow.value = false;
            currentQuery.value = null;
            // 檢查是否需要縮小窗口
            await checkAndShrinkWindow();
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

            // 調整視窗大小並置中
            try {
                WindowSetSize(450, 500);
                // 等待大小調整完成
                await new Promise((resolve) => setTimeout(resolve, 50));
                // 置中視窗
                WindowCenter();
            } catch (error) {
                console.log("Failed to resize/center window:", error);
            }

            // 顯示 ResponseWindow
            showResponseWindow.value = true;
            currentQuery.value = null;
        };

        const closeResponseWindow = async () => {
            showResponseWindow.value = false;
            // 檢查是否需要縮小窗口
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
