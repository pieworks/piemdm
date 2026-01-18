<template>
  <div class="upload-container">
    <!-- 文件列表 -->
    <div v-if="fileList.length > 0" class="file-list mb-2">
      <div v-for="(file, index) in fileList" :key="index" class="file-item">
        <!-- 文件图标 -->
        <i class="bi bi-file-earmark file-icon"></i>

        <!-- 文件名 -->
        <span class="file-name" :title="file.name">{{ file.name }}</span>

        <!-- 文件大小 -->
        <span class="file-size">{{ formatSize(file.size) }}</span>

        <!-- 删除按钮 -->
        <button
          @click="removeFile(index)"
          class="btn-remove"
          type="button"
          :disabled="disabled"
        >
          <i class="bi bi-x-lg"></i>
        </button>
      </div>
    </div>

    <!-- 上传按钮 -->
    <div v-if="!maxReached" class="upload-button">
      <input
        ref="fileInput"
        type="file"
        :accept="accept"
        :multiple="multiple"
        @change="handleFileChange"
        style="display: none"
        :disabled="disabled"
      />
      <button
        @click="$refs.fileInput.click()"
        class="btn btn-outline-primary btn-sm"
        type="button"
        :disabled="disabled || uploading"
      >
        <span v-if="uploading" class="spinner-border spinner-border-sm me-1"></span>
        <i v-else class="bi bi-upload me-1"></i>
        {{ uploading ? '上传中...' : '选择文件' }}
      </button>
      <small v-if="accept" class="text-muted ms-2">
        支持: {{ accept }}
      </small>
      <small v-if="maxSize" class="text-muted ms-2">
        最大: {{ formatSize(maxSize) }}
      </small>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { uploadFile } from '@/api/upload';
import { AppModal } from '@/components/Modal/modal.js';

const props = defineProps({
  modelValue: {
    type: [String, Array],
    default: ''
  },
  multiple: {
    type: Boolean,
    default: false
  },
  accept: {
    type: String,
    default: ''
  },
  maxSize: {
    type: Number,
    default: 10 * 1024 * 1024 // 10MB
  },
  maxFiles: {
    type: Number,
    default: 10
  },
  disabled: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue']);

const fileList = ref([]);
const fileInput = ref(null);
const uploading = ref(false);

const maxReached = computed(() => {
  return props.multiple && fileList.value.length >= props.maxFiles;
});

// 辅助函数 - 必须在 watch 之前定义
const getFilenameFromUrl = (url) => {
  if (!url) return '';
  const parts = url.split('/');
  return parts[parts.length - 1] || 'file';
};

const formatSize = (bytes) => {
  if (!bytes) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
};

const validateFileType = (file) => {
  if (!props.accept) return true;

  const acceptTypes = props.accept.split(',').map(t => t.trim());
  const fileName = file.name.toLowerCase();
  const fileType = file.type.toLowerCase();

  return acceptTypes.some(type => {
    // 通配符匹配 (如 image/*)
    if (type.includes('*')) {
      const prefix = type.split('/')[0];
      return fileType.startsWith(prefix + '/');
    }
    // 扩展名匹配 (如 .pdf)
    if (type.startsWith('.')) {
      return fileName.endsWith(type.toLowerCase());
    }
    // MIME 类型匹配
    return fileType === type;
  });
};

// 初始化文件列表
watch(() => props.modelValue, (newValue) => {
  if (!newValue) {
    fileList.value = [];
    return;
  }

  // 处理多文件
  if (Array.isArray(newValue)) {
    fileList.value = newValue.map(url => ({
      url,
      name: getFilenameFromUrl(url),
      size: 0
    }));
  }
  // 处理单文件
  else if (typeof newValue === 'string' && newValue) {
    fileList.value = [{
      url: newValue,
      name: getFilenameFromUrl(newValue),
      size: 0
    }];
  }
}, { immediate: true });

const handleFileChange = async (event) => {
  const files = Array.from(event.target.files);

  // 单文件模式:清空现有文件列表(替换模式)
  // 多文件模式:保留现有文件列表(累加模式)
  if (!props.multiple) {
    fileList.value = [];
  }

  for (const file of files) {
    // 验证文件大小
    if (props.maxSize && file.size > props.maxSize) {
      AppModal.alert({
        title: '文件大小超限',
        bodyContent: `文件 ${file.name} 超过大小限制 (${formatSize(props.maxSize)})`,
      });
      continue;
    }

    // 验证文件类型
    if (props.accept && !validateFileType(file)) {
      AppModal.alert({
        title: '文件类型不支持',
        bodyContent: `文件 ${file.name} 类型不支持`,
      });
      continue;
    }

    // 上传文件
    try {
      uploading.value = true;
      const formData = new FormData();
      formData.append('file', file);

      const res = await uploadFile(formData);

      console.log('Upload response:', res);
      console.log('Upload response data:', res.data);

      // 后端返回: res.data = {url, filename, size}
      if (!res.data || !res.data.url) {
        throw new Error('服务器返回数据格式错误');
      }

      fileList.value.push({
        name: file.name,
        url: res.data.url,
        size: file.size
      });

      updateValue();

      // 单文件模式:只上传第一个文件后就退出
      if (!props.multiple) {
        break;
      }
    } catch (error) {
      console.error('Upload error:', error);
      AppModal.alert({
        title: '文件上传失败',
        bodyContent: `文件 ${file.name} 上传失败: ${error.response?.data?.error || error.message || '未知错误'}`,
      });
    } finally {
      uploading.value = false;
    }
  }

  // 清空 input
  event.target.value = '';
};

const removeFile = (index) => {
  fileList.value.splice(index, 1);
  updateValue();
};

const updateValue = () => {
  if (props.multiple) {
    emit('update:modelValue', fileList.value.map(f => f.url));
  } else {
    emit('update:modelValue', fileList.value[0]?.url || '');
  }
};
</script>

<style scoped>
.upload-container {
  width: 100%;
}

.file-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 0.5rem;
  border: 1px solid #dee2e6;
  border-radius: 0.25rem;
  background-color: #f8f9fa;
}

.file-icon {
  font-size: 1.5rem;
  color: #6c757d;
  margin-right: 0.75rem;
  width: 40px;
  text-align: center;
}

.file-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.875rem;
}

.file-size {
  color: #6c757d;
  font-size: 0.75rem;
  margin-left: 0.5rem;
  margin-right: 0.5rem;
}

.btn-remove {
  border: none;
  background: none;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  color: #dc3545;
  font-size: 0.875rem;
}

.btn-remove:hover:not(:disabled) {
  color: #bb2d3b;
}

.btn-remove:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.upload-button {
  display: flex;
  align-items: center;
}
</style>
