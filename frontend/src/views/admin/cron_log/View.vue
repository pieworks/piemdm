<template>
  <div class="mt-3" style="min-height: calc(100vh - 105px); font-size: 0.9rem">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Cron Log') }}{{ $t('View') }}</div>
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
                  <div class="form-group row">
                    <legend
                      for="Method"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Method') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.Method }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Param"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Param') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.Param }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="ErrMsg"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('ErrMsg') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.ErrMsg }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="StartTime"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Start Time') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.StartTime }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="EndTime"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('End Time') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.EndTime }}
                    </div>
                  </div>
                </div>

                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="ExecTime"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Exec Time') }}:
                    </legend>
                    <div class="col-sm-3">
                      {{ dataInfo.ExecTime }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Status"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Status') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.Status }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="CreatedAt"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Created At') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.CreatedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="UpdatedAt"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Updated At') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.UpdatedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="DeletedAt"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Deleted At') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.DeletedAt }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div class="col-sm-12 my-3">
              <div class="form-group row">
                <div class="col-sm-10 col-form-label">
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
</template>

<script setup>
  import { findCronLog } from '@/api/cron_log';
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const dataInfo = ref({});

  onMounted(() => {
    getCronLogInfo();
  });

  // get cron log info
  const getCronLogInfo = async () => {
    const params = router.currentRoute.value.query;
    const res = await findCronLog({
      id: params.id,
    });
    dataInfo.value = res.data;
  };

  // goto index page
  function goIndex() {
    router.push('/admin/cron_log/index');
  }
</script>

<style scoped></style>
