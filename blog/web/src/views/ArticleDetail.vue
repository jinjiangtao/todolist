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
          <span v-if="article.category">分类：{{ article.category.name }}</span>
        </div>
        <div v-if="article.tags && article.tags.length" class="tags">
          <el-tag v-for="tag in article.tags" :key="tag.id" size="small" style="margin-right: 8px">{{ tag.name }}</el-tag>
        </div>
        <div class="content" v-html="safeHtml" ref="contentRef"></div>
      </div>
    </div>

    <!-- 图片查看器 -->
    <el-image-viewer
      v-if="imageViewerVisible"
      :url-list="imageUrls"
      :initial-index="currentImageIndex"
      @close="imageViewerVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DOMPurify from 'dompurify'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const article = ref(null)
const loading = ref(true)
const contentRef = ref(null)
const imageViewerVisible = ref(false)
const imageUrls = ref([])
const currentImageIndex = ref(0)

const safeHtml = computed(() => {
  if (!article.value) return ''
  return DOMPurify.sanitize(article.value.content || '')
})

async function fetchArticle() {
  loading.value = true
  try {
    const res = await request.get(`/articles/${route.params.id}`)
    article.value = res.data
    await nextTick()
    setupImageClickHandler()
  } finally {
    loading.value = false
  }
}

function setupImageClickHandler() {
  if (!contentRef.value) return
  const images = contentRef.value.querySelectorAll('img')
  imageUrls.value = Array.from(images).map(img => img.src)
  
  images.forEach((img, index) => {
    img.style.cursor = 'pointer'
    img.style.maxWidth = '100%'
    img.style.height = 'auto'
    img.addEventListener('click', () => {
      currentImageIndex.value = index
      imageViewerVisible.value = true
    })
  })

  const videos = contentRef.value.querySelectorAll('video')
  videos.forEach(video => {
    video.style.maxWidth = '100%'
    video.controls = true
  })
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
  margin-bottom: 16px;
  display: flex;
  gap: 20px;
}

.tags {
  margin-bottom: 32px;
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

.content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 16px 0;
}

.content :deep(video) {
  max-width: 100%;
  border-radius: 4px;
  margin: 16px 0;
}

.content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 16px 0;
}

.content :deep(th),
.content :deep(td) {
  border: 1px solid #ddd;
  padding: 8px 12px;
  text-align: left;
}

.content :deep(th) {
  background: #f5f5f5;
}

.content :deep(blockquote) {
  border-left: 4px solid #409eff;
  padding-left: 16px;
  margin: 16px 0;
  color: #666;
  background: #f9f9f9;
  padding: 16px;
  border-radius: 0 4px 4px 0;
}
</style>
