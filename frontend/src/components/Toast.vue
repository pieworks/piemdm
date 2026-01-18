<!-- text-bg-primary -->
<!-- text-bg-secondary -->
<!-- text-bg-success -->
<!-- text-bg-danger -->
<!-- text-bg-warning -->
<!-- text-bg-info -->
<!-- text-bg-light -->
<!-- text-bg-dark -->

<template>

  <div class="toast-container position-fixed top-0 end-0 p-3">
    <div
      :class="['toast', color ? 'text-bg-' + color : 'text-bg-light', 'border-0', 'fade', show ? 'show' : '']"
      role="alert"
      aria-live="assertive"
      aria-atomic="true"
    >
      <div
        v-if="header"
        class="toast-header"
      >
        <i :class="['bi', 'bi-square-fill', 'me-2', color ? 'text-' + color : 'text-primary']"></i>
        <strong class="me-auto">{{ headerTitle || 'Notification' }}</strong>
        <button
          type="button"
          class="btn-close"
          data-bs-dismiss="toast"
          aria-label="Close"
        ></button>
      </div>
      <div class="d-flex">
        <div class="toast-body">
          {{ message }}
        </div>
        <button
          v-if="!header"
          type="button"
          class="btn-close me-2 m-auto"
          :class="color ? 'btn-close-white' : ''"
          data-bs-dismiss="toast"
          aria-label="Close"
        ></button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  show: Boolean,
  header: Boolean,
  headerTitle: String,
  color: String,
  message: String,
})
</script>

<style></style>
