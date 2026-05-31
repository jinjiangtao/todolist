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
            <div class="article-content">
              <div v-if="getFirstImage(article.content)" class="thumbnail">
                <img :src="getFirstImage(article.content)" alt="文章缩略图" />
              </div>
              <div class="text-content">
                <h2 class="title">{{ article.title }}</h2>
                <div class="meta">
                  <span>{{ formatDate(article.created_at) }}</span>
                  <span>阅读：{{ article.view_count }}</span>
                  <span v-if="article.category">分类：{{ article.category.name }}</span>
                </div>
                <p class="summary">{{ getSummary(article.content) }}</p>
                <div v-if="article.tags && article.tags.length" class="tags">
                  <el-tag v-for="tag in article.tags.slice(0, 3)" :key="tag.id" size="small">{{ tag.name }}</el-tag>
                </div>
              </div>
            </div>
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

function getFirstImage(content) {
  if (!content) return null
  const imgRegex = /<img[^>]+src="([^">]+)"/
  const match = content.match(imgRegex)
  return match ? match[1] : null
}

function removeHtmlTags(str) {
  if (!str) return ''
  return str.replace(/<[^>]*>/g, '').replace(/&nbsp;/g, ' ')
}

function getSummary(content) {
  const text = removeHtmlTags(content)
  return text ? text.substring(0, 150) + (text.length > 150 ? '...' : '') : ''
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
  max-width: 900px;
  margin: 40px auto;
  padding: 0 20px;
}

.loading,
.empty {
  text-align: center;
  padding: 60px 0;
  color: #999;
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.article-item {
  background: white;
  padding: 24px;
  border-radius: 8px;
  cursor: pointer;
  transition: box-shadow 0.3s;
}

.article-item:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.article-content {
  display: flex;
  gap: 20px;
}

.thumbnail {
  flex-shrink: 0;
  width: 150px;
  height: 100px;
  overflow: hidden;
  border-radius: 8px;
}

.thumbnail img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.text-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.title {
  font-size: 20px;
  color: #333;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.meta {
  font-size: 14px;
  color: #999;
  margin-bottom: 8px;
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.summary {
  color: #666;
  line-height: 1.6;
  flex: 1;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.tags {
  margin-top: 12px;
  flex-shrink: 0;
}

.tags :deep(.el-tag) {
  margin-right: 8px;
  margin-bottom: 4px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

@media (max-width: 600px) {
  .article-content {
    flex-direction: column;
  }
  
  .thumbnail {
    width: 100%;
    height: 150px;
  }
}
</style>
