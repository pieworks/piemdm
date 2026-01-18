<template>
  <div class="mt-3" style="min-height: calc(100vh - 105px); font-size: 0.9rem">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('User') }}{{ $t('View') }}</div>
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
                      for="ID"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('ID') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.ID }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="EmployeeID"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('EmployeeID') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.EmployeeID }}
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
                    <div class="col-sm-auto">
                      {{ dataInfo.Username }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="Password"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Password') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.Password }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="FirstName"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('FirstName') }}:
                    </legend>
                    <div class="col-sm-auto">
                      {{ dataInfo.FirstName }}
                    </div>
                  </div>
                </div>

                <div class="col-sm-12">
                  <div class="form-group row">
                    <legend
                      for="LastName"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('LastName') }}:
                    </legend>
                    <div class="col-sm-3">
                      {{ dataInfo.LastName }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="DisplayName"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('DisplayName') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.DisplayName }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Email"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Email') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.Email }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Phone"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Phone') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.Phone }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Language"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Language') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.Language }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Sex"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Sex') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.Sex }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Avatar"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Avatar') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.Avatar }}
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
                    <div class="col-sm-6">
                      {{ dataInfo.Desc }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="Admin"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('Admin') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.Admin }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 mb-1">
                  <div class="form-group row">
                    <legend
                      for="LastLogin"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('LastLogin') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.LastLogin }}
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
                      for="CreatedBy"
                      class="col-form-label col-sm-2"
                    >
                      {{ $t('CreatedBy') }}:
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
                      {{ $t('UpdatedBy') }}:
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
                      {{ $t('CreatedAt') }}:
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
                      {{ $t('UpdatedAt') }}:
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
                      {{ $t('DeletedAt') }}:
                    </legend>
                    <div class="col-sm-6">
                      {{ dataInfo.DeletedAt }}
                    </div>
                  </div>
                </div>
                <div class="col-sm-12 my-3 mb-3">
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
  import { getUser } from '@/api/user';
  import { onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const dataInfo = ref({});

  onMounted(() => {
    getUserInfo();
  });

  // get user info
  const getUserInfo = async () => {
    const params = router.currentRoute.value.query;
    const res = await getUser({
      id: params.id,
    });
    dataInfo.value = res.data;
  };

  // goto index page
  function goIndex() {
    router.push('/admin/user/index');
  }
</script>

<style scoped></style>
