<template>
    <transition name="fade">
        <div v-if="message" class="error-toast">
            {{ message }}
        </div>
    </transition>
</template>

<script>
export default {
    name: 'ErrorToast',
    props: {
        message: String
    }
};
</script>

<style scoped>
.error-toast {
    position: fixed;
    top: 50%;
    right: 100px;
    transform: translateY(-50%);
    background: #fee2e2;
    color: #991b1b;
    padding: 12px 16px;
    border-radius: 8px;
    border: 1px solid #ef4444;
    font-size: 12px;
    font-weight: 600;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    max-width: 250px;
    word-break: break-word;
    animation: toast-slide-in 0.3s ease-out;
}

@keyframes toast-slide-in {
    from {
        transform: translateX(100%);
        opacity: 0;
    }
    to {
        transform: translateX(0);
        opacity: 1;
    }
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}
</style>
