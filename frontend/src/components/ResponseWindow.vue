<template>
    <div
        class="fixed bottom-6 right-6 w-96 bg-white rounded-lg shadow-2xl border border-slate-200 z-50 max-h-96 flex flex-col"
    >
        <div
            class="flex items-center justify-between px-4 py-3 border-b border-slate-200 bg-gradient-to-r from-purple-50 to-pink-50"
        >
            <div class="flex items-center space-x-2">
                <span class="font-semibold text-slate-800">AI Response</span>
            </div>
            <div class="flex items-center space-x-1">
                <button
                    @click="copyResponse"
                    class="p-1.5 hover:bg-white rounded"
                    title="Copy"
                >
                    Copy
                </button>
                <button
                    @click="pinResponse"
                    class="p-1.5 hover:bg-white rounded"
                    title="Pin"
                >
                    Pin
                </button>
                <button
                    @click="close"
                    class="p-1.5 hover:bg-white rounded"
                    title="Close"
                >
                    ×
                </button>
            </div>
        </div>
        <div class="flex-1 overflow-auto p-4">
            <div v-if="loading">Loading...</div>
            <div v-else>
                <div v-if="screenshot" class="mb-4">
                    <img
                        :src="screenshot"
                        alt="Screenshot"
                        class="max-w-full h-auto rounded border border-slate-300"
                    />
                </div>
                <div class="text-sm text-slate-700 whitespace-pre-wrap">
                    {{ response }}
                </div>
            </div>
        </div>
    </div>
</template>
<script>
export default {
    name: "ResponseWindow",
    props: { response: String, loading: Boolean, screenshot: String },
    emits: ["close", "pin"],
    setup(props, { emit }) {
        const close = () => emit("close");
        const pinResponse = () => emit("pin");
        const copyResponse = () =>
            navigator.clipboard.writeText(props.response);
        return { close, pinResponse, copyResponse };
    },
};
</script>
