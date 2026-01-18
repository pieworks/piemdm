<template>
  <form name="permissionForm" id="permissionForm" method="post">
    <div class="col-12 col-sm-12">
      <div class="row">
        <!-- Code -->
        <div class="col-sm-6">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-4', { required: !isView }]">
              {{ $t('Code') }}:
            </legend>
            <div class="col-sm-8">
              <input v-if="!isView" type="text" class="form-control form-control-sm" v-model="formData.code"
                placeholder="e.g. user:list" required />
              <input v-else type="text" class="form-control form-control-sm" :value="formData.code" disabled />
            </div>
          </div>
        </div>

        <!-- Name -->
        <div class="col-sm-6">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-4', { required: !isView }]">
              {{ $t('Name') }}:
            </legend>
            <div class="col-sm-8">
              <input v-if="!isView" type="text" class="form-control form-control-sm" v-model="formData.name" required />
              <input v-else type="text" class="form-control form-control-sm" :value="formData.name" disabled />
            </div>
          </div>
        </div>

        <!-- Resource -->
        <div class="col-sm-6 mt-2">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-4', { required: !isView }]">
              {{ $t('Resource') }}:
            </legend>
            <div class="col-sm-8">
              <input v-if="!isView" type="text" class="form-control form-control-sm" v-model="formData.resource"
                placeholder="e.g. user" required />
              <input v-else type="text" class="form-control form-control-sm" :value="formData.resource" disabled />
            </div>
          </div>
        </div>

        <!-- Action -->
        <div class="col-sm-6 mt-2">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-4', { required: !isView }]">
              {{ $t('Action') }}:
            </legend>
            <div class="col-sm-8">
              <input v-if="!isView" type="text" class="form-control form-control-sm" v-model="formData.action"
                placeholder="e.g. list" required />
              <input v-else type="text" class="form-control form-control-sm" :value="formData.action" disabled />
            </div>
          </div>
        </div>

        <!-- Parent ID -->
        <div class="col-sm-6 mt-2">
          <div class="form-group row">
            <legend class="col-form-label col-sm-4">
              {{ $t('ParentID') }}:
            </legend>
            <div class="col-sm-8">
              <input v-if="!isView" type="number" class="form-control form-control-sm" v-model="formData.parent_id" />
              <input v-else type="number" class="form-control form-control-sm" :value="formData.parent_id" disabled />
            </div>
          </div>
        </div>

        <!-- Description -->
        <div class="col-sm-12 mt-2">
          <div class="form-group row">
            <legend class="col-form-label col-sm-2">
              {{ $t('Description') }}:
            </legend>
            <div class="col-sm-10">
              <textarea v-if="!isView" class="form-control form-control-sm" v-model="formData.description"
                rows="3"></textarea>
              <textarea v-else class="form-control form-control-sm" :value="formData.description" rows="3"
                disabled></textarea>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="form-group row mt-3">
      <div class="col-sm-auto">
        <button v-if="!isView" type="button" class="btn btn-outline-primary btn-sm me-2" @click="onSubmit">
          {{ $t('Submit') }}
        </button>
        <button type="button" class="btn btn-outline-secondary btn-sm" @click="router.back()">
          {{ $t(isView ? 'Back' : 'Cancel') }}
        </button>
      </div>
    </div>
  </form>
</template>

<script setup>
import { createPermission, getPermission, updatePermission } from '@/api/permission';
import { AppToast } from '@/components/toast.js';
import { computed, onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const router = useRouter();
const route = useRoute();
const formData = ref({
  parent_id: 0
});
const isUpdate = ref(false);
const isView = computed(() => route.path.includes('/view'));

onMounted(async () => {
  if (route.query.id) {
    isUpdate.value = true;
    const res = await getPermission(route.query.id);
    if (res && res.data) {
      formData.value = res.data;
    }
  }
});

const onSubmit = async () => {
  if (!formData.value.code || !formData.value.name) {
    AppToast.show({ message: 'Code and Name are required', color: 'danger' });
    return;
  }

  let res;
  if (isUpdate.value) {
    res = await updatePermission(formData.value.ID, formData.value);
  } else {
    res = await createPermission(formData.value);
  }

  if (res) {
    AppToast.show({
      message: isUpdate.value ? 'Updated successfully' : 'Created successfully',
      color: 'success',
    });
    router.push('/admin/permission/index');
  }
};
</script>

<style scoped>
.required:after {
  content: ' *';
  color: #dc3545;
  font-weight: bold;
}
</style>
