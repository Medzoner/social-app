<!-- components/ui/ImagePreview.vue -->
<template>
  <div v-if="urls.length" class="mt-4 grid grid-cols-1 gap-4 sm:grid-cols-2">
    <div v-for="(url, index) in urls" :key="index" class="rounded border bg-gray-50 p-2 shadow-sm">
      <img
        v-if="files[index].type.startsWith(MediaType.Image)"
        :src="url"
        class="mx-auto max-h-64 rounded"
        alt="image preview"
      />
      <video
        v-else-if="files[index].type.startsWith(MediaType.Video)"
        controls
        class="mx-auto max-w-full"
      >
        <source :src="url" />
        Votre navigateur ne prend pas en charge la vidéo.
      </video>
      <audio v-else-if="files[index].type.startsWith(MediaType.Audio)" controls class="w-full">
        <source :src="url" />
        Votre navigateur ne prend pas en charge l'audio.
      </audio>
      <p v-else class="text-sm text-gray-500">Type non pris en charge</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, onBeforeUnmount, watch, ref } from 'vue'
import { MediaType } from '@/types'

const props = defineProps<{ files: File[] }>()
const urls = ref<string[]>([])

const revokeAll = () => {
  urls.value.forEach((url) => URL.revokeObjectURL(url))
  urls.value = []
}

// Re-génère les URLs à chaque changement de fichiers
watch(
  () => props.files,
  (newFiles) => {
    revokeAll()
    urls.value = newFiles.map((file) => URL.createObjectURL(file))
  },
  { immediate: true, deep: true }
)

onBeforeUnmount(() => {
  revokeAll()
})
</script>
