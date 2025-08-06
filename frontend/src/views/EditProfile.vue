<template>
  <div class="mx-auto max-w-xl p-4">
    <h2 class="mb-4 text-xl font-bold">Modifier mon profil</h2>

    <div
      v-if="notification"
      class="mt-4 rounded p-2"
      :class="
        notificationType === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
      "
    >
      {{ notification }}
    </div>
    <AvatarImage
      v-if="media && media.file_type.includes(MediaType.Image)"
      :src="'/api/' + media.file_path"
      size="32"
      alt="Avatar de l'utilisateur"
    />
    <form @submit.prevent="updateProfile" enctype="multipart/form-data" class="flex flex-col gap-4">
      <textarea
        v-model="bio"
        placeholder="Votre bio"
        class="rounded border p-2"
        rows="4"
      ></textarea>

      <input type="file" @change="handleFileUpload" class="p-2" />
      <button type="submit" class="rounded bg-blue-600 p-2 text-white">Mettre à jour</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { useProfileStore } from '@/stores/profile'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import AvatarImage from '@/components/profile/AvatarImage.vue'
import { Media, MediaType } from '@/types'

const profileStore = useProfileStore()
const router = useRouter()
const auth = useAuthStore()
const { user } = storeToRefs(auth)

const bio = ref('')
const media = ref<Media | null>(null)
const avatar = ref(null)
const notification = ref(null)
const notificationType = ref('success')
const profile = ref(null)

const showNotification = (message: string, type = 'success'): void => {
  notification.value = message
  notificationType.value = type

  setTimeout(() => {
    notification.value = null
  }, 5000)
}

const fetchProfile = async (): Promise<void> => {
  if (!user.value) return
  try {
    profile.value = await profileStore.fetchProfile(user.value.id)
    bio.value = profile.value.bio || ''
    media.value = profile.value.avatar_media || ''
  } catch (err) {
    console.error('Erreur fetch profil', err)
  }
}

const handleFileUpload = (e: any): void => {
  avatar.value = e.target.files[0]
}

const updateProfile = async (): Promise<void> => {
  try {
    let avatarUuid = ''

    if (avatar.value) {
      const formData = new FormData()
      formData.append(MediaType.Image, avatar.value)
      const res = await axios.post('/api/upload', formData, {
        headers: { ...auth.getAuthHeader() }
      })
      avatarUuid = res.data.media.uuid
    }

    const res = await axios.patch(
      `/api/users/${user.value.id}`,
      {
        bio: bio.value,
        avatar: avatarUuid
      },
      {
        headers: { ...auth.getAuthHeader() }
      }
    )
    profile.value = res.data
    media.value = profile.value.avatar_media || ''
    showNotification('Profil mis à jour avec succès ✅')
  } catch (err) {
    console.error('Erreur update profile', err)
    showNotification('Échec de la mise à jour du profil ❌', 'error')
  }
}

onMounted((): void => {
  const { isAuthenticated } = storeToRefs(auth)
  if (!isAuthenticated.value) {
    router.push('/login')
  } else {
    fetchProfile()
  }
})
</script>
