<template>
  <nav class="bg-gray-800 text-white">
    <div class="container mx-auto flex items-center justify-between px-4 py-3">
      <!-- Logo / Titre -->
      <RouterLink to="/feed" class="text-xl font-bold text-white">RéseauSocial</RouterLink>

      <!-- Bouton menu burger (mobile) -->
      <!--      <button @click="toggleMenu" class="rounded p-2 hover:bg-gray-700">-->
      <!--        <MenuIcon />-->
      <!--      </button>-->

      <!-- Menu desktop -->
      <div class="hidden items-center gap-6 md:flex">
        <RouterLink to="/feed" class="hover:underline">Feed</RouterLink>

        <template v-if="isAuthenticated">
          <RouterLink :to="`/chats`" class="hover:underline">Chats</RouterLink>
          <RouterLink :to="`/profile/${user.id}`" class="hover:underline">Mon profil</RouterLink>
          <RouterLink to="/edit-profile" class="hover:underline">Modifier profil</RouterLink>
          <RouterLink v-if="user.role === 'admin'" to="/admin" class="hover:underline"
            >Admin</RouterLink
          >
          <NotificationDropdown />
          <button @click="logout" class="text-left hover:underline">Déconnexion</button>
        </template>

        <template v-else>
          <RouterLink to="/login" class="hover:underline">Connexion</RouterLink>
          <RouterLink to="/register" class="hover:underline">Inscription</RouterLink>
        </template>
      </div>
    </div>

    <!-- Menu mobile (affiché quand menuOpen) -->
    <transition name="fade">
      <div v-if="menuOpen" class="flex flex-col gap-3 px-4 pb-4 md:hidden">
        <RouterLink to="/feed" class="hover:underline">Feed</RouterLink>

        <template v-if="isAuthenticated">
          <RouterLink :to="`/chats`" class="hover:underline">Chats</RouterLink>
          <RouterLink :to="`/profile/${user.id}`" class="hover:underline">Mon profil</RouterLink>
          <RouterLink to="/edit-profile" class="hover:underline">Modifier profil</RouterLink>
          <RouterLink v-if="user.role === 'admin'" to="/admin" class="hover:underline"
            >Admin</RouterLink
          >
          <NotificationDropdown />
          <button @click="logout" class="text-left hover:underline">Déconnexion</button>
        </template>

        <template v-else>
          <RouterLink to="/login" class="hover:underline">Connexion</RouterLink>
          <RouterLink to="/register" class="hover:underline">Inscription</RouterLink>
        </template>
      </div>
    </transition>
  </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
// import MenuIcon from '@/icons/MenuIcon.vue'
import NotificationDropdown from '@/components/notification/NotificationDropdown.vue'

const auth = useAuthStore()
const { user, isAuthenticated } = storeToRefs(auth)
const router = useRouter()

const menuOpen = ref(false)

async function logout() {
  await auth.logout()
  await router.push('/login')
  menuOpen.value = false
}
</script>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
