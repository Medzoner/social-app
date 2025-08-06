<!-- src/components/auth/sso/Form.vue -->
<template>
  <div v-if="message" class="flex items-center justify-between">
    <span>{{ message }}</span>
    <button @click="$emit('close')" class="ml-2 text-sm text-red-600">X</button>
  </div>
</template>

<script setup lang="ts">
defineProps<{ message: string }>()
</script>

<style scoped></style>
