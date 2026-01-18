<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Approval') }}{{ $t('View') }}</div>
    <div class="card-body overlay-wrapper">
      <div id="create_wrapper">
        <div
          class="tab-content"
          id="myTabContent"
        >
          <div
            class="tab-pane fade show active"
            id="baseinfo"
          >
            <div id="create_wrapper container">
              <div class="row">
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="ID"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('ID') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.ID }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="Code"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Code') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Code }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="Title"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Title') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Title }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="ApprovalDefCode"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Approval Define Code') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.ApprovalDefCode }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="EntityCode"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Entity Code') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.EntityCode }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="SerialNumber"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Serial Number') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.SerialNumber }}
                    </div>
                  </div>
                </div>

                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="CurrentTaskID"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Current Task ID') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.CurrentTaskID }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="CurrentTaskName"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Current Task Name') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.CurrentTaskName }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="FormData"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Form Data') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.FormData }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="FormSchema"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Form Schema') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.FormSchema }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="Priority"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Priority') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Priority }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="Urgency"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Urgency') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Urgency }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="Description"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Description') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Description }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="Status"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Status') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Status }}
                    </div>
                  </div>
                </div>

                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="CreatedBy"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Created By') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.CreatedBy }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="UpdatedBy"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Updated By') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.UpdatedBy }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="CreatedAt"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Created At') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.CreatedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="UpdatedAt"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Updated At') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.UpdatedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row h-100">
                    <legend
                      for="DeletedAt"
                      class="col-form-label col-sm-2 my-auto"
                    >
                      {{ $t('Deleted At') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.DeletedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 my-3">
                  <div class="form-group row">
                    <div class="col-sm-10">
                      <button
                        class="btn btn-outline-secondary btn-sm"
                        type="button"
                        @click="goIndex"
                      >
                        {{ $t('Back') }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { findApproval } from '@/api/approval';
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const dataInfo = ref({});

  onMounted(() => {
    getTableData();
  });

  // 查询
  const getTableData = async () => {
    const params = router.currentRoute.value.query;
    const res = await findApproval({
      id: params.id,
    });
    if (res) {
      dataInfo.value = res.data;
    }
  };

  function goIndex() {
    router.push('/admin/approval/index');
  }
</script>

<style scoped></style>
