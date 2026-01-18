<script setup>
const props = defineProps({
  show: Boolean
})
</script>

<template>
  <!-- operation list-->
  <div class="form-group row py-2  px-1">
    <div class="col-sm-10">
      <a class="btn btn-outline-primary btn-sm" :href="[['/web/entity/create?table_code=' + tableCode]]" role="button">
        <i class="bi bi-file-earmark-plus"></i>
        New
      </a>
      <button type="button" class="btn btn-outline-primary btn-sm" data-bs-toggle="modal"
        data-bs-target="#massOperateModal" data-bs-operation="MF" data-bs-title="Mass Freeze">
        Freeze
      </button>
      <button type="button" class="btn btn-outline-primary btn-sm" data-bs-toggle="modal"
        data-bs-target="#massOperateModal" data-bs-operation="MT" data-bs-title="Mass Unfreeze">
        Unfreeze
      </button>
      <button type="button" class="btn btn-outline-primary btn-sm" data-bs-toggle="modal"
        data-bs-target="#massOperateModal" data-bs-operation="MD" data-bs-title="Mass Delete">
        <i class="bi bi-trash3"></i>
      </button>
      <!--
      <button type="button" class="btn btn-outline-primary btn-sm">
        Change Documents
      </button>
      <button type="button" class="btn btn-outline-primary btn-sm">
        Replication Status
      </button>
      <button type="button" class="btn btn-outline-primary btn-sm">
        Replicate
      </button>
      <button type="button" class="btn btn-outline-primary btn-sm">
        Create Cleansing Case
      </button> -->
      <!-- <button type="button" class="btn btn-outline-primary btn-sm" @click="massOperate('Copy')">
        <i class="bi bi-files"></i>
        Copy
      </button> -->
      <!-- <button type="button" class="btn btn-outline-primary btn-sm">
        <i class="bi bi-pencil"></i>
        Multi-Processsing
      </button> -->
      <!-- <button type="button" class="btn btn-outline-primary btn-sm">
        Change Requests
      </button> -->
      <!-- <button type="button" class="btn btn-outline-primary btn-sm" @click="selectFile()">
        Import
      </button> -->
      <button type="button" class="btn btn-outline-primary btn-sm" data-bs-toggle="modal" data-bs-target="#importModal">
        Import
      </button>
      <div class="btn-group">
        <button type="button" class="btn btn-outline-primary btn-sm dropdown-toggle" data-bs-toggle="dropdown"
          aria-expanded="false">
          Export
        </button>
        <ul class="dropdown-menu dropdown-menu-end">
          <li>
            <button type="button" class="dropdown-item" @click="exportFile('selected')">
              Selected
            </button>
          </li>
          <li>
            <button type="button" class="dropdown-item" @click="exportFile('filtered')">
              Filtered
            </button>
          </li>
          <li>
            <button type="button" class="dropdown-item" @click="exportFile('all')">
              All
            </button>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<style scoped>
.btn {
  margin-right: 5px;
}
</style>
