import { useNotificationStore } from '@/stores/notif'
import { useAuthStore } from '@/stores/auth'
import { Notification } from '@/types'
import { useChatStore } from '@/stores/chat'

const reconnectDelay = 2000
let typingTimeout = 0

export function initNotificationSocket(): void {
  const auth = useAuthStore()

  const chatStore = useChatStore()

  const token = auth.getAccessToken()
  if (!token) return

  window.ws = new WebSocket('ws://localhost:3222/ws?token=' + token)

  window.ws.onmessage = (e): void => {
    if (e?.data) {
      const notif = JSON.parse(e.data)
      if (!notif) return

      const notification = notif as Notification
      notification.timestamp = notif.timestamp || new Date().toISOString()

      if (notification.type === 'ping') {
        if (!window.ws || window.ws.readyState !== 1) return
        window.ws.send(JSON.stringify({ type: 'pong' }))
        return
      }

      if (notification.type === 'message') {
        const payload = JSON.parse(notif.payload)

        chatStore.receiveMessage(payload)

        window.dispatchEvent(new CustomEvent('chat:message', { detail: payload }))
      }

      if (notification.type === 'seen') {
        chatStore.markAsRead(notification.receiver_id, notification.user_id)
      }

      if (notification.type === 'typing') {
        if (notification.user_id === chatStore.currentReceiver) {
          chatStore.setTyping(true)

          if (typingTimeout) {
            clearTimeout(typingTimeout)
          }

          typingTimeout = window.setTimeout(() => {
            chatStore.setTyping(false)
          }, 3000)
        }
      }

      const store = useNotificationStore()
      store.receive(notification)
    }
  }

  window.ws.onclose = (e) => {
    console.log('Disconnected. Reconnecting in', reconnectDelay, 'ms')
    setTimeout(initNotificationSocket, reconnectDelay)
  }

  console.log('✅ WebSocket connected')
}
