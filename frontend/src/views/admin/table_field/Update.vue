<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Table Field') }} {{ $t('Update') }}</div>
    <div class="card-body overlay-wrapper mt-2" style="min-height: calc(100vh - 150px); font-size: 0.9rem">
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
import { findTableField, updateTableField } from '@/api/table_field';
import { AppToast } from '@/components/toast.js';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();
const dataInfo = ref({});
const params = ref({});

onMounted(() => {
  params.value = router.currentRoute.value.query;
  getTableFieldInfo();
});

// get table field info
const getTableFieldInfo = async () => {
  const res = await findTableField({
    id: params.value.id,
  });
  dataInfo.value = res.data;
};

// update table field
async function submitForm(data) {
  const res = await updateTableField({
    ...data,
  });
  if (res) {
    AppToast.show({
      message: 'Update Table Field success。',
      color: 'success',
    });
    router.push('/admin/table_field/index?table_code=' + params.value.table_code);
  }
}

// go to table field index page
function goIndex() {
  router.push('/admin/table_field/index?table_code=' + params.value.table_code);
}
</script>

<style scoped></style>
