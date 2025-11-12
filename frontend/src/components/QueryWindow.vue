<template>
    <div
        class="bg-white rounded-lg shadow-lg border border-slate-200 p-6 max-w-2xl mx-auto"
    >
        <h3 class="text-lg font-semibold text-slate-800 mb-4">
            Ask about this screenshot
        </h3>

        <!-- Screenshot Preview -->
        <div class="mb-4">
            <img
                :src="screenshot"
                alt="Screenshot"
                class="w-full rounded border border-slate-200 max-h-64 object-contain bg-slate-50"
            />
        </div>

        <!-- Query Input -->
        <div class="mb-4">
            <label class="block text-sm font-medium text-slate-700 mb-2">
                Your Question
            </label>
            <textarea
                v-model="queryText"
                @keydown.ctrl.enter="submit"
                @keydown.meta.enter="submit"
                class="w-full px-4 py-3 border border-slate-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                rows="4"
                placeholder="e.g., 'Explain this code', 'What's wrong with this error?', 'Summarize this chart'..."
                autofocus
            ></textarea>
            <div class="flex items-center justify-between mt-2">
                <span class="text-xs text-slate-500">
                    Press
                    <kbd
                        class="px-1.5 py-0.5 bg-slate-100 rounded text-xs border border-slate-300"
                        >{{ submitKey }}</kbd
                    >
                    to submit
                </span>
                <span class="text-xs text-slate-400"
                    >{{ queryText.length }} / 1000</span
                >
            </div>
        </div>

        <!-- Quick Prompts -->
        <div class="mb-6">
            <p class="text-xs font-medium text-slate-600 mb-2">
                Quick prompts:
            </p>
            <div class="flex flex-wrap gap-2">
                <button
                    v-for="prompt in quickPrompts"
                    :key="prompt"
                    @click="queryText = prompt"
                    class="px-3 py-1.5 text-xs bg-slate-100 hover:bg-slate-200 text-slate-700 rounded-full transition-colors"
                >
                    {{ prompt }}
                </button>
            </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end space-x-3">
            <button
                @click="cancel"
                class="px-4 py-2 text-slate-600 hover:text-slate-800 font-medium transition-colors"
            >
                Cancel
            </button>
            <button
                @click="submit"
                :disabled="!queryText.trim()"
                class="px-6 py-2 bg-blue-500 hover:bg-blue-600 disabled:bg-slate-300 disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors flex items-center space-x-2"
            >
                <svg
                    class="w-5 h-5"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M13 10V3L4 14h7v7l9-11h-7z"
                    />
                </svg>
                <span>Ask AI</span>
            </button>
        </div>
    </div>
</template>

<script>
import { ref, computed } from "vue";

export default {
    name: "QueryWindow",
    props: {
        screenshot: {
            type: String,
            required: true,
        },
    },
    emits: ["submit", "cancel"],
    setup(props, { emit }) {
        const queryText = ref("");
        const quickPrompts = [
            "Explain this code",
            "What does this mean?",
            "Find the error",
            "Summarize this",
            "Translate this",
            "Improve this code",
        ];

        const isMac = computed(() => {
            return navigator.userAgent.toLowerCase().includes("mac");
        });

        const submitKey = computed(() => {
            return isMac.value ? "Cmd+Enter" : "Ctrl+Enter";
        });

        const submit = () => {
            const trimmed = queryText.value.trim();
            if (trimmed) {
                const clamped = trimmed.slice(0, 1000);
                emit("submit", clamped);
                queryText.value = "";
            }
        };

        const cancel = () => {
            emit("cancel");
        };

        return {
            queryText,
            quickPrompts,
            submitKey,
            submit,
            cancel,
        };
    },
};
</script>
