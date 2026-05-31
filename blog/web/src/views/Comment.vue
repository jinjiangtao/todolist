<template>
  <div class="comment-section">
    <h3 class="comment-title">评论 ({{ comments.length }})</h3>
    
    <!-- 评论列表 -->
    <div v-if="comments.length > 0" class="comment-list">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="comment-avatar">
          <el-avatar :src="comment.author_email" :size="48">
            {{ comment.author_name.charAt(0).toUpperCase() }}
          </el-avatar>
        </div>
        <div class="comment-content">
          <div class="comment-header">
            <span class="comment-author">
              {{ comment.author_name }}
              <el-tag v-if="comment.is_admin" type="success" size="small" style="margin-left: 8px">博主</el-tag>
            </span>
            <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
          </div>
          <div class="comment-body">{{ comment.content }}</div>
          <div class="comment-actions">
            <el-button 
              link 
              size="small" 
              @click="handleLike(comment)"
              :disabled="hasLiked(comment.id)"
            >
              <span style="margin-right: 4px">👍</span>
              <span>{{ comment.likes }}</span>
            </el-button>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="no-comments">
      暂无评论，来发表第一篇评论吧！
    </div>

    <!-- 评论表单 -->
    <div class="comment-form">
      <h4 class="form-title">发表评论</h4>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="0">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item prop="author_name">
              <el-input 
                v-model="form.author_name" 
                placeholder="昵称（必填）"
                prefix-icon="User"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item prop="author_email">
              <el-input 
                v-model="form.author_email" 
                placeholder="邮箱（必填）"
                prefix-icon="Message"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="4"
            placeholder="评论内容（必填）"
          />
        </el-form-item>
        
        <!-- 验证码 -->
        <div v-if="captchaEnabled" class="captcha-section">
          <el-form-item prop="captcha" style="display: inline-block; width: 200px;">
            <el-input
              v-model="form.captcha"
              placeholder="验证码"
              @keyup.enter="handleSubmit"
            />
          </el-form-item>
          <span class="captcha-question" @click="refreshCaptcha">
            {{ captchaQuestion }}
            <i class="el-icon-refresh" style="cursor: pointer; margin-left: 8px;"></i>
          </span>
        </div>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">提交评论</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 提示消息 -->
    <el-message
      v-if="showMessage"
      :type="messageType"
      :message="messageText"
      :duration="3000"
      @close="showMessage = false"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const props = defineProps({
  articleId: {
    type: [Number, String],
    required: true
  }
})

const comments = ref([])
const loading = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const form = ref({
  author_name: '',
  author_email: '',
  content: '',
  captcha: ''
})
const captchaId = ref('')
const captchaQuestion = ref('')
const captchaEnabled = ref(true)
const likedComments = ref(new Set())
const showMessage = ref(false)
const messageType = ref('success')
const messageText = ref('')

const rules = {
  author_name: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  author_email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入评论内容', trigger: 'blur' }
  ],
  captcha: [
    { required: true, message: '请输入验证码', trigger: 'blur' }
  ]
}

async function fetchComments() {
  loading.value = true
  try {
    const res = await request.get(`/articles/${props.articleId}/comments`)
    comments.value = res.data || []
  } finally {
    loading.value = false
  }
}

async function fetchCaptcha() {
  try {
    const res = await request.get('/captcha')
    captchaId.value = res.data.captcha_id
    captchaQuestion.value = res.data.question
    captchaEnabled.value = true
  } catch (error) {
    console.error('Failed to fetch captcha:', error)
    captchaEnabled.value = false
  }
}

async function refreshCaptcha() {
  await fetchCaptcha()
  form.value.captcha = ''
}

function hasLiked(commentId) {
  return likedComments.value.has(commentId)
}

async function handleLike(comment) {
  if (hasLiked(comment.id)) {
    ElMessage.warning('您已经点过赞了')
    return
  }

  try {
    const res = await request.post(`/comments/${comment.id}/like`)
    comment.likes = res.data.likes
    likedComments.value.add(comment.id)
    ElMessage.success('点赞成功')
  } catch (error) {
    if (error.response?.data?.message) {
      ElMessage.warning(error.response.data.message)
    } else {
      ElMessage.error('点赞失败')
    }
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const payload = {
          article_id: parseInt(props.articleId),
          author_name: form.value.author_name,
          author_email: form.value.author_email,
          content: form.value.content
        }

        if (captchaEnabled.value) {
          payload.captcha = form.value.captcha
          payload.captcha_id = captchaId.value
        }

        await request.post('/comments', payload)
        
        ElMessage.success('评论提交成功')
        
        form.value.content = ''
        form.value.captcha = ''
        
        await fetchComments()
        
        if (captchaEnabled.value) {
          await fetchCaptcha()
        }
      } catch (error) {
        if (error.response?.data?.message) {
          ElMessage.error(error.response.data.message)
        } else {
          ElMessage.error('评论提交失败')
        }
      } finally {
        submitting.value = false
      }
    }
  })
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchComments()
  fetchCaptcha()
})
</script>

<style scoped>
.comment-section {
  margin-top: 60px;
  padding-top: 40px;
  border-top: 2px solid #e8e8e8;
}

.comment-title {
  font-size: 24px;
  color: #333;
  margin-bottom: 30px;
}

.comment-list {
  margin-bottom: 40px;
}

.comment-item {
  display: flex;
  gap: 16px;
  padding: 20px 0;
  border-bottom: 1px solid #f0f0f0;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-avatar {
  flex-shrink: 0;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.comment-author {
  font-weight: 600;
  color: #333;
}

.comment-date {
  font-size: 12px;
  color: #999;
}

.comment-body {
  color: #666;
  line-height: 1.6;
  margin-bottom: 12px;
  white-space: pre-wrap;
  word-break: break-word;
}

.comment-actions {
  display: flex;
  gap: 16px;
}

.no-comments {
  text-align: center;
  padding: 40px 0;
  color: #999;
  font-size: 14px;
}

.comment-form {
  background: #f9f9f9;
  padding: 30px;
  border-radius: 8px;
  margin-top: 30px;
}

.form-title {
  font-size: 18px;
  color: #333;
  margin-bottom: 20px;
}

.captcha-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.captcha-question {
  font-size: 16px;
  color: #409eff;
  cursor: pointer;
  user-select: none;
}

:deep(.el-avatar img) {
  width: 100%;
}
</style>
