<template>
  <div
    v-if="total"
    class="form-group row mt-2 px-1"
  >
    <div
      class="col-sm-10"
      style="font-size:0.8rem"
    >
      {{ $t("Result List") }}: {{ startNumber }} to {{ endNumber }} of {{ total }}
      {{ $t("items") }}
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
const total = ref(0)
const startNumber = ref(0)
const endNumber = ref(0)

const props = defineProps({
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 2
  },
  total: {
    type: Number,
    default: 0,
  }
})

// defineEmits(['pageChange'])

// onMounted(() => {
//   changData(props)
// })

watch(props, (newProps) => {
  changData(newProps)
})

function changData(val) {
  total.value = val.total
  startNumber.value = (val.page - 1) * val.pageSize + 1
  endNumber.value = startNumber.value + val.pageSize - 1
  if (endNumber.value >= val.total) endNumber.value = val.total
}
</script>

<style></style>
