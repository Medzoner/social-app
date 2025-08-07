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
import axios from 'axios'

const router = useRouter()
const auth = useAuthStore()

onMounted(async () => {
  if (auth.isAuthenticated) {
    await router.push('/feed')
    return
  }
  const code = new URLSearchParams(window.location.search).get('code')

  let cb = {
    access_token: '',
    id_token: '',
    refresh_token: '',
    verified: false,
    id: ''
  }
  try {
    const res = await axios.get(`/api/oauth/google/callback?code=${code}`)
    cb = res.data || cb
  } catch (err) {
    console.error('Error fetching redirect URL:', err)
  }

  if (!cb.verified && !cb.access_token) {
    localStorage.setItem('tmp_user_id', cb.id)
    await router.push('/verify')
    return
  }

  if (cb.access_token) {
    const jwt: JwtToken = {
      access_token: cb.access_token || '',
      refresh_token: cb.refresh_token || '',
      id_token: cb.id_token || ''
    }

    auth.loginWithToken(jwt)

    await router.push('/feed')
  }
})
</script>
