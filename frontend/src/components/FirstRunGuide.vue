<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
    <div class="bg-white rounded-lg shadow-2xl max-w-2xl mx-4 p-6">
      <h2 class="text-2xl font-bold text-slate-800 mb-4">Welcome to SnapAsk!</h2>
      <div class="mb-6">
        <p class="text-slate-600 mb-4">On Linux, please manually set up the global hotkey:</p>
        <ol class="list-decimal list-inside space-y-2 text-slate-700">
          <li>Open your system settings (GNOME Settings or KDE System Settings)</li>
          <li>Navigate to Keyboard Shortcuts</li>
          <li>Add a new custom shortcut with command:</li>
        </ol>
        <div class="mt-3 p-3 bg-slate-100 rounded font-mono text-sm">
          flatpak run com.snapask.SnapAsk --screenshot
        </div>
        <p class="mt-3 text-sm text-slate-600">Recommended hotkey: Ctrl+Alt+Q</p>
      </div>
      <button @click="close" class="w-full px-4 py-2 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-medium">
        Got it!
      </button>
    </div>
  </div>
</template>
<script>
export default {
  name: 'FirstRunGuide',
  emits: ['close'],
  setup(props, { emit }) {
    const close = () => emit('close')
    return { close }
  }
}
</script>