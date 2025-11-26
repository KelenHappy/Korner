<template>
    <transition name="pie-container-fade">
        <div v-if="visible" class="pie-menu-container" @click.self="$emit('hide')">
            
            <div 
                class="pie-menu-content"
                :style="{ top: clickY + 'px', left: clickX + 'px' }"
            >
                <transition-group name="pie-pop" tag="div" class="pie-items">
                    <div
                        key="item0"
                        class="pie-item photo-theme"
                        :style="{ ...getPieItemStyle(0), '--delay': '0.05s' }"
                        @click.stop="$emit('screenshot')"
                        @mouseenter="activeItem = 0"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">📸</span>
                         <span class="pie-label" v-show="activeItem === 0">拍照</span>
                    </div>

                    <div
                        key="item1"
                        class="pie-item talk-theme"
                        :style="{ ...getPieItemStyle(1), '--delay': '0.1s' }"
                        @click.stop="$emit('ask-question')"
                        @mouseenter="activeItem = 1"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">💬</span>
                        <span class="pie-label" v-show="activeItem === 1">說話</span>
                    </div>

                    <div
                        key="item2"
                        class="pie-item settings-theme"
                        :style="{ ...getPieItemStyle(2), '--delay': '0.15s' }"
                        @click.stop="$emit('settings')"
                        @mouseenter="activeItem = 2"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">⚙️</span>
                        <span class="pie-label" v-show="activeItem === 2">設定</span>
                    </div>

                    <div
                        key="item3"
                        class="pie-item hide-theme"
                        :style="{ ...getPieItemStyle(3), '--delay': '0.2s' }"
                        @click.stop="$emit('hide-pet')"
                        @mouseenter="activeItem = 3"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">💤</span>
                        <span class="pie-label" v-show="activeItem === 3">休息</span>
                    </div>
                </transition-group>
            </div>
        </div>
    </transition>
</template>

<script>
import { ref } from "vue";

export default {
    name: "DesktopPetPieMenu",
    props: {
        visible: {
            type: Boolean,
            default: false,
        },
        // 改用點擊座標
        clickX: {
            type: Number,
            default: 0,
        },
        clickY: {
            type: Number,
            default: 0,
        },
    },
    emits: ["screenshot", "ask-question", "settings", "hide", "hide-pet"],
    setup() {
        const activeItem = ref(null);
        // 半徑可以稍微加大一點，讓卡通圖標不那麼擁擠
        const radius = 40; 

        const getPieItemStyle = (index) => {
            const angle = (index * 360) / 4 - 90;
            const radians = (angle * Math.PI) / 180;
            // 這裡計算的是「最終位置」的偏移量
            const x = Math.cos(radians) * radius;
            const y = Math.sin(radians) * radius;

            // 利用 CSS 變數存儲最終位置，供動畫使用
            return {
                '--end-x': `${x}px`,
                '--end-y': `${y}px`,
            };
        };

        return {
            activeItem,
            getPieItemStyle,
        };
    },
};
</script>

<style scoped>
/* 容器淡入淡出，不影響內部彈跳 */
.pie-menu-container {
    position: fixed;
    top: 0; left: 0; width: 100vw; height: 100vh;
    z-index: 9999;
    background: transparent; /* 桌寵通常不需要背景遮罩 */
}
.pie-container-fade-enter-active, .pie-container-fade-leave-active { transition: opacity 0.2s; }
.pie-container-fade-enter-from, .pie-container-fade-leave-to { opacity: 0; }


/* 核心內容定位點，寬高為 0，確保是絕對中心 */
.pie-menu-content {
    position: absolute;
    width: 0; height: 0;
    /* 這裡不需要 transform translate，因為 top/left 已經是精確點擊位置 */
    pointer-events: none; /* 讓點擊穿透到 items */
}

.pie-items {
    position: absolute;
    /* 讓 items 的中心點對齊 content 的中心點 */
    top: 0; left: 0;
    width: 0; height: 0;
}

/* --- 卡通化按鈕樣式 --- */
.pie-item {
    position: absolute;
    /* 讓按鈕中心對齊定位點 */
    left: -20px; top: -20px; 
    width: 40px; height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center; justify-content: center;
    cursor: pointer;
    pointer-events: auto;
    user-select: none;

    /* 桌寵風格：粗邊框、鮮明陰影 */
    background: #fff;
    border: 3px solid #4a4a4a;
    box-shadow: 2px 4px 0px rgba(0,0,0,0.3); 
    
    /* 這是最終靜止狀態的位置，從 CSS 變數讀取 */
    transform: translate(var(--end-x), var(--end-y)) scale(1);
    transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1), background 0.2s, box-shadow 0.2s; /* 懸停時的彈性 */
}

.pie-icon { font-size: 18px; }

/* Hover 效果：更 Q 彈 */
.pie-item:hover {
    transform: translate(var(--end-x), var(--end-y)) scale(1.2) rotate(5deg);
    box-shadow: 4px 6px 0px rgba(0,0,0,0.4);
    z-index: 10;
}

/* --- 不同功能的配色主題 (可選) --- */
.photo-theme:hover { border-color: #ff6b6b; color: #ff6b6b; }
.talk-theme:hover { border-color: #51cf66; color: #51cf66; }
.settings-theme:hover { border-color: #339af0; color: #339af0; }
.hide-theme:hover { border-color: #fcc419; color: #fcc419; }


/* --- Hover 文字標籤 --- */
.pie-label {
    position: absolute;
    bottom: -25px; /* 顯示在圓圈下方 */
    white-space: nowrap;
    background: rgba(50, 50, 50, 0.9);
    color: #fff;
    padding: 3px 8px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: bold;
    pointer-events: none;
    box-shadow: 1px 2px 4px rgba(0,0,0,0.2);
}

/* --- 核心動畫：Q 彈噴射 (Pop Animation) --- */
/* 進場前狀態：在中心點，縮小為 0 */
.pie-pop-enter-from {
    transform: translate(0px, 0px) scale(0.1) !important; /* 強制覆蓋原本的 translate */
    opacity: 0;
}

/* 進場動畫過程 */
.pie-pop-enter-active {
    /* 使用貝茲曲線製造「衝過頭再拉回」的彈性效果 */
    transition: transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease-out;
    /* 應用 JavaScript 計算出的延遲 */
    transition-delay: var(--delay);
}

/* 離場狀態 (可選：讓它們縮回中心，或者直接淡出) */
.pie-pop-leave-to {
     transform: translate(0px, 0px) scale(0.1) !important;
     opacity: 0;
     transition: transform 0.3s ease-in, opacity 0.2s ease-in;
}
/* 離場時取消延遲，一起消失比較俐落 */
.pie-pop-leave-active { transition-delay: 0s !important; }

</style>