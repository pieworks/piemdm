<template>
  <div class="p-1">
    <form id="search">
      <div class="form-group row">
        <div class="form-group row mb-1 align-items-center" v-for="(criteria, index) in criterias" :key="index">
          <div class="col-2">
            <v-select label="Code" v-model="criteria.field" :options="fieldData" :placeholder="$t('Please Select')"
              :reduce="option => option.Code">
              <template v-slot:option="option">{{ option.Code }} {{ option.Name }}</template>
            </v-select>
          </div>

          <div class="col-1">
            <select class="form-select form-select-sm" v-model="criteria.symbol">
              <option value="=">=</option>
              <option value="!=">!=</option>
              <option value="like">like</option>
              <option value=">">&gt;</option>
              <option value="<">&lt;</option>
              <option value=">=">&gt;=</option>
              <option value="<=">&lt;=</option>
            </select>
          </div>

          <div class="col-3">
            <textarea class="form-control form-control-sm" v-model="criteria.value" rows="1"
              :placeholder="$t('Please enter search value')" style="resize: vertical; min-height: 31px;"></textarea>
          </div>

          <div class="col-2">
            <button type="button" class="btn btn-sm btn-outline-primary me-1" @click="addCriteria" v-if="index === 0">
              <i class="bi bi-plus"></i>
            </button>
            <button type="button" class="btn btn-sm btn-outline-danger" @click="deleteCriteria(index)" v-if="index > 0">
              <i class="bi bi-dash"></i>
            </button>
          </div>
        </div>
      </div>

      <div class="row mt-2">
        <div class="col-12">
          <button type="button" class="btn btn-sm btn-outline-primary me-2" @click="$emit('search')">
            <i class="bi bi-search"></i>
            {{ $t('Search') }}
          </button>
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleReset">
            <i class="bi bi-arrow-counterclockwise"></i>
            {{ $t('Reset') }}
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import vSelect from 'vue-select';
import 'vue-select/dist/vue-select.css';

const emit = defineEmits(['search']);

const props = defineProps({
  fieldData: {
    type: Array,
    default: () => [],
  },
  criterias: {
    type: Array,
    default: () => [
      {
        field: 'ID',
        symbol: '=',
        value: '',
      },
    ],
  },
});

const addCriteria = () => {
  props.criterias.push({
    symbol: '=',
  });
};

const deleteCriteria = key => {
  if (key <= 0) return;
  props.criterias.splice(key, 1);
};

const handleReset = () => {
  props.criterias.splice(1, props.criterias.length - 1);
  props.criterias[0].value = '';
};
</script>
