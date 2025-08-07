<template>
  <ul class="mt-6 space-y-4">
    <li v-for="post in filteredPosts" :key="post.id" class="post-content rounded border bg-white p-4 shadow-sm">
      <small>
        <RouterLink :to="`/profile/${post.user.id}`" class="text-blue-600 hover:underline">
          <AvatarImage
            v-if="post?.user?.avatar_media?.file_type.includes(MediaType.Image)"
            :src="post.user.avatar_media.file_path"
            size="28"
            alt="Avatar de {{ post.user.username }}"
            class="avatar mr-2 inline-block h-7 w-7 rounded-full bg-gray-300 text-center text-xs leading-7 text-white"
          />
          <span
            v-else
            class="avatar mr-2 inline-block h-7 w-7 rounded-full bg-gray-300 text-center text-xs leading-7 text-white"
          >
            {{ post.user.username.charAt(0).toUpperCase() }}
          </span>
          #{{ post.user.username }}
        </RouterLink>
      </small>

      <div v-if="post.medias && post.medias.length" class="mt-2 space-y-2">
        <template v-for="media in post.medias" :key="media.id">
          <img
            v-if="media.file_type.includes(MediaType.Image)"
            :src="media.file_path"
            class="max-w-full rounded"
            alt=""
          />
          <video
            v-else-if="media.file_type.includes(MediaType.Video)"
            controls
            class="max-w-full rounded"
          >
            <source :src="media.file_path" :type="`video/${media.file_ext}`" />
            Your browser does not support the video tag.
          </video>
          <audio v-else-if="media.file_type.includes(MediaType.Audio)" controls class="w-full">
            <source :src="media.file_path" :type="`audio/${media.file_ext}`" />
            Your browser does not support the audio element.
          </audio>
          <!-- Optionally handle other file types -->
          <div v-else class="text-sm text-gray-500">Unsupported media type</div>
        </template>
      </div>

      <div class="text-lg">{{ post.content }}</div>
      <LikeButton
        :post-id="post.id"
        :initial-like-count="post.like_count"
        :initially-liked="post.has_like"
        @update:like="onLikeUpdate"
      />

      <CommentSection
        :post-id="post.id"
        :comment-count="post.comment_count"
        @notify="emit('notify', $event)"
      />
    </li>

    <li v-if="loading" class="py-4 text-center text-gray-500">Chargement...</li>
    <li v-if="!nextCursor && !loading" class="py-4 text-center text-gray-400">Fin du fil 🎉</li>
  </ul>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import axios from 'axios'
import qs from 'qs'
import CommentSection from './CommentSection.vue'
import LikeButton from './LikeButton.vue'
import AvatarImage from '@/components/profile/AvatarImage.vue'
import { useAuthStore } from '@/stores/auth'
import { MediaType } from '@/types'

const auth = useAuthStore()
const emit = defineEmits(['notify'])
const posts = ref([])
const loading = ref(false)
const nextCursor = ref(null)

const props = defineProps({
  profileId: {
    type: Number,
    default: null
  },
  searchQuery: {
    type: String,
    default: ''
  }
})

const fetchPosts = async () => {
  if (loading.value || nextCursor.value === false) return
  loading.value = true

  let url = props.profileId ? `/api/profile/${props.profileId}/posts` : '/api/posts'
  if (nextCursor.value) {
    url += `?cursor=${nextCursor.value}`
  }

  try {
    const res = await axios.get(url, {
      headers: { ...auth.getAuthJSONHeader() }
    })

    posts.value.push(...res.data.posts)
    nextCursor.value = res.data.next_cursor || false
  } catch (err) {
    console.error('Error fetchPosts:', err)
  } finally {
    loading.value = false
  }

  const postIds = posts.value.map((p) => p.id)
  // const countComms = async () => {
  try {
    const res = await axios.get(`/api/posts/comments/counts`, {
      headers: { ...auth.getAuthJSONHeader() },
      params: { ids: postIds },
      paramsSerializer: (params) => qs.stringify(params, { arrayFormat: 'repeat' })
    })
    res.data.counts.forEach((comm: any) => {
      const post = posts.value.find((p) => p.id === comm.post_id)
      if (post) {
        post.comment_count = comm.total
      }
    })
  } catch (err) {
    console.error(`Error fetchCommentCounts for posts: ${postIds}:`, err)
  }
  // }
  //
  // const countLikes = async () => {
  try {
    const res = await axios.get(`/api/posts/like/counts`, {
      headers: { ...auth.getAuthJSONHeader() },
      params: { ids: postIds },
      paramsSerializer: (params) => qs.stringify(params, { arrayFormat: 'repeat' })
    })
    res.data.forEach((likeData: any) => {
      const post = posts.value.find((p) => p.id === likeData.post_id)
      if (post) {
        post.like_count = likeData.total_likes
        post.has_like = likeData.liked_by_current_user
      }
    })
  } catch (err) {
    console.error(`Error fetchLikeCounts for posts: ${postIds}:`, err)
  }
  // }
  //
  // await Promise.all([countComms, countLikes])
}

const filteredPosts = computed(() =>
  posts.value.filter((post) => post.content.toLowerCase().includes(props.searchQuery.toLowerCase()))
)

const handleScroll = () => {
  const bottom = window.innerHeight + window.scrollY >= document.body.offsetHeight - 300
  if (bottom && !loading.value && nextCursor.value !== false) {
    fetchPosts()
  }
}

const onLikeUpdate = ({ postId, likeCount, hasLike, error }) => {
  if (error) {
    emit('notify', { type: 'error', message: 'Error updating like' })
    return
  }

  const post = posts.value.find((p) => p.id === postId)
  if (post) {
    post.like_count = likeCount
    post.has_like = hasLike
  }
}

onMounted(() => {
  fetchPosts()
  window.addEventListener('scroll', handleScroll)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>
