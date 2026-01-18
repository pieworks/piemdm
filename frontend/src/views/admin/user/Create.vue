<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('User') }}{{ $t('Create') }}</div>
    <div class="card-body" style="min-height: calc(100vh - 150px); font-size: 0.9rem">
      <div id="create_wrapper">
        <div class="row">
          <div class="col-sm-12">
            <AppForm :data-info="dataInfo" @submit-form="submitForm" @go-index="goIndex" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { createUser } from '@/api/user';
import { AppToast } from '@/components/toast.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();
const dataInfo = ref({
  Status: 'Normal',
  Language: 'zh-cn',
  Sex: 'Male',
  Admin: 'No',
});

onMounted(() => { });

async function submitForm(data) {
  const res = await createUser({
    ...data,
  });
  if (res) {
    AppToast.show({
      message: 'Create user success。',
      color: 'success',
    });
    router.push('/admin/user/index');
  }
}

function goIndex() {
  router.push('/admin/user/index');
}
</script>

<style scoped></style>
