<script setup>
  const props = defineProps({
    show: Boolean,
  });
</script>

<template>
  <div
    v-if="show"
    class="spinner-mask"
  >
    <div class="spinner-container">
      <div
        class="spinner-border text-primary"
        role="status"
      >
        <span class="visually-hidden">Loading...</span>
      </div>
    </div>
  </div>
</template>

<style>
  .spinner-mask {
    position: fixed;
    z-index: 9998;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    /* background-color: rgba(0, 0, 0, 0.5); */
    display: flex;
    transition: opacity 0.3s ease;
  }

  .spinner-container {
    width: 45px;
    margin: auto;
  }
</style>
