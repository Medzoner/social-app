<template>
  <div class="mt-4">
    <h4 class="mb-2 font-semibold">Commentaires</h4>

    <div v-if="!show">
      <button
        @click="load(nextCursor)"
        class="rounded bg-gray-100 px-3 py-1 text-sm hover:bg-gray-200"
      >
        Afficher les commentaires ({{ commentCount }})
      </button>
    </div>
    <div v-else>
      <button @click="toggle" class="mb-2 rounded bg-gray-100 px-3 py-1 text-sm hover:bg-gray-200">
        Masquer les commentaires ({{ commentCount }})
      </button>
    </div>

    <div v-if="loading" class="flex items-center gap-2 text-sm text-gray-500">
      <svg class="h-4 w-4 animate-spin text-gray-500" viewBox="0 0 24 24" fill="none">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z" />
      </svg>
      Chargement...
    </div>

    <ul v-if="show" class="space-y-1">
      <li v-for="comment in comments" :key="comment.id">
        <small>#{{ comment.user.username }}</small> : {{ comment.content }}
      </li>
    </ul>

    <div v-if="nextCursor" class="mt-2">
      <button @click="load(nextCursor)" class="text-sm text-blue-600 hover:underline">
        Voir plus de commentaires
      </button>
    </div>

    <input
      v-model="input"
      placeholder="Ajouter un commentaire"
      class="mt-2 w-full rounded border p-1"
    />
    <button @click="send" class="mt-1 rounded bg-gray-200 px-2 py-1 text-sm">Envoyer</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const props = defineProps({
  postId: {
    type: Number,
    required: true
  },
  commentCount: {
    type: Number,
    default: 0
  }
})
const emit = defineEmits(['notify'])

const show = ref(false)
const loading = ref(false)
const nextCursor = ref(null)
const comments = ref([])
const input = ref('')

const load = async (cursor = null) => {
  if (loading.value) return
  loading.value = true

  const url = `/api/posts/${props.postId}/comments` + (cursor ? `?cursor=${cursor}` : '')
  const res = await axios.get(url, {
    headers: { ...auth.getAuthHeader() }
  })

  comments.value.push(...res.data.comments)
  nextCursor.value = res.data.next_cursor
  show.value = true
  loading.value = false
}

const toggle = () => {
  show.value = !show.value
  comments.value = []
}

const send = async () => {
  await axios.post(
    `/api/posts/${props.postId}/comments`,
    { content: input.value },
    { headers: { ...auth.getAuthHeader() } }
  )
  emit('notify', 'Comment added!')
  input.value = ''
  comments.value = []
  nextCursor.value = null
  await load()
}
</script>
