<template>
    <div
        ref="petRef"
        class="pet-container"
        :class="{ dragging: isDragging }"
        @pointerdown="startDrag"
    >
        <div class="pet-character">
            <img src="/icon.png" alt="Korner Pet" class="pet-logo" />
            
            <button
                @click.stop="$emit('screenshot')"
                class="action-btn screenshot-btn"
                title="Take Screenshot"
            >
                📸
            </button>
        </div>

        <button
            @click.stop="$emit('toggle-menu')"
            class="menu-toggle"
            :class="{ active: menuActive }"
        >
            {{ menuActive ? "✕" : "☰" }}
        </button>
    </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue';

export default {
    name: 'PetAvatar',
    props: {
        menuActive: {
            type: Boolean,
            default: false
        }
    },
    emits: ['screenshot', 'toggle-menu', 'position-change'],
    setup(props, { emit }) {
        const petRef = ref(null);
        const isDragging = ref(false);
        const position = ref({ x: 20, y: 20 });
        
        let dragOffsetX = 0;
        let dragOffsetY = 0;

        const startDrag = (e) => {
            if (e.target.closest('button')) return;
            
            e.preventDefault();
            const el = petRef.value;
            if (!el) return;

            const style = window.getComputedStyle(el);
            const matrix = new DOMMatrix(style.transform);
            
            isDragging.value = true;
            dragOffsetX = e.clientX - matrix.m41;
            dragOffsetY = e.clientY - matrix.m42;

            el.setPointerCapture(e.pointerId);
            el.addEventListener('pointermove', onDrag);
            el.addEventListener('pointerup', stopDrag);
            el.addEventListener('pointercancel', stopDrag);
        };

        const onDrag = (e) => {
            if (!isDragging.value || !petRef.value) return;

            let newX = e.clientX - dragOffsetX;
            let newY = e.clientY - dragOffsetY;

            const petWidth = 80;
            const petHeight = 80;
            newX = Math.max(0, Math.min(newX, window.innerWidth - petWidth));
            newY = Math.max(0, Math.min(newY, window.innerHeight - petHeight));

            petRef.value.style.transform = `translate(${newX}px, ${newY}px)`;
            
            position.value = { x: newX, y: newY };
            emit('position-change', position.value);
        };

        const stopDrag = (e) => {
            isDragging.value = false;
            const el = petRef.value;
            if (el) {
                el.releasePointerCapture(e.pointerId);
                el.removeEventListener('pointermove', onDrag);
                el.removeEventListener('pointerup', stopDrag);
                el.removeEventListener('pointercancel', stopDrag);
            }
        };

        onMounted(() => {
            if (petRef.value) {
                petRef.value.style.transform = `translate(${position.value.x}px, ${position.value.y}px)`;
                emit('position-change', position.value);
            }
        });

        onUnmounted(() => {
            if (petRef.value) {
                petRef.value.removeEventListener('pointermove', onDrag);
                petRef.value.removeEventListener('pointerup', stopDrag);
            }
        });

        return {
            petRef,
            isDragging,
            startDrag
        };
    }
};
</script>

<style scoped>
.pet-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 80px;
    height: 80px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: grab;
    user-select: none;
    pointer-events: auto;
    will-change: transform;
    touch-action: none;
}

.pet-container.dragging {
    cursor: grabbing;
}

.pet-container.dragging .pet-character {
    animation: none;
    transition: none;
}

.pet-character {
    position: relative;
    width: 60px;
    height: 60px;
    background: transparent;
    border-radius: 50%;
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: transform 0.3s ease, box-shadow 0.3s ease;
    animation: float 3s ease-in-out infinite;
    overflow: hidden;
}

.pet-character:hover {
    transform: scale(1.1);
    box-shadow: 0 6px 20px rgba(102, 126, 234, 0.5);
}

.pet-logo {
    width: 100%;
    height: 100%;
    object-fit: cover;
    pointer-events: none;
}

@keyframes float {
    0%, 100% {
        transform: translateY(0px);
    }
    50% {
        transform: translateY(-10px);
    }
}

.action-btn {
    position: absolute;
    bottom: -3px;
    right: -3px;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #000000;
    border: 2px solid white;
    box-shadow: 0 3px 10px rgba(0, 0, 0, 0.25);
    font-size: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    pointer-events: auto;
    z-index: 10;
}

.action-btn:hover {
    transform: scale(1.15) rotate(15deg);
    box-shadow: 0 4px 15px rgba(102, 126, 234, 0.6);
}

.action-btn:active {
    transform: scale(0.95);
}

.menu-toggle {
    position: absolute;
    bottom: -3px;
    left: -3px;
    width: 28px;
    height: 28px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: 2px solid rgba(102, 126, 234, 0.3);
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    pointer-events: auto;
    z-index: 10;
}

.menu-toggle:hover {
    background: #000000;
    color: white;
    transform: scale(1.1) rotate(90deg);
    border-color: transparent;
}

.menu-toggle.active {
    background: #000000;
    color: white;
    transform: rotate(180deg);
}
</style>
