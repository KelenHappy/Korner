<template>
    <Teleport to="body">
        <div class="progress-overlay">
            <div class="progress-modal">
                <div class="progress-content">
                    <div class="progress-title">{{ status }}</div>
                    
                    <ProgressBar :progress="progress" />
                    
                    <div class="progress-text">{{ message }}</div>
                    
                    <StepIndicator :current-step="currentStep" />
                </div>
            </div>
        </div>
    </Teleport>
</template>

<script>
import { computed } from 'vue';
import ProgressBar from './MeetingSummary/ProgressBar.vue';
import StepIndicator from './MeetingSummary/StepIndicator.vue';

export default {
    name: 'MeetingSummaryProgress',
    components: {
        ProgressBar,
        StepIndicator
    },
    props: {
        status: {
            type: String,
            default: '處理中...'
        },
        progress: {
            type: Number,
            default: 0
        },
        message: {
            type: String,
            default: ''
        }
    },
    setup(props) {
        const currentStep = computed(() => {
            if (props.progress < 30) return 1;
            if (props.progress < 70) return 2;
            return 3;
        });

        return {
            currentStep
        };
    }
};
</script>

<style scoped>
.progress-overlay {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 10000;
    animation: fadeIn 0.2s ease-out;
    pointer-events: auto;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        backdrop-filter: blur(0px);
    }
    to {
        opacity: 1;
        backdrop-filter: blur(8px);
    }
}

.progress-modal {
    background: #ffffff;
    border-radius: 20px;
    box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.05), 0 10px 40px rgba(0, 0, 0, 0.15);
    width: 700px;
    max-width: 85vw;
    display: flex;
    flex-direction: column;
    animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    overflow: hidden;
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(30px) scale(0.96);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

.progress-content {
    padding: 40px 50px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
}

.progress-title {
    font-size: 28px;
    font-weight: 700;
    color: #1f2937;
    margin-bottom: 24px;
    text-align: center;
}

.progress-text {
    font-size: 16px;
    color: #6b7280;
    text-align: center;
    font-weight: 500;
    line-height: 1.5;
    margin-bottom: 30px;
}
</style>
