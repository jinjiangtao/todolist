<template>
  <div class="app">
    <header class="header">
      <h1>智能文件批量重命名工具</h1>
      <p class="subtitle">Smart Renamer</p>
    </header>

    <div class="main-content">
      <div class="left-panel">
        <DirectorySelector @directory-selected="onDirectorySelected" :directory="directory" />
        <RuleManager :rules="rules" @rules-changed="onRulesChanged" />
      </div>

      <div class="right-panel">
        <PreviewTable 
          :files="files" 
          :preview-items="previewItems" 
          :conflicts="conflicts"
          :loading="loading"
        />
        <ActionPanel 
          :can-execute="canExecute" 
          :has-undo="hasUndo"
          @execute="onExecute"
          @undo="onUndo"
        />
      </div>
    </div>

    <div v-if="message" :class="['message', message.type]">{{ message.text }}</div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import DirectorySelector from './components/DirectorySelector.vue'
import RuleManager from './components/RuleManager.vue'
import PreviewTable from './components/PreviewTable.vue'
import ActionPanel from './components/ActionPanel.vue'

const directory = ref('')
const files = ref([])
const rules = ref([])
const previewItems = ref([])
const conflicts = ref([])
const loading = ref(false)
const lastUndoId = ref('')
const message = ref(null)

const canExecute = computed(() => previewItems.value.length > 0 && !loading.value)
const hasUndo = computed(() => lastUndoId.value !== '')

function showMessage(text, type = 'info') {
  message.value = { text, type }
  setTimeout(() => {
    message.value = null
  }, 3000)
}

async function onDirectorySelected(dir) {
  directory.value = dir
  loading.value = true
  try {
    const req = JSON.stringify({ directory: dir })
    const resp = JSON.parse(await window.go.main.App.ScanDirectory(req))
    if (resp.success) {
      files.value = resp.files
      await generatePreview()
    } else {
      showMessage(resp.error || '扫描目录失败', 'error')
    }
  } catch (e) {
    showMessage('扫描目录失败: ' + e.message, 'error')
  } finally {
    loading.value = false
  }
}

function onRulesChanged(newRules) {
  rules.value = newRules
  generatePreview()
}

async function generatePreview() {
  if (files.value.length === 0 || rules.value.length === 0) {
    previewItems.value = []
    conflicts.value = []
    return
  }

  try {
    const req = JSON.stringify({ files: files.value, rules: rules.value })
    const resp = JSON.parse(await window.go.main.App.GeneratePreview(req))
    if (resp.success) {
      previewItems.value = resp.results
      conflicts.value = resp.conflicts
    } else {
      showMessage(resp.error || '生成预览失败', 'error')
    }
  } catch (e) {
    showMessage('生成预览失败: ' + e.message, 'error')
  }
}

async function onExecute() {
  if (!confirm('确定要重命名这些文件吗？')) return

  try {
    const req = JSON.stringify({ items: previewItems.value, conflictResolution: 'autoNum' })
    const resp = JSON.parse(await window.go.main.App.ExecuteRename(req))
    if (resp.success) {
      showMessage(`成功重命名 ${resp.renamed} 个文件，跳过 ${resp.skipped} 个`, 'success')
      lastUndoId.value = resp.undoId
      await onDirectorySelected(directory.value)
    } else {
      showMessage(resp.error || '重命名失败', 'error')
    }
  } catch (e) {
    showMessage('重命名失败: ' + e.message, 'error')
  }
}

async function onUndo() {
  if (!lastUndoId.value) return

  try {
    const req = JSON.stringify({ undoId: lastUndoId.value })
    const resp = JSON.parse(await window.go.main.App.UndoRename(req))
    if (resp.success) {
      showMessage(`成功恢复 ${resp.restored} 个文件`, 'success')
      lastUndoId.value = ''
      await onDirectorySelected(directory.value)
    } else {
      showMessage(resp.error || '撤销失败', 'error')
    }
  } catch (e) {
    showMessage('撤销失败: ' + e.message, 'error')
  }
}

onMounted(() => {
  rules.value = [
    { type: 'sequence', enabled: false, sequence: { start: 1, digits: 3, prefix: '', suffix: '' } },
    { type: 'replace', enabled: false, replace: { search: '', replace: '', useRegex: false } },
    { type: 'date', enabled: false, date: { source: 'modified', format: 'YYYYMMDD', position: 'prefix', separator: '_' } },
    { type: 'case', enabled: false, case: { mode: 'lower' } },
    { type: 'extension', enabled: false, extension: { newExt: '', unifyCase: false, targetCase: 'lower' } },
  ]
})
</script>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: #f8fafc;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px 30px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.header h1 {
  font-size: 24px;
  font-weight: 600;
}

.subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin-top: 4px;
}

.main-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.left-panel {
  width: 400px;
  background: white;
  border-right: 1px solid #e2e8f0;
  overflow-y: auto;
  padding: 20px;
}

.right-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 20px;
  gap: 16px;
}

.message {
  position: fixed;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  border-radius: 8px;
  color: white;
  font-weight: 500;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  z-index: 1000;
}

.message.success {
  background: #10b981;
}

.message.error {
  background: #ef4444;
}

.message.info {
  background: #3b82f6;
}
</style>
