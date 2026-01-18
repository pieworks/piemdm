      style="min-height: calc(100vh - 160px); font-size: 0.9rem"
<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Webhook') }}{{ $t('Update') }}</div>
    <div class="card-body mt-3" style="min-height: calc(100vh - 160px); font-size: 0.9rem">
      <div id="create_wrapper">
        <div class="row">
          <div class="col-sm-12">
            <app-form :data-info="dataInfo" @submit-form="submitForm" @go-index="goIndex" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { findWebhook, updateWebhook } from '@/api/webhook';
import { AppToast } from '@/components/toast.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();
const dataInfo = ref({});

onMounted(() => {
  getWebhookInfo();
});

const getWebhookInfo = async () => {
  const params = router.currentRoute.value.query;
  const res = await findWebhook({
    id: params.id,
  });
  dataInfo.value = res.data;
};

async function submitForm(data) {
  const res = await updateWebhook({
    ...data,
  });
  if (res) {
    AppToast.show({
      message: 'Update Webhook success。',
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
