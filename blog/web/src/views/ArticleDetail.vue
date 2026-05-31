<template>
  <div class="article-detail">
    <div class="header">
      <h1>个人博客</h1>
      <router-link to="/" class="back-link">返回首页</router-link>
    </div>
    <div class="container">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="article">
        <h1 class="title">{{ article.title }}</h1>
        <div class="meta">
          <span>{{ formatDate(article.created_at) }}</span>
          <span>阅读：{{ article.view_count }}</span>
        </div>
        <div class="content" v-html="renderedContent"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked } from 'marked'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const article = ref(null)
const loading = ref(true)

const renderedContent = computed(() => {
  if (!article.value) return ''
  return marked(article.value.content || '')
})

async function fetchArticle() {
  loading.value = true
  try {
    const res = await request.get(`/articles/${route.params.id}`)
    article.value = res.data
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchArticle()
})
</script>

<style scoped>
.article-detail {
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

.back-link {
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

.loading {
  text-align: center;
  padding: 60px 0;
  color: #999;
}

.title {
  font-size: 32px;
  color: #333;
  margin-bottom: 16px;
}

.meta {
  font-size: 14px;
  color: #999;
  margin-bottom: 32px;
  display: flex;
  gap: 20px;
}

.content {
  background: white;
  padding: 40px;
  border-radius: 8px;
  line-height: 1.8;
}

.content :deep(h1),
.content :deep(h2),
.content :deep(h3) {
  margin-top: 24px;
  margin-bottom: 12px;
  color: #333;
}

.content :deep(p) {
  margin-bottom: 16px;
  color: #666;
}

.content :deep(code) {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
}

.content :deep(pre) {
  background: #2d2d2d;
  color: #f8f8f2;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin-bottom: 16px;
}

.content :deep(pre code) {
  background: none;
  padding: 0;
}
</style>
