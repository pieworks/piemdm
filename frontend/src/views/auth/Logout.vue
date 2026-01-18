<template>
  <div>正在登出...</div>
</template>

<script setup>
import { useUserStore } from '@/pinia/modules/user';
import { onMounted } from 'vue';

const userStore = useUserStore();

onMounted(() => {
  userStore.logout();
});
</script>

<style scoped>
</style>
