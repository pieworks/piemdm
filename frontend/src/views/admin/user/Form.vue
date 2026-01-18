<template>
  <form
    name="entityForm"
    id="entityForm"
    method="post"
    enctype="multipart/form-data"
  >
    <div class="col-12 col-sm-12 mb-2">
      <div class="row">
        <div class="col-sm-12">
          <div class="form-group row">
            <legend
              :class="['col-form-label', 'col-sm-2', { required: requiredFields.employeeid }]"
            >
              {{ $t('EmployeeID') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="employeeID"
                v-bind="employeeIDAttrs"
                name="employeeID"
                :placeholder="$t('EmployeeID')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.EmployeeID"
                class="text-danger small mt-1"
              >
                {{ errors.EmployeeID }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.username }]">
              {{ $t('Username') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                autocomplete="username"
                class="form-control form-control-sm"
                v-model="username"
                v-bind="usernameAttrs"
                name="username"
                :placeholder="$t('Username')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.Username"
                class="text-danger small mt-1"
              >
                {{ errors.Username }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.password }]">
              {{ $t('Password') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="password"
                autocomplete="current-password"
                class="form-control form-control-sm"
                v-model="password"
                v-bind="passwordAttrs"
                name="password"
                :placeholder="$t('Password')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.Password"
                class="text-danger small mt-1"
              >
                {{ errors.Password }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.firstname }]">
              {{ $t('FirstName') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="firstName"
                v-bind="firstNameAttrs"
                name="firstName"
                :placeholder="$t('FirstName')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.FirstName"
                class="text-danger small mt-1"
              >
                {{ errors.FirstName }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.lastname }]">
              {{ $t('LastName') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="lastName"
                v-bind="lastNameAttrs"
                name="lastName"
                :placeholder="$t('LastName')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.LastName"
                class="text-danger small mt-1"
              >
                {{ errors.LastName }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend
              :class="['col-form-label', 'col-sm-2', { required: requiredFields.displayName }]"
            >
              {{ $t('DisplayName') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="displayName"
                v-bind="displayNameAttrs"
                name="displayName"
                :placeholder="$t('DisplayName')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.DisplayName"
                class="text-danger small mt-1"
              >
                {{ errors.DisplayName }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.fullname }]">
              {{ $t('FullName') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="fullName"
                v-bind="fullNameAttrs"
                name="fullName"
                :placeholder="$t('FullName')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.FullName"
                class="text-danger small mt-1"
              >
                {{ errors.FullName }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.email }]">
              {{ $t('Email') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="email"
                class="form-control form-control-sm"
                v-model="email"
                v-bind="emailAttrs"
                name="email"
                :placeholder="$t('Email')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.Email"
                class="text-danger small mt-1"
              >
                {{ errors.Email }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.phone }]">
              {{ $t('Phone') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="phone"
                v-bind="phoneAttrs"
                name="phone"
                :placeholder="$t('Phone')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.Phone"
                class="text-danger small mt-1"
              >
                {{ errors.Phone }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.language }]">
              {{ $t('Language') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect
                v-model="language"
                v-bind="languageAttrs"
                name="language"
                :reduce="option => option.value"
                :placeholder="$t('Please Select')"
                :options="languageOptions"
              ></VSelect>
              <div
                v-if="errors.Language"
                class="text-danger small mt-1"
              >
                {{ errors.Language }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.sex }]">
              {{ $t('Sex') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect
                v-model="sex"
                v-bind="sexAttrs"
                name="sex"
                :reduce="option => option.value"
                :placeholder="$t('Please Select')"
                :options="sexOptions"
              ></VSelect>
              <div
                v-if="errors.Sex"
                class="text-danger small mt-1"
              >
                {{ errors.Sex }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.avatar }]">
              {{ $t('Avatar') }}:
            </legend>
            <div class="col-sm-auto">
              <input
                type="text"
                class="form-control form-control-sm"
                v-model="avatar"
                v-bind="avatarAttrs"
                name="avatar"
                :placeholder="$t('Avatar')"
                maxlength="64"
                size="64"
              />
              <div
                v-if="errors.Avatar"
                class="text-danger small mt-1"
              >
                {{ errors.Avatar }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.admin }]">
              {{ $t('Admin') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect
                v-model="admin"
                v-bind="adminAttrs"
                name="admin"
                :reduce="option => option.value"
                :placeholder="$t('Please Select')"
                :options="yesNoOptions"
              ></VSelect>
              <div
                v-if="errors.Admin"
                class="text-danger small mt-1"
              >
                {{ errors.Admin }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12">
          <div class="form-group row">
            <legend :class="['col-form-label', 'col-sm-2', { required: requiredFields.status }]">
              {{ $t('Status') }}:
            </legend>
            <div class="col-sm-3">
              <VSelect
                v-model="status"
                v-bind="statusAttrs"
                name="status"
                :reduce="option => option.value"
                :placeholder="$t('Please Select')"
                :options="statusOptions"
              ></VSelect>
              <div
                v-if="errors.Status"
                class="text-danger small mt-1"
              >
                {{ errors.Status }}
              </div>
            </div>
          </div>
        </div>
        <div class="col-sm-12 mb-1">
          <div class="form-group row">
            <legend
              :class="['col-form-label', 'col-sm-2', { required: requiredFields.description }]"
            >
              {{ $t('Description') }}:
            </legend>
            <div class="col-sm-6">
              <textarea
                class="form-control form-control-sm"
                v-model="description"
                v-bind="descriptionAttrs"
                name="description"
                :placeholder="$t('Description')"
                maxlength="255"
                rows="3"
              ></textarea>
              <div
                v-if="errors.Description"
                class="text-danger small mt-1"
              >
                {{ errors.Description }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="form-group row mt-3">
      <div class="col-sm-auto">
        <input
          v-if="props.dataInfo && props.dataInfo.ID"
          type="hidden"
          name="id"
          :value="props.dataInfo.ID"
        />
        <button
          type="button"
          class="btn btn-outline-primary btn-sm me-2"
          @click="onSubmit()"
        >
          {{ $t('Submit') }}
        </button>
        <button
          type="button"
          class="btn btn-outline-secondary btn-sm"
          @click="$emit('goIndex')"
        >
          {{ $t('Cancel') }}
        </button>
      </div>
    </div>
  </form>
</template>

<script setup>
  import { useFormOptions } from '@/composables/useFormOptions';
  import { useForm } from 'vee-validate';
  import { computed, watch } from 'vue';
  import VSelect from 'vue-select';
  import 'vue-select/dist/vue-select.css';
  import * as yup from 'yup';

  const { statusOptions, languageOptions, sexOptions, yesNoOptions } = useFormOptions();

  // Define props to receive dataInfo from parent component
  const props = defineProps({
    dataInfo: {
      type: Object,
      default: () => ({}),
    },
  });

  // 验证模式 - 直接使用平面结构
  const validationSchema = yup.object({
    EmployeeID: yup.string().required(),
    Username: yup.string().required(),
    Password: yup.string().required().min(6),
    FirstName: yup.string(),
    LastName: yup.string(),
    DisplayName: yup.string(),
    Email: yup.string().email(),
    Phone: yup.string(),
    Language: yup.string(),
    Sex: yup.string(),
    Avatar: yup.string(),
    Admin: yup.string().required(),
    Status: yup.string().required(),
    Description: yup.string(),
  });

  // required字段映射
  const requiredFields = computed(() => {
    const requiredMap = {};
    Object.keys(validationSchema.fields).forEach(key => {
      const field = validationSchema.fields[key];
      requiredMap[key.toLowerCase()] = field.tests.some(test => test.OPTIONS?.name === 'required');
    });
    return requiredMap;
  });

  // 表单初始化
  const { values, errors, defineField, handleSubmit, setValues } = useForm({
    validationSchema,
    initialValues: props.dataInfo,
  });

  // 字段定义
  const [employeeID, employeeIDAttrs] = defineField('EmployeeID');
  const [username, usernameAttrs] = defineField('Username');
  const [password, passwordAttrs] = defineField('Password');
  const [firstName, firstNameAttrs] = defineField('FirstName');
  const [lastName, lastNameAttrs] = defineField('LastName');
  const [displayName, displayNameAttrs] = defineField('DisplayName');
  const [email, emailAttrs] = defineField('Email');
  const [phone, phoneAttrs] = defineField('Phone');
  const [language, languageAttrs] = defineField('Language');
  const [sex, sexAttrs] = defineField('Sex');
  const [avatar, avatarAttrs] = defineField('Avatar');
  const [admin, adminAttrs] = defineField('Admin');
  const [status, statusAttrs] = defineField('Status');
  const [description, descriptionAttrs] = defineField('Description');

  const emit = defineEmits(['submitForm', 'goIndex']);

  const onSubmit = handleSubmit(values => {
    emit('submitForm', values);
  });

  // 监听 dataInfo 变化并更新表单值
  watch(
    () => props.dataInfo,
    newDataInfo => {
      if (newDataInfo && Object.keys(newDataInfo).length > 0) {
        setValues(newDataInfo);
      }
    },
    { immediate: true, deep: true }
  );
</script>

<style scoped>
  .required:after {
    content: ' *';
    color: #dc3545;
    font-weight: bold;
  }
</style>
