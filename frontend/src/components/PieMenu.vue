<template>
    <transition name="pie-container-fade">
        <div v-if="visible" class="pie-menu-container" @click.self="$emit('hide')">
            <div 
                class="pie-menu-content"
                :class="{ 'is-closing': isClosing }"
                :style="contentStyle"
            >
                <transition-group name="pie-pop" tag="div" class="pie-items">
                    <div
                        v-if="showItems"
                        key="item3"
                        class="pie-item hide-theme"
                        :style="{ '--delay': '0.05s' }"
                        @click.stop="$emit('hide-pet')"
                        @mouseenter="activeItem = 3"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">💤</span>
                        <span class="pie-label" v-show="activeItem === 3">{{ t('menu.minimize') }}</span>
                    </div>

                    <div
                        v-if="showItems"
                        key="item0"
                        class="pie-item photo-theme"
                        :style="{ '--delay': '0.1s' }"
                        @click.stop="$emit('screenshot')"
                        @mouseenter="activeItem = 0"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">📸</span>
                        <span class="pie-label" v-show="activeItem === 0">{{ t('menu.screenshot') }}</span>
                    </div>

                    <div
                        v-if="showItems"
                        key="item1"
                        class="pie-item talk-theme"
                        :style="{ '--delay': '0.15s' }"
                        @click.stop="$emit('ask-question')"
                        @mouseenter="activeItem = 1"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">💬</span>
                        <span class="pie-label" v-show="activeItem === 1">{{ t('menu.askQuestion') }}</span>
                    </div>

                    <div
                        v-if="showItems"
                        key="item2"
                        class="pie-item settings-theme"
                        :style="{ '--delay': '0.2s' }"
                        @click.stop="$emit('settings')"
                        @mouseenter="activeItem = 2"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">⚙️</span>
                        <span class="pie-label" v-show="activeItem === 2">{{ t('menu.settings') }}</span>
                    </div>

                    <div
                        v-if="showItems"
                        key="item4"
                        class="pie-item history-theme"
                        :style="{ '--delay': '0.25s' }"
                        @click.stop="$emit('history')"
                        @mouseenter="activeItem = 4"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">📜</span>
                        <span class="pie-label" v-show="activeItem === 4">{{ t('menu.history') }}</span>
                    </div>

                    <div
                        v-if="showItems"
                        key="item5"
                        class="pie-item voice-theme"
                        :style="{ '--delay': '0.3s' }"
                        @click.stop="$emit('voice-meeting')"
                        @mouseenter="activeItem = 5"
                        @mouseleave="activeItem = null"
                    >
                        <span class="pie-icon">🎤</span>
                        <span class="pie-label" v-show="activeItem === 5">{{ t('menu.voiceMeeting') }}</span>
                    </div>
                </transition-group>
            </div>
        </div>
    </transition>
</template>

<script>
import { ref, computed, watch } from "vue";
import { useI18n } from "vue-i18n";

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
        // 桌寵位置（關閉時移動到這裡）
        petX: {
            type: Number,
            default: 0,
        },
        petY: {
            type: Number,
            default: 0,
        },
    },
    emits: ["screenshot", "ask-question", "settings", "history", "hide", "hide-pet", "voice-meeting"],
    setup(props, { emit }) {
        const { t } = useI18n();
        const activeItem = ref(null);
        const isClosing = ref(false);
        const showItems = ref(false);

        // 計算內容位置：關閉時移動到桌寵位置
        const contentStyle = computed(() => {
            if (isClosing.value) {
                return {
                    top: props.petY + 'px',
                    left: props.petX + 'px',
                    transition: 'top 0.3s ease-in, left 0.3s ease-in'
                };
            }
            return {
                top: props.clickY + 'px',
                left: props.clickX + 'px'
            };
        });

        // 監聽 visible 變化
        watch(() => props.visible, (newVal, oldVal) => {
            if (oldVal && !newVal) {
                // 關閉：先隱藏 items，等縮小動畫完成後再移動
                showItems.value = false;
                setTimeout(() => {
                    isClosing.value = true;
                }, 100);
            } else if (newVal) {
                // 打開：先重置狀態，等位置設定好後再顯示 items
                isClosing.value = false;
                showItems.value = false;
                // 等待下一幀讓位置更新，然後顯示 items
                setTimeout(() => {
                    showItems.value = true;
                }, 50);
            }
        });

        const handleSettingsClick = () => {
            console.log("[Korner][PieMenu] Settings clicked");
            emit('settings');
        };

        return {
            t,
            activeItem,
            isClosing,
            showItems,
            contentStyle,
            handleSettingsClick,
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


/* 核心內容定位點 */
.pie-menu-content {
    position: absolute;
    /* 水平居中，垂直從頂部開始 */
    transform: translateX(-50%);
    pointer-events: none; /* 讓點擊穿透到 items */
}

/* 關閉時的過渡效果 */
.pie-menu-content.is-closing {
    transition: top 0.3s ease-in, left 0.3s ease-in;
}

.pie-items {
    display: flex;
    flex-direction: column;
    gap: 8px;
    align-items: center;
    justify-content: center;
    flex-wrap: nowrap;
}

/* --- 卡通化按鈕樣式 --- */
.pie-item {
    position: relative;
    width: 32px; height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center; justify-content: center;
    cursor: pointer;
    pointer-events: auto;
    user-select: none;
    flex-shrink: 0;

    /* 桌寵風格：粗邊框、鮮明陰影 */
    background: #fff;
    border: 2px solid #4a4a4a;
    box-shadow: 1px 3px 0px rgba(0,0,0,0.3); 
    
    transform: scale(1);
    transition: transform 0.2s cubic-bezier(0.34, 1.56, 0.64, 1), background 0.2s, box-shadow 0.2s; /* 懸停時的彈性 */
}

.pie-icon { font-size: 14px; }

/* Hover 效果：更 Q 彈 */
.pie-item:hover {
    transform: scale(1.15) rotate(5deg);
    box-shadow: 3px 5px 0px rgba(0,0,0,0.4);
    z-index: 10;
}

.photo-theme:hover { border-color: #ff6b6b; color: #ff6b6b; }
.talk-theme:hover { border-color: #51cf66; color: #51cf66; }
.settings-theme:hover { border-color: #339af0; color: #339af0; }
.history-theme:hover { border-color: #9775fa; color: #9775fa; }
.hide-theme:hover { border-color: #fcc419; color: #fcc419; }
.voice-theme:hover { border-color: #20c997; color: #20c997; }


.pie-label {
    position: absolute;
    bottom: -24px; 
    left: 50%;
    transform: translateX(-50%);
    white-space: nowrap;
    background: rgba(50, 50, 50, 0.9);
    color: #fff;
    padding: 2px 6px;
    border-radius: 10px;
    font-size: 11px;
    font-weight: bold;
    pointer-events: none;
    box-shadow: 1px 2px 4px rgba(0,0,0,0.2);
}

.pie-pop-enter-from {
    transform: scale(0.1) !important;
    opacity: 0;
}

.pie-pop-enter-active {
    
    transition: transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease-out;
    transition-delay: var(--delay);
}

.pie-pop-leave-to {
     transform: scale(0.1) !important;
     opacity: 0;
     transition: transform 0.15s ease-in, opacity 0.1s ease-in;
}

.pie-pop-leave-active { 
    transition-delay: 0s !important;
}

</style>