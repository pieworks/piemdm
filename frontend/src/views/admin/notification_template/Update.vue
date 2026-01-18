<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-3 pb-2">
      {{ $t('Edit') }} {{ $t('Notification Template') }}
    </div>
    <div class="card-body mt-3" style="min-height: calc(100vh - 160px); font-size: 0.9rem">
      <div id="create_wrapper">
        <div class="row">
          <div class="col-sm-12">
            <AppForm v-if="dataInfo" :dataInfo="dataInfo" @submitForm="submitForm" @goIndex="goIndex" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { findNotificationTemplate, updateNotificationTemplate } from '@/api/notification_template';
import { AppToast } from '@/components/toast.js';
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import AppForm from './Form.vue';

const route = useRoute();
const router = useRouter();
const dataInfo = ref(null);
const id = route.query.id;

onMounted(async () => {
  if (id) {
    const res = await findNotificationTemplate(id);
    if (res) {
      dataInfo.value = res.data;
    }
  }
});

const submitForm = async formData => {
  let variables = {};
  if (formData.Variables) {
    try {
      variables = JSON.parse(formData.Variables);
    } catch (e) {
      console.error('Failed to parse variables JSON:', e);
      // Fallback or error handling if needed, but validation should catch this
    }
  }

  const res = await updateNotificationTemplate({
    ...formData,
    Variables: variables,
    ID: id,
  });
  if (res) {
    AppToast.show({
      message: '更新成功',
      color: 'success',
    });
    router.push('/admin/notification_template/index');
  }
};

const goIndex = () => {
  router.push('/admin/notification_template/index');
};
</script>
