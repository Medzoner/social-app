<template>
  <div class="mx-auto max-w-md p-4">
    <h2 class="mb-4 text-xl font-bold">OAuth Callback</h2>
    <p>Processing your login...</p>
  </div>
</template>
<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { JwtToken } from '@/types'

const router = useRouter()
const auth = useAuthStore()

onMounted(() => {
  const token = new URLSearchParams(window.location.search).get('access_token')
  const refresh = new URLSearchParams(window.location.search).get('refresh_token')
  const idToken = new URLSearchParams(window.location.search).get('id_token')
  if (token) {
    console.log('OAuth Callback:', token)
    const jwt: JwtToken = {
      access_token: token || '',
      refresh_token: refresh || '',
      id_token: idToken || ''
    }

    auth.loginWithToken(jwt)

    window.close()
  }
})
</script>
