<template>
  <div class="preview-table">
    <div class="table-header">
      <h3>预览 <span v-if="previewItems.length > 0">({{ previewItems.length }} 个文件)</span></h3>
      <div v-if="conflicts.length > 0" class="conflict-warning">
        ⚠️ 发现 {{ conflicts.length }} 个命名冲突
      </div>
    </div>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="previewItems.length === 0" class="empty">
      <p>选择目录并添加规则后，此处将显示重命名预览</p>
    </div>
    <div v-else class="table-container">
      <table>
        <thead>
          <tr>
            <th>原文件名</th>
            <th>新文件名</th>
            <th>状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in previewItems" :key="index" :class="{ conflict: item.conflict }">
            <td class="old-name">{{ item.originalFile.name }}{{ item.originalFile.ext }}</td>
            <td class="new-name">{{ item.newName }}</td>
            <td class="status">
              <span v-if="item.conflict" class="badge conflict">冲突</span>
              <span v-else-if="item.originalFile.name + item.originalFile.ext === item.newName" class="badge unchanged">不变</span>
              <span v-else class="badge changed">将重命名</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
defineProps(['files', 'previewItems', 'conflicts', 'loading'])
</script>

<style scoped>
.preview-table {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  overflow: hidden;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.table-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.conflict-warning {
  color: #d97706;
  font-size: 13px;
  font-weight: 500;
}

.loading, .empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748b;
  font-size: 14px;
}

.table-container {
  flex: 1;
  overflow: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  position: sticky;
  top: 0;
  background: #f8fafc;
  z-index: 10;
}

th {
  padding: 12px 20px;
  text-align: left;
  font-size: 13px;
  font-weight: 600;
  color: #475569;
  border-bottom: 2px solid #e2e8f0;
}

td {
  padding: 10px 20px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 14px;
}

tr:hover {
  background: #f8fafc;
}

tr.conflict {
  background: #fef3c7;
}

tr.conflict:hover {
  background: #fde68a;
}

.old-name {
  color: #64748b;
  text-decoration: line-through;
}

.new-name {
  color: #10b981;
  font-weight: 500;
}

.badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 100px;
  font-size: 12px;
  font-weight: 500;
}

.badge.conflict {
  background: #fee2e2;
  color: #dc2626;
}

.badge.unchanged {
  background: #f1f5f9;
  color: #64748b;
}

.badge.changed {
  background: #d1fae5;
  color: #059669;
}
</style>
