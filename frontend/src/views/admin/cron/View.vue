<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Cron') }}{{ $t('View') }}</div>
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
                      for="Code"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Code') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.Code }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Expression"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Expression') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.Expression }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Name"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Name') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.Name }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="EntityCode"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Entity Code') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.EntityCode }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="System"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('System') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.System }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Url"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Url') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.Url }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Protocol"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Protocol') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.Protocol }}
                    </div>
                  </div>
                </div>
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
                      for="AppId"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('AppId') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.AppId }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="AppKey"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('AppKey') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.AppKey }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="SignType"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('SignType') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.SignType }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Status"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Status') }}:
                    </legend>
                    <div class="col-sm-3">
                      {{ dataInfo.Status }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Description"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Description') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.Description }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="CreatedBy"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Created By') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.CreatedBy }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="UpdatedBy"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Updated By') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.UpdatedBy }}
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
  import { findCron } from '@/api/cron';
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const dataInfo = ref({});

  onMounted(() => {
    getCronInfo();
  });

  // get Cron info
  const getCronInfo = async () => {
    const params = router.currentRoute.value.query;
    const res = await findCron({
      id: params.id,
    });
    dataInfo.value = res.data;
  };

  // goto index page
  function goIndex() {
    router.push('/admin/cron/index');
  }
</script>

<style scoped></style>
