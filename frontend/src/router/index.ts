import { createRouter, createWebHistory, RouteRecordRaw, useRoute } from 'vue-router'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import Feed from '../views/Feed.vue'
import Profile from '../views/Profile.vue'
import EditProfile from '../views/EditProfile.vue'
import Admin from '../views/Admin.vue'
import Chat from '../views/Chat.vue'
import ChatList from '../views/ChatList.vue'
import Verify from '../views/Verify.vue'
import { useAuthStore } from '@/stores/auth'
import About from '../views/About.vue'
import OAuthLogin from '../views/OAuthLogin.vue'
import OAuthCallback from '../views/OAuthCallback.vue'
import Page404 from '../views/404.vue'

export const routeNames = {
  login: 'login',
  register: 'register',
  feed: 'feed',
  profile: 'profile',
  editProfile: 'edit-profile',
  admin: 'admin',
  chat: 'chat',
  chatList: 'chats',
  verify: 'verify',
  about: 'about',
  oauthLogin: 'oauth-login',
  oauthCallback: 'oauth-callback'
}

export const routes: Readonly<RouteRecordRaw[]> = [
  { path: '/404', component: Page404 },
  { path: '/', redirect: `/${routeNames.feed}` },
  { path: `/${routeNames.about}`, component: About },
  { path: `/${routeNames.login}`, component: Login },
  { path: `/${routeNames.register}`, component: Register },
  { path: `/${routeNames.oauthLogin}`, component: OAuthLogin, meta: { requiresAuth: false } },
  { path: `/${routeNames.oauthCallback}`, component: OAuthCallback, meta: { requiresAuth: false } },
  { path: `/${routeNames.feed}`, component: Feed, meta: { requiresAuth: true } },
  { path: `/${routeNames.chatList}`, component: ChatList, meta: { requiresAuth: true } },
  { path: `/${routeNames.chat}/:id`, component: Chat, meta: { requiresAuth: true } },
  { path: `/${routeNames.profile}/:id`, component: Profile, meta: { requiresAuth: true } },
  { path: `/${routeNames.editProfile}`, component: EditProfile, meta: { requiresAuth: true } },
  {
    path: `/${routeNames.admin}`,
    component: Admin,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: `/${routeNames.verify}`,
    name: `${routeNames.verify}`,
    component: Verify,
    meta: { requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next): void => {
  const auth = useAuthStore()

  if (!to.meta.requiresAuth && !to.meta.requiresAdmin) {
    return next()
  }

  if (
    (!auth?.user?.verified && to.meta.requiresAuth && to.path !== `/${routeNames.login}`) ||
    (!auth.isAuthenticated && to.path !== `/${routeNames.verify}`)
  ) {
    return next(`/${routeNames.login}`)
  }

  if (
    auth?.user?.verified &&
    to.meta.requiresAdmin &&
    auth.user?.role !== 'admin'
  ) {
    return next(`/${routeNames.feed}`)
  }

  if (auth.user && !auth.user.verified && to.meta.requiresAuth && to.path !== `/${routeNames.verify}`) {
    return next(`/${routeNames.verify}`)
  }

  next()
})

export default router
