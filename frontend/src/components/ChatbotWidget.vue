<template>
  <div>
    <button class="chat-toggle" @click="isOpen = !isOpen">💬</button>

    <div v-if="isOpen" class="chat-window">
      <div class="chat-messages">
        <div v-for="(msg, i) in messages" :key="i" :class="msg.role">
          <template v-if="msg.role === 'user'"> <strong>Vous:</strong> {{ msg.content }} </template>
          <template v-else>
            <strong>Bot:</strong>
            <div class="bot-html" v-html="msg.content"></div>
          </template>
        </div>
      </div>

      <form @submit.prevent="sendMessage">
        <input v-model="userInput" placeholder="Écris quelque chose..." />
        <button type="submit">Envoyer</button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const isOpen = ref(false)
const userInput = ref('')
const messages = ref([])

const sendMessage = async () => {
  if (!userInput.value.trim()) return

  const content = userInput.value.trim()
  messages.value.push({ role: 'user', content })
  userInput.value = ''

  try {
    const res = await fetch('http://localhost:3222/llm/chat', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message: content })
    })

    const data = await res.json()
    messages.value.push({ role: 'bot', content: data.reply || 'Réponse vide.' })
  } catch (err) {
    messages.value.push({ role: 'bot', content: 'Erreur de connexion.' })
  }
}
</script>

<style scoped>
.chat-toggle {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background: #007bff;
  color: white;
  border-radius: 50%;
  padding: 1rem;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  z-index: 9999;
}

.chat-window {
  position: fixed;
  bottom: 80px;
  right: 20px;
  width: 400px;
  height: 500px;
  background: white;
  border: 1px solid #ccc;
  display: flex;
  flex-direction: column;
  z-index: 9999;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.user {
  text-align: right;
  margin-bottom: 5px;
}

.bot {
  text-align: left;
  margin-bottom: 5px;
}

form {
  display: flex;
  border-top: 1px solid #ccc;
}

input {
  flex: 1;
  border: none;
  padding: 10px;
}

button[type='submit'] {
  background: #007bff;
  color: white;
  border: none;
  padding: 0 1rem;
}
</style>
