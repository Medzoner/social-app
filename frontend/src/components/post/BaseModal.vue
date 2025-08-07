<template>
  <transition name="fade">
    <div
      v-if="show"
      id="post-modal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
      @click.self="emit('close')"
    >
      <div
        class="relative w-full max-w-md overflow-y-auto rounded-lg bg-white p-4 shadow-lg sm:max-h-[90vh] sm:w-[90%] sm:max-w-md sm:rounded"
      >
        <button
          name="close-button"
          @click="emit('close')"
          class="absolute right-2 top-2 text-gray-500 hover:text-black"
        >
          ✖
        </button>

        <slot name="header">
          <h3 class="mb-4 text-lg font-semibold">Modal</h3>
        </slot>

        <slot />
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
defineProps<{ show: boolean }>()
const emit = defineEmits(['close'])
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
