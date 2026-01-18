<template>
  <div
    class="message"
    :class="type ? 'message-' + type : ''"
  >
    <span>{{ message }}</span>
    <span
      class="close"
      @click="close"
    >x</span>
  </div>
</template>

<script setup>
import { useAttrs } from "vue";

defineOptions({
  inheritAttrs: false,
});

const attrs = useAttrs();

const props = defineProps({
  message: {
    type: String,
  },
  type: {
    type: String,
  },
  duration: {
    type: Number,
    default: 8000,
  },
});

function close() {
  if (props.duration) {
    attrs.close && attrs.close();
  }
}
</script>

<style scoped>
.message {
  padding: 8px 16px;
  border-radius: 4px;
  background-color: #ddd;
  display: flex;
  justify-content: space-between;
  position: absolute;
  top: 0;
  left: 0;
  justify-content: center;
  align-items: center;
  z-index: 99999;
}

.close {
  cursor: pointer;
}

.message-success {
  background-color: #67c23a;
}

.message-error {
  background-color: #f56c6c;
}</style>
