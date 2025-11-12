<template>
    <div id="app" class="h-screen w-screen overflow-hidden bg-transparent">
        <!-- Window drag handle (invisible overlay in top area) -->
        <div class="window-drag-handle" @mousedown="startWindowDrag"></div>

        <!-- Fullscreen screenshot overlay -->
        <ScreenshotOverlay
            v-if="showScreenshotOverlay"
            @screenshot-captured="handleScreenshotCaptured"
            @cancel="cancelScreenshot"
        />

        <!-- Desktop Pet - always visible except during screenshot -->
        <DesktopPet
            v-if="!showScreenshotOverlay"
            :status="petStatus"
            :chatMessage="latestResponse"
            :chatLoading="isLoadingResponse"
            :lastScreenshot="lastScreenshot"
            @screenshot="triggerScreenshot"
            @settings="showSettings"
            @minimize="minimizeToTray"
        />

        <!-- Query Window Modal - shown after screenshot captured -->
        <QueryWindow
            v-if="showQueryWindow && currentQuery"
            :screenshot="currentQuery.screenshot"
            @submit="handleQuerySubmit"
            @cancel="cancelQueryWindow"
        />

        <!-- Response Window - shown after query submitted -->
        <ResponseWindow
            v-if="showResponseWindow"
            :response="latestResponse"
            :loading="isLoadingResponse"
            :screenshot="lastScreenshot"
            @close="closeResponseWindow"
            @pin="pinResponseWindow"
        />
    </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from "vue";
import ScreenshotOverlay from "./components/ScreenshotOverlay.vue";
import QueryWindow from "./components/QueryWindow.vue";
import DesktopPet from "./components/DesktopPet.vue";
import ResponseWindow from "./components/ResponseWindow.vue";

import {
    WindowSetAlwaysOnTop,
    WindowFullscreen,
    WindowUnfullscreen,
    EventsOn,
    EventsOff,
    WindowShow,
} from "../wailsjs/runtime/runtime";

export default {
    name: "App",
    components: {
        ScreenshotOverlay,
        QueryWindow,
        DesktopPet,
        ResponseWindow,
    },
    setup() {
        const showScreenshotOverlay = ref(false);
        const showQueryWindow = ref(false);
        const showResponseWindow = ref(false);
        const currentQuery = ref(null);
        const latestResponse = ref("");
        const isLoadingResponse = ref(false);
        const lastScreenshot = ref(null);
        const platform = ref("unknown");
        const petStatus = ref("");

        onMounted(async () => {
            // Detect platform
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    platform.value = await window.go.main.App.GetPlatform();
                } else {
                    const ua = navigator.userAgent.toLowerCase();
                    if (ua.includes("mac")) platform.value = "darwin";
                    else if (ua.includes("win")) platform.value = "windows";
                    else platform.value = "unknown";
                }
            } catch {
                const ua = navigator.userAgent.toLowerCase();
                if (ua.includes("mac")) platform.value = "darwin";
                else if (ua.includes("win")) platform.value = "windows";
                else platform.value = "unknown";
            }

            // Keep window always on top
            try {
                WindowSetAlwaysOnTop(true);
            } catch {
                // ignore if not running under Wails runtime (dev mode)
            }

            // Listen for system tray screenshot trigger
            try {
                EventsOn("trigger-screenshot", () => {
                    triggerScreenshot();
                });
            } catch {
                // ignore if not running under Wails runtime (dev mode)
            }
        });

        onUnmounted(() => {
            try {
                EventsOff("trigger-screenshot");
            } catch {
                // ignore if not running under Wails runtime (dev mode)
            }
        });

        const hotkey = computed(() =>
            platform.value === "darwin" ? "Cmd+Option+Q" : "Ctrl+Alt+Q",
        );

        const triggerScreenshot = async () => {
            console.log("[Korner] Triggering screenshot flow");
            showResponseWindow.value = false;
            latestResponse.value = "";
            isLoadingResponse.value = false;
            try {
                WindowShow();
                WindowSetAlwaysOnTop(true);
            } catch {}

            // Try to maximize window first, then fullscreen
            try {
                if (window.wails && window.wails.Window) {
                    await window.wails.Window.Maximise();
                }
            } catch (e) {
                console.log("[Korner] Maximise not available:", e);
            }

            try {
                WindowFullscreen();
            } catch {}

            showScreenshotOverlay.value = true;
        };

        const handleScreenshotCaptured = (screenshotData) => {
            console.log("[Korner] Screenshot captured", {
                hasData: !!screenshotData,
                length: screenshotData ? screenshotData.length : 0,
            });
            // Screenshot captured - store it and prepare for query window

            currentQuery.value = {
                screenshot: screenshotData,

                timestamp: new Date(),
            };

            lastScreenshot.value = screenshotData;

            console.log("[Korner] Stored screenshot for query window", {
                length: screenshotData ? screenshotData.length : 0,
            });

            showScreenshotOverlay.value = false;

            try {
                WindowUnfullscreen();
            } catch {}

            // Show query window with a small delay for smooth transition

            setTimeout(() => {
                console.log("[Korner] Opening QueryWindow modal");
                showQueryWindow.value = true;
            }, 200);
        };

        const cancelScreenshot = () => {
            showScreenshotOverlay.value = false;
            try {
                WindowUnfullscreen();
            } catch {}
        };

        const cancelQueryWindow = () => {
            showQueryWindow.value = false;
            currentQuery.value = null;
            latestResponse.value = "";
        };

        const handleQuerySubmit = async (queryText) => {
            if (!currentQuery.value) return;
            console.log("[Korner] Submitting query", {
                textLength: queryText.length,
                hasScreenshot: !!currentQuery.value?.screenshot,
            });

            // Prepare screenshot data
            let screenshotB64 = currentQuery.value.screenshot || "";
            if (screenshotB64.startsWith("data:image")) {
                const comma = screenshotB64.indexOf(",");
                if (comma !== -1) {
                    screenshotB64 = screenshotB64.substring(comma + 1);
                }
            }

            // Close query window and show loading in pet's chatbox
            showQueryWindow.value = false;
            isLoadingResponse.value = true;
            latestResponse.value = ""; // Will trigger chatbox in pet

            try {
                if (window.go && window.go.main && window.go.main.App) {
                    const response = await window.go.main.App.QueryLLM(
                        queryText,
                        screenshotB64,
                    );
                    latestResponse.value = response;
                    isLoadingResponse.value = false;
                } else {
                    // Dev fallback
                    setTimeout(() => {
                        latestResponse.value = `Response to: "${queryText}"\n\nThis connects to AMD GPT OSS 120B model.\n\n(Running in dev mode - backend not connected)`;
                        isLoadingResponse.value = false;
                    }, 1200);
                }
            } catch (error) {
                const msg = error?.message || "Failed to query LLM";
                latestResponse.value = `Error: ${msg}\n\nPlease check:\n1. AMD_LLM_ENDPOINT is set\n2. AMD_API_KEY is valid\n3. Network connection`;
                isLoadingResponse.value = false;
            }

            // Show response window

            showResponseWindow.value = true;

            console.log("[Korner] Response received", {
                hasResponse: !!latestResponse.value,
                loading: isLoadingResponse.value,
            });

            // Clear current query

            currentQuery.value = null;
        };

        const showSettings = () => {
            petStatus.value = "Settings coming soon!";
            setTimeout(() => {
                petStatus.value = "";
            }, 2000);
        };

        const startWindowDrag = (e) => {
            try {
                if (window.wails && window.wails.Window) {
                    window.wails.Window.Drag();
                }
            } catch (err) {
                // Fallback for frameless or if Drag not available
                console.log("Window drag not available");
            }
        };

        const minimizeToTray = () => {
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    window.go.main.App.HideWindow();
                }
            } catch (e) {
                console.error("Failed to hide window:", e);
            }
        };

        const closeResponseWindow = () => {
            showResponseWindow.value = false;
        };

        const pinResponseWindow = () => {
            // For now, just keep it open; could add pinning logic later
        };

        return {
            showScreenshotOverlay,
            showQueryWindow,
            showResponseWindow,
            currentQuery,
            latestResponse,
            isLoadingResponse,
            lastScreenshot,
            platform,
            petStatus,
            hotkey,
            triggerScreenshot,
            cancelScreenshot,
            handleScreenshotCaptured,
            cancelQueryWindow,
            handleQuerySubmit,
            showSettings,
            minimizeToTray,
            startWindowDrag,
            closeResponseWindow,
            pinResponseWindow,
        };
    },
};
</script>

<style>
#app {
    background: transparent !important;
    background: linear-gradient(
        135deg,
        rgba(102, 126, 234, 0.05) 0%,
        rgba(118, 75, 162, 0.05) 100%
    ) !important;
}

body {
    background: transparent !important;
    overflow: hidden;
}

/* Window drag handle - invisible draggable area at the top */
.window-drag-handle {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 50px;
    cursor: move;
    z-index: 1;
    pointer-events: auto;
    user-select: none;
}

.window-drag-handle:hover {
    background: rgba(102, 126, 234, 0.02);
}
</style>
