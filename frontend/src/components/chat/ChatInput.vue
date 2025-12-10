<template>
    <div class="input-area">
        <textarea
            v-model="inputText"
            @keydown.ctrl.enter="handleSubmit"
            @keydown.meta.enter="handleSubmit"
            class="chat-input"
            rows="3"
            :placeholder="placeholder"
            :disabled="disabled"
        ></textarea>
        <div class="input-actions">
            <div class="left-actions">
                <!-- 只在沒有截圖時顯示文件上傳和聯網搜尋 -->
                <template v-if="!hasScreenshot">
                    <button
                        @click="handleFileSelect"
                        :disabled="disabled"
                        class="file-btn"
                        :class="{ 'has-file': selectedFiles.length > 0 }"
                        :title="selectedFiles.length > 0 ? `已選擇 ${selectedFiles.length} 個檔案` : '上傳文件 (支援 TXT, MD, PDF, JSON, CSV)'"
                    >
                        📄
                        <span v-if="selectedFiles.length > 0" class="file-count">{{ selectedFiles.length }}</span>
                    </button>
                    <button
                        @click="webSearchEnabled = !webSearchEnabled"
                        :disabled="disabled"
                        class="web-search-btn"
                        :class="{ 'active': webSearchEnabled }"
                        :title="webSearchEnabled ? '關閉旅遊搜尋' : '開啟旅遊搜尋'"
                    >
                        🚆
                    </button>

                    <div v-if="selectedFiles.length > 0" class="file-list">
                        <div 
                            v-for="(file, index) in selectedFiles" 
                            :key="index"
                            class="file-tag"
                            :title="file.name"
                        >
                            <span>{{ file.name }}</span>
                            <button @click="removeFile(index)" class="remove-file">✕</button>
                        </div>
                    </div>
                </template>
            </div>
            <div class="right-actions">
                <span class="char-count">{{ inputText.length }} / 1000 {{ charCountLabel }}</span>
                <button
                    @click="handleSubmit"
                    :disabled="!inputText.trim() || disabled"
                    class="send-btn"
                >
                    <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
                        <path d="M2 10L18 2L10 18L8 11L2 10Z" fill="currentColor" />
                    </svg>
                </button>
            </div>
        </div>

    </div>
</template>

<script>
import { ref, watch } from 'vue';

export default {
    name: 'ChatInput',
    props: {
        modelValue: {
            type: String,
            default: ''
        },
        placeholder: {
            type: String,
            default: 'Type your message...'
        },
        disabled: {
            type: Boolean,
            default: false
        },
        charCountLabel: {
            type: String,
            default: 'characters'
        },
        screenshot: {
            type: String,
            default: ''
        },
        hasScreenshot: {
            type: Boolean,
            default: false
        }
    },
    emits: ['update:modelValue', 'submit', 'extract-text'],
    setup(props, { emit }) {
        const inputText = ref(props.modelValue);
        const selectedFiles = ref([]);
        const webSearchEnabled = ref(false);

        watch(() => props.modelValue, (newValue) => {
            inputText.value = newValue;
        });

        watch(inputText, (newValue) => {
            emit('update:modelValue', newValue);
        });

        const handleSubmit = () => {
            const trimmed = inputText.value.trim();
            if (!trimmed || props.disabled) return;
            
            // 保存原始用戶輸入（用於顯示）
            const userInput = trimmed;
            
            // 合併用戶輸入和所有檔案內容（用於發送給 AI）
            let finalText = trimmed;
            
            if (selectedFiles.value.length > 0) {
                finalText += '\n\n--- 檔案內容 ---\n';
                selectedFiles.value.forEach((file, index) => {
                    if (index > 0) finalText += '\n\n';
                    finalText += `[檔案 ${index + 1}: ${file.name}]\n${file.content}`;
                });
            }
            
            // 提交數據：text 是完整內容（給 AI），userInput 是用戶輸入（用於顯示）
            const submitData = {
                text: finalText.slice(0, 10000),
                userInput: userInput,
                webSearch: webSearchEnabled.value
            };
            
            emit('submit', submitData);
            inputText.value = '';
            selectedFiles.value = [];
            webSearchEnabled.value = false;
        };

        const handleFileSelect = async () => {
            if (props.disabled) return;
            
            try {
                // 使用後端文件選擇對話框（支援所有文件類型包括 PDF）
                const filePaths = await window.go.main.App.SelectDocumentFiles();
                
                if (filePaths.length === 0) return;
                
                for (const filePath of filePaths) {
                    const fileName = filePath.split(/[/\\]/).pop();
                    
                    try {
                        const content = await window.go.main.App.ReadDocumentFile(filePath);
                        selectedFiles.value.push({
                            name: fileName,
                            content: content.slice(0, 3000)
                        });
                    } catch (error) {
                        console.error(`讀取檔案 ${fileName} 失敗:`, error);
                        alert(`無法讀取檔案 ${fileName}: ${error}`);
                    }
                }
            } catch (error) {
                console.error('選擇文件失敗:', error);
            }
        };

        const removeFile = (index) => {
            selectedFiles.value.splice(index, 1);
        };

        return {
            inputText,
            selectedFiles,
            webSearchEnabled,
            handleSubmit,
            handleFileSelect,
            removeFile
        };
    }
};
</script>

<style scoped>
.input-area {
    border-top: 1px solid #cccccc;
    padding: 16px 24px;
    background: #ffffff;
}

.chat-input {
    width: 100%;
    padding: 12px 16px;
    border: 1.5px solid rgba(0, 0, 0, 0.08);
    border-radius: 12px;
    font-size: 14px;
    line-height: 1.6;
    color: #1a1a1a;
    resize: none;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    box-sizing: border-box;
    background: white;
    margin-bottom: 12px;
}

.chat-input:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 1px #667eea, 0 0 0 4px rgba(102, 126, 234, 0.1);
}

.chat-input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.chat-input::placeholder {
    color: #9ca3af;
}

.input-actions {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.left-actions {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    min-width: 0;
    overflow: hidden;
}

.right-actions {
    display: flex;
    align-items: center;
    gap: 12px;
}

.file-btn {
    width: 36px;
    height: 36px;
    border: 1.5px solid rgba(0, 0, 0, 0.08);
    background: white;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    transition: all 0.2s;
    position: relative;
}

.file-btn:hover:not(:disabled) {
    border-color: #667eea;
    background: rgba(102, 126, 234, 0.05);
}

.file-btn.has-file {
    border-color: #667eea;
    background: rgba(102, 126, 234, 0.1);
}

.file-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.web-search-btn {
    width: 36px;
    height: 36px;
    border: 1.5px solid rgba(0, 0, 0, 0.08);
    background: white;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    transition: all 0.2s;
}

.web-search-btn:hover:not(:disabled) {
    border-color: #3b82f6;
    background: rgba(59, 130, 246, 0.05);
}

.web-search-btn.active {
    border-color: #3b82f6;
    background: rgba(59, 130, 246, 0.15);
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.web-search-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.screenshot-btn {
    width: 36px;
    height: 36px;
    border: 1.5px solid rgba(0, 0, 0, 0.08);
    background: white;
    border-radius: 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 18px;
    transition: all 0.2s;
}

.screenshot-btn:hover:not(:disabled) {
    border-color: #10b981;
    background: rgba(16, 185, 129, 0.05);
}

.screenshot-btn.has-screenshot {
    border-color: #10b981;
    background: rgba(16, 185, 129, 0.1);
}

.screenshot-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.file-count {
    position: absolute;
    top: -6px;
    right: -6px;
    background: #667eea;
    color: white;
    font-size: 10px;
    font-weight: 700;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid white;
}

.file-list {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    flex: 1;
    min-width: 0;
    overflow: hidden;
}

.file-tag {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 4px 8px;
    background: rgba(102, 126, 234, 0.1);
    border: 1px solid rgba(102, 126, 234, 0.2);
    border-radius: 6px;
    font-size: 11px;
    color: #667eea;
    font-weight: 500;
    max-width: 150px;
    min-width: 0;
    flex-shrink: 1;
}

.file-tag span {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
}

.remove-file {
    background: none;
    border: none;
    color: #667eea;
    cursor: pointer;
    padding: 0;
    font-size: 12px;
    line-height: 1;
    opacity: 0.6;
    transition: opacity 0.2s;
}

.remove-file:hover {
    opacity: 1;
}

.char-count {
    font-size: 12px;
    color: #9ca3af;
    font-weight: 500;
}

.send-btn {
    width: 40px;
    height: 40px;
    border: none;
    background: #667eea;
    color: white;
    border-radius: 50%;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    box-shadow: 0 0 0 1px rgba(102, 126, 234, 0.1), 0 4px 12px rgba(102, 126, 234, 0.3);
}

.send-btn:hover:not(:disabled) {
    transform: scale(1.1);
    box-shadow: 0 0 0 1px rgba(102, 126, 234, 0.2), 0 6px 16px rgba(102, 126, 234, 0.4);
}

.send-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: scale(1);
}
</style>
