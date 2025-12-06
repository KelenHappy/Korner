import { ref } from 'vue';

export function useRecording() {
    const isRecording = ref(false);
    const duration = ref(0);
    const savedFile = ref('');
    let durationInterval = null;

    const startRecording = async () => {
        savedFile.value = '';
        await window.go.main.App.StartRecording();
        isRecording.value = true;
        duration.value = 0;
        
        durationInterval = setInterval(async () => {
            duration.value = await window.go.main.App.GetRecordingDuration();
        }, 100);
    };

    const stopRecording = async () => {
        if (durationInterval) {
            clearInterval(durationInterval);
            durationInterval = null;
        }
        
        const filePath = await window.go.main.App.StopRecording();
        isRecording.value = false;
        savedFile.value = filePath;
    };

    const toggleRecording = () => {
        isRecording.value ? stopRecording() : startRecording();
    };

    const cleanup = () => {
        if (durationInterval) {
            clearInterval(durationInterval);
        }
        if (isRecording.value) {
            stopRecording();
        }
    };

    const formatDuration = (seconds) => {
        const mins = Math.floor(seconds / 60);
        const secs = Math.floor(seconds % 60);
        return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    };

    return {
        isRecording,
        duration,
        savedFile,
        startRecording,
        stopRecording,
        toggleRecording,
        formatDuration,
        cleanup
    };
}
