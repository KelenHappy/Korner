<template>
    <div class="modal-header">
        <h2 class="modal-title">{{ title }}</h2>
        <button v-if="showClose" class="close-btn" @click="$emit('close')">✕</button>
    </div>
</template>

<script>
export default {
    name: 'ModalHeader',
    props: {
        title: {
            type: String,
            required: true
        },
        showClose: {
            type: Boolean,
            default: true
        }
    },
    emits: ['close']
};
</script>

<style scoped>
.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 24px 32px;
    border-bottom: 1px solid #e5e7eb;
}

.modal-title {
    font-size: 24px;
    font-weight: 700;
    color: #1f2937;
    margin: 0;
}

.close-btn {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: none;
    background: #f3f4f6;
    color: #6b7280;
    font-size: 20px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
}

.close-btn:hover {
    background: #e5e7eb;
    color: #1f2937;
}
</style>
