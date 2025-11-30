<template>
    <div
        v-if="visible"
        class="selection-box"
        :style="{
            left: rect.x + 'px',
            top: rect.y + 'px',
            width: rect.width + 'px',
            height: rect.height + 'px',
        }"
    >
        <!-- Corner indicators -->
        <div class="corner corner-tl"></div>
        <div class="corner corner-tr"></div>
        <div class="corner corner-bl"></div>
        <div class="corner corner-br"></div>

        <!-- Dimension display -->
        <div class="dimension-label">
            {{ Math.abs(rect.width) }} × {{ Math.abs(rect.height) }}
        </div>
    </div>
</template>

<script>
export default {
    name: 'SelectionBox',
    props: {
        visible: {
            type: Boolean,
            default: false
        },
        rect: {
            type: Object,
            default: () => ({ x: 0, y: 0, width: 0, height: 0 })
        }
    }
};
</script>

<style scoped>
.selection-box {
    position: absolute;
    border: 2px solid rgb(59, 130, 246);
    background-color: rgba(59, 130, 246, 0.1);
    pointer-events: none;
}

.corner {
    position: absolute;
    width: 8px;
    height: 8px;
    background-color: rgb(59, 130, 246);
    border-radius: 50%;
}

.corner-tl {
    top: -4px;
    left: -4px;
}

.corner-tr {
    top: -4px;
    right: -4px;
}

.corner-bl {
    bottom: -4px;
    left: -4px;
}

.corner-br {
    bottom: -4px;
    right: -4px;
}

.dimension-label {
    position: absolute;
    top: -32px;
    left: 0;
    background-color: rgb(59, 130, 246);
    color: white;
    font-size: 12px;
    padding: 4px 8px;
    border-radius: 4px;
    white-space: nowrap;
}
</style>
