<template>
  <div class="mt-3" style="min-height: calc(100vh - 105px); font-size: 0.9rem">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Webhook') }}{{ $t('View') }}</div>
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
                      for="Url"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Url') }}:
                    </legend>
                    <div class="col-sm-10 my-auto">
                      {{ dataInfo.Url }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="TableCode"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Table Code') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.TableCode }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Username"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Username') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Username }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="ContentType"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Content Type') }}:
                    </legend>
                    <div class="col-sm-3 my-auto">
                      {{ dataInfo.ContentType }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Secret"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Secret') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Secret }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Events"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Events') }}:
                    </legend>
                    <div class="col-sm-6 my-auto">
                      {{ dataInfo.Events }}
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
                    <div class="col-sm-3 my-auto">
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
                    <div class="col-sm-6 my-auto">
                      {{ dataInfo.Description }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 my-3 mb-5">
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
    </div>
  </div>
</template>

<script setup>
  import { getWebhook } from '@/api/webhook';
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const dataInfo = ref({});

  onMounted(() => {
    getWebhookInfo();
  });

  // get webhook info
  const getWebhookInfo = async () => {
    const params = router.currentRoute.value.query;
    const res = await getWebhook({
      id: params.id,
    });
    dataInfo.value = res.data;
  };

  // goto index page
  function goIndex() {
    router.push('/admin/webhook/index');
  }
</script>

<style scoped></style>
