<template>
  <div class="flex min-h-screen flex-col bg-gray-100 text-gray-900">
    <NavBar />

    <main class="container mx-auto flex-1 px-4 py-6">
      <RouterView />
    </main>
  </div>
  <ChatbotWidget />
</template>

<script setup>
import NavBar from './components/NavBar.vue'
import ChatbotWidget from './components/ChatbotWidget.vue'
</script>
