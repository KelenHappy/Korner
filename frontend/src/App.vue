<template>
    <div id="app" class="h-screen w-screen overflow-hidden">
        <ScreenshotOverlay
            v-if="showScreenshotOverlay"
            @screenshot-captured="handleScreenshotCaptured"
            @cancel="cancelScreenshot"
        />
        <FirstRunGuide v-if="showFirstRunGuide" @close="closeFirstRunGuide" />
        <div
            v-if="!showScreenshotOverlay"
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
                                SnapAsk
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
                    <QueryWindow
                        v-if="currentQuery"
                        :screenshot="currentQuery.screenshot"
                        @submit="handleQuerySubmit"
                        @cancel="currentQuery = null"
                    />
                </div>
            </main>
        </div>
        <ResponseWindow
            v-if="showResponseWindow && latestResponse"
            :response="latestResponse"
            :loading="isLoadingResponse"
            @close="showResponseWindow = false"
            @pin="pinResponse"
        />
    </div>
</template>

<script>
import { ref, computed, onMounted } from "vue";
import ScreenshotOverlay from "./components/ScreenshotOverlay.vue";
import QueryWindow from "./components/QueryWindow.vue";
import ResponseWindow from "./components/ResponseWindow.vue";
import FirstRunGuide from "./components/FirstRunGuide.vue";
import {
    WindowSetAlwaysOnTop,
    WindowFullscreen,
    WindowUnfullscreen,
} from "@wailsapp/runtime";

export default {
    name: "App",
    components: {
        ScreenshotOverlay,
        QueryWindow,
        ResponseWindow,
        FirstRunGuide,
    },
    setup() {
        const showScreenshotOverlay = ref(false);
        const showFirstRunGuide = ref(false);
        const currentQuery = ref(null);
        const conversationHistory = ref([]);
        const showResponseWindow = ref(false);
        const latestResponse = ref("");
        const isLoadingResponse = ref(false);
        const platform = ref("unknown");

        onMounted(async () => {
            // Try to get platform from Wails backend, fallback to userAgent detection
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    platform.value = await window.go.main.App.GetPlatform();
                } else {
                    // Fallback for dev mode or if backend not ready
                    const ua = navigator.userAgent.toLowerCase();
                    if (ua.includes("linux")) platform.value = "linux";
                    else if (ua.includes("mac")) platform.value = "darwin";
                    else if (ua.includes("win")) platform.value = "windows";
                    else platform.value = "unknown";
                }
            } catch (err) {
                console.warn(
                    "Failed to get platform from backend, using userAgent:",
                    err,
                );
                const ua = navigator.userAgent.toLowerCase();
                if (ua.includes("linux")) platform.value = "linux";
                else if (ua.includes("mac")) platform.value = "darwin";
                else if (ua.includes("win")) platform.value = "windows";
                else platform.value = "unknown";
            }

            const isFirstRun = localStorage.getItem("snapask_first_run");
            if (!isFirstRun && platform.value === "linux") {
                showFirstRunGuide.value = true;
            }
            try {
                WindowSetAlwaysOnTop(true);
            } catch (e) {
                // ignore if not running under Wails runtime
            }
        });

        const hotkey = computed(() =>
            platform.value === "darwin" ? "Cmd+Option+Q" : "Ctrl+Alt+Q",
        );

        const triggerScreenshot = async () => {
            try {
                WindowSetAlwaysOnTop(true);
            } catch (e) {}
            try {
                WindowFullscreen();
            } catch (e) {}
            showScreenshotOverlay.value = true;
        };

        const handleScreenshotCaptured = (screenshotData) => {
            currentQuery.value = {
                screenshot: screenshotData,
                timestamp: new Date(),
            };
            showScreenshotOverlay.value = false;
            try {
                WindowUnfullscreen();
            } catch (e) {}
        };

        const cancelScreenshot = () => {
            showScreenshotOverlay.value = false;
            try {
                WindowUnfullscreen();
            } catch (e) {}
        };

        const handleQuerySubmit = async (queryText) => {
            if (!currentQuery.value) return;
            const query = {
                query: queryText,
                screenshot: currentQuery.value.screenshot,
                response: "",
                loading: true,
            };
            conversationHistory.value.push(query);
            currentQuery.value = null;
            showResponseWindow.value = true;
            isLoadingResponse.value = true;

            try {
                // Try to call Wails backend
                if (window.go && window.go.main && window.go.main.App) {
                    // Extract base64 from data URL if present
                    let screenshotB64 = currentQuery.value?.screenshot || "";
                    if (screenshotB64.startsWith("data:image")) {
                        const commaIndex = screenshotB64.indexOf(",");
                        if (commaIndex !== -1) {
                            screenshotB64 = screenshotB64.substring(
                                commaIndex + 1,
                            );
                        }
                    }

                    const response = await window.go.main.App.QueryLLM(
                        queryText,
                        screenshotB64,
                    );
                    query.response = response;
                    query.loading = false;
                    latestResponse.value = response;
                    isLoadingResponse.value = false;
                } else {
                    // Fallback simulation for dev mode
                    setTimeout(() => {
                        const response = `Response to: "${queryText}"\n\nThis connects to AMD GPT OSS 120B model.\n\n(Running in dev mode - backend not connected)`;
                        query.response = response;
                        query.loading = false;
                        latestResponse.value = response;
                        isLoadingResponse.value = false;
                    }, 2000);
                }
            } catch (error) {
                console.error("QueryLLM failed:", error);
                query.response = `Error: ${error.message || "Failed to query LLM"}\n\nPlease check:\n1. AMD_LLM_ENDPOINT is set\n2. AMD_API_KEY is valid\n3. Network connection`;
                query.loading = false;
                latestResponse.value = query.response;
                isLoadingResponse.value = false;
            }
        };

        const pinResponse = () => {
            showResponseWindow.value = false;
        };
        const closeFirstRunGuide = () => {
            showFirstRunGuide.value = false;
            localStorage.setItem("snapask_first_run", "true");
        };

        return {
            showScreenshotOverlay,
            showFirstRunGuide,
            currentQuery,
            conversationHistory,
            showResponseWindow,
            latestResponse,
            isLoadingResponse,
            hotkey,
            triggerScreenshot,
            cancelScreenshot,
            handleScreenshotCaptured,
            handleQuerySubmit,
            pinResponse,
            closeFirstRunGuide,
        };
    },
};
</script>
