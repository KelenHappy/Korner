<template>
    <div id="app" class="h-screen w-screen overflow-hidden bg-transparent">
        <!-- Fullscreen screenshot overlay -->
        <ScreenshotOverlay
            v-if="showScreenshotOverlay"
            @screenshot-captured="handleScreenshotCaptured"
            @cancel="cancelScreenshot"
        />

        <!-- Desktop Pet Mode (floating transparent pet) - ALWAYS SHOW UNLESS SCREENSHOTTING -->
        <DesktopPet
            v-if="!showScreenshotOverlay && !currentQuery"
            :status="petStatus"
            :chatMessage="latestResponse"
            :chatLoading="isLoadingResponse"
            @screenshot="triggerScreenshot"
            @settings="showSettings"
            @minimize="minimizeToTray"
            @close="closeApp"
        />

        <!-- Query window - shown after screenshot is captured -->
        <QueryWindow
            v-if="currentQuery"
            :screenshot="currentQuery.screenshot"
            @submit="handleQuerySubmit"
            @cancel="currentQuery = null"
        />

        <!-- Main app UI - DISABLED, always use desktop pet mode -->
        <div
            v-if="false"
            class="h-full w-full flex flex-col bg-gradient-to-br from-slate-50 to-slate-100"
        >
            <header
                class="bg-white shadow-sm border-b border-slate-200 px-6 py-4"
            >
                <div class="flex items-center justify-between">
                    <div class="flex items-center space-x-3">
                        <div
                            class="w-10 h-10 bg-blue-500 rounded-lg flex items-center justify-center"
                        >
                            <span class="text-white font-bold">S</span>
                        </div>
                        <div>
                            <h1 class="text-xl font-bold text-slate-800">
                                Korner
                            </h1>
                            <p class="text-xs text-slate-500">
                                AI Screenshot Assistant
                            </p>
                        </div>
                    </div>
                    <button
                        @click="triggerScreenshot"
                        class="px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-medium transition-colors"
                    >
                        New Screenshot
                    </button>
                </div>
            </header>

            <main class="flex-1 overflow-auto p-6">
                <div class="max-w-4xl mx-auto">
                    <div
                        v-if="!currentQuery && conversationHistory.length === 0"
                        class="text-center py-12"
                    >
                        <h2 class="text-2xl font-bold text-slate-800 mb-2">
                            Ready to Ask
                        </h2>
                        <p class="text-slate-600 mb-6">
                            Press
                            <kbd
                                class="px-2 py-1 bg-slate-200 rounded text-sm"
                                >{{ hotkey }}</kbd
                            >
                            or click "New Screenshot" to start
                        </p>
                    </div>

                    <div
                        v-if="conversationHistory.length > 0"
                        class="space-y-4"
                    >
                        <div
                            v-for="(item, index) in conversationHistory"
                            :key="index"
                            class="bg-white rounded-lg shadow-sm border border-slate-200 p-4"
                        >
                            <div class="mb-4">
                                <p
                                    class="text-sm font-medium text-slate-700 mb-2"
                                >
                                    {{ item.query }}
                                </p>
                                <img
                                    v-if="item.screenshot"
                                    :src="item.screenshot"
                                    class="rounded border border-slate-200 max-w-sm"
                                />
                            </div>
                            <div class="bg-slate-50 rounded-lg p-3">
                                <div v-if="item.loading">Loading...</div>
                                <div
                                    v-else
                                    class="text-sm text-slate-700 whitespace-pre-wrap"
                                >
                                    {{ item.response }}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>

        <!-- Floating response window (visible in both normal & compact modes) -->
        <ResponseWindow
            v-if="showResponseWindow && latestResponse"
            :response="latestResponse"
            :loading="isLoadingResponse"
            @close="closeResponseWindow"
            @pin="pinResponse"
        />
    </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from "vue";
import ScreenshotOverlay from "./components/ScreenshotOverlay.vue";
import QueryWindow from "./components/QueryWindow.vue";
import ResponseWindow from "./components/ResponseWindow.vue";
import DesktopPet from "./components/DesktopPet.vue";

// Wails Runtime: used for window management (always-on-top, resize, move, fullscreen)
import {
    WindowSetAlwaysOnTop,
    WindowFullscreen,
    WindowUnfullscreen,
    WindowGetSize,
    WindowGetPosition,
    WindowSetSize,
    WindowSetPosition,
    ScreenGetAll,
    EventsOn,
    EventsOff,
    WindowShow,
} from "../wailsjs/runtime/runtime";

export default {
    name: "App",
    components: {
        ScreenshotOverlay,
        QueryWindow,
        ResponseWindow,
        DesktopPet,
    },
    setup() {
        const showScreenshotOverlay = ref(false);
        const currentQuery = ref(null);
        const conversationHistory = ref([]);
        const showResponseWindow = ref(false);
        const latestResponse = ref(""); // Empty on startup - chatbox hidden
        const isLoadingResponse = ref(false);
        const platform = ref("unknown");

        // Desktop pet mode - ALWAYS true (removed toggle logic)
        const compactMode = ref(true);
        const prevSize = ref({ w: 0, h: 0 });
        const prevPos = ref({ x: 0, y: 0 });
        const petStatus = ref("");

        onMounted(async () => {
            // Detect platform (from backend if available)
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

            // Desktop pet mode is always on (compactMode always true)

            // Try to keep window always on top
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
            // Show window if hidden
            try {
                WindowShow();
                WindowSetAlwaysOnTop(true);
            } catch {}
            try {
                WindowFullscreen();
            } catch {}
            showScreenshotOverlay.value = true;
            // Keep compactMode true - don't change it
        };

        const handleScreenshotCaptured = (screenshotData) => {
            currentQuery.value = {
                screenshot: screenshotData,
                timestamp: new Date(),
            };
            showScreenshotOverlay.value = false;
            try {
                WindowUnfullscreen();
            } catch {}
        };

        const cancelScreenshot = () => {
            showScreenshotOverlay.value = false;
            try {
                WindowUnfullscreen();
            } catch {}
        };

        const handleQuerySubmit = async (queryText) => {
            if (!currentQuery.value) return;

            // Capture screenshot data BEFORE clearing currentQuery to avoid losing it
            let screenshotDataUrl = currentQuery.value.screenshot || "";
            let screenshotB64 = screenshotDataUrl;
            if (screenshotDataUrl.startsWith("data:image")) {
                const comma = screenshotDataUrl.indexOf(",");
                if (comma !== -1) {
                    screenshotB64 = screenshotDataUrl.substring(comma + 1);
                }
            }

            const query = {
                query: queryText,
                screenshot: screenshotDataUrl,
                response: "",
                loading: true,
            };
            conversationHistory.value.push(query);

            // Stay in desktop pet mode (always true)
            showResponseWindow.value = false;
            isLoadingResponse.value = true;
            latestResponse.value = ""; // Will show loading animation in chatbox

            // Clear current query after we saved the screenshot
            currentQuery.value = null;

            // Exit fullscreen and show pet with loading chat
            try {
                WindowUnfullscreen();
            } catch {}

            try {
                if (window.go && window.go.main && window.go.main.App) {
                    const response = await window.go.main.App.QueryLLM(
                        queryText,
                        screenshotB64,
                    );
                    query.response = response;
                    query.loading = false;
                    latestResponse.value = response;
                    isLoadingResponse.value = false;
                } else {
                    // Dev fallback
                    setTimeout(() => {
                        const response = `Response to: "${queryText}"\n\nThis connects to AMD GPT OSS 120B model.\n\n(Running in dev mode - backend not connected)`;
                        query.response = response;
                        query.loading = false;
                        latestResponse.value = response;
                        isLoadingResponse.value = false;
                    }, 1200);
                }
            } catch (error) {
                const msg = error?.message || "Failed to query LLM";
                query.response = `Error: ${msg}\n\nPlease check:\n1. AMD_LLM_ENDPOINT is set\n2. AMD_API_KEY is valid\n3. Network connection`;
                query.loading = false;
                latestResponse.value = query.response;
                isLoadingResponse.value = false;
            }
        };

        // Make the app window compact and keep it floating on all windows
        const pinResponse = async () => {
            // compactMode is always true
            showResponseWindow.value = true;
            try {
                prevSize.value = await WindowGetSize();
                prevPos.value = await WindowGetPosition();

                const screens = await ScreenGetAll();
                const primary = screens.find((s) => s.isPrimary) || screens[0];

                const width = 420;
                const height = 380;
                const margin = 16;
                const x = Math.max(
                    0,
                    (primary?.width || 1280) - width - margin,
                );
                const y = Math.max(
                    0,
                    (primary?.height || 800) - height - margin,
                );

                WindowSetSize(width, height);
                WindowSetPosition(x, y);
                WindowSetAlwaysOnTop(true);
            } catch {
                // ignore when runtime not available (dev)
            }
        };

        // Close the floating window and restore the window size/position if compact
        const closeResponseWindow = async () => {
            showResponseWindow.value = false;
            // compactMode is always true - no need to restore

            // Try to restore window size if needed
            if (true) {
                try {
                    if (prevSize.value.w && prevSize.value.h) {
                        WindowSetSize(prevSize.value.w, prevSize.value.h);
                    }
                    if (
                        typeof prevPos.value.x === "number" &&
                        typeof prevPos.value.y === "number"
                    ) {
                        WindowSetPosition(prevPos.value.x, prevPos.value.y);
                    }
                } catch {
                    // ignore when runtime not available (dev)
                }
            }
        };

        const showSettings = () => {
            // TODO: Implement settings
            petStatus.value = "Settings coming soon!";
            setTimeout(() => {
                petStatus.value = "";
            }, 2000);
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

        const closeApp = () => {
            try {
                if (window.wails && window.wails.Quit) {
                    window.wails.Quit();
                }
            } catch (e) {
                console.error("Failed to quit:", e);
            }
        };

        return {
            // state
            showScreenshotOverlay,
            currentQuery,
            conversationHistory,
            showResponseWindow,
            latestResponse,
            isLoadingResponse,
            platform,
            compactMode,
            petStatus,

            // computed
            hotkey,

            // actions
            triggerScreenshot,
            cancelScreenshot,
            handleScreenshotCaptured,
            handleQuerySubmit,
            pinResponse,
            closeResponseWindow,
            showSettings,
            minimizeToTray,
            closeApp,
        };
    },
};
</script>

<style>
#app {
    background: transparent !important;
    /* Fallback if transparency doesn't work */
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

/* Force transparency for desktop pet container */
.desktop-pet-container {
    background: transparent !important;
}
</style>
