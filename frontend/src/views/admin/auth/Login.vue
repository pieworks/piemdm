<template>
  <div class="form-signin h-100 w-100 m-auto">
    <div class="w-100 h-100 d-flex align-items-center">
      <div class="w-100">
        <form @submit.prevent="onSubmit">
          <img class="mb-4" src="@/assets/piemdm.png" alt="" width="72" />
          <h1 class="h3 mb-3 fw-normal">Admin Login</h1>

          <div class="form-floating">
            <input type="text" class="form-control" :class="{ 'is-invalid': errors.Username }" id="Username"
              name="Username" v-model="username" v-bind="usernameAttrs" autocomplete="username"
              :placeholder="$t('Username')" />
            <label>{{ $t('Username') }}</label>
            <div v-if="errors.Username" class="invalid-feedback text-start">
              {{ errors.Username }}
            </div>
          </div>
          <div class="form-floating">
            <input type="password" class="form-control" :class="{ 'is-invalid': errors.Password }" id="Password"
              name="Password" v-model="password" v-bind="passwordAttrs" autocomplete="current-password"
              :placeholder="$t('Password')" />
            <label>{{ $t('Password') }}</label>
            <div v-if="errors.Password" class="invalid-feedback text-start">
              {{ errors.Password }}
            </div>
          </div>
          <div class="checkbox mb-3 ps-1 float-start">
            <label>
              <input type="checkbox" value="remember-me" />
              {{ $t('Remember me') }}
            </label>
          </div>
          <button class="w-100 btn btn-lg btn-primary" type="submit">
            {{ $t('Sign in') }}
          </button>
        </form>
        <div class="mt-5">
          <p class="mt-5 mb-3 text-muted">&copy; 2017–2022</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useUserStore } from '@/pinia/modules/user';
import { yup } from '@/utils/yup-config';
import { useForm } from 'vee-validate';
import { useRouter } from 'vue-router';

const router = useRouter();
const userStore = useUserStore();

// 1. 定义验证 Schema
const validationSchema = yup.object({
  Username: yup.string().required(),
  Password: yup.string().required().min(6),
});

// 2. 初始化 useForm
const { errors, defineField, handleSubmit, setFieldError } = useForm({
  validationSchema,
  initialValues: {},
});

// 3. 定义字段
const [username, usernameAttrs] = defineField('Username');
const [password, passwordAttrs] = defineField('Password');

// 4. 定义提交函数
const onSubmit = handleSubmit(async values => {
  const result = await userStore.adminLogin(values);
  if (!result.success) {
    // 如果登录失败，使用 setFieldError 在特定字段上显示后端返回的错误
    setFieldError('Password', 'Incorrect username or password');
  } else {
    // 登录成功，跳转页面
    await router.push('/admin/dashboard/index');
  }
});
</script>

<style scoped>
body {
  background-color: #f5f5f5;
}

.form-signin {
  max-width: 330px;
  padding: 15px;
  text-align: center;
}

.form-signin #Username {
  margin-bottom: -1px;
  border-bottom-right-radius: 0;
  border-bottom-left-radius: 0;
}

.form-signin #Password {
  margin-bottom: 10px;
  border-top-left-radius: 0;
  border-top-right-radius: 0;
}
</style>
