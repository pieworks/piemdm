<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Cron') }}{{ $t('Create') }}</div>
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
import { createCron } from '@/api/cron';
import { AppToast } from '@/components/toast.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();
const dataInfo = ref({
  Status: 'Normal',
  Protocol: 'Http',
  Method: 'GET',
});

onMounted(() => { });

async function submitForm(data) {
  const res = await createCron({
    ...data,
  });
  if (res) {
    AppToast.show({
      message: 'Create cron success。',
      color: 'success',
    });
    router.push('/admin/cron/index');
  }
}

function goIndex() {
  router.push('/admin/cron/index');
}
</script>

<style scoped></style>
