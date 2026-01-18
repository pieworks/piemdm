<template>
</template>

<script setup>
import { onMounted } from 'vue';
import request from '@/utils/request'
import { getUser } from '@/utils';
import { useRouter } from 'vue-router';
import { AppToast } from '@/components/toast.js';

const router = useRouter();

onMounted(
    () => {
        const params = new URLSearchParams(window.location.search);
        const state = params.get("state");
        const code = params.get("code");
        const queryString = atob(state);
        const newParams = new URLSearchParams(queryString);
        const authType = newParams.get("oauth");

    request.post("/api/v1/auth/token", {
        authType : authType,
        authCode: code,
        setCookie: true,
      }).then((response) => {
        let user = getUser();
        AppToast.show({
          message: 'Login Success: Hi~ ' + user.name,
          color: 'success',
        });
        router.push('/');
      })
    }
)

</script>

<style scoped>
</style>
