<template>
    <div class="message message-ai">
        <div class="message-avatar">
            <span>✨</span>
        </div>
        <div class="message-content">
            <div class="loading-dots">
                <span></span>
                <span></span>
                <span></span>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'LoadingIndicator'
};
</script>

<style scoped>
.message {
    display: flex;
    gap: 12px;
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
    background: #cccccc;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
}

.message-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.loading-dots {
    display: flex;
    gap: 6px;
    padding: 12px 16px;
}

.loading-dots span {
    width: 8px;
    height: 8px;
    background: #f5576c;
    border-radius: 50%;
    animation: bounce 1.4s infinite ease-in-out both;
}

.loading-dots span:nth-child(1) {
    animation-delay: -0.32s;
}

.loading-dots span:nth-child(2) {
    animation-delay: -0.16s;
}

@keyframes bounce {
    0%, 80%, 100% {
        transform: scale(0);
    }
    40% {
        transform: scale(1);
    }
}
</style>
