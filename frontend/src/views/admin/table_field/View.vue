<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Table Field') }} {{ $t('View') }}</div>
    <div
      class="card-body mt-3"
      style="min-height: calc(100vh - 160px); font-size: 0.9rem"
    >
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
                      for="Code"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Code') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Code }}
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
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Name }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="NameEn"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('English Name') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.NameEn }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Type"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Type') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Type }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Length"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Length') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Length }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="IsUnique"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Is Unique') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.IsUnique }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Desc"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Desc') }}:
                    </legend>
                    <div class="col-sm-6 my-auto">
                      {{ dataInfo.Desc }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Required"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Required') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Required }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="IsFilter"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Is Filter') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.IsFilter }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="IsShow"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Is Show') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.IsShow }}
                    </div>
                  </div>
                </div>
                <!-- 废弃字段已移除,新配置存储在 Options JSON 中 -->
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Sort"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Sort') }}:
                    </legend>
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Sort }}
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
                    <div class="col-sm-auto my-auto">
                      {{ dataInfo.Status }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 my-3 mb-5">
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
  import { getTableField } from '@/api/table_field';
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const dataInfo = ref({});
  const params = ref({});

  onMounted(() => {
    params.value = router.currentRoute.value.query;
    getTableFieldInfo();
  });

  // get table field info
  const getTableFieldInfo = async () => {
    const res = await getTableField({
      id: params.value.id,
    });
    dataInfo.value = res.data;
  };

  // goto index page
  function goIndex() {
    router.push('/admin/table_field/index?table_code=' + params.value.table_code);
  }
</script>

<style scoped></style>
