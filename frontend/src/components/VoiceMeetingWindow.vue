<template>
    <div class="voice-meeting-window">
        <div class="voice-meeting-content">
            <div class="voice-header">
                <h2>🎤 {{ t('voiceMeeting.title') }}</h2>
                <button class="close-btn" @click="$emit('close')">✕</button>
            </div>

            <div class="voice-body">
                <div class="recording-status">
                    <div class="status-indicator" :class="{ recording: isRecording }">
                        <div class="pulse"></div>
                        <span class="status-icon">{{ isRecording ? '🔴' : '⚪' }}</span>
                    </div>
                    <div class="status-text">
                        <h3>{{ isRecording ? t('voiceMeeting.recording') : t('voiceMeeting.ready') }}</h3>
                        <p v-if="isRecording" class="duration">{{ formatDuration(duration) }}</p>
                    </div>
                </div>

                <div class="controls">
                    <button 
                        v-if="!isRecording" 
                        class="record-btn start"
                        @click="startRecording"
                    >
                        <span class="btn-icon">🎙️</span>
                        <span>{{ t('voiceMeeting.start') }}</span>
                    </button>
                    <button 
                        v-else 
                        class="record-btn stop"
                        @click="stopRecording"
                    >
                        <span class="btn-icon">⏹️</span>
                        <span>{{ t('voiceMeeting.stop') }}</span>
                    </button>
                </div>

                <div v-if="savedFile" class="saved-info">
                    <p class="success-msg">✅ {{ t('voiceMeeting.saved') }}</p>
                    <p class="file-path">{{ savedFile }}</p>
                    <button class="open-folder-btn" @click="openFolder">
                        📁 {{ t('voiceMeeting.openFolder') }}
                    </button>
                </div>

                <div v-if="errorMsg" class="error-info">
                    <p class="error-msg">❌ {{ errorMsg }}</p>
                </div>

                <div class="info-section">
                    <h4>{{ t('voiceMeeting.infoTitle') }}</h4>
                    <ul>
                        <li>{{ t('voiceMeeting.info1') }}</li>
                        <li>{{ t('voiceMeeting.info2') }}</li>
                        <li>{{ t('voiceMeeting.info3') }}</li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { ref, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';

export default {
    name: 'VoiceMeetingWindow',
    emits: ['close'],
    setup(props, { emit }) {
        const { t } = useI18n();
        const isRecording = ref(false);
        const duration = ref(0);
        const savedFile = ref('');
        const errorMsg = ref('');
        let durationInterval = null;

        const startRecording = async () => {
            try {
                errorMsg.value = '';
                savedFile.value = '';
                
                if (window.go && window.go.main && window.go.main.App) {
                    await window.go.main.App.StartRecording();
                    isRecording.value = true;
                    duration.value = 0;
                    
                    // Start duration counter
                    durationInterval = setInterval(async () => {
                        if (window.go && window.go.main && window.go.main.App) {
                            duration.value = await window.go.main.App.GetRecordingDuration();
                        }
                    }, 100);
                } else {
                    errorMsg.value = 'Backend not available';
                }
            } catch (error) {
                console.error('Failed to start recording:', error);
                errorMsg.value = String(error);
            }
        };

        const stopRecording = async () => {
            try {
                if (durationInterval) {
                    clearInterval(durationInterval);
                    durationInterval = null;
                }
                
                if (window.go && window.go.main && window.go.main.App) {
                    const filePath = await window.go.main.App.StopRecording();
                    isRecording.value = false;
                    savedFile.value = filePath;
                } else {
                    errorMsg.value = 'Backend not available';
                }
            } catch (error) {
                console.error('Failed to stop recording:', error);
                errorMsg.value = String(error);
                isRecording.value = false;
            }
        };

        const formatDuration = (seconds) => {
            const mins = Math.floor(seconds / 60);
            const secs = Math.floor(seconds % 60);
            return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
        };

        onUnmounted(() => {
            if (durationInterval) {
                clearInterval(durationInterval);
            }
            // Stop recording if still recording
            if (isRecording.value) {
                stopRecording();
            }
        });

        const openFolder = async () => {
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    await window.go.main.App.OpenRecordingFolder();
                }
            } catch (error) {
                console.error('Failed to open folder:', error);
                errorMsg.value = String(error);
            }
        };

        return {
            t,
            isRecording,
            duration,
            savedFile,
            errorMsg,
            startRecording,
            stopRecording,
            formatDuration,
            openFolder,
        };
    },
};
</script>

<style scoped>
.voice-meeting-window {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 10000;
}

.voice-meeting-content {
    background: white;
    border-radius: 16px;
    width: 90%;
    max-width: 600px;
    max-height: 90vh;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.voice-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 24px;
    border-bottom: 1px solid #e5e7eb;
}

.voice-header h2 {
    font-size: 20px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
}

.close-btn {
    background: none;
    border: none;
    font-size: 24px;
    color: #6b7280;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 4px;
    transition: all 0.2s;
}

.close-btn:hover {
    background: #f3f4f6;
    color: #1f2937;
}

.voice-body {
    padding: 32px 24px;
}

.recording-status {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
    margin-bottom: 32px;
}

.status-indicator {
    position: relative;
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.status-indicator .pulse {
    position: absolute;
    width: 100%;
    height: 100%;
    border-radius: 50%;
    background: #e5e7eb;
    opacity: 0;
}

.status-indicator.recording .pulse {
    animation: pulse 2s ease-out infinite;
}

@keyframes pulse {
    0% {
        transform: scale(0.8);
        opacity: 0.8;
    }
    100% {
        transform: scale(1.5);
        opacity: 0;
    }
}

.status-icon {
    font-size: 48px;
    z-index: 1;
}

.status-text h3 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 8px 0;
}

.duration {
    font-size: 24px;
    font-weight: 700;
    color: #ef4444;
    font-family: 'Courier New', monospace;
    margin: 0;
}

.controls {
    display: flex;
    justify-content: center;
    margin-bottom: 24px;
}

.record-btn {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 32px;
    border: none;
    border-radius: 12px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}

.record-btn.start {
    background: #10b981;
    color: white;
}

.record-btn.start:hover {
    background: #059669;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.record-btn.stop {
    background: #ef4444;
    color: white;
}

.record-btn.stop:hover {
    background: #dc2626;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.btn-icon {
    font-size: 20px;
}

.saved-info, .error-info {
    padding: 16px;
    border-radius: 8px;
    margin-bottom: 24px;
}

.saved-info {
    background: #d1fae5;
    border: 1px solid #10b981;
}

.error-info {
    background: #fee2e2;
    border: 1px solid #ef4444;
}

.success-msg, .error-msg {
    font-weight: 600;
    margin: 0 0 8px 0;
}

.success-msg {
    color: #065f46;
}

.error-msg {
    color: #991b1b;
    margin: 0;
}

.file-path {
    font-size: 12px;
    color: #047857;
    word-break: break-all;
    margin: 0 0 12px 0;
}

.open-folder-btn {
    background: #10b981;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}

.open-folder-btn:hover {
    background: #059669;
    transform: translateY(-1px);
}

.info-section {
    background: #f9fafb;
    padding: 16px;
    border-radius: 8px;
}

.info-section h4 {
    font-size: 14px;
    font-weight: 600;
    color: #374151;
    margin: 0 0 12px 0;
}

.info-section ul {
    margin: 0;
    padding-left: 20px;
}

.info-section li {
    font-size: 13px;
    color: #6b7280;
    margin-bottom: 8px;
}
</style>
