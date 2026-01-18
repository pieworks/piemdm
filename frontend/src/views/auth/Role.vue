<template>
  <div class="p-4">
    <div class="flex flex-row w-full space-x-[2rem]">
      <div class="input-group">
        <span class="input-group-text bg-white border-end-0">
          <i class="bi bi-search"></i>
        </span>
        <input
          type="text"
          class="form-control border-start-0 ps-0"
          v-model="search"
          placeholder="Type to search"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
  import request from '@/utils/request';
  import { deleteItem, updateItem } from '@/utils';

  import { AppToast } from '@/components/toast.js';
  import { computed, onMounted, ref, unref } from 'vue';

  const roles = ref([]);
  const groups = ref([]);

  const showCreate = ref(false);
  const showUpdate = ref(false);
  const showDelete = ref(-1);

  const newRole = ref({
    rules: [{}],
  });
  const updatedRole = ref({});
  const updateRow = ref({});

  const createFormRef = ref();
  const updateFormRef = ref();

  const search = ref('');
  const filterRoles = computed(() =>
    roles.value.filter(
      data => !search.value || data.name.toLowerCase().includes(search.value.toLowerCase())
    )
  );

  onMounted(() => {
    request.get(`/api/v1/roles`).then(response => {
      roles.value = Array.from(response.data.data);
    });

    request.get(`/api/v1/groups`).then(response => {
      groups.value = Array.from(response.data.data);
    });
  });

  const createRole = () => {
    const form = unref(createFormRef);
    if (!form) {
      return;
    }
    form.validate((valid, err) => {
      if (valid) {
        request.post('/api/v1/roles', newRole.value).then(response => {
          AppToast.show({ message: 'Create success', color: 'success' });
          roles.value.push(response.data.data);
          showCreate.value = false;
        });
      } else {
        AppToast.show({ message: 'Input invalid, all fields required', color: 'danger' });
      }
    });
  };

  const editRole = row => {
    updatedRole.value = row;
    updateRow.value = row;
    showUpdate.value = true;
  };

  const updateRole = () => {
    const form = unref(updateFormRef);
    if (!form) {
      return;
    }

    form.validate((valid, err) => {
      if (valid) {
        request.put(`/api/v1/roles/${updatedRole.value.id}`, updatedRole.value).then(response => {
          AppToast.show({ message: 'Update success', color: 'success' });
          updateItem(roles, updateRow, updatedRole.value);
          showUpdate.value = false;
        });
      } else {
        AppToast.show({
          message: `Input invalid: ${err}`,
          color: 'danger',
        });
      }
    });
  };

  const deleteRole = row => {
    request.delete(`/api/v1/roles/${row.id}`).then(() => {
      AppToast.show({ message: 'Delete success', color: 'success' });
      deleteItem(roles.value, row);
      showDelete.value = -1;
    });
  };

  const removeRule = (role, item) => {
    const index = role.rules.indexOf(item);
    if (index !== -1) {
      role.rules.splice(index, 1);
    }
  };

  const addRule = role => {
    role.rules.push({
      key: '',
      value: '',
    });
  };
</script>

<style scoped></style>
