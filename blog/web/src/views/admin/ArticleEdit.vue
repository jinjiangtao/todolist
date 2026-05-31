<template>
  <div class="article-edit">
    <div class="header">
      <h1>{{ isEdit ? '编辑文章' : '新建文章' }}</h1>
      <div class="header-right">
        <el-button @click="handlePreview" type="info">预览</el-button>
        <el-button @click="handleSaveDraft" type="default">保存草稿</el-button>
        <router-link to="/admin" class="nav-link">返回</router-link>
      </div>
    </div>
    <div class="container">
      <el-tabs v-if="isEdit" v-model="activeTab">
        <el-tab-pane label="编辑文章" name="edit">
          <ArticleForm
            ref="articleFormRef"
            :form="form"
            :rules="rules"
            :categories="categories"
            :tags="tags"
            :editor-config="editorConfig"
            :toolbar-config="toolbarConfig"
            @editor-created="handleCreated"
          />
          <el-form-item>
            <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
            <el-button @click="goBack">取消</el-button>
          </el-form-item>
        </el-tab-pane>
        <el-tab-pane label="历史版本" name="history">
          <ArticleHistory
            :article-id="route.params.id"
            @restore="handleRestore"
          />
        </el-tab-pane>
      </el-tabs>

      <div v-else>
        <ArticleForm
          ref="articleFormRef"
          :form="form"
          :rules="rules"
          :categories="categories"
          :tags="tags"
          :editor-config="editorConfig"
          :toolbar-config="toolbarConfig"
          @editor-created="handleCreated"
        />
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </div>
    </div>

    <el-dialog v-model="previewVisible" title="文章预览" width="80%" top="5vh">
      <div class="preview-content" v-html="safeHtml"></div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import * as wangEditor from '@wangeditor/editor'
import DOMPurify from 'dompurify'
import request from '@/utils/request'
import { useAuthStore } from '@/store/auth'
import axios from 'axios'
import ArticleForm from './ArticleEditForm.vue'
import ArticleHistory from './ArticleHistory.vue'

const route = useRoute()
const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const editorRef = ref()
const previewVisible = ref(false)
const categories = ref([])
const tags = ref([])
const activeTab = ref('edit')
const authStore = useAuthStore()

const isEdit = computed(() => !!route.params.id)

const form = ref({
  title: '',
  content: '',
  status: 0,
  is_top: false,
  password: '',
  category_id: null,
  tag_ids: [],
  published_at: null
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }]
}

const DOMPURIFY_CONFIG = {
  ALLOWED_TAGS: ['img', 'p', 'br', 'span', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'ul', 'ol', 'li', 'blockquote', 'pre', 'code', 'strong', 'em', 'a', 'table', 'thead', 'tbody', 'tr', 'th', 'td', 'hr', 'div', 'video', 'source'],
  ALLOWED_ATTR: ['src', 'href', 'class', 'style', 'id', 'target', 'controls', 'type', 'width', 'height', 'alt', 'title', 'data-value']
}

const safeHtml = computed(() => {
  return DOMPurify.sanitize(form.value.content || '', DOMPURIFY_CONFIG)
})

const toolbarConfig = {}

async function customUploadImage(file, insertFn) {
  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await axios.post('http://localhost:8080/api/v1/admin/upload', formData, {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    const fullUrl = `http://localhost:8080${res.data.data.url}`
    insertFn(fullUrl)
  } catch (err) {
    ElMessage.error('图片上传失败')
    console.error(err)
  }
}

async function customUploadVideo(file, insertFn) {
  const formData = new FormData()
  formData.append('file', file)

  try {
    const res = await axios.post('http://localhost:8080/api/v1/admin/upload', formData, {
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      }
    })
    const fullUrl = `http://localhost:8080${res.data.data.url}`
    insertFn(fullUrl)
  } catch (err) {
    ElMessage.error('视频上传失败')
    console.error(err)
  }
}

const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {
    uploadImage: {
      customUpload: customUploadImage,
      fieldName: 'file',
      maxFileSize: 10 * 1024 * 1024,
      allowedFileTypes: ['image/*']
    },
    uploadVideo: {
      customUpload: customUploadVideo,
      fieldName: 'file',
      maxFileSize: 100 * 1024 * 1024,
      allowedFileTypes: ['video/*']
    }
  }
}

function handleCreated(editor) {
  editorRef.value = editor
}

async function loadCategories() {
  const res = await request.get('/categories')
  categories.value = res.data || []
}

async function loadTags() {
  const res = await request.get('/tags')
  tags.value = res.data || []
}

async function fetchArticle() {
  if (!route.params.id) return
  const res = await request.get('/admin/articles')
  const article = res.data.list.find(a => a.id === parseInt(route.params.id))
  if (article) {
    form.value = {
      title: article.title,
      content: article.content,
      status: article.status,
      is_top: article.is_top,
      password: article.password || '',
      category_id: article.category_id,
      tag_ids: article.tags ? article.tags.map(t => t.id) : [],
      published_at: article.published_at
    }
  }
}

async function handleSave(publish = false) {
  if (!articleFormRef.value) return
  await articleFormRef.value.formRef.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (publish) {
          form.value.status = 1
        }
        const data = {
          ...form.value,
          published_at: form.value.published_at ? new Date(form.value.published_at).toISOString() : ''
        }
        if (isEdit.value) {
          await request.put(`/admin/articles/${route.params.id}`, data)
          ElMessage.success('更新成功')
        } else {
          await request.post('/admin/articles', data)
          ElMessage.success('创建成功')
        }
        router.push('/admin')
      } finally {
        loading.value = false
      }
    }
  })
}

async function handleSaveDraft() {
  await handleSave(false)
}

function handlePreview() {
  previewVisible.value = true
}

function goBack() {
  router.push('/admin')
}

function handleRestore(article) {
  form.value.title = article.title
  form.value.content = article.content
  activeTab.value = 'edit'
}

onMounted(() => {
  loadCategories()
  loadTags()
  fetchArticle()
})

onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null) return
  editor.destroy()
})
</script>

<style src="@wangeditor/editor/dist/css/style.css"></style>

<style scoped>
.article-edit {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  background: #409eff;
  color: white;
  padding: 20px 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h1 {
  font-size: 24px;
  margin: 0;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 8px 16px;
  border: 1px solid white;
  border-radius: 4px;
}

.header-right {
  display: flex;
  gap: 10px;
  align-items: center;
}

.container {
  max-width: 1200px;
  margin: 40px auto;
  padding: 40px;
  background: white;
  border-radius: 8px;
}

.editor-container {
  border: 1px solid #ccc;
  border-radius: 4px;
}

.preview-content {
  padding: 20px;
  line-height: 1.8;
}

.preview-content :deep(img) {
  max-width: 100%;
  height: auto;
}

.preview-content :deep(video) {
  max-width: 100%;
}
</style>
