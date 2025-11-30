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
                    <button @click="close" class="close-btn" :title="t('common.close')">✕</button>
                </div>
            </div>

            <div class="modal-body">
                <FilterTabs
                    :modelValue="filter"
                    :tabs="filterTabs"
                    @change="filter = $event"
                />

                <div class="history-list" v-if="conversations.length > 0">
                    <ConversationItem
                        v-for="conv in conversations"
                        :key="conv.id"
                        :conversation="conv"
                        :questionLabel="t('history.question')"
                        :answerLabel="t('history.answer')"
                        :screenshotLabel="t('history.screenshot')"
                        @delete="deleteItem(conv.id)"
                    />
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
import { ref, computed, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import FilterTabs from './history/FilterTabs.vue';
import ConversationItem from './history/ConversationItem.vue';

export default {
    name: 'HistoryWindow',
    components: {
        FilterTabs,
        ConversationItem
    },
    emits: ['close'],
    setup(props, { emit }) {
        const { t } = useI18n();
        const conversations = ref([]);
        const filter = ref('all');
        const loading = ref(false);

        const filterTabs = computed(() => [
            { value: 'all', label: t('history.all') },
            { value: 'today', label: t('history.today') },
            { value: 'recent', label: t('history.recent') }
        ]);

        const loadHistory = async () => {
            loading.value = true;
            try {
                let result;
                if (window.go && window.go.main && window.go.main.App) {
                    if (filter.value === 'today') {
                        result = await window.go.main.App.GetTodayHistory();
                    } else if (filter.value === 'recent') {
                        result = await window.go.main.App.GetRecentHistory(10);
                    } else {
                        result = await window.go.main.App.GetAllHistory();
                    }
                    conversations.value = result || [];
                }
            } catch (error) {
                console.error('Failed to load history:', error);
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
                console.error('Failed to delete item:', error);
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
                console.error('Failed to clear history:', error);
                alert(t('history.clearFailed') + error);
            }
        };

        const exportHistory = async () => {
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
                    const filename = `korner_history_${timestamp}.txt`;
                    await window.go.main.App.ExportHistoryToText(filename);
                    alert(t('history.exportSuccess') + filename);
                }
            } catch (error) {
                console.error('Failed to export history:', error);
                alert(t('history.exportFailed') + error);
            }
        };

        const close = () => {
            emit('close');
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
            filterTabs,
            loading,
            deleteItem,
            clearHistory,
            exportHistory,
            close
        };
    }
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

.history-list {
    flex: 1;
    overflow-y: auto;
    padding: 16px 24px;
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
