<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-3 pb-2">
      {{ $t('New') }} {{ $t('Notification Template') }}
    </div>
    <div class="card-body mt-3" style="min-height: calc(100vh - 160px); font-size: 0.9rem">
      <div id="create_wrapper">
        <div class="row">
          <div class="col-sm-12">
            <AppForm @submitForm="submitForm" @goIndex="goIndex" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { createNotificationTemplate } from '@/api/notification_template';
import { AppToast } from '@/components/toast.js';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();

const submitForm = async formData => {
  let variables = {};
  if (formData.Variables) {
    try {
      variables = JSON.parse(formData.Variables);
    } catch (e) {
      console.error('Failed to parse variables JSON:', e);
    }
  }

  const res = await createNotificationTemplate({
    ...formData,
    Variables: variables,
  });
  if (res) {
    AppToast.show({
      message: '创建成功',
      color: 'success',
    });
    router.push('/admin/notification_template/index');
  }
};

const goIndex = () => {
  router.push('/admin/notification_template/index');
};
</script>
