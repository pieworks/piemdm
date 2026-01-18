<script setup>
  import { onMounted, watch } from 'vue';

  const props = defineProps({
    show: Boolean,
    display: String,
  });

  onMounted(() => {});

  watch(props, newProps => {});
</script>

<template>
  <!-- Modal -->
  <div
    :class="['modal', 'fade', show ? 'show' : '']"
    id="exampleModal"
    tabindex="-1"
    :style="{ display: display }"
  >
    <div class="modal-dialog modal-dialog-centered">
      <div class="modal-content">
        <div class="modal-header">
          <h1
            class="modal-title fs-5"
            id="exampleModalLabel"
          >
            title
          </h1>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
          ></button>
        </div>
        <div class="modal-body"></div>
        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-secondary"
            data-bs-dismiss="modal"
          >
            Close
          </button>
          <button
            type="button"
            class="btn btn-primary"
          >
            Save changes
          </button>
        </div>
      </div>
    </div>
  </div>
  <div :class="['modal-backdrop', 'fade', show ? 'show' : '']"></div>
</template>

<style></style>
