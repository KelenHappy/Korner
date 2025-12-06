<template>
    <div class="progress-steps">
        <div 
            v-for="(step, index) in steps" 
            :key="index"
            class="step-wrapper"
        >
            <div 
                class="step" 
                :class="{ 
                    active: currentStep >= step.number, 
                    completed: currentStep > step.number 
                }"
            >
                <div class="step-icon">{{ step.icon }}</div>
                <div class="step-label">{{ step.label }}</div>
            </div>
            <div v-if="index < steps.length - 1" class="step-divider"></div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'StepIndicator',
    props: {
        currentStep: {
            type: Number,
            default: 1
        },
        steps: {
            type: Array,
            default: () => [
                { number: 1, icon: '🎤', label: '音訊轉錄' },
                { number: 2, icon: '🤖', label: 'AI 分析' },
                { number: 3, icon: '📝', label: '生成摘要' }
            ]
        }
    }
};
</script>

<style scoped>
.progress-steps {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0;
    width: 100%;
    max-width: 800px;
}

.step-wrapper {
    display: flex;
    align-items: center;
}

.step {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    opacity: 0.3;
    transition: all 0.3s ease;
}

.step.active {
    opacity: 1;
}

.step.completed {
    opacity: 0.6;
}

.step-icon {
    font-size: 32px;
    width: 56px;
    height: 56px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f3f4f6;
    border-radius: 50%;
    transition: all 0.3s ease;
}

.step.active .step-icon {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    transform: scale(1.05);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.step.completed .step-icon {
    background: #10b981;
}

.step-label {
    font-size: 14px;
    font-weight: 600;
    color: #6b7280;
}

.step.active .step-label {
    color: #1f2937;
}

.step-divider {
    width: 50px;
    height: 3px;
    background: #e5e7eb;
    border-radius: 2px;
    margin: 0 24px;
}
</style>
