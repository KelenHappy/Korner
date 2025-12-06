<template>
    <Teleport to="body">
        <div class="voice-meeting-widget">
            <!-- 返回按鈕 (灰色) -->
            <button 
                class="bubble-btn back" 
                @click="$emit('close')"
                :title="t('voiceMeeting.back') || '返回'"
            >
                <span class="bubble-icon">✕</span>
            </button>

            <!-- 錄音指示器 -->
            <div 
                class="record-indicator" 
                :class="{ recording: isRecording, stopped: savedFile }"
                @click="toggleRecording"
                :title="isRecording ? t('voiceMeeting.stop') : t('voiceMeeting.start')"
            >
                <div v-if="isRecording" class="pulse-ring"></div>
                <div class="record-dot">
                    <span v-if="isRecording" class="duration-text">{{ formatDuration(duration) }}</span>
                </div>
            </div>

            <!-- 操作按鈕 (錄音完成後顯示) -->
            <transition name="bubble-fade">
                <div v-if="savedFile && !isRecording" class="action-bubbles">
                    <button 
                        class="bubble-btn summary" 
                        @click="generateMeetingSummary"
                        :disabled="isProcessing"
                        :title="t('voiceMeeting.meetingSummary')"
                    >
                        <span class="bubble-icon">📝</span>
                    </button>
                    <button 
                        class="bubble-btn folder" 
                        @click="openFolder"
                        :title="t('voiceMeeting.openFolder')"
                    >
                        <span class="bubble-icon">📁</span>
                    </button>
                </div>
            </transition>

            <!-- 錯誤提示 -->
            <transition name="fade">
                <div v-if="errorMsg" class="error-toast">
                    {{ errorMsg }}
                </div>
            </transition>
        </div>
    </Teleport>
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
        const isProcessing = ref(false);
        let durationInterval = null;

        const startRecording = async () => {
            console.log('[VoiceMeeting] Start recording button clicked');
            try {
                errorMsg.value = '';
                savedFile.value = '';
                
                if (window.go && window.go.main && window.go.main.App) {
                    console.log('[VoiceMeeting] Calling backend StartRecording...');
                    await window.go.main.App.StartRecording();
                    isRecording.value = true;
                    duration.value = 0;
                    console.log('[VoiceMeeting] Recording started successfully');
                    
                    // Start duration counter
                    durationInterval = setInterval(async () => {
                        if (window.go && window.go.main && window.go.main.App) {
                            duration.value = await window.go.main.App.GetRecordingDuration();
                        }
                    }, 100);
                } else {
                    console.error('[VoiceMeeting] Backend not available');
                    errorMsg.value = 'Backend not available';
                }
            } catch (error) {
                console.error('[VoiceMeeting] Failed to start recording:', error);
                errorMsg.value = String(error);
            }
        };

        const stopRecording = async () => {
            console.log('[VoiceMeeting] Stop recording button clicked');
            try {
                if (durationInterval) {
                    clearInterval(durationInterval);
                    durationInterval = null;
                }
                
                if (window.go && window.go.main && window.go.main.App) {
                    console.log('[VoiceMeeting] Calling backend StopRecording...');
                    const filePath = await window.go.main.App.StopRecording();
                    isRecording.value = false;
                    savedFile.value = filePath;
                    console.log('[VoiceMeeting] Recording stopped, file:', filePath);
                } else {
                    console.error('[VoiceMeeting] Backend not available');
                    errorMsg.value = 'Backend not available';
                }
            } catch (error) {
                console.error('[VoiceMeeting] Failed to stop recording:', error);
                errorMsg.value = String(error);
                isRecording.value = false;
            }
        };

        const toggleRecording = () => {
            if (isRecording.value) {
                stopRecording();
            } else {
                startRecording();
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
            if (isRecording.value) {
                stopRecording();
            }
        });

        const generateMeetingSummary = async () => {
            console.log('[VoiceMeeting] Generate meeting summary button clicked');
            try {
                isProcessing.value = true;
                errorMsg.value = '';
                
                if (window.go && window.go.main && window.go.main.App) {
                    console.log('[VoiceMeeting] Calling backend GenerateMeetingSummary...');
                    const summary = await window.go.main.App.GenerateMeetingSummary(savedFile.value);
                    console.log('[VoiceMeeting] Summary generated:', summary.substring(0, 100));
                    
                    // 發送摘要到聊天窗口
                    if (window.sendMessageToChat) {
                        window.sendMessageToChat(summary);
                    }
                    
                    // 關閉錄音窗口並返回
                    emit('close');
                } else {
                    errorMsg.value = 'Backend not available';
                }
            } catch (error) {
                console.error('Failed to generate meeting summary:', error);
                errorMsg.value = String(error);
            } finally {
                isProcessing.value = false;
            }
        };

        const openFolder = async () => {
            console.log('[VoiceMeeting] Open folder button clicked');
            try {
                if (window.go && window.go.main && window.go.main.App) {
                    await window.go.main.App.OpenRecordingFolder();
                    console.log('[VoiceMeeting] Folder opened');
                }
                
                // 打開資料夾後也關閉窗口並返回
                emit('close');
            } catch (error) {
                console.error('[VoiceMeeting] Failed to open folder:', error);
                errorMsg.value = String(error);
            }
        };

        return {
            t,
            isRecording,
            duration,
            savedFile,
            errorMsg,
            isProcessing,
            toggleRecording,
            formatDuration,
            openFolder,
            generateMeetingSummary,
        };
    },
};
</script>

<style scoped>
.voice-meeting-widget {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 9999;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
}

/* 錄音指示器 */
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

/* 脈衝環 */
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

/* 紅點 */
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

/* 時間顯示 */
.duration-text {
    font-size: 8px;
    font-weight: 700;
    color: white;
    font-family: 'Courier New', monospace;
}

/* 操作按鈕 */
.action-bubbles {
    display: flex;
    gap: 10px;
    animation: bubble-slide-up 0.3s ease-out;
}

@keyframes bubble-slide-up {
    from {
        transform: translateY(20px);
        opacity: 0;
    }
    to {
        transform: translateY(0);
        opacity: 1;
    }
}

.bubble-btn {
    width: 35px;
    height: 35px;
    border-radius: 50%;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    transition: all 0.2s;
    pointer-events: auto;
}

.bubble-btn:hover:not(:disabled) {
    transform: translateY(-3px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.bubble-btn:active:not(:disabled) {
    transform: translateY(-1px);
}

.bubble-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.bubble-btn.back {
    background: #6b7280;
}

.bubble-btn.back:hover:not(:disabled) {
    background: #4b5563;
}

.bubble-btn.summary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.bubble-btn.folder {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.bubble-icon {
    font-size: 16px;
}

/* 錯誤提示 */
.error-toast {
    position: fixed;
    top: 50%;
    right: 100px;
    transform: translateY(-50%);
    background: #fee2e2;
    color: #991b1b;
    padding: 12px 16px;
    border-radius: 8px;
    border: 1px solid #ef4444;
    font-size: 12px;
    font-weight: 600;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    max-width: 250px;
    word-break: break-word;
    animation: toast-slide-in 0.3s ease-out;
}

@keyframes toast-slide-in {
    from {
        transform: translateX(100%);
        opacity: 0;
    }
    to {
        transform: translateX(0);
        opacity: 1;
    }
}

/* 過渡動畫 */
.bubble-fade-enter-active,
.bubble-fade-leave-active {
    transition: all 0.3s ease;
}

.bubble-fade-enter-from {
    transform: translateY(20px);
    opacity: 0;
}

.bubble-fade-leave-to {
    transform: translateY(-20px);
    opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
</style>
