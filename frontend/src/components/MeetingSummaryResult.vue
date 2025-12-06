<template>
    <Teleport to="body">
        <div class="summary-overlay" @click="$emit('close')">
            <div class="summary-modal" @click.stop>
                <ModalHeader 
                    title="📝 會議摘要" 
                    @close="$emit('close')" 
                />
                
                <div class="summary-content">
                    <pre class="summary-text">{{ summary }}</pre>
                </div>
                
                <div class="summary-footer">
                    <button class="btn btn-primary" @click="copySummary">
                        📋 複製摘要
                    </button>
                    <button class="btn btn-secondary" @click="$emit('close')">
                        關閉
                    </button>
                </div>
            </div>
        </div>
    </Teleport>
</template>

<script>
import ModalHeader from './MeetingSummary/ModalHeader.vue';

export default {
    name: 'MeetingSummaryResult',
    components: {
        ModalHeader
    },
    props: {
        summary: {
            type: String,
            required: true
        }
    },
    emits: ['close'],
    setup(props) {
        const copySummary = () => {
            navigator.clipboard.writeText(props.summary)
                .then(() => alert('摘要已複製到剪貼簿！'))
                .catch(err => console.error('複製失敗:', err));
        };

        return {
            copySummary
        };
    }
};
</script>

<style scoped>
.summary-overlay {
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

.summary-modal {
    background: #ffffff;
    border-radius: 20px;
    box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.05), 0 10px 40px rgba(0, 0, 0, 0.15);
    width: 900px;
    max-width: 85vw;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
    animation: slideUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    overflow: hidden;
}

@keyframes slideUp {
    from {
        opacity: 0;
        transform: translateY(50px) scale(0.96);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

.summary-content {
    flex: 1;
    overflow-y: auto;
    padding: 32px;
}

.summary-text {
    font-family: 'Microsoft JhengHei', 'PingFang TC', 'Noto Sans TC', sans-serif;
    font-size: 15px;
    line-height: 1.7;
    color: #1f2937;
    white-space: pre-wrap;
    word-wrap: break-word;
    margin: 0;
}

.summary-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 20px 32px;
    border-top: 1px solid #e5e7eb;
}

.btn {
    padding: 10px 24px;
    border-radius: 10px;
    border: none;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}

.btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
}

.btn-primary:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-secondary {
    background: #f3f4f6;
    color: #6b7280;
}

.btn-secondary:hover {
    background: #e5e7eb;
    color: #1f2937;
}
</style>
