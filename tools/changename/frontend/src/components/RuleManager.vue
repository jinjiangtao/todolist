<template>
  <div class="rule-manager">
    <h3>重命名规则</h3>
    <div class="rules-list">
      <div v-for="(rule, index) in localRules" :key="index" :class="['rule-item', { active: rule.enabled }]">
        <div class="rule-header" @click="toggleRule(index)">
          <input 
            type="checkbox" 
            :checked="rule.enabled" 
            @click.stop
            @change="toggleRule(index)"
          />
          <span class="rule-name">{{ getRuleName(rule.type) }}</span>
        </div>
        <div v-if="rule.enabled" class="rule-content">
          <template v-if="rule.type === 'sequence'">
            <div class="field">
              <label>起始数字</label>
              <input type="number" v-model.number="rule.sequence.start" min="0" @input="notifyChange" />
            </div>
            <div class="field">
              <label>位数</label>
              <input type="number" v-model.number="rule.sequence.digits" min="1" max="10" @input="notifyChange" />
            </div>
            <div class="field">
              <label>前缀</label>
              <input type="text" v-model="rule.sequence.prefix" @input="notifyChange" />
            </div>
            <div class="field">
              <label>后缀</label>
              <input type="text" v-model="rule.sequence.suffix" @input="notifyChange" />
            </div>
          </template>
          <template v-else-if="rule.type === 'replace'">
            <div class="field">
              <label>查找内容</label>
              <input type="text" v-model="rule.replace.search" @input="notifyChange" />
            </div>
            <div class="field">
              <label>替换为</label>
              <input type="text" v-model="rule.replace.replace" @input="notifyChange" />
            </div>
            <div class="field checkbox">
              <label>
                <input type="checkbox" v-model="rule.replace.useRegex" @change="notifyChange" />
                使用正则表达式
              </label>
            </div>
          </template>
          <template v-else-if="rule.type === 'date'">
            <div class="field">
              <label>日期来源</label>
              <select v-model="rule.date.source" @change="notifyChange">
                <option value="modified">文件修改时间</option>
                <option value="exif">EXIF信息（图片）</option>
              </select>
            </div>
            <div class="field">
              <label>日期格式</label>
              <select v-model="rule.date.format" @change="notifyChange">
                <option value="YYYYMMDD">YYYYMMDD</option>
                <option value="YYYY-MM-DD">YYYY-MM-DD</option>
                <option value="YYYY_MM_DD">YYYY_MM_DD</option>
                <option value="YYYYMMDD_HHMMSS">YYYYMMDD_HHMMSS</option>
              </select>
            </div>
            <div class="field">
              <label>位置</label>
              <select v-model="rule.date.position" @change="notifyChange">
                <option value="prefix">前缀</option>
                <option value="suffix">后缀</option>
              </select>
            </div>
            <div class="field">
              <label>分隔符</label>
              <input type="text" v-model="rule.date.separator" @input="notifyChange" />
            </div>
          </template>
          <template v-else-if="rule.type === 'case'">
            <div class="field">
              <label>转换方式</label>
              <select v-model="rule.case.mode" @change="notifyChange">
                <option value="upper">全大写</option>
                <option value="lower">全小写</option>
                <option value="title">首字母大写</option>
                <option value="camel">驼峰式</option>
              </select>
            </div>
          </template>
          <template v-else-if="rule.type === 'extension'">
            <div class="field">
              <label>新扩展名（留空不修改）</label>
              <input type="text" v-model="rule.extension.newExt" @input="notifyChange" placeholder="例如: jpg" />
            </div>
            <div class="field checkbox">
              <label>
                <input type="checkbox" v-model="rule.extension.unifyCase" @change="notifyChange" />
                统一扩展名大小写
              </label>
            </div>
            <div v-if="rule.extension.unifyCase" class="field">
              <label>目标格式</label>
              <select v-model="rule.extension.targetCase" @change="notifyChange">
                <option value="lower">小写</option>
                <option value="upper">大写</option>
              </select>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps(['rules'])
const emit = defineEmits(['rules-changed'])

const localRules = ref([])

watch(() => props.rules, (newRules) => {
  localRules.value = JSON.parse(JSON.stringify(newRules))
}, { immediate: true, deep: true })

function getRuleName(type) {
  const names = {
    sequence: '序号填充',
    replace: '查找替换',
    date: '日期标记',
    case: '大小写转换',
    extension: '扩展名修改'
  }
  return names[type] || type
}

function toggleRule(index) {
  localRules.value[index].enabled = !localRules.value[index].enabled
  notifyChange()
}

function notifyChange() {
  emit('rules-changed', JSON.parse(JSON.stringify(localRules.value)))
}
</script>

<style scoped>
.rule-manager h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 12px;
}

.rules-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.rule-item {
  background: #f8fafc;
  border: 2px solid #e2e8f0;
  border-radius: 10px;
  overflow: hidden;
  transition: all 0.2s;
}

.rule-item.active {
  border-color: #667eea;
  background: #f0f4ff;
}

.rule-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  cursor: pointer;
  user-select: none;
}

.rule-header input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: #667eea;
}

.rule-name {
  font-weight: 500;
  font-size: 14px;
  color: #334155;
}

.rule-content {
  padding: 12px 16px 16px;
  border-top: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field.checkbox {
  flex-direction: row;
  align-items: center;
}

.field label {
  font-size: 13px;
  color: #475569;
  font-weight: 500;
}

.field input,
.field select {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 14px;
}

.field input:focus,
.field select:focus {
  outline: none;
  border-color: #667eea;
}
</style>
