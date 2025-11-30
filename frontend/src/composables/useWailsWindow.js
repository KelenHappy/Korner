// Wails 視窗操作的簡單封裝
import {
    WindowSetAlwaysOnTop,
    WindowFullscreen,
    WindowUnfullscreen,
    WindowSetSize,
    WindowSetPosition,
    WindowGetPosition,
    WindowGetSize,
    WindowCenter,
    EventsOn,
    EventsOff,
} from '../wailsjs/runtime/runtime';

export {
    WindowSetAlwaysOnTop,
    WindowFullscreen,
    WindowUnfullscreen,
    WindowSetSize,
    WindowSetPosition,
    WindowGetPosition,
    WindowGetSize,
    WindowCenter,
    EventsOn,
    EventsOff,
};

// 延遲函數
export const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

// 安全執行視窗操作
export const safeWindowOp = async (fn) => {
    try {
        return await fn();
    } catch (error) {
        console.log('Window operation failed:', error);
        return null;
    }
};
