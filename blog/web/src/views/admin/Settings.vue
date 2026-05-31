<template>
  <div class="settings">
    <div class="header">
      <h1>设置</h1>
      <div class="header-right">
        <router-link to="/admin" class="nav-link">返回</router-link>
      </div>
    </div>
    <div class="container">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="修改密码" name="password">
          <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
            <el-form-item label="旧密码" prop="oldPassword">
              <el-input
                v-model="passwordForm.oldPassword"
                type="password"
                placeholder="请输入旧密码"
              />
            </el-form-item>
            <el-form-item label="新密码" prop="newPassword">
              <el-input
                v-model="passwordForm.newPassword"
                type="password"
                placeholder="请输入新密码（至少6位）"
              />
            </el-form-item>
            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input
                v-model="passwordForm.confirmPassword"
                type="password"
                placeholder="请再次输入新密码"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="passwordLoading" @click="handleSavePassword">保存</el-button>
              <el-button @click="goBack">取消</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="评论设置" name="comment">
          <el-form :model="commentSettings" label-width="200px">
            <el-form-item label="评论需要审核">
              <el-switch
                v-model="commentSettings.comment_moderation"
                @change="handleSaveCommentSettings"
              />
              <div class="form-tip">开启后，用户提交的评论需要管理员审核后才能显示</div>
            </el-form-item>
            <el-form-item label="开启验证码">
              <el-switch
                v-model="commentSettings.captcha_enabled"
                @change="handleSaveCommentSettings"
              />
              <div class="form-tip">开启后，用户发表评论需要输入验证码</div>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)
const activeTab = ref('password')

const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const commentSettings = ref({
  comment_moderation: true,
  captcha_enabled: true
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== passwordForm.value.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  oldPassword: [{ required: true, message: '请输入旧密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

async function fetchCommentSettings() {
  try {
    const res = await request.get('/admin/settings')
    commentSettings.value = res.data
  } catch (error) {
    console.error('Failed to fetch settings:', error)
  }
}

async function handleSaveCommentSettings() {
  try {
    await request.put('/admin/settings', commentSettings.value)
    ElMessage.success('设置已保存')
  } catch (error) {
    ElMessage.error('保存失败')
    fetchCommentSettings()
  }
}

async function handleSavePassword() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await request.put('/admin/password', {
          old_password: passwordForm.value.oldPassword,
          new_password: passwordForm.value.newPassword
        })
        ElMessage.success('密码修改成功')
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
  fetchCommentSettings()
})
</script>

<style scoped>
.settings {
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
  max-width: 600px;
  margin: 40px auto;
  padding: 0 20px;
  background: white;
  padding: 40px;
  border-radius: 8px;
}

.form-tip {
  font-size: 12px;
  color: #999;
  margin-left: 16px;
  margin-top: 4px;
}
</style>
