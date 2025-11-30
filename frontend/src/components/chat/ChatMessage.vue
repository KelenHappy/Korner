<template>
    <div :class="['message', isUser ? 'message-user' : 'message-ai']">
        <div class="message-avatar">
            <span v-if="isUser">👤</span>
            <span v-else>✨</span>
        </div>
        <div class="message-content">
            <div class="message-text">{{ content }}</div>
            <div class="message-time">{{ formattedTime }}</div>
        </div>
    </div>
</template>

<script>
import { computed } from 'vue';

export default {
    name: 'ChatMessage',
    props: {
        content: {
            type: String,
            required: true
        },
        role: {
            type: String,
            required: true,
            validator: (value) => ['user', 'assistant'].includes(value)
        },
        timestamp: {
            type: Date,
            required: true
        }
    },
    setup(props) {
        const isUser = computed(() => props.role === 'user');
        
        const formattedTime = computed(() => {
            return props.timestamp.toLocaleTimeString([], {
                hour: '2-digit',
                minute: '2-digit'
            });
        });

        return {
            isUser,
            formattedTime
        };
    }
};
</script>

<style scoped>
.message {
    display: flex;
    gap: 12px;
    animation: messageSlideIn 0.3s ease;
}

@keyframes messageSlideIn {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.message-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    flex-shrink: 0;
}

.message-user .message-avatar {
    background: #ffffff;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
}

.message-ai .message-avatar {
    background: #cccccc;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
}

.message-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.message-text {
    padding: 12px 16px;
    border-radius: 12px;
    font-size: 14px;
    line-height: 1.6;
    white-space: pre-wrap;
    word-wrap: break-word;
}

.message-user .message-text {
    background: #ffffff;
    color: #000000;
    border-top-left-radius: 4px;
}

.message-ai .message-text {
    background: #ffffff;
    color: #6b7280;
    border: 1px solid rgba(0, 0, 0, 0.06);
    border-top-left-radius: 4px;
}

.message-time {
    font-size: 11px;
    color: #94a3b8;
    padding-left: 4px;
}
</style>
