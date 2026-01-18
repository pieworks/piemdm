<template>
  <div
    :class="['modal', animation ? 'fade' : '']"
    tabindex="-1"
    ref="modalEle"
    @keydown.esc="onCancel"
  >
    <div
      :class="[
        'modal-dialog',
        centered ? 'modal-dialog-centered' : '',
        'modal-' + size,
        scrollable ? 'modal-dialog-scrollable' : '',
      ]"
    >
      <div class="modal-content">
        <div class="modal-header">
          <slot name="header">
            <h5 class="modal-title">{{ title }}</h5>
          </slot>
          <button
            type="button"
            class="btn-close"
            @click="onCancel"
            aria-label="Close"
          ></button>
        </div>
        <div class="modal-body">
          <div
            v-if="bodyHtml"
            v-html="bodyContent"
          ></div>
          <slot v-else>{{ bodyContent }}</slot>
        </div>
        <div class="modal-footer">
          <slot name="footer">
            <button
              type="button"
              class="btn btn-secondary btn-sm"
              @click="onCancel"
              v-if="showCancelButton"
            >
              {{ cancelTitle }}
            </button>
            <button
              v-if="okTitle"
              type="button"
              class="btn btn-primary btn-sm"
              @click="onConfirm"
            >
              {{ okTitle }}
            </button>
          </slot>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { Modal } from 'bootstrap';
  import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';

  const props = defineProps({
    show: Boolean,
    title: {
      type: String,
      default: 'Title',
    },
    bodyContent: {
      type: String,
      default: 'Content',
    },
    animation: {
      type: Boolean,
      default: true,
    },
    cancelTitle: {
      type: String,
      default: '取消',
    },
    okTitle: {
      type: String,
      default: '确定',
    },
    centered: {
      type: Boolean,
      default: true,
    },
    scrollable: {
      type: Boolean,
      default: true,
    },
    size: {
      type: String,
      default: 'md',
    },
    bodyHtml: {
      type: Boolean,
      default: false,
    },
    showCancelButton: {
      type: Boolean,
      default: true,
    },
  });

  const emit = defineEmits(['update:show', 'confirm', 'cancel']);
  const modalEle = ref(null);
  let modalObj = null;

  function onCancel() {
    emit('cancel');
    nextTick(() => {
      if (modalObj) modalObj.hide();
    });
  }

  function onConfirm() {
    emit('confirm');
    nextTick(() => {
      if (modalObj) modalObj.hide();
    });
  }

  watch(
    () => props.show,
    val => {
      if (modalObj) {
        if (val) modalObj.show();
        else modalObj.hide();
      }
    }
  );

  onMounted(() => {
    // 使用 Bootstrap 默认配置,让它自己管理 backdrop
    modalObj = new Modal(modalEle.value, {
      backdrop: 'static',
      keyboard: true
    });

    if (props.show) modalObj.show();

    modalEle.value.addEventListener('hidden.bs.modal', () => {
      emit('update:show', false);
    });
  });

  onBeforeUnmount(() => {
    // 确保在组件销毁前清理 Modal 实例
    if (modalObj) {
      modalObj.dispose();
      modalObj = null;
    }
  });
</script>

<style scoped>
  /* 移除所有自定义样式,完全使用 Bootstrap 原生样式 */
</style>
