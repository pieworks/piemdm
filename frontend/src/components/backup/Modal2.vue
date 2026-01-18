<script setup>
  import { ref, watch } from 'vue';
  const list = ref([]);

  const props = defineProps({
    show: Boolean,
    title: {
      type: String,
      default: 'Title',
    },
    list: {
      type: Object,
      default: {},
    },
  });

  // onMounted(() => {
  //   changData(props)
  // })

  watch(props, newProps => {
    // changData(newProps)
  });

  // 获取审批流任务
  // TODO 如果code为空，增加提示
  function getProcessInfo() {
    // spinner.style.display = "block"
    // if (code === "") return;
    // let url = '/web/workflow/task?resultType=json';
    // url = url + "&code=" + code;
    // const statusToast = document.getElementById('statusToast')
    // const toastBody = statusToast.querySelector('.toast-body')
    // var toast = new bootstrap.Toast(statusToast)
    // const response = await fetch(url);
    // if (response.ok) {
    //   const res = await response.json()
    //   if (res.code == 200) {
    //     this.tasks = res.data.data
    // const tm = new bootstrap.Modal(document.getElementById('taskModal'))
    // tm.show()
    //   } else {
    //     toast.show()
    //   }
    // } else {
    //   toastBody.innerHTML = response.statusText
    //   toast.show()
    // }
    // spinner.style.display = "none"
  }
</script>

<template>
  <!-- Modal -->
  <div
    v-if="show"
    class="modal-mask"
  >
    <div
      class="modal"
      data-bs-backdrop="static"
      data-bs-keyboard="true"
      tabindex="-1"
      style="display: block"
    >
      <div class="modal-dialog modal-lg modal-dialog-centered modal-dialog-scrollable">
        <div class="modal-content">
          <div class="modal-header">
            <h1
              class="modal-title fs-5"
              id="taskModalLabel"
            >
              Process Information
            </h1>
            <button
              type="button"
              class="btn-close"
              data-bs-dismiss="modal"
              @click="$emit('close')"
            ></button>
          </div>
          <div class="modal-body">
            <table class="table table-sm table-bordered table-hover">
              <thead>
                <tr>
                  <th>#</th>
                  <th>NodeName</th>
                  <th>UserName</th>
                  <th>OperateType</th>
                  <th>CreatedAt</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(t, i) in tasks">
                  <td>[[ t.ID ]]</td>
                  <td>[[t.NodeName]]</td>
                  <td>[[t.UserName]]</td>
                  <td>[[t.OperateType]]</td>
                  <td>[[formatDatetime(t.CreatedAt)]]</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-outline-primary btn-sm"
              data-bs-dismiss="modal"
              @click="$emit('close')"
            >
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
  .modal-mask {
    position: fixed;
    z-index: 9998;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    transition: opacity 0.3s ease;
  }
</style>
