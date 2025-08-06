<template>
  <div class="flex items-center gap-2">
    <button
      @click="toggleLike"
      class="flex items-center gap-1 text-sm"
      :class="{ 'font-semibold text-red-600': hasLike, 'text-gray-600': !hasLike }"
    >
      <i :class="hasLike ? 'fas fa-heart' : 'far fa-heart'"></i>
      {{ hasLike ? "Je n'aime plus" : "J'aime" }}
    </button>
    <span class="text-sm text-gray-600">{{ likeCount }} J'aime</span>
  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const props = defineProps({
  postId: { type: Number, required: true },
  initialLikeCount: { type: Number, default: 0 },
  initiallyLiked: { type: Boolean, default: false }
})

const emit = defineEmits(['update:like'])

const likeCount = ref(0)
const hasLike = ref(false)

// 🔁 Sync props -> state
watchEffect(() => {
  likeCount.value = props.initialLikeCount
  hasLike.value = props.initiallyLiked
})

const toggleLike = async () => {
  const url = `/api/posts/${props.postId}/like`

  try {
    if (!hasLike.value) {
      await axios.post(url, {}, { headers: { ...auth.getAuthHeader() } })
      likeCount.value++
      hasLike.value = true
    } else {
      await axios.delete(url, { headers: { ...auth.getAuthHeader() } })
      likeCount.value = Math.max(0, likeCount.value - 1)
      hasLike.value = false
    }

    emit('update:like', {
      postId: props.postId,
      likeCount: likeCount.value,
      hasLike: hasLike.value
    })
  } catch (err) {
    console.error('Error toggleLike:', err)
    emit('update:like', { postId: props.postId, error: err })
  }
}
</script>
