<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Table Field') }} {{ $t('Create') }}</div>
    <div class="card-body overlay-wrapper mt-2" style="min-height: calc(100vh - 150px);; font-size: 0.9rem">
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
import { createTableField, getTableFields } from '@/api/table_field';
import { AppToast } from '@/components/toast.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();
const dataInfo = ref({
  // 基本配置由 _form.vue 中的 createDefaultFormData 提供
  // 废弃字段的默认值已移除
});
const params = ref({});
const fieldCount = ref(0); // 现有字段数量

onMounted(async () => {
  params.value = router.currentRoute.value.query;
  dataInfo.value.TableCode = params.value.table_code;

  // 获取现有字段数量
  try {
    const res = await getTableFields({ table_code: params.value.table_code });
    if (res && res.data) {
      fieldCount.value = res.data.length || 0;
      // 设置默认排序值为字段数量+1
      dataInfo.value.Sort = fieldCount.value + 1;
    }
  } catch (error) {
    console.error('Failed to get field count:', error);
    dataInfo.value.Sort = 0;
  }
});

const submitForm = async data => {
  const res = await createTableField({
    table_code: params.value.table_code,
    ...data,
  });
  if (res) {
    AppToast.show({
      message: 'Create table field success。',
      color: 'success',
    });
    router.push('/admin/table_field/index?table_code=' + params.value.table_code);
  }
};

function goIndex() {
  router.push('/admin/table_field/index?table_code=' + params.value.table_code);
}
</script>

<style scoped></style>
