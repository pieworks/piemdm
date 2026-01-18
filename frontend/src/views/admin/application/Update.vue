<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Application') }}{{ $t('Update') }}</div>
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
import { findApplication, updateApplication } from '@/api/application';
import { AppToast } from '@/components/toast.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();
const dataInfo = ref({});

onMounted(() => {
  getApplicationInfo();
});

const getApplicationInfo = async () => {
  const params = router.currentRoute.value.query;
  const res = await findApplication({
    id: params.id,
  });
  dataInfo.value = res.data;
};

async function submitForm(data) {
  const res = await updateApplication({
    ...data,
  });
  if (res) {
    AppToast.show({
      message: 'Update Application success。',
      color: 'success',
    });
    router.push('/admin/application/index');
  }
}

function goIndex() {
  router.push('/admin/application/index');
}
</script>

<style scoped></style>
