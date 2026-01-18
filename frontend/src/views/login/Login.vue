<template>
  <div class="login-container">
    <div class="login-form">
      <h2 class="text-center mb-4">{{ t('Please sign in') }}</h2>
      <form @submit.prevent="submitForm">
        <div class="mb-3">
          <label class="form-label">{{ t('Username') }}</label>
          <input
            v-model="loginForm.username"
            type="text"
            class="form-control"
            :class="{ 'is-invalid': errors.username }"
            :placeholder="t('Username')"
            autocomplete="username"
            required
          />
          <div class="invalid-feedback">{{ errors.username }}</div>
        </div>
        <div class="mb-3">
          <label class="form-label">{{ t('Password') }}</label>
          <input
            v-model="loginForm.password"
            type="password"
            class="form-control"
            :class="{ 'is-invalid': errors.password }"
            :placeholder="t('Password')"
            autocomplete="current-password"
            required
          />
          <div class="invalid-feedback">{{ errors.password }}</div>
        </div>
        <div class="mb-3 form-check">
          <input
            v-model="loginForm.remember"
            type="checkbox"
            class="form-check-input"
            id="remember"
          />
          <label
            class="form-check-label"
            for="remember"
          >
            {{ t('Remember me') }}
          </label>
        </div>
        <button
          type="submit"
          class="btn btn-primary w-100"
          :disabled="loading"
        >
          <span
            v-if="loading"
            class="spinner-border spinner-border-sm me-1"
            role="status"
          ></span>
          {{ t('Sign in') }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
  import { useUserStore } from '@/pinia/modules/user';
  import { reactive, ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { useRouter } from 'vue-router';

  const { t } = useI18n();
  const router = useRouter();
  const userStore = useUserStore();
  const loading = ref(false);

  const loginForm = reactive({
    remember: false,
  });

  const errors = reactive({
    username: '',
    password: '',
  });

  const validateForm = () => {
    let isValid = true;

    if (!loginForm.username) {
      errors.username = t('Please input username');
      isValid = false;
    }

    if (!loginForm.password) {
      errors.password = t('Please input password');
      isValid = false;
    }

    return isValid;
  };

  const submitForm = async () => {
    if (!validateForm()) return;

    try {
      loading.value = true;
      await userStore.LoginIn(loginForm);
      router.push({ path: '/admin' });
    } catch (error) {
      console.error('登录失败', error);
    } finally {
      loading.value = false;
    }
  };
</script>

<style scoped>
  .login-container {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f5f5;
  }

  .login-form {
    width: 100%;
    max-width: 400px;
    padding: 2rem;
    background-color: white;
    border-radius: 0.5rem;
    box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.1);
  }
</style>
