<template>
  <!-- pagination -->
  <nav
    v-if="total"
    class="Page navigation p-1"
  >
    <ul class="pagination pagination-sm justify-content-start">
      <li class="page-item">
        <a
          class="page-link"
          href="javascript:;"
          @click="$emit('pageChange', prevPage)"
        >
          <span aria-hidden="true">&laquo;</span>
        </a>
      </li>
      <li
        class="page-item"
        v-if="minPage > 1"
      >
        <a
          class="page-link"
          href="javascript:;"
          @click="$emit('pageChange', 1)"
        >
          1
        </a>
      </li>
      <li
        class="page-item disabled"
        v-if="minPage > 1"
      >
        <a class="page-link">...</a>
      </li>
      <li
        :class="['page-item', page == p ? ' active' : '']"
        v-for="(p, index) in pages"
      >
        <a
          class="page-link"
          href="javascript:;"
          @click="$emit('pageChange', p)"
        >
          {{ p }}
        </a>
      </li>
      <li
        class="page-item"
        v-if="maxPage < pageCount"
      >
        <a class="page-link disabled">...</a>
      </li>
      <li
        class="page-item"
        v-if="maxPage < pageCount"
      >
        <a
          class="page-link"
          href="javascript:;"
          @click="$emit('pageChange', pageCount)"
        >
          {{ pageCount }}
        </a>
      </li>
      <li class="page-item">
        <a
          class="page-link"
          href="javascript:;"
          @click="$emit('pageChange', nextPage)"
        >
          <span aria-hidden="true">&raquo;</span>
        </a>
      </li>
      <li
        class="page-item my-auto ms-3"
        style="font-size: 0.8rem"
      >
        {{ $t('Showing') }} {{ startNumber }} {{ $t('to') }} {{ endNumber }} {{ $t('of') }}
        {{ total }} {{ $t('items') }}
      </li>
    </ul>
  </nav>
</template>

<script setup>
  import { ref, watch } from 'vue';
  const page = ref(1);
  const total = ref(0);
  const pageCount = ref(0);
  const minPage = ref(0);
  const maxPage = ref(0);
  const prevPage = ref(1);
  const nextPage = ref(0);
  const startNumber = ref(0);
  const endNumber = ref(0);
  const pages = ref([]);

  const props = defineProps({
    page: {
      type: Number,
      default: 1,
    },
    pageSize: {
      type: Number,
      default: 2,
    },
    total: {
      type: Number,
      default: 0,
    },
  });

  watch(props, newProps => {
    changData(newProps);
  });

  function changData(val) {
    pageCount.value = Math.ceil(val.total / val.pageSize);
    total.value = val.total;
    page.value = val.page;
    minPage.value = val.page - 5;
    if (minPage.value <= 1) minPage.value = 1;
    maxPage.value = minPage.value + 9;
    if (maxPage.value >= pageCount.value) maxPage.value = pageCount.value;
    startNumber.value = (val.page - 1) * val.pageSize + 1;
    endNumber.value = startNumber.value + val.pageSize - 1;
    if (endNumber.value >= val.total) endNumber.value = val.total;
    prevPage.value = val.page - 1;
    if (prevPage.value <= 1) prevPage.value = 1;
    nextPage.value = page.value + 1;
    if (nextPage.value >= pageCount.value) nextPage.value = pageCount.value;

    pages.value = [];
    for (var i = minPage.value; i < maxPage.value + 1; i++) {
      pages.value[i - minPage.value] = i;
    }
  }
</script>

<style></style>
