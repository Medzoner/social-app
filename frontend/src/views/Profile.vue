<template>
  <div v-if="profile" class="mx-auto max-w-3xl px-4 py-8">
    <div class="flex flex-col items-center gap-4 text-center">
      <AvatarImage
        v-if="media && media.file_type.includes(MediaType.Image)"
        :src="media.file_path"
        size="96"
        alt="Avatar de l'utilisateur"
      />

      <div>
        <h2 class="text-2xl font-semibold text-gray-800">@{{ profile.username }}</h2>
        <p class="mt-1 max-w-md text-gray-600">{{ profile.bio }}</p>
      </div>
    </div>

    <hr class="my-8 border-gray-300" />

    <div>
      <div>
        <!-- Menu -->
        <div class="mb-6 border-b border-gray-200">
          <nav class="flex gap-6 text-gray-600">
            <span class="cursor-default border-b-2 border-blue-500 pb-2 text-blue-600">
              Publications
            </span>

            <RouterLink
              :to="`/chat/${profile.id}`"
              class="border-b-2 border-transparent pb-2 hover:border-blue-300 hover:text-blue-600"
            >
              Chat
            </RouterLink>
          </nav>
        </div>

        <!-- Publications -->
        <div>
          <h3 class="mb-4 text-lg font-medium text-gray-700">Ses publications</h3>
          <PostComponent
            @notify="notification = $event"
            :profileId="Number(profile.id)"
            :searchQuery="searchQuery"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRoute, useRouter } from 'vue-router'
import PostComponent from '../components/post/PostComponent.vue'
import AvatarImage from '../components/profile/AvatarImage.vue'
import { useAuthStore } from '@/stores/auth'
import { useProfileStore } from '@/stores/profile'
import { Profile, Post, Media, MediaType } from '@/types'

const auth = useAuthStore()
const profileStore = useProfileStore()
const route = useRoute()
const router = useRouter()

const profile = ref<Profile | null>(null)
const searchQuery = ref<string>('')
const notification = ref<string>('')
const media = ref<Media | null>(null)

const fetchProfile = async (): Promise<void> => {
  try {
    const data = await profileStore.fetchProfile(parseInt(route?.params?.id as string))
    profile.value = data
    media.value = data.avatar_media || ''
  } catch (error) {
    console.error('Erreur lors de la récupération du profil:', error)
    await router.push('/404')
    return
  }

  const posts = await axios.get<Post[]>(`/api/profile/${route.params.id}/posts`, {
    headers: { ...auth.getAuthJSONHeader() }
  })

  if (profile.value) {
    profile.value.posts = posts.data
  }
}

onMounted(fetchProfile)
</script>
