<template>
  <div class="relative inline-block">
    <img
      v-if="src"
      :src="src"
      :alt="alt"
      :style="{ width: computedSize, height: computedSize }"
      :class="[
        'rounded-full object-cover shadow transition duration-300',
        ring ? 'ring-2 ring-blue-500 hover:ring-4' : ''
      ]"
    />
    <div
      v-else
      :style="{ width: computedSize, height: computedSize }"
      class="flex items-center justify-center rounded-full bg-gray-200 text-gray-500 shadow"
    >
      <span class="text-sm">?</span>
    </div>

    <span
      v-if="online"
      class="absolute bottom-0 right-0 block h-3 w-3 rounded-full border-2 border-white bg-green-500"
    ></span>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  src: String,
  alt: {
    type: String,
    default: 'Avatar'
  },
  size: {
    type: [String, Number],
    default: '128'
  },
  ring: {
    type: Boolean,
    default: true
  },
  online: { type: Boolean, default: false }
})

const computedSize = computed(() => {
  return typeof props.size === 'number' || /^\d+$/.test(props.size) ? `${props.size}px` : props.size
})
</script>
