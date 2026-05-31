<template>
  <div class="settings">
    <div class="header">
      <h1>设置</h1>
      <div class="header-right">
        <router-link to="/admin" class="nav-link">返回</router-link>
      </div>
    </div>
    <div class="container">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input
            v-model="form.oldPassword"
            type="password"
            placeholder="请输入旧密码"
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="form.newPassword"
            type="password"
            placeholder="请输入新密码（至少6位）"
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
          />
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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const formRef = ref(null)
const loading = ref(false)

const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.value.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
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

async function handleSave() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await request.put('/admin/password', {
          old_password: form.value.oldPassword,
          new_password: form.value.newPassword
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
</style>
