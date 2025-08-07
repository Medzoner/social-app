<!-- src/views/LoginView.vue -->
<template>
  <div class="mx-auto max-w-md p-4">
    <h2 class="mb-4 text-xl font-bold">Connexion</h2>
    <form @submit.prevent="login" class="flex flex-col gap-3" message="">
      <input
        v-model="username"
        placeholder="Username"
        class="rounded border p-2"
        autocomplete="username"
      />
      <input
        v-model="password"
        type="password"
        placeholder="Password"
        class="rounded border p-2"
        autocomplete="current-password"
      />
      <button type="submit" class="rounded bg-blue-600 p-2 text-white">Login</button>
    </form>

    <!-- Redirect Oauth Login -->
    <div class="mt-4">
      <RouterLink to="oauth-login" class="mt-4 rounded bg-red-500 p-2 text-white">
        Se connecter avec Google
      </RouterLink>
    </div>

    <!-- Message de notification -->
    <Form
      class="mt-4 rounded bg-green-100 p-2 text-green-800"
      :message="notification"
      @close="notification = ''"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Google from '@/components/auth/sso/Google.vue'
import Form from '@/components/auth/sso/Form.vue'

const username = ref('')
const password = ref('')
const notification = ref<string>('')

const router = useRouter()
const auth = useAuthStore()

onMounted(() => {
  if (auth.isAuthenticated) {
    router.push('/feed')
  }
})

const login = async (): Promise<void> => {
  try {
    const res = await axios.post('/api/login', {
      username: username.value,
      password: password.value
    })

    if (res.data.need2fa) {
      localStorage.setItem('tmp_user_id', res.data.id)
      await router.push('/verify')
      return
    }

    auth.loginWithToken(res.data)
    await router.push('/feed')
  } catch (err) {
    console.error('Erreur de connexion :', err)
    notification.value = 'Identifiants incorrects ou problème de serveur.'
  }
}
</script>

<style scoped></style>
