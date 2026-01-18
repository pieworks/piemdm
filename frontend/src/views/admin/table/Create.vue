<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Table') }} {{ $t('Create') }}</div>
    <ul class="nav nav-tabs mh-30" id="myTab" role="tablist">
      <li class="nav-item mh-30" role="presentation">
        <a class="nav-link mh-30 active" id="baseinfo-tab" data-toggle="tab" href="#baseinfo">
          {{ $t('Basic Info') }}
        </a>
      </li>
    </ul>
    <div class="card-body overlay-wrapper mt-2" style="min-height: calc(100vh - 200px); font-size: 0.9rem">
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
import { createTable } from '@/api/table';
import { AppToast } from '@/components/toast';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import AppForm from './Form.vue';

const router = useRouter();
const dataInfo = ref({
  DisplayMode: 'List',
  Status: 'Normal',
  TableType: 'Entity',
});

onMounted(() => { });

async function submitForm(data) {
  const res = await createTable({
    ...data,
  });
  if (res) {
    AppToast.show({
      message: 'Create table success。',
      color: 'success',
    });
    router.push('/admin/table/index');
  }
}

function goIndex() {
  router.push('/admin/table/index');
}
</script>

<style scoped></style>
