<template>
  <button @click="openGooglePopup" class="mt-4 rounded bg-red-500 p-2 text-white">
    Se connecter avec Google
  </button>
</template>

<script setup lang="ts">
import router from '@/router'
import { useAuthStore } from '@/stores/auth'
import { defineEmits } from 'vue'
const auth = useAuthStore()

const emit = defineEmits(['notify'])

const openGooglePopup = () => {
  const popup = window.open('/api/oauth/google/login', 'googleLogin', 'width=500,height=600')

  emit('notify', 'Connexion via Google...')

  const popupChecker = setInterval(() => {
    if (!popup || popup.closed) {
      clearInterval(popupChecker)
      emit('notify', 'Connection canceled or window closed.')
    }
  }, 500)

  window.addEventListener('message', (event) => {
    if (event.origin !== window.location.origin) return

    const { token } = event.data || {}
    if (token) {
      clearInterval(popupChecker)
    }
  })
}
</script>
