<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Webhook') }}{{ $t('Create') }}</div>
    <div class="card-body mt-3" style="min-height: calc(100vh - 160px); font-size: 0.9rem">
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
import { createWebhook } from '@/api/webhook';
import { AppToast } from '@/components/toast.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();
const dataInfo = ref({
  Status: 'Normal',
  ContentType: 'application/json',
});

onMounted(() => { });

async function submitForm(data) {
  const res = await createWebhook({
    ...data,
  });
  if (res) {
    AppToast.show({
      message: 'Create webhook success。',
      color: 'success',
    });
    router.push('/admin/webhook/index');
  }
}

function goIndex() {
  router.push('/admin/webhook/index');
}
</script>

<style scoped></style>
