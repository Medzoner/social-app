<template>
  <div class="mx-auto max-w-xl p-4">
    <div v-if="profile" class="mb-6 flex items-center gap-4">
      <AvatarImage
        v-if="profile.avatar_media?.file_type.includes(MediaType.Image)"
        :src="profile.avatar_media.file_path"
        size="28"
        :online="isOnline"
        alt="Avatar de {{ profile.username }}"
      />
      <div>
        <h2 class="text-xl font-bold text-gray-800">Discussion avec {{ profile.username }}</h2>
      </div>
    </div>

    <div v-if="error" class="mb-4 text-red-600">{{ error }}</div>

    <div
      class="mb-2 h-64 overflow-y-scroll rounded border bg-white p-2 shadow"
      ref="chatBox"
      @scroll="handleScroll"
    >
      <transition-group name="fade" tag="div">
        <div v-for="group in groupedMessages" :key="group.date" class="mb-6">
          <div class="mb-4 text-center text-xs text-gray-500">
            {{ group.date }}
          </div>

          <div
            v-for="m in group.messages"
            :key="m.id"
            class="mb-2 flex gap-2"
            :class="m.sender_id === myId ? 'justify-end' : 'justify-start'"
          >
            <!-- Avatar -->
            <AvatarImage
              v-if="
                m.sender_id !== myId && profile?.avatar_media?.file_type.includes(MediaType.Image)
              "
              :src="profile.avatar_media.file_path"
              size="24"
              class="mt-1"
            />

            <!-- Bulle -->
            <div
              :class="[
                'max-w-[75%] rounded-xl px-4 py-2 text-sm shadow',
                m.sender_id === myId
                  ? 'rounded-br-none bg-green-100 text-right text-green-900'
                  : 'rounded-bl-none border bg-white text-left text-gray-900'
              ]"
            >
              <div
                class="whitespace-pre-line break-words"
                :data-user-id="m.user_id"
                :data-my-id="myId"
              >
                {{ m.content }}
              </div>
              <div class="mt-1 flex items-center justify-end gap-1 text-xs text-gray-500">
                {{
                  new Date(m.created_at).toLocaleTimeString([], {
                    hour: '2-digit',
                    minute: '2-digit'
                  })
                }}
                <template v-if="m.sender_id === myId">
                  <i
                    class="text-sm"
                    :class="{
                      'fas fa-check': !m.read,
                      'fas fa-check-double text-blue-500': m.read
                    }"
                  ></i>

                  <span v-if="!m.read">⌛</span>
                  <span v-else class="text-blue-500">✓</span>
                </template>
              </div>
            </div>
          </div>
        </div>
      </transition-group>
      <div v-if="!messages || messages.length === 0" class="text-center text-gray-500">
        Aucun message pour le moment.
      </div>
      <p class="text-sm text-gray-500">
        {{ isTyping ? '...' : 'En ligne récemment' }}
      </p>
    </div>

    <form @submit.prevent="send" class="flex gap-2">
      <input
        ref="input"
        @input="onInput"
        v-model="content"
        placeholder="Écrire un message..."
        class="w-full rounded border p-2"
        autocomplete="off"
      />
      <button type="submit" class="rounded bg-blue-600 px-4 py-2 text-white">Envoyer</button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import AvatarImage from '@/components/profile/AvatarImage.vue'
import { useProfileStore } from '@/stores/profile'
import { useAuthStore } from '@/stores/auth'
import { useChatStore } from '@/stores/chat'
import { MediaType, Message, MessageGoup } from '@/types'
const isOnline = ref(false)

const chatStore = useChatStore()

const messages = computed(() => chatStore.messages)
const isTyping = computed(() => chatStore.typing)

const auth = useAuthStore()
const profileStore = useProfileStore()

const route = useRoute()
const to: number = parseInt(route.params.id as string)
const content = ref('')
const chatBox = ref(null)
const input = ref(null)
const myId = ref(null)
const profile = ref(null)
const error = ref(null)
const loading = ref(false)
const nextCursor = ref(null)
let typingInterval = null
let timeout = null

const notifyTyping = (): void => {
  if (!window.ws || window.ws.readyState !== 1) return

  window.ws.send(
    JSON.stringify({
      type: 'typing',
      sender_id: myId.value,
      receiver_id: to
    })
  )
}

const onInput = (): void => {
  if (typingInterval) return
  notifyTyping()
  typingInterval = setTimeout(() => {
    typingInterval = null
  }, 3500)
}

const fetchProfile = async (): Promise<void> => {
  try {
    profile.value = await profileStore.fetchProfile(to)
  } catch (err) {
    error.value = 'Error fetching profile.'
  }
}

const sortedMessages = computed((): any[] => {
  return [...messages.value].sort(
    (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
  )
})

const fetchMessages = async (): Promise<void> => {
  loading.value = true
  try {
    const res = await axios.get(`/api/messages/${to}`, {
      headers: { ...auth.getAuthJSONHeader() }
    })
    nextCursor.value = res.data.next_cursor || false
    chatStore.setMessages(res.data.messages)
  } catch (err) {
    error.value = 'Error loading messages.'
  } finally {
    loading.value = false
  }
}

const fetchPreviousMessages = async (): Promise<void> => {
  if (loading.value || nextCursor.value === false) return

  const el = chatBox.value
  const previousScrollHeight = el.scrollHeight

  loading.value = true
  let url = `/api/messages/${to}`
  if (nextCursor.value) {
    url += `?cursor=${nextCursor.value}`
  }

  try {
    const res = await axios.get(url, {
      headers: { ...auth.getAuthJSONHeader() }
    })
    chatStore.setMessages([...res.data.messages, ...messages.value])
    nextCursor.value = res.data.next_cursor || false

    await nextTick()
    const newScrollHeight = el.scrollHeight
    el.scrollTop = newScrollHeight - previousScrollHeight
  } catch (err) {
    error.value = 'Error loading previous messages.'
  } finally {
    loading.value = false
  }
}

const send = async (): Promise<void> => {
  const trimmed = content.value.trim()
  if (!trimmed) return

  const tempId = Date.now()

  const msg: Message = {
    id: tempId,
    created_at: new Date().toISOString(),
    content: trimmed,
    user_id: myId.value,
    sender_id: myId.value,
    receiver_id: to,
    read: false,
    timestamp: new Date(),
    error: false
  }
  chatStore.receiveMessage(msg)

  scrollToBottom()
  content.value = ''

  try {
    const { data } = await axios.post(
      '/api/messages',
      {
        to: to,
        content: trimmed
      },
      {
        headers: { ...auth.getAuthJSONHeader() }
      }
    )

    const index = messages.value.findIndex((m) => m.id === tempId)
    if (index !== -1) {
      messages.value[index] = { ...messages.value[index], id: data.id, created_at: data.created_at }
    }
  } catch (err) {
    console.error(err)
    error.value = 'Error sending message.'

    const index = messages.value.findIndex((m) => m.id === tempId)
    if (index !== -1) {
      messages.value[index].error = true
    }
  }
}

const handleScroll = (e): void => {
  const el = e.target
  const top = el.scrollTop <= 100
  const atBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 10
  if (top && !loading.value && nextCursor.value !== false) {
    fetchPreviousMessages()
  }

  if (atBottom) {
    notifySeen()
  }
}

const groupedMessages = computed(() => {
  const groups = {}

  for (const msg of sortedMessages.value) {
    const dateKey = new Date(msg.created_at).toDateString()
    if (!groups[dateKey]) groups[dateKey] = []
    groups[dateKey] = [...groups[dateKey], msg]
  }

  return Object.entries(groups).map(([dateKey, msgs]) => {
    const group: MessageGoup = {
      date: formatDateHeader(msgs[0].created_at),
      messages: msgs as Message[],
      id: parseInt(dateKey)
    }

    return group
  })
})

const formatDateHeader = (isoDate: string): string => {
  const date = new Date(isoDate)
  const now = new Date()
  const yesterday = new Date(now)
  yesterday.setDate(now.getDate() - 1)

  const isToday = date.toDateString() === now.toDateString()
  const isYesterday = date.toDateString() === yesterday.toDateString()

  if (isToday) return 'Today'
  if (isYesterday) return 'Yesterday'

  return date.toLocaleDateString('fr-FR', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  })
}

const scrollToBottom = (): void => {
  nextTick(() => {
    chatBox.value.scrollTop = chatBox.value.scrollHeight
  })
}

const checkOnlineStatus = async () => {
  try {
    const { data } = await axios.get(`/api/users/${to}/online`, {
      headers: { ...auth.getAuthJSONHeader() }
    })
    isOnline.value = data.status === 'online'
  } catch (err) {
    console.error('Erreur vérification status online', err)
  }

  if (timeout) clearTimeout(timeout)

  timeout = window.setTimeout(() => {
    checkOnlineStatus()
  }, 10000)
}

onBeforeUnmount((): void => {
  if (typingInterval) clearTimeout(typingInterval)
  if (timeout) clearTimeout(timeout)
  window.removeEventListener('chat:message', () => {})
})

onMounted(async (): Promise<void> => {
  const token = auth.getAccessToken()
  const payload = JSON.parse(atob(token.split('.')[1]))
  myId.value = payload.sub

  chatStore.setCurrentReceiver(to)
  chatStore.setTyping(false)

  await Promise.all([fetchProfile(), fetchMessages(), checkOnlineStatus()])
  scrollToBottom()

  window.addEventListener('chat:message', (e: CustomEvent) => {
    const newMsg = e.detail
    if (newMsg.sender_id === to || newMsg.receiver_id === to) {
      scrollToBottom()
      notifySeen()
    }
  })

  input.value?.focus()
})

const notifySeen = async (): Promise<void> => {
  if (!window.ws || window.ws.readyState !== 1) return

  window.ws.send(
    JSON.stringify({
      type: 'seen',
      user_id: myId.value,
      sender_id: myId.value,
      receiver_id: to
    })
  )

  if (to) {
    try {
      await axios.post(`/api/messages/${to}/read`, null, {
        headers: { ...auth.getAuthJSONHeader() }
      })
    } catch (err) {
      console.error('Erreur lors de la mise à jour du message lu', err)
    }
  }
}
</script>

<style scoped>
.fade-enter-active {
  transition: all 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-enter-to {
  opacity: 1;
  transform: translateY(0);
}
</style>
