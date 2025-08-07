<template>
  <div class="mx-auto max-w-md p-4">
    <h2 class="mb-4 text-xl font-bold">OAuth Login</h2>
    <p>Processing your login...</p>
  </div>
</template>

<script setup lang="ts">
import router from '@/router'
import { useAuthStore } from '@/stores/auth'
import { onMounted } from 'vue'
import axios from 'axios'
const auth = useAuthStore()

onMounted(async () => {
  if (auth.isAuthenticated) {
    await router.push('/feed')
    return
  }

  try {
    const res = await axios.get('/api/oauth/google/login')
    if (res.data.oauth_url) {
      window.location.href = res.data.oauth_url
    }
  } catch (err) {
    console.error('Error fetching redirect URL:', err)
  }
})
</script>
