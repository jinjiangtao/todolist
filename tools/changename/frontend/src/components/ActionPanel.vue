<template>
  <div class="action-panel">
    <button 
      v-if="hasUndo" 
      @click="$emit('undo')" 
      class="btn btn-undo"
    >
      ↶ 撤销上次操作
    </button>
    <button 
      :disabled="!canExecute" 
      @click="$emit('execute')" 
      class="btn btn-execute"
    >
      开始重命名
    </button>
  </div>
</template>

<script setup>
defineProps(['canExecute', 'hasUndo'])
defineEmits(['execute', 'undo'])
</script>

<style scoped>
.action-panel {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 4px;
}

.btn {
  padding: 12px 28px;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 600;
  transition: all 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-undo {
  background: #f1f5f9;
  color: #475569;
}

.btn-undo:hover:not(:disabled) {
  background: #e2e8f0;
}

.btn-execute {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-execute:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}
</style>
