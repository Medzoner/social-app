import { defineStore } from 'pinia'
import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import { Notification } from '@/types'

export interface NotificationState {
  list: Array<any>
  unreadCount: number
}

export const useNotificationStore = defineStore('notifications', {
  state: (): NotificationState => ({
    list: [],
    unreadCount: 0
  }),
  actions: {
    async fetchNotifications(): Promise<void> {
      const auth = useAuthStore()
      const res = await axios.get('/api/notifications', {
        headers: { ...auth.getAuthHeader() }
      })
      this.list = res.data
      this.unreadCount = res.data.length
    },
    async markAllRead(): Promise<void> {
      const auth = useAuthStore()
      await axios.post(
        '/api/notifications/all-read',
        {},
        {
          headers: { ...auth.getAuthHeader() }
        }
      )
      this.unreadCount = 0
    },
    async markRead(notif): Promise<void> {
      const auth = useAuthStore()
      await axios.post(
        '/api/notifications/read',
        { id: notif.id },
        {
          headers: { ...auth.getAuthHeader() }
        }
      )
      this.list = this.list.filter((n) => n.id !== notif.id)
      this.unreadCount--
    },
    receive(notification: Notification): void {
      this.list.push(notification)
      this.unreadCount++
    }
  }
})
