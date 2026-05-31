<template>
  <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
    <el-form-item label="标题" prop="title">
      <el-input v-model="form.title" placeholder="请输入文章标题" />
    </el-form-item>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="分类">
          <el-select v-model="form.category_id" placeholder="请选择分类" clearable style="width: 100%">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="标签">
          <el-select v-model="form.tag_ids" multiple placeholder="请选择标签" style="width: 100%">
            <el-option v-for="tag in tags" :key="tag.id" :label="tag.name" :value="tag.id" />
          </el-select>
        </el-form-item>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="发布时间">
          <el-date-picker
            v-model="form.published_at"
            type="datetime"
            placeholder="选择发布时间"
            style="width: 100%"
          />
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="访问密码">
          <el-input v-model="form.password" placeholder="留空则不需要密码" />
        </el-form-item>
      </el-col>
    </el-row>

    <el-form-item>
      <el-checkbox v-model="form.is_top">置顶文章</el-checkbox>
    </el-form-item>

    <el-form-item label="内容">
      <div class="editor-container">
        <Toolbar
          :editor="editorRef"
          :defaultConfig="toolbarConfig"
          mode="default"
          style="border-bottom: 1px solid #ccc"
        />
        <Editor
          v-model="form.content"
          :defaultConfig="editorConfig"
          mode="default"
          style="height: 500px; overflow-y: hidden"
          @on-created="handleCreated"
        />
      </div>
    </el-form-item>

    <el-form-item label="状态">
      <el-radio-group v-model="form.status">
        <el-radio :label="0">草稿</el-radio>
        <el-radio :label="1">发布</el-radio>
      </el-radio-group>
    </el-form-item>
  </el-form>
</template>

<script setup>
import { ref } from 'vue'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'

const props = defineProps({
  form: {
    type: Object,
    required: true
  },
  rules: {
    type: Object,
    required: true
  },
  categories: {
    type: Array,
    required: true
  },
  tags: {
    type: Array,
    required: true
  },
  formRef: {
    type: Object,
    required: true
  },
  editorConfig: {
    type: Object,
    required: true
  },
  toolbarConfig: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['editor-created'])

const editorRef = ref()

function handleCreated(editor) {
  editorRef.value = editor
  emit('editor-created', editor)
}
</script>

<style src="@wangeditor/editor/dist/css/style.css"></style>

<style scoped>
.editor-container {
  border: 1px solid #ccc;
  border-radius: 4px;
}
</style>
