<template>
  <div class="mx-auto max-w-md p-4">
    <h2 class="mb-4 text-xl font-bold">Register</h2>
    <form @submit.prevent="register" class="flex flex-col gap-3">
      <input v-model="username" placeholder="Username" class="rounded border p-2" />
      <input v-model="email" type="email" placeholder="Email" class="rounded border p-2" />
      <input v-model="password" type="password" placeholder="Password" class="rounded border p-2" />
      <button type="submit">Register</button>
    </form>
    <button @click="signinWithGoogle" class="mt-4 rounded bg-red-500 p-2 text-white">
      Se connecter avec Google
    </button>

    <div v-if="notification" class="mt-4 rounded bg-green-100 p-2 text-green-800">
      {{ notification }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

const username = ref('')
const password = ref('')
const email = ref('')
const notification = ref('')
const router = useRouter()

const register = async (): Promise<void> => {
  try {
    await axios.post('/api/register', {
      username: username.value,
      password: password.value,
      email: email.value
    })
  } catch (error) {
    if (error.response && error.response.data) {
      notification.value = error.response.data.message || 'Registration failed'
    } else {
      notification.value = 'An unexpected error occurred'
    }
    console.error('Registration error:', error)
    return
  }
  await router.push('/login')
}

const signinWithGoogle = (): void => {
  window.location.href = '/api/oauth/google/login'
}
</script>
