<template>
  <div class="mx-auto max-w-xl p-4">
    <h2 class="mb-4 text-xl font-bold">Notifications</h2>
    <ul>
      <li v-for="n in notifications" :key="n.id" class="mb-2">
        <span :class="{ 'font-bold': !n.is_read }">{{ n.content }}</span>
        <button @click="markRead(n.id)" class="ml-2 text-sm text-blue-600 underline">
          Mark as read
        </button>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const notifications = ref([])

onMounted(async (): Promise<void> => {
  const token = auth.getAccessToken()
  const payload = JSON.parse(atob(token.split('.')[1]))
  const res = await axios.get('/api/notifications?user_id=' + payload.sub)
  notifications.value = res.data
})

const markRead = async (id: number): Promise<void> => {
  await axios.post('/api/notifications/mark-read', { id })
  const n = notifications.value.find((n) => n.id === id)
  if (n) n.is_read = true
}
</script>
