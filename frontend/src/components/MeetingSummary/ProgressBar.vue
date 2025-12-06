<template>
    <div class="progress-bar-container">
        <div class="progress-bar">
            <div class="progress-fill" :style="{ width: progress + '%' }"></div>
        </div>
        <div class="progress-percentage">{{ progress }}%</div>
    </div>
</template>

<script>
export default {
    name: 'ProgressBar',
    props: {
        progress: {
            type: Number,
            default: 0
        }
    }
};
</script>

<style scoped>
.progress-bar-container {
    width: 100%;
    margin-bottom: 20px;
}

.progress-bar {
    width: 100%;
    height: 20px;
    background: #e5e7eb;
    border-radius: 10px;
    overflow: hidden;
    box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 12px;
}

.progress-fill {
    height: 100%;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    border-radius: 10px;
    transition: width 0.5s ease;
}

.progress-percentage {
    font-size: 20px;
    font-weight: 700;
    color: #667eea;
    text-align: center;
}
</style>
