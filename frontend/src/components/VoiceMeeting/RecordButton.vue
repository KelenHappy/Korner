<template>
    <div 
        class="record-indicator" 
        :class="{ recording: isRecording, stopped: hasSavedFile }"
        @click="$emit('toggle')"
        :title="isRecording ? (stopTitle || '停止') : (startTitle || '開始')"
    >
        <div v-if="isRecording" class="pulse-ring"></div>
        <div class="record-dot">
            <span v-if="isRecording" class="duration-text">{{ formattedDuration }}</span>
        </div>
    </div>
</template>

<script>
export default {
    name: 'RecordButton',
    props: {
        isRecording: Boolean,
        hasSavedFile: Boolean,
        formattedDuration: String,
        startTitle: String,
        stopTitle: String
    },
    emits: ['toggle']
};
</script>

<style scoped>
.record-indicator {
    position: relative;
    width: 35px;
    height: 35px;
    cursor: pointer;
    transition: transform 0.2s;
}

.record-indicator:hover {
    transform: scale(1.1);
}

.record-indicator:active {
    transform: scale(0.95);
}

.pulse-ring {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 35px;
    height: 35px;
    border-radius: 50%;
    background: rgba(239, 68, 68, 0.3);
    animation: pulse-animation 1.5s ease-out infinite;
}

@keyframes pulse-animation {
    0% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 1;
    }
    100% {
        transform: translate(-50%, -50%) scale(2);
        opacity: 0;
    }
}

.record-dot {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 35px;
    height: 35px;
    border-radius: 50%;
    background: #ef4444;
    box-shadow: 0 4px 12px rgba(239, 68, 68, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s;
}

.record-indicator:not(.recording) .record-dot {
    background: #10b981;
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.5);
}

.record-indicator.stopped .record-dot {
    background: #6b7280;
    box-shadow: 0 4px 12px rgba(107, 114, 128, 0.5);
}

.duration-text {
    font-size: 8px;
    font-weight: 700;
    color: white;
    font-family: 'Courier New', monospace;
}
</style>
