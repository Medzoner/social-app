import { defineStore } from 'pinia'
import { jwtDecode } from 'jwt-decode'
import axios from 'axios'
import router, {routeNames, routes} from '@/router'
import { AuthHeader, JwtToken, JwtUser } from '@/types'
import { initNotificationSocket } from '@/notifications/ws'

let expireInterval = null
const accessTokenKey = 'access_token'
const idTokenKey = 'id_token'
const refreshTokenKey = 'refresh_token'

export interface AuthState {
  user: JwtUser | null
  token: JwtToken | object
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: {}
  }),

  getters: {
    isAuthenticated: (state): boolean => !!state.user
  },

  actions: {
    expireIn: async function (): Promise<number | null> {
      if (!this.token?.access_token) return null

      const decoded = jwtDecode(this.token.access_token)
      const now = Date.now() / 1000
      const expiresIn = decoded.exp - now

      if (expiresIn < 30) {
        await this.refreshToken()
        return null
      }

      return expiresIn > 0 ? expiresIn : null
    },

    startTokenWatcher(): void {
      if (expireInterval) return

      expireInterval = setInterval(async () => {
        await this.expireIn()
      }, 30_000)
    },

    stopTokenWatcher(): void {
      if (expireInterval) {
        clearInterval(expireInterval)
        expireInterval = null
      }
    },

    loginWithToken(token: JwtToken | null): void {
      if (!token?.access_token) {
        console.error('❌ Invalid token provided')
        return
      }

      if (window.ws?.readyState !== 1) {
        initNotificationSocket()
      }

      try {
        const decoded: JwtUser = jwtDecode(token.access_token)
        this.token = token
        this.user = {
          id: decoded.sub,
          username: decoded.username,
          role: decoded.role,
          verified: decoded.verified,
          email: decoded.email
        }
        localStorage.setItem(accessTokenKey, token.access_token || '')
        localStorage.setItem(idTokenKey, token.id_token || '')
        localStorage.setItem(refreshTokenKey, token.refresh_token || '')
        this.startTokenWatcher()
      } catch (e) {
        console.error('JWT invalid', e)
      }
    },

    async logout(): Promise<void> {
      this.stopTokenWatcher()
      if (window.ws) {
        window.ws.onclose = () => {
          console.log('WebSocket closed during logout')
        }
      }

      try {
        await axios.post(
          '/api/logout',
          {},
          {
            headers: { ...this.getAuthHeader() }
          }
        )
      } catch (error) {
        console.error('Error during logout:', error)
      }
      this.user = null
      this.token = null

      localStorage.removeItem(accessTokenKey)
      localStorage.removeItem(idTokenKey)
      localStorage.removeItem(refreshTokenKey)

      if (window.ws && window.ws.readyState === WebSocket.OPEN) {
        window.ws.close()
        window.ws = null
      }
    },

    getAuthHeader(type = 'access'): AuthHeader {
      if (type === 'refresh') {
        return { Authorization: `Bearer ${this.token.refresh_token}` }
      }

      return { Authorization: `Bearer ${this.token.access_token}` }
    },

    getAuthJSONHeader(type = 'access'): AuthHeader {
      return this.getJSONHeader(this.getAuthHeader(type))
    },

    getJSONHeader(headers: AuthHeader): AuthHeader {
      const contentType = 'application/json'

      return {
        ...headers,
        'Content-Type': contentType,
      }
    },

    getAccessToken(): string | null {
      if (!this.token?.access_token) {
        console.error('❌ No access token available')
        return null
      }
      return this.token.access_token
    },

    async refreshToken(): Promise<void> {
      try {
        const res = await axios.post(
          '/api/refresh',
          {},
          {
            headers: { ...this.getAuthHeader('refresh') }
          }
        )

        const newToken = res.data
        if (!newToken?.refresh_token) {
          console.error('❌ No new token received')
          await this.logout()
          await router.push(`/${routeNames.login}`)
          return
        }

        this.loginWithToken(newToken)
        console.info('🔄 Token refreshed')
      } catch (err) {
        console.error('❌ Error refreshing token', err)
        await this.logout()
        await router.push(`/${routeNames.login}`)
      }
    },

    restoreSession(): void {
      const token = {
        access_token: localStorage.getItem(accessTokenKey),
        id_token: localStorage.getItem(idTokenKey),
        refresh_token: localStorage.getItem(refreshTokenKey)
      }
      if (token?.access_token && token?.refresh_token && token?.id_token) {
        this.loginWithToken(token)
      }
    }
  }
})
