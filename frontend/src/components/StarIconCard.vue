<template>
    <div class="star-icon-card">
        <div class="gradient-circle">
            <svg class="star-icon" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path
                    d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z"
                    fill="#FF6B9D"
                />
            </svg>
        </div>
    </div>
</template>

<script>
export default {
    name: 'StarIconCard',
};
</script>

<style scoped>
.star-icon-card {
    width: 170px;
    height: 130px;
    background: linear-gradient(145deg, #3a3d47, #2f3239);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 
        0 8px 16px rgba(0, 0, 0, 0.4),
        inset 0 1px 1px rgba(255, 255, 255, 0.05);
    position: relative;
}

.star-icon-card::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: 12px;
    padding: 1px;
    background: linear-gradient(145deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.02));
    -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
}

.gradient-circle {
    width: 70px;
    height: 70px;
    border-radius: 50%;
    background: linear-gradient(135deg, #7B68EE 0%, #6366F1 50%, #8B7AE8