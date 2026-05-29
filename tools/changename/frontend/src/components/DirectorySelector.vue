<template>
  <div class="directory-selector">
    <h3>选择目录</h3>
    <div class="input-group">
      <input 
        type="text" 
        :value="directory" 
        placeholder="拖放文件夹或点击选择..."
        readonly
        @click="selectDirectory"
        @dragover.prevent
        @drop="onDrop"
      />
      <button @click="selectDirectory" class="btn-select">浏览</button>
    </div>
    <p class="hint">提示：可以直接拖放文件夹到输入框</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['directory-selected'])
const props = defineProps(['directory'])

function selectDirectory() {
  const path = window.prompt('请输入目录路径：', props.directory || '')
  if (path && path.trim()) {
    emit('directory-selected', path.trim())
  }
}

function onDrop(e) {
  e.preventDefault()
  const files = Array.from(e.dataTransfer.files)
  if (files.length > 0) {
    const file = files[0]
    if (file.type === '' || file.name.indexOf('.') === -1) {
      emit('directory-selected', file.path)
    } else {
      const path = file.path.substring(0, file.path.lastIndexOf('\\'))
      emit('directory-selected', path)
    }
  }
}
</script>

<style scoped>
.directory-selector {
  margin-bottom: 24px;
}

.directory-selector h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 12px;
}

.input-group {
  display: flex;
  gap: 8px;
}

.input-group input {
  flex: 1;
  padding: 10px 14px;
  border: 2px solid #e2e8f0;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s;
  cursor: pointer;
}

.input-group input:focus {
  border-color: #667eea;
  outline: none;
}

.btn-select {
  padding: 10px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 8px;
  font-weight: 500;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-select:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.hint {
  font-size: 12px;
  color: #64748b;
  margin-top: 8px;
}
</style>
