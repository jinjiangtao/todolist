<template>
  <div class="article-history">
    <el-table :data="histories" style="width: 100%" v-loading="loading">
      <el-table-column prop="version" label="版本号" width="100" />
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="created_at" label="修改时间" width="200">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250">
        <template #default="{ row }">
          <el-button size="small" @click="viewVersion(row)">查看</el-button>
          <el-button size="small" @click="compareWithCurrent(row)">与当前对比</el-button>
          <el-button size="small" @click="handleCompare(row)">选择对比</el-button>
          <el-button type="primary" size="small" @click="confirmRestore(row)">恢复</el-button>
          <el-button type="danger" size="small" @click="confirmDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="loadHistories"
      @current-change="loadHistories"
      style="margin-top: 20px; justify-content: flex-end"
    />

    <el-dialog v-model="viewVisible" title="查看历史版本" width="80%" top="5vh">
      <div v-if="viewVersionData">
        <h3>{{ viewVersionData.title }}</h3>
        <p><small>{{ formatDate(viewVersionData.created_at) }}</small></p>
        <div class="version-content" v-html="viewVersionData.content"></div>
      </div>
    </el-dialog>

    <el-dialog v-model="compareVisible" title="版本对比" width="90%" top="5vh">
      <div class="compare-container">
        <div v-if="compareMode === 'inline'" class="compare-inline">
          <div v-html="inlineDiff" class="diff-content"></div>
        </div>
        <div v-else class="compare-split">
          <div class="compare-side">
            <h4>{{ leftVersion ? `版本 ${leftVersion.version}` : '当前版本' }}</h4>
            <div class="version-content" v-html="leftContent"></div>
          </div>
          <div class="compare-side">
            <h4>{{ rightVersion ? `版本 ${rightVersion.version}` : '当前版本' }}</h4>
            <div class="version-content" v-html="rightContent"></div>
          </div>
        </div>
        <div style="margin-top: 20px">
          <el-radio-group v-model="compareMode">
            <el-radio value="inline">高亮对比</el-radio>
            <el-radio value="split">双栏对比</el-radio>
          </el-radio-group>
        </div>
      </div>
    </el-dialog>

    <el-dialog v-model="selectCompareVisible" title="选择版本进行对比" width="50%">
      <el-table :data="histories" style="width: 100%" @row-click="selectForCompare">
        <el-table-column prop="version" label="版本号" width="100" />
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="created_at" label="修改时间" width="200">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import DOMPurify from 'dompurify'
import { diffLines } from 'diff'

const props = defineProps({
  articleId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits(['restore'])

const loading = ref(false)
const histories = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const viewVisible = ref(false)
const viewVersionData = ref(null)

const compareVisible = ref(false)
const selectCompareVisible = ref(false)
const compareMode = ref('inline')
const compareBaseVersion = ref(null)
const leftVersion = ref(null)
const rightVersion = ref(null)

const DOMPURIFY_CONFIG = {
  ALLOWED_TAGS: ['img', 'p', 'br', 'span', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'ul', 'ol', 'li', 'blockquote', 'pre', 'code', 'strong', 'em', 'a', 'table', 'thead', 'tbody', 'tr', 'th', 'td', 'hr', 'div', 'video', 'source', 'del', 'ins'],
  ALLOWED_ATTR: ['src', 'href', 'class', 'style', 'id', 'target', 'controls', 'type', 'width', 'height', 'alt', 'title', 'data-value']
}

const currentArticle = ref(null)

const leftContent = computed(() => {
  if (leftVersion) {
    return DOMPurify.sanitize(leftVersion.content || '', DOMPURIFY_CONFIG)
  } else if (currentArticle.value) {
    return DOMPurify.sanitize(currentArticle.value.content || '', DOMPURIFY_CONFIG)
  }
  return ''
})

const rightContent = computed(() => {
  if (rightVersion) {
    return DOMPurify.sanitize(rightVersion.content || '', DOMPURIFY_CONFIG)
  } else if (currentArticle.value) {
    return DOMPurify.sanitize(currentArticle.value.content || '', DOMPURIFY_CONFIG)
  }
  return ''
})

const inlineDiff = computed(() => {
  const oldText = leftContent.value || ''
  const newText = rightContent.value || ''

  const diff = diffLines(oldText, newText)
  let html = ''

  diff.forEach((part) => {
    if (part.added) {
      html += `<ins style="background: #e6ffec; color: #22863a;">${escapeHtml(part.value)}</ins>`
    } else if (part.removed) {
      html += `<del style="background: #ffeef0; color: #cb2431; text-decoration: line-through;">${escapeHtml(part.value)}</del>`
    } else {
      html += escapeHtml(part.value)
    }
  })

  return DOMPurify.sanitize(html, { ...DOMPURIFY_CONFIG, ALLOWED_TAGS: [...DOMPURIFY_CONFIG.ALLOWED_TAGS, 'ins', 'del'] })
})

function escapeHtml(text) {
  const div = document.createElement('div')
  div.textContent = text
  return div.innerHTML
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

async function loadHistories() {
  loading.value = true
  try {
    const res = await request.get(`/admin/articles/${props.articleId}/histories?page=${page.value}&page_size=${pageSize.value}`)
    histories.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (err) {
    ElMessage.error('加载历史版本失败')
    console.error(err)
  } finally {
    loading.value = false
  }
}

async function loadCurrentArticle() {
  try {
    const res = await request.get('/admin/articles')
    currentArticle.value = res.data.list.find(a => a.id === parseInt(props.articleId))
  } catch (err) {
    console.error(err)
  }
}

async function viewVersion(version) {
  viewVersionData.value = version
  viewVisible.value = true
}

async function compareWithCurrent(version) {
  await loadCurrentArticle()
  leftVersion.value = version
  rightVersion.value = null
  compareMode.value = 'inline'
  compareVisible.value = true
}

function handleCompare(version) {
  compareBaseVersion.value = version
  selectCompareVisible.value = true
}

async function selectForCompare(row) {
  if (!compareBaseVersion.value) return

  if (row.id === compareBaseVersion.value.id) {
    ElMessage.warning('请选择不同的版本进行对比')
    return
  }

  selectCompareVisible.value = false
  await loadCurrentArticle()

  if (compareBaseVersion.value.version < row.version) {
    leftVersion.value = compareBaseVersion.value
    rightVersion.value = row
  } else {
    leftVersion.value = row
    rightVersion.value = compareBaseVersion.value
  }

  compareMode.value = 'inline'
  compareVisible.value = true
}

async function confirmRestore(version) {
  try {
    await ElMessageBox.confirm(
      `确定要恢复到版本 ${version.version} 吗？`,
      '恢复确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await request.post(`/admin/articles/${props.articleId}/histories/${version.id}/restore`)
    ElMessage.success('恢复成功')
    
    emit('restore', {
      title: version.title,
      content: version.content
    })
    loadHistories()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('恢复失败')
      console.error(err)
    }
  }
}

async function confirmDelete(version) {
  try {
    await ElMessageBox.confirm(
      `确定要删除版本 ${version.version} 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await request.delete(`/admin/articles/${props.articleId}/histories/${version.id}`)
    ElMessage.success('删除成功')
    loadHistories()
  } catch (err) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
      console.error(err)
    }
  }
}

onMounted(() => {
  loadHistories()
  loadCurrentArticle()
})
</script>

<style scoped>
.article-history {
  padding: 10px 0;
}

.version-content {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 20px;
  margin-top: 10px;
  max-height: 400px;
  overflow-y: auto;
}

.compare-container {
  padding: 10px 0;
}

.compare-split {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.compare-side {
  flex: 1;
}

.diff-content {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 20px;
  max-height: 500px;
  overflow-y: auto;
  white-space: pre-wrap;
}

.diff-content ins {
  background: #e6ffec;
  color: #22863a;
  padding: 2px 4px;
  border-radius: 3px;
}

.diff-content del {
  background: #ffeef0;
  color: #cb2431;
  padding: 2px 4px;
  border-radius: 3px;
}
</style>
