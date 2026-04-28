<template>
  <div class="todo-app">
    <h1>代办清单</h1>
    
    <!-- 用户选择 -->
    <div class="user-selector">
      <label>选择用户：</label>
      <select v-model="selectedUserId" @change="fetchTodos">
        <option value="">-- 请选择用户 --</option>
        <option v-for="user in users" :key="user.id" :value="user.id">
          {{ user.username }}
        </option>
      </select>
      <button @click="showUserForm = !showUserForm">
        {{ showUserForm ? '取消' : '添加用户' }}
      </button>
    </div>
    
    <!-- 添加用户表单 -->
    <div v-if="showUserForm" class="user-form">
      <h2>添加用户</h2>
      <form @submit.prevent="addUser">
        <input type="text" v-model="newUser.username" placeholder="用户名" required />
        <input type="password" v-model="newUser.password" placeholder="密码" required />
        <textarea v-model="newUser.bio" placeholder="简介"></textarea>
        <button type="submit">添加</button>
      </form>
    </div>
    
    <!-- TODO表单 -->
    <form @submit.prevent="addTodo" v-if="selectedUserId">
      <input type="text" v-model="newTodo" placeholder="输入新的代办事项" required />
      <button type="submit">添加</button>
    </form>
    
    <!-- TODO列表 -->
    <ul v-if="selectedUserId">
      <li v-for="todo in todos" :key="todo.id" :class="{ completed: todo.completed }">
        <div v-if="!todo.editing">
          <input type="checkbox" v-model="todo.completed" @change="updateTodo(todo)" />
          <span>{{ todo.title }}</span>
          <button class="edit" @click="todo.editing = true">编辑</button>
          <button @click="deleteTodo(todo.id)">删除</button>
        </div>
        <div v-else class="edit-form">
          <input type="text" v-model="todo.title" @keyup.enter="updateTodo(todo)" @keyup.esc="todo.editing = false" />
          <button @click="updateTodo(todo)">保存</button>
          <button class="cancel" @click="todo.editing = false">取消</button>
        </div>
      </li>
    </ul>
    
    <!-- 提示信息 -->
    <div v-if="!selectedUserId" class="info">
      请选择一个用户来查看和管理代办事项
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',
  data() {
    return {
      todos: [],
      users: [],
      newTodo: '',
      newUser: {
        username: '',
        password: '',
        bio: ''
      },
      selectedUserId: '',
      showUserForm: false
    };
  },
  mounted() {
    this.fetchUsers();
  },
  methods: {
    // 用户相关方法
    async fetchUsers() {
      try {
        const response = await axios.get('/api/v1/users');
        this.users = response.data.data;
      } catch (error) {
        console.error('获取用户列表失败:', error);
      }
    },
    async addUser() {
      if (this.newUser.username.trim() === '' || this.newUser.password.trim() === '') return;
      
      try {
        const response = await axios.post('/api/v1/users', this.newUser);
        this.users.push(response.data.data);
        this.selectedUserId = response.data.data.id;
        this.newUser = {
          username: '',
          password: '',
          bio: ''
        };
        this.showUserForm = false;
        this.fetchTodos();
      } catch (error) {
        console.error('添加用户失败:', error);
      }
    },
    
    // TODO相关方法
    async fetchTodos() {
      if (!this.selectedUserId) {
        this.todos = [];
        return;
      }
      
      try {
        const response = await axios.get(`/api/v1/todos?uid=${this.selectedUserId}`);
        this.todos = response.data.data.map(todo => ({
          ...todo,
          completed: todo.status === 1
        }));
      } catch (error) {
        console.error('获取代办事项失败:', error);
      }
    },
    async addTodo() {
      if (this.newTodo.trim() === '' || !this.selectedUserId) return;
      
      try {
        const response = await axios.post('/api/v1/todos', {
          title: this.newTodo,
          status: 0,
          uid: parseInt(this.selectedUserId)
        });
        this.todos.push({
          ...response.data.data,
          completed: response.data.data.status === 1
        });
        this.newTodo = '';
      } catch (error) {
        console.error('添加代办事项失败:', error);
      }
    },
    async updateTodo(todo) {
      try {
        await axios.put(`/api/v1/todos/${todo.id}`, {
          title: todo.title,
          status: todo.completed ? 1 : 0,
          uid: parseInt(this.selectedUserId)
        });
        todo.editing = false;
      } catch (error) {
        console.error('更新代办事项失败:', error);
      }
    },
    async deleteTodo(id) {
      try {
        await axios.delete(`/api/v1/todos/${id}`);
        this.todos = this.todos.filter(todo => todo.id !== id);
      } catch (error) {
        console.error('删除代办事项失败:', error);
      }
    }
  }
};
</script>

<style scoped>
.todo-app {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.user-selector {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-selector select {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  flex: 1;
}

.user-form {
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: #f9f9f9;
}

.user-form h2 {
  margin-bottom: 10px;
  font-size: 18px;
}

.user-form input,
.user-form textarea {
  width: 100%;
  padding: 8px;
  margin-bottom: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}

.user-form textarea {
  resize: vertical;
  min-height: 80px;
}

.info {
  text-align: center;
  padding: 20px;
  color: #666;
  background-color: #f9f9f9;
  border-radius: 4px;
  margin-top: 20px;
}
</style>