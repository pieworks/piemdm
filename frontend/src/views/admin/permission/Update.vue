<template>
  <div class="mt-3">
    <div class="border-bottom fs-5 mb-2 pb-2">{{ $t('Permission') }} {{ $t('Update') }}</div>
    <div class="card-body mt-3" style="min-height: calc(100vh - 160px); font-size: 0.9rem">
      <div id="create_wrapper">
        <div class="row">
          <div class="col-sm-12">
            <AppForm />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import AppForm from './Form.vue';
</script>
