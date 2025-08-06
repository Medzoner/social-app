<template>
  <div>
    <h2>Admin Dashboard</h2>
    <ul>
      <li v-for="user in users" :key="user.id">
        <strong>{{ user.username }}</strong> - {{ user.role }} - {{ user.bio }}
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const users = ref([])

const fetchUsers = async (): Promise<void> => {
  const token = auth.getAccessToken()
  if (!token) return

  const payload = JSON.parse(atob(token.split('.')[1]))
  if (payload.role !== 'admin') return alert('Access denied')

  const res = await axios.get('/api/admin/users')
  users.value = res.data
}

onMounted(fetchUsers)
</script>
