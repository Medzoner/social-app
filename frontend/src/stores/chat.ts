// stores/chat.ts
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { Message } from '@/types'

export const useChatStore = defineStore('chat', () => {
  const messages = ref<Message[]>([])
  const currentReceiver = ref<number | null>(null)
  const typing = ref<boolean>(false)

  function setMessages(msgs: Message[]) {
    messages.value = msgs
  }

  function receiveMessage(msg: Message) {
    if (!messages.value.find((m) => m.id === msg.id)) {
      messages.value.unshift(msg)
    }
  }

  function markAsRead(senderId: number, receiverId: number) {
    messages.value = messages.value.map((m) =>
      m.sender_id === senderId && m.receiver_id === receiverId ? { ...m, read: true } : m
    )
  }

  function setTyping(value: boolean) {
    typing.value = value
  }

  function setCurrentReceiver(receiverId: number | null) {
    currentReceiver.value = receiverId
  }

  return {
    messages,
    setMessages,
    receiveMessage,
    markAsRead,
    setTyping,
    currentReceiver,
    setCurrentReceiver,
    typing
  }
})
