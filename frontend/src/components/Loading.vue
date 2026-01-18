<template>
  <div
    class="loading-container"
    v-show="show"
  >
    <div class="loading-box">
      <div class="d-flex justify-content-center">
        <div
          class="spinner-border"
          role="status"
        >
          <span class="visually-hidden">Loading...</span>
        </div>
      </div>
      <span>{{ title }}</span>
    </div>
  </div>
</template>

<script setup>
import { } from "vue";

const props = defineProps({
  show: Boolean,
  title: String,
})

</script>

<style scoped>
.loading-container {
  position: fixed;
  z-index: 999;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-box {
  background-color: rgba(0, 0, 0, .5);
  width: 100px;
  height: 100px;
  border-radius: 5px;
  color: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}</style>
