<template>
  <div class="article-edit">
    <div class="header">
      <h1>{{ isEdit ? '编辑文章' : '新建文章' }}</h1>
      <div class="header-right">
        <router-link to="/admin" class="nav-link">返回</router-link>
      </div>
    </div>
    <div class="container">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入文章标题" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="20"
            placeholder="请输入文章内容（支持 Markdown）"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :label="0">草稿</el-radio>
            <el-radio :label="1">发布</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSave">保存</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const formRef = ref(null)
const loading = ref(false)

const isEdit = computed(() => !!route.params.id)

const form = ref({
  title: '',
  content: '',
  status: 0
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }]
}

async function fetchArticle() {
  if (!route.params.id) return
  const res = await request.get(`/admin/articles`)
  const article = res.data.list.find(a => a.id === parseInt(route.params.id))
  if (article) {
    form.value = {
      title: article.title,
      content: article.content,
      status: article.status
    }
  }
}

async function handleSave() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        if (isEdit.value) {
          await request.put(`/admin/articles/${route.params.id}`, form.value)
          ElMessage.success('更新成功')
        } else {
          await request.post('/admin/articles', form.value)
          ElMessage.success('创建成功')
        }
        router.push('/admin')
      } finally {
        loading.value = false
      }
    }
  })
}

function goBack() {
  router.push('/admin')
}

onMounted(() => {
  fetchArticle()
})
</script>

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
}

.nav-link {
  color: white;
  text-decoration: none;
}

.container {
  max-width: 1000px;
  margin: 40px auto;
  padding: 0 20px;
  background: white;
  padding: 40px;
  border-radius: 8px;
}
</style>
