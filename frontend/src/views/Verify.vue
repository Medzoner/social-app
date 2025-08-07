<template>
  <div class="mx-auto mt-10 max-w-md rounded-lg bg-white p-6 shadow">
    <h1 class="mb-4 text-center text-2xl font-bold">Vérification de votre compte</h1>

    <div class="mb-6 flex justify-center gap-4">
      <button :class="mode === 'email' ? activeBtn : inactiveBtn" @click="mode = 'email'">
        Par Email
      </button>
      <button :class="mode === 'phone' ? activeBtn : inactiveBtn" @click="mode = 'phone'">
        Par Téléphone
      </button>
    </div>

    <div v-if="successMessage" class="mb-4 text-center text-green-600">
      {{ successMessage }}
    </div>

    <div class="mb-4">
      <label class="mb-1 block text-sm font-medium">
        Code de vérification {{ mode === 'email' ? 'email' : 'SMS' }}
      </label>
      <input v-model="code" class="w-full rounded border p-2" placeholder="Entrez le code reçu" />
    </div>

    <div class="flex items-center justify-between">
      <button
        class="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
        @click="verifyCode"
        :disabled="loading"
      >
        Vérifier
      </button>
      <button class="text-sm text-blue-600 hover:underline" @click="resendCode" :disabled="loading">
        Renvoyer le code
      </button>
    </div>

    <div v-if="errorMessage" class="mt-4 text-center text-red-600">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const mode = ref('email')
const code = ref('')
const loading = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const tmpUserIdKey = 'tmp_user_id'

const activeBtn = 'px-4 py-2 bg-blue-600 text-white rounded'
const inactiveBtn = 'px-4 py-2 bg-gray-200 text-gray-600 rounded'

const verifyCode = async (): Promise<void> => {
  if (!code.value.trim()) return

  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const endpoint = '/api/profile/verify'
    const res = await axios.post(endpoint, {
      code: code.value,
      id: parseInt(localStorage.getItem(tmpUserIdKey)) || auth.user.id,
      type: mode.value === 'email' ? 'email' : 'phone'
    })
    successMessage.value = `✅ ${mode.value === 'email' ? 'Email' : 'Téléphone'} vérifié avec succès`

    if (!res.data.access_token) {
      await auth.refreshToken()
      return
    }

    auth.loginWithToken(res.data)
    setTimeout(async () => {
      await router.push('/')
    }, 2000)
    return
  } catch (err) {
    errorMessage.value = 'Code invalide ou expiré.'
  } finally {
    loading.value = false
  }
}

const resendCode = async (): Promise<void> => {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const endpoint = '/api/profile/request-code'
    await axios.post(endpoint, {
      id: parseInt(localStorage.getItem(tmpUserIdKey)),
      type: mode.value === 'email' ? 'email' : 'phone'
    })
    successMessage.value = `Le code de vérification a été renvoyé par ${mode.value === 'email' ? 'email' : 'SMS'} 📩`
  } catch (err) {
    errorMessage.value = `Erreur lors de l'envoi du code ${mode.value === 'email' ? 'email' : 'SMS'}`
  } finally {
    loading.value = false
  }
}
</script>
