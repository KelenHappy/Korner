<template>
    <Teleport to="body">
        <div class="voice-meeting-widget">
            <!-- 返回按鈕 -->
            <button 
                class="bubble-btn back" 
                @click="$emit('close')"
                :title="t('voiceMeeting.back') || '返回'"
            >
                <span class="bubble-icon">✕</span>
            </button>

            <!-- 錄音按鈕 -->
            <RecordButton
                :is-recording="recording.isRecording.value"
                :has-saved-file="!!recording.savedFile.value"
                :formatted-duration="recording.formatDuration(recording.duration.value)"
                :start-title="t('voiceMeeting.start')"
                :stop-title="t('voiceMeeting.stop')"
                @toggle="handleToggleRecording"
            />

            <!-- 操作按鈕 (錄音完成後) -->
            <ActionButtons
                v-if="recording.savedFile.value && !recording.isRecording.value && !summary.isProcessing.value"
                :buttons="recordedButtons"
                @summary="handleGenerateSummary"
                @folder="handleOpenFolder"
            />

            <!-- 選擇本地音檔按鈕 -->
            <ActionButtons
                v-if="!recording.isRecording.value && !recording.savedFile.value"
                :buttons="selectFileButtons"
                @selectFile="handleSelectFile"
            />

            <!-- 處理進度 -->
            <MeetingSummaryProgress 
                v-if="summary.isProcessing.value"
                :status="summary.processingStatus.value"
                :progress="summary.processingProgress.value"
                :message="summary.processingText.value"
            />

            <!-- 摘要結果 -->
            <MeetingSummaryResult
                v-if="summary.showSummaryResult.value"
                :summary="summary.summaryResult.value"
                @close="handleCloseSummary"
            />

            <!-- 錯誤提示 -->
            <ErrorToast :message="errorMsg" />
        </div>
    </Teleport>
</template>

<script>
import { ref, computed, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRecording } from '../composables/useRecording';
import { useMeetingSummary } from '../composables/useMeetingSummary';
import RecordButton from './VoiceMeeting/RecordButton.vue';
import ActionButtons from './VoiceMeeting/ActionButtons.vue';
import ErrorToast from './VoiceMeeting/ErrorToast.vue';
import MeetingSummaryProgress from './MeetingSummaryProgress.vue';
import MeetingSummaryResult from './MeetingSummaryResult.vue';

export default {
    name: 'VoiceMeetingWindow',
    components: {
        RecordButton,
        ActionButtons,
        ErrorToast,
        MeetingSummaryProgress,
        MeetingSummaryResult
    },
    emits: ['close'],
    setup(props, { emit }) {
        const { t } = useI18n();
        const recording = useRecording();
        const summary = useMeetingSummary();
        const errorMsg = ref('');

        // 按鈕配置
        const recordedButtons = computed(() => [
            { name: 'summary', icon: '📝', title: t('voiceMeeting.meetingSummary'), event: 'summary' },
            { name: 'folder', icon: '📁', title: t('voiceMeeting.openFolder'), event: 'folder' }
        ]);

        const selectFileButtons = computed(() => [
            { name: 'file-select', icon: '🎵', title: '選擇本地音檔', event: 'selectFile', disabled: summary.isProcessing.value }
        ]);

        // 事件處理
        const handleToggleRecording = async () => {
            try {
                errorMsg.value = '';
                await recording.toggleRecording();
            } catch (error) {
                errorMsg.value = String(error);
            }
        };

        const handleSelectFile = async () => {
            try {
                errorMsg.value = '';
                const filePath = await window.go.main.App.SelectAudioFile();
                if (filePath) {
                    recording.savedFile.value = filePath;
                }
                // If filePath is empty, user cancelled - no error needed
            } catch (error) {
                errorMsg.value = String(error);
            }
        };

        const handleOpenFolder = async () => {
            try {
                await window.go.main.App.OpenRecordingFolder();
                emit('close');
            } catch (error) {
                errorMsg.value = String(error);
            }
        };

        const handleGenerateSummary = async () => {
            try {
                errorMsg.value = '';
                await summary.generateSummary(recording.savedFile.value);
            } catch (error) {
                errorMsg.value = String(error);
            }
        };

        const handleCloseSummary = () => {
            summary.closeSummary();
            emit('close');
        };

        onUnmounted(() => {
            recording.cleanup();
        });

        return {
            t,
            recording,
            summary,
            errorMsg,
            recordedButtons,
            selectFileButtons,
            handleToggleRecording,
            handleSelectFile,
            handleOpenFolder,
            handleGenerateSummary,
            handleCloseSummary
        };
    }
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
    pointer-events: auto;
}

.bubble-btn.back {
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
    background: #6b7280;
}

.bubble-btn.back:hover {
    background: #4b5563;
    transform: translateY(-3px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.bubble-icon {
    font-size: 16px;
}
</style>
