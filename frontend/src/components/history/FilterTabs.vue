<template>
    <div class="filter-tabs">
        <button
            v-for="tab in tabs"
            :key="tab.value"
            @click="$emit('change', tab.value)"
            :class="['tab-btn', { active: modelValue === tab.value }]"
        >
            {{ tab.label }}
        </button>
    </div>
</template>

<script>
export default {
    name: 'FilterTabs',
    props: {
        modelValue: {
            type: String,
            required: true
        },
        tabs: {
            type: Array,
            required: true
        }
    },
    emits: ['change']
};
</script>

<style scoped>
.filter-tabs {
    display: flex;
    gap: 8px;
    padding: 16px 24px;
    border-bottom: 1px solid #e2e8f0;
}

.tab-btn {
    padding: 8px 16px;
    background: transparent;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    font-size: 13px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}

.tab-btn:hover {
    background: #f1f5f9;
}

.tab-btn.active {
    background: #000;
    color: white;
    border-color: #000;
}
</style>
