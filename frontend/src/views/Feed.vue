<template>
  <div class="mx-auto max-w-2xl p-4">
    <h2 class="mb-4 text-xl font-bold">Feed</h2>

    <input
      v-model="searchQuery"
      name="search"
      placeholder="Search posts..."
      class="mb-4 w-full rounded border p-2"
    />

    <button @click="showModal = true" class="mb-4 rounded bg-blue-600 px-4 py-2 text-white">
      Créer un post
    </button>

    <BaseModal :show="showModal" @close="closeModal">
      <template #header>
        <h3 class="mb-4 text-lg font-semibold">Créer un post</h3>
      </template>

      <form @submit.prevent="createPost" enctype="multipart/form-data" class="space-y-3">
        <textarea
          v-model="content"
          name="content"
          placeholder="Quoi de neuf ?"
          class="w-full resize-none rounded border p-2"
          rows="4"
        ></textarea>

        <div
          class="w-full cursor-pointer rounded border-2 border-dashed border-gray-300 p-4 text-center hover:bg-gray-50"
          @dragover.prevent
          @drop.prevent="onDrop"
          @click="triggerFileInput"
        >
          Glissez et déposez une image, vidéo ou audio ici, ou cliquez pour choisir un fichier.
          <input
            type="file"
            class="hidden"
            ref="fileInput"
            @change="handleFileUpload"
            accept="image/*,video/*,audio/*"
            multiple
          />
        </div>

        <div v-if="uploadProgress !== null" class="mt-2 text-sm text-gray-700">
          Upload en cours : {{ uploadProgress }}%
          <div class="mt-1 h-2 rounded bg-gray-200">
            <div class="h-full rounded bg-blue-500" :style="{ width: uploadProgress + '%' }"></div>
          </div>
        </div>

        <ImagePreview :files="files" />

        <button name="publish" type="submit" class="w-full rounded bg-blue-600 py-2 text-white">
          Publier
        </button>
      </form>
    </BaseModal>

    <!-- ✅ Notification -->
    <div id="notification" v-if="notification" class="mt-4 rounded bg-green-100 p-2 text-green-800">
      {{ notification }}
    </div>

    <PostComponent @notify="notification = $event" :searchQuery="searchQuery" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import PostComponent from '@/components/post/PostComponent.vue'
import { useAuthStore } from '@/stores/auth'
import BaseModal from '@/components/post/BaseModal.vue'
import ImagePreview from '@/components/ui/ImagePreview.vue'
import { MediaTypes, MediaType } from '@/types'

const auth = useAuthStore()
const content = ref('')
const files = ref<File[]>([])
const previewFiles = ref<
  { url: string; type: MediaType.Image | MediaType.Video | MediaType.Audio | 'other' }[]
>([])
const searchQuery = ref('')
const notification = ref('')
const showModal = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const uploadProgress = ref<number | null>(null)

const onDrop = (e: DragEvent) => {
  const dropped = Array.from(e.dataTransfer?.files || [])
  handleSelectedFiles(dropped)
}

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileUpload = (e: Event): void => {
  const selected = Array.from((e.target as HTMLInputElement).files || [])
  handleSelectedFiles(selected)
}

const handleSelectedFiles = (selected: File[]) => {
  const accepted = selected.filter((file) => MediaTypes.some((type) => file.type.startsWith(type)))
  files.value.push(...accepted)
  accepted.forEach((file) => {
    const url = URL.createObjectURL(file)
    const type = file.type.startsWith(MediaType.Image)
      ? MediaType.Image
      : file.type.startsWith(MediaType.Video)
        ? MediaType.Video
        : file.type.startsWith(MediaType.Audio)
          ? MediaType.Audio
          : 'other'
    previewFiles.value.push({ url, type })
  })
}

const createPost = async (): Promise<void> => {
  let mediaUuids: string[] = []

  if (files.value.length > 0) {
    const formData = new FormData()
    files.value.forEach((file) => {
      formData.append('images', file)
    })

    const uploadRes = await axios.post('/api/upload', formData, {
      headers: { ...auth.getAuthHeader() },
      onUploadProgress: (e) => {
        if (e.total) {
          uploadProgress.value = Math.round((e.loaded / e.total) * 100)
        }
      }
    })

    uploadProgress.value = null
    mediaUuids = uploadRes.data.medias.map((media: { uuid: string }) => media.uuid)
  }

  const res = await axios.post(
    '/api/posts',
    { content: content.value, media_uuids: mediaUuids },
    { headers: { ...auth.getAuthJSONHeader() } }
  )

  if (res.status !== 201) {
    notification.value = 'Erreur lors de la création du post.'
    return
  }

  closeModal()
  notification.value = 'Post publié !'
}

const closeModal = () => {
  showModal.value = false
  content.value = ''
  files.value = []
  previewFiles.value = []
  uploadProgress.value = null
}
</script>
