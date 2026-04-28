<template>
  <div class="todo-app">
    <h1>代办清单</h1>
    <form @submit.prevent="addTodo">
      <input type="text" v-model="newTodo" placeholder="输入新的代办事项" required />
      <button type="submit">添加</button>
    </form>
    <ul>
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
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',
  data() {
    return {
      todos: [],
      newTodo: ''
    };
  },
  mounted() {
    this.fetchTodos();
  },
  methods: {
    async fetchTodos() {
      try {
        const response = await axios.get('/api/todos');
        this.todos = response.data;
      } catch (error) {
        console.error('获取代办事项失败:', error);
      }
    },
    async addTodo() {
      if (this.newTodo.trim() === '') return;
      
      try {
        const response = await axios.post('/api/todos', {
          title: this.newTodo,
          completed: false
        });
        this.todos.push(response.data);
        this.newTodo = '';
      } catch (error) {
        console.error('添加代办事项失败:', error);
      }
    },
    async updateTodo(todo) {
      try {
        await axios.put(`/api/todos/${todo.id}`, todo);
        todo.editing = false;
      } catch (error) {
        console.error('更新代办事项失败:', error);
      }
    },
    async deleteTodo(id) {
      try {
        await axios.delete(`/api/todos/${id}`);
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
</style>