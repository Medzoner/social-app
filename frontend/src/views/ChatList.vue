<template>
  <div class="mx-auto max-w-md p-4">
    <h2 class="mb-4 text-2xl font-bold">Mes conversations</h2>

    <ul v-if="chats.length > 0" class="divide-y divide-gray-200">
      <li
        v-for="chat in chats"
        :key="chat.user.id"
        @click="goToChat(chat.user.id)"
        class="cursor-pointer rounded p-3 transition hover:bg-gray-100"
      >
        <RouterLink :to="`/chat/${chat.user.id}`" class="flex items-center space-x-3">
          <div class="font-semibold text-gray-800">{{ chat.user.username }}</div>
          <div class="truncate text-sm text-gray-600">{{ chat.content }}</div>
          <div class="text-xs text-gray-400">{{ formatDate(chat.created_at) }}</div>
        </RouterLink>
      </li>
    </ul>

    <p v-else class="text-gray-500">Aucune conversation encore.</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const chats = ref([])
const router = useRouter()

const fetchChats = async (): Promise<void> => {
  try {
    const res = await axios.get('/api/chats', {
      headers: { ...auth.getAuthJSONHeader() }
    })
    chats.value = res.data
  } catch (err) {
    console.error('Erreur récupération chats:', err)
  }
}

const formatDate = (dateStr: string): string => {
  const d = new Date(dateStr)
  return d.toLocaleString()
}

const goToChat = (userId: string): void => {
  router.push(`/chat/${userId}`)
}

onMounted(fetchChats)
</script>
