<template>
    <div class="conversation-item">
        <div class="conv-header">
            <span class="conv-time">{{ formattedTime }}</span>
            <span class="conv-provider">{{ conversation.provider }}</span>
            <button @click="$emit('delete')" class="delete-btn" title="刪除">✕</button>
        </div>
        <div class="conv-question">
            <strong>{{ questionLabel }}</strong>{{ conversation.question }}
        </div>
        <div class="conv-answer">
            <strong>{{ answerLabel }}</strong>{{ conversation.answer }}
        </div>
        <div v-if="conversation.screenshot_path" class="conv-screenshot">
            📷 {{ screenshotLabel }}{{ fileName }}
        </div>
    </div>
</template>

<script>
import { computed } from 'vue';

export default {
    name: 'ConversationItem',
    props: {
        conversation: {
            type: Object,
            required: true
        },
        questionLabel: {
            type: String,
            default: 'Q:'
        },
        answerLabel: {
            type: String,
            default: 'A:'
        },
        screenshotLabel: {
            type: String,
            default: 'Screenshot:'
        }
    },
    emits: ['delete'],
    setup(props) {
        const formattedTime = computed(() => {
            const date = new Date(props.conversation.timestamp);
            return date.toLocaleString('zh-TW', {
                year: 'numeric',
                month: '2-digit',
                day: '2-digit',
                hour: '2-digit',
                minute: '2-digit'
            });
        });

        const fileName = computed(() => {
            if (!props.conversation.screenshot_path) return '';
            return props.conversation.screenshot_path.split(/[\\/]/).pop();
        });

        return {
            formattedTime,
            fileName
        };
    }
};
</script>

<style scoped>
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
</style>
