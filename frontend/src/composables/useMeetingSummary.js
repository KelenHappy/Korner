import { ref } from 'vue';

export function useMeetingSummary() {
    const isProcessing = ref(false);
    const processingStatus = ref('');
    const processingProgress = ref(0);
    const processingText = ref('');
    const showSummaryResult = ref(false);
    const summaryResult = ref('');

    const generateSummary = async (audioPath) => {
        isProcessing.value = true;
        processingProgress.value = 10;
        processingStatus.value = '正在轉錄音訊...';
        processingText.value = '使用 Whisper 模型處理中';
        
        const progressInterval = setInterval(() => {
            if (processingProgress.value < 60) {
                processingProgress.value += 5;
            }
        }, 500);
        
        const summaryPromise = window.go.main.App.GenerateMeetingSummary(audioPath);
        
        setTimeout(() => {
            processingStatus.value = '正在生成會議摘要...';
            processingText.value = '使用 AI 分析會議內容';
            processingProgress.value = 70;
        }, 3000);
        
        const summary = await summaryPromise;
        clearInterval(progressInterval);
        
        processingProgress.value = 100;
        processingStatus.value = '完成！';
        summaryResult.value = summary;
        showSummaryResult.value = true;
        
        isProcessing.value = false;
        processingProgress.value = 0;
    };

    const closeSummary = () => {
        showSummaryResult.value = false;
        summaryResult.value = '';
    };

    return {
        isProcessing,
        processingStatus,
        processingProgress,
        processingText,
        showSummaryResult,
        summaryResult,
        generateSummary,
        closeSummary
    };
}
