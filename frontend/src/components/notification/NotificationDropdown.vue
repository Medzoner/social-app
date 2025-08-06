<template>
  <div class="relative" @click.stop="toggle">
    <div class="relative cursor-pointer">
      <span v-if="count > 0" class="text-2xl">🔔</span>
      <span v-else class="text-2xl">🔕</span>
      <span
        v-if="count > 0"
        class="absolute -right-2 -top-1 rounded-full bg-red-600 px-1 text-xs text-white"
      >
        {{ count }}
      </span>
    </div>

    <transition v-if="count > 0" name="fade">
      <div
        v-if="open"
        class="absolute right-0 z-50 mt-2 w-64 rounded bg-white text-black shadow-lg"
      >
        <div class="border-b p-2 font-semibold">Notifications</div>
        <ul class="max-h-64 overflow-y-auto">
          <li v-for="notif in reversed" :key="notif.id" class="px-4 py-2 text-sm hover:bg-gray-100">
            <RouterLink :to="notif.link" class="hover:underline">{{ notif.content }}</RouterLink>
            <br />
            <small class="text-gray-500">{{ new Date(notif.created_at).toLocaleString() }}</small>
            <span @click.stop="markRead(notif)">vu</span>
          </li>
        </ul>
        <div class="border-t text-center">
          <button
            @click.stop="markAllRead"
            class="w-full py-2 text-sm text-blue-600 hover:underline"
          >
            Marquer comme lues
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useNotificationStore } from '@/stores/notif'
import { storeToRefs } from 'pinia'

const store = useNotificationStore()
const { list, unreadCount: count } = storeToRefs(store)
const open = ref(false)

const reversed = computed(() => (Array.isArray(list.value) ? [...list.value].reverse() : []))

function toggle() {
  open.value = !open.value
}

function markAllRead() {
  store.markAllRead()
  open.value = false
}

function markRead(notif) {
  store.markRead(notif)
}

onMounted(() => {
  store.fetchNotifications()
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
