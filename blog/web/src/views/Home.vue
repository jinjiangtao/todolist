<template>
  <div class="home">
    <div class="header">
      <h1>个人博客</h1>
      <router-link to="/login" class="login-link">登录</router-link>
    </div>
    <div class="container">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else>
        <div v-if="articles.length === 0" class="empty">暂无文章</div>
        <div v-else class="article-list">
          <div
            v-for="article in articles"
            :key="article.id"
            class="article-item"
            @click="goToDetail(article.id)"
          >
            <h2 class="title">{{ article.title }}</h2>
            <div class="meta">
              <span>{{ formatDate(article.created_at) }}</span>
              <span>阅读：{{ article.view_count }}</span>
            </div>
            <p class="summary">{{ getSummary(article.content) }}</p>
          </div>
        </div>
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
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/utils/request'

const router = useRouter()
const articles = ref([])
const loading = ref(true)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const totalPage = ref(0)

async function fetchArticles() {
  loading.value = true
  try {
    const res = await request.get('/articles', {
      params: { page: page.value, page_size: pageSize.value }
    })
    articles.value = res.data.list
    total.value = res.data.total
    totalPage.value = res.data.total_page
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

function getSummary(content) {
  return content ? content.substring(0, 150) + (content.length > 150 ? '...' : '') : ''
}

function goToDetail(id) {
  router.push(`/article/${id}`)
}

onMounted(() => {
  fetchArticles()
})
</script>

<style scoped>
.home {
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

.login-link {
  color: white;
  text-decoration: none;
  padding: 8px 16px;
  border: 1px solid white;
  border-radius: 4px;
}

.container {
  max-width: 800px;
  margin: 40px auto;
  padding: 0 20px;
}

.loading,
.empty {
  text-align: center;
  padding: 60px 0;
  color: #999;
}

.article-item {
  background: white;
  padding: 24px;
  margin-bottom: 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: box-shadow 0.3s;
}

.article-item:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.title {
  font-size: 20px;
  color: #333;
  margin-bottom: 12px;
}

.meta {
  font-size: 14px;
  color: #999;
  margin-bottom: 12px;
  display: flex;
  gap: 20px;
}

.summary {
  color: #666;
  line-height: 1.6;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}
</style>
