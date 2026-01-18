<template>
  <div class="mt-3" style="min-height: calc(100vh - 105px); font-size: 0.9rem">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Application') }}{{ $t('View') }}</div>
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
                      for="AppId"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('App ID') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.AppId }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="AppSecret"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('App Secret') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.AppSecret }}
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
                      for="IP"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('IP') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.IP }}
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
                <div class="col-sm-12">
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
  import { findApplication } from '@/api/application';
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const dataInfo = ref({});

  onMounted(() => {
    getApplicationInfo();
  });

  // get application info
  const getApplicationInfo = async () => {
    const params = router.currentRoute.value.query;
    const res = await findApplication({
      id: params.id,
    });
    dataInfo.value = res.data;
  };

  // goto index page
  function goIndex() {
    router.push('/admin/application/index');
  }
</script>

<style scoped></style>
