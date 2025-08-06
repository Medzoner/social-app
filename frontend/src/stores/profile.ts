import { defineStore } from 'pinia'
import axios from 'axios'
import { useAuthStore } from './auth.js'
import { Profile } from '@/types'

export interface ProfileState {
  profiles: Profile[]
}

export const useProfileStore = defineStore('profile', {
  state: (): ProfileState => ({
    profiles: []
  }),

  actions: {
    async fetchProfile(id: number): Promise<any> {
      const auth = useAuthStore()
      if (this.profiles[id]) {
        return this.profiles[id] // ✅ en cache
      }

      try {
        const res = await axios.get(`/api/users/${id}`, {
          headers: { ...auth.getAuthHeader() }
        })
        this.profiles[id] = res.data
        return res.data
      } catch (err) {
        console.error(`Erreur fetch du profil ID ${id}:`, err)
        throw err
      }
    }
  }
})
