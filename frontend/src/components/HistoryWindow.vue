<template>
    <div class="history-overlay" @click.self="close">
        <div class="history-modal">
            <div class="modal-header">
                <h2 class="modal-title">📜 {{ t('history.title') }}</h2>
                <div class="header-actions">
                    <button @click="exportHistory" class="action-btn" :title="t('history.export')">
                        💾 {{ t('history.export') }}
                    </button>
                    <button @click="clearHistory" class="action-btn danger" :title="t('history.clear')">
                        🗑️ {{ t('history.clear') }}
                    </button>
                    <button @click="close" class="close-btn" :title="t('common.close')">
                        ✕
                    </button>
                </div>
            </div>

            <div class="modal-body">
                <div class="filter-tabs">
                    <button
                        @click="filter = 'all'"
                        :class="['tab-btn', { active: filter === 'all' }]"
                    >
                        {{ t('history.all') }}
                    </button>
                    <button
                        @click="filter = 'today'"
                        :class="['tab-btn', { active: filter === 'today' }]"
                    >
                        {{ t('history.today') }}
                    </button>
                    <button
                        @click="filter = 'recent'"
                        :class="['tab-btn', { active: filter === 'recent' }]"
                    >
                        {{ t('history.recent') }}
                    </button>
                </div>

                <div class="history-list" v-if="conversations.length > 0">
                    <div
                        v-for="conv in conversations"
                        :key="conv.id"
                        class="conversation-item"
                    >
                        <div class="conv-header">
                            <span class="conv-time">{{ formatTime(conv.timestamp) }}</span>
                            <span class="conv-provider">{{ conv.provider }}</span>
                            <button
                                @click="deleteItem(conv.id)"
                                class="delete-btn"
                                title="刪除"
                            >
                                ✕
                            </button>
                        </div>
                        <div class="conv-question">
                            <strong>{{ t('history.question') }}</strong>{{ conv.question }}
                        </div>
                        <div class="conv-answer">
                            <strong>{{ t('history.answer') }}</strong>{{ conv.answer }}
                        </div>
                        <div v-if="conv.screenshot_path" class="conv-screenshot">
                            📷 {{ t('history.screenshot') }}{{ getFileName(conv.screenshot_path) }}
                        </div>
                    </div>
                </div>

                <div v-else class="empty-state">
                    <div class="empty-icon">📭</div>
                    <p>{{ t('history.empty') }}</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, onMounted, watch } from "vue";
import { useI18n } from "vue-i18n";

export default {
    name: "HistoryWindow",
    emits: ["close"],
    setup(props, { emit }) {
        const { t } = useI18n();
        const conversations = ref([]);
        const filter = ref("all");
        const loading = ref(false);

        const loadHistory = async () => {
            loading.value = true;
            try {
                let result;
                if (window.go && window.go.main && window.go.main.App) {
                    if (filter.value === "today") {
                        result = await window.go.main.App.GetTodayHistory();
                    } else if (filter.value === "recent") {
                        result = await window.go.main.App.GetRecentHistory(10);
                    } else {
                        result = await window.go.main.App.GetAllHistory();
                    }
                    conversations.value = result || [];
                }
            } catch (error) {
                console.error("Failed to load history:", error);
                conversations.value = [];
            } finally {
                loading.value = false;
            }
        };

        const deleteItem = async (id) => {
            if (!confirm(t('history.deleteConfirm'))) return;

            try {
                if (window.go && window.go.main && window.go.main.App) {
                    await window.go.main.App.DeleteHistoryItem(id);
                    await loadHistory();
                }
            } catch (error) {
                console.error("Failed to delete item:", error);
                alert(t('history.deleteFailed') + error);
            }
        };

        const clearHistory = async () => {
            if (!confirm(t('history.clearConfirm'))) return;

            try {
                if (window.go && window.go.main && window.go.main.App) {
                    await window.go.main.App.ClearHistory();
                    conversations.value = [];
                }
            } catch (error) {
                console.error("Failed to clear history:", error);
                alert(t('history.clearFailed') + error);
            }
        };

        const exportHistory = async () => {
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    const timestamp = new Date().toISOString().replace(/[:.]/g, "-");
                    const filename = `korner_history_${timestamp}.txt`;
                    await window.go.main.App.ExportHistoryToText(filename);
                    alert(t('history.exportSuccess') + filename);
                }
            } catch (error) {
                console.error("Failed to export history:", error);
                alert(t('history.exportFailed') + error);
            }
        };

        const formatTime = (timestamp) => {
            const date = new Date(timestamp);
            return date.toLocaleString("zh-TW", {
                year: "numeric",
                month: "2-digit",
                day: "2-digit",
                hour: "2-digit",
                minute: "2-digit",
            });
        };

        const getFileName = (path) => {
            return path.split(/[\\/]/).pop();
        };

        const close = () => {
            emit("close");
        };

        watch(filter, () => {
            loadHistory();
        });

        onMounted(() => {
            loadHistory();
        });

        return {
            t,
            conversations,
            filter,
            loading,
            deleteItem,
            clearHistory,
            exportHistory,
            formatTime,
            getFileName,
            close,
        };
    },
};
</script>

<style scoped>
.history-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10000;
    animation: fadeIn 0.2s ease;
}

.history-modal {
    width: 90%;
    max-width: 800px;
    max-height: 90vh;
    background: white;
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    display: flex;
    flex-direction: column;
    animation: slideUp 0.3s ease;
}

.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 24px;
    border-bottom: 1px solid #e2e8f0;
}

.modal-title {
    font-size: 20px;
    font-weight: 700;
    color: #1e293b;
    margin: 0;
}

.header-actions {
    display: flex;
    gap: 8px;
}

.action-btn {
    padding: 8px 16px;
    background: #f1f5f9;
    border: none;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}

.action-btn:hover {
    background: #e2e8f0;
}

.action-btn.danger {
    background: #fee2e2;
    color: #dc2626;
}

.action-btn.danger:hover {
    background: #fecaca;
}

.close-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: #f1f5f9;
    border-radius: 8px;
    font-size: 18px;
    cursor: pointer;
    transition: all 0.2s;
}

.close-btn:hover {
    background: #e2e8f0;
    transform: scale(1.1);
}

.modal-body {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
}

.filter-tabs {
    display: flex;
    gap: 8px;
    padding: 16px 24px;
    border-bottom: 1px solid #e2e8f0;
}

.tab-btn {
    padding: 8px 16px;
    background: transparent;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}

.tab-btn:hover {
    background: #f1f5f9;
}

.tab-btn.active {
    background: #000;
    color: white;
    border-color: #000;
}

.history-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px 24px;
}

.conversation-item {
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    border-radius: 12px;
    padding: 16px;
    margin-bottom: 12px;
    transition: all 0.2s;
}

.conversation-item:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.conv-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
}

.conv-time {
    font-size: 12px;
    color: #64748b;
    font-weight: 600;
}

.conv-provider {
    font-size: 11px;
    padding: 2px 8px;
    background: #000;
    color: white;
    border-radius: 4px;
    font-weight: 600;
}

.delete-btn {
    margin-left: auto;
    width: 24px;
    height: 24px;
    border: none;
    background: #fee2e2;
    color: #dc2626;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.2s;
}

.delete-btn:hover {
    background: #fecaca;
    transform: scale(1.1);
}

.conv-question,
.conv-answer {
    font-size: 14px;
    line-height: 1.6;
    margin-bottom: 8px;
    color: #334155;
}

.conv-question strong,
.conv-answer strong {
    color: #1e293b;
    font-weight: 700;
}

.conv-screenshot {
    font-size: 12px;
    color: #64748b;
    margin-top: 8px;
    padding-top: 8px;
    border-top: 1px solid #e2e8f0;
}

.empty-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
}

.empty-icon {
    font-size: 64px;
    margin-bottom: 16px;
    opacity: 0.5;
}

.empty-state p {
    font-size: 16px;
    color: #64748b;
}

@keyframes fadeIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
</style>
