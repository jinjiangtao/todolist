<template>
  <div class="admin">
    <div class="header">
      <h1>管理后台</h1>
      <div class="header-right">
        <router-link to="/" class="nav-link">查看博客</router-link>
        <router-link to="/admin/settings" class="nav-link">设置</router-link>
        <a @click="handleLogout" class="nav-link">退出</a>
      </div>
    </div>
    <div class="container">
      <el-tabs v-model="activeTab">
        <!-- 文章管理 -->
        <el-tab-pane label="文章管理" name="articles">
          <div class="toolbar">
            <el-button type="primary" @click="goToCreate">新建文章</el-button>
          </div>
          <div v-if="loading" class="loading">加载中...</div>
          <div v-else>
            <el-table :data="articles" stripe>
              <el-table-column prop="id" label="ID" width="80" />
              <el-table-column prop="title" label="标题" />
              <el-table-column label="分类" width="120">
                <template #default="{ row }">
                  {{ row.category ? row.category.name : '-' }}
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 1 ? 'success' : 'info'">
                    {{ row.status === 1 ? '已发布' : '草稿' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column label="置顶" width="80">
                <template #default="{ row }">
                  <el-tag v-if="row.is_top" type="warning">是</el-tag>
                  <span v-else>-</span>
                </template>
              </el-table-column>
              <el-table-column prop="view_count" label="阅读" width="100" />
              <el-table-column prop="created_at" label="创建时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="200" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" size="small" @click="goToEdit(row.id)">编辑</el-button>
                  <el-button type="danger" size="small" @click="handleDelete(row.id)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-pagination
              v-if="totalPage > 1"
              v-model:current-page="page"
              :page-size="pageSize"
              :total="total"
              layout="prev, pager, next"
              @current-change="fetchArticles"
              class="pagination"
            />
          </div>
        </el-tab-pane>

        <!-- 分类管理 -->
        <el-tab-pane label="分类管理" name="categories">
          <div class="toolbar">
            <el-input v-model="newCategory" placeholder="输入新分类名称" style="width: 300px; margin-right: 10px" />
            <el-button type="primary" @click="handleAddCategory">添加分类</el-button>
          </div>
          <el-table :data="categories" stripe style="margin-top: 20px">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="分类名称" />
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button type="danger" size="small" @click="handleDeleteCategory(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 标签管理 -->
        <el-tab-pane label="标签管理" name="tags">
          <div class="toolbar">
            <el-input v-model="newTag" placeholder="输入新标签名称" style="width: 300px; margin-right: 10px" />
            <el-button type="primary" @click="handleAddTag">添加标签</el-button>
          </div>
          <el-table :data="tags" stripe style="margin-top: 20px">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="标签名称" />
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button type="danger" size="small" @click="handleDeleteTag(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const authStore = useAuthStore()
const activeTab = ref('articles')
const articles = ref([])
const categories = ref([])
const tags = ref([])
const loading = ref(true)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const totalPage = ref(0)
const newCategory = ref('')
const newTag = ref('')

async function fetchArticles() {
  loading.value = true
  try {
    const res = await request.get('/admin/articles', {
      params: { page: page.value, page_size: pageSize.value }
    })
    articles.value = res.data.list
    total.value = res.data.total
    totalPage.value = res.data.total_page
  } finally {
    loading.value = false
  }
}

async function fetchCategories() {
  const res = await request.get('/admin/categories')
  categories.value = res.data || []
}

async function fetchTags() {
  const res = await request.get('/admin/tags')
  tags.value = res.data || []
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

function goToCreate() {
  router.push('/admin/article/edit')
}

function goToEdit(id) {
  router.push(`/admin/article/edit/${id}`)
}

async function handleDelete(id) {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/admin/articles/${id}`)
    ElMessage.success('删除成功')
    fetchArticles()
  } catch (error) {
    if (error !== 'cancel') {
      throw error
    }
  }
}

async function handleAddCategory() {
  if (!newCategory.value.trim()) {
    ElMessage.warning('请输入分类名称')
    return
  }
  await request.post('/admin/categories', { name: newCategory.value })
  ElMessage.success('添加成功')
  newCategory.value = ''
  fetchCategories()
}

async function handleDeleteCategory(id) {
  try {
    await ElMessageBox.confirm('确定要删除这个分类吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/admin/categories/${id}`)
    ElMessage.success('删除成功')
    fetchCategories()
  } catch (error) {
    if (error !== 'cancel') {
      throw error
    }
  }
}

async function handleAddTag() {
  if (!newTag.value.trim()) {
    ElMessage.warning('请输入标签名称')
    return
  }
  await request.post('/admin/tags', { name: newTag.value })
  ElMessage.success('添加成功')
  newTag.value = ''
  fetchTags()
}

async function handleDeleteTag(id) {
  try {
    await ElMessageBox.confirm('确定要删除这个标签吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/admin/tags/${id}`)
    ElMessage.success('删除成功')
    fetchTags()
  } catch (error) {
    if (error !== 'cancel') {
      throw error
    }
  }
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

onMounted(() => {
  fetchArticles()
  fetchCategories()
  fetchTags()
})
</script>

<style scoped>
.admin {
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

.header-right {
  display: flex;
  gap: 24px;
}

.nav-link {
  color: white;
  text-decoration: none;
  cursor: pointer;
}

.container {
  max-width: 1200px;
  margin: 40px auto;
  padding: 0 20px;
  background: white;
  padding: 40px;
  border-radius: 8px;
}

.toolbar {
  margin-bottom: 20px;
}

.loading {
  text-align: center;
  padding: 60px 0;
  color: #999;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
