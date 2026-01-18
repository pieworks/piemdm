<template>
  <span :class="badgeClass">{{ status }}</span>
</template>

<script setup>
  import { computed } from 'vue';

  const props = defineProps({
    status: {
      type: String,
      required: true,
    },
  });

  const badgeClass = computed(() => {
    switch (props.status) {
      case 'Normal':
        return 'badge text-bg-light';
      case 'Frozen':
        return 'badge text-bg-warning';
      case 'Pending':
        return 'badge text-bg-warning';
      case 'Approved':
      case 'Resolved':
        return 'badge text-bg-success';
      case 'Deleted':
        return 'badge text-bg-danger';
      default:
        return 'badge text-bg-secondary';
    }
  });
</script>
