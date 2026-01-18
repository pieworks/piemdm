<template>
  <!-- 横向滚动容器 -->
  <div class="d-flex py-4 border-top border-bottom">
    <!-- 第一个审批节点 -->
    <div
      class="d-flex flex-column align-items-center text-center position-relative pe-4"
      v-for="(task, index) in approvalTasks"
      :key="task.id"
    >
      <!-- 连接线到上一个节点 -->
      <div
        class="position-absolute bg-success"
        style="top: 1rem; left: 0rem; width: calc(50% - 1.7rem); height: 2px; z-index: 1"
        v-if="index > 0"
      ></div>
      <button
        type="button"
        :class="
          task.Status === 'Approved' || task.Status === 'Done'
            ? 'btn btn-sm btn-success rounded-pill'
            : task.Status === 'Pending'
              ? 'btn btn-sm btn-warning rounded-pill'
              : 'btn btn-sm btn-secondary rounded-pill'
        "
        style="width: 2rem; height: 2rem"
      >
        {{ index + 1 }}
      </button>
      <div class="small">
        <div class="fw-bold">{{ task.NodeName }}</div>
        <div v-if="index < approvalTasks.length">
          <div
            class="mt-1"
            v-if="task.IsAgentTask"
          >
            <span
              class="badge"
              :class="
                task.Status === 'Approved' || task.Status === 'Done'
                  ? 'bg-success'
                  : task.Status === 'Pending'
                    ? 'bg-warning'
                    : 'bg-secondary'
              "
            >
              {{ task.AgentName + '(代)' || '&nbsp;' }}
            </span>
          </div>
          <div
            class="mt-1"
            v-else
          >
            <span
              class="badge"
              :class="
                task.Status === 'Approved' || task.Status === 'Done'
                  ? 'bg-success'
                  : task.Status === 'Pending'
                    ? 'bg-warning'
                    : 'bg-secondary'
              "
            >
              {{ task.AssigneeName || '&nbsp;' }}
            </span>
          </div>
          <div>
            {{
              task.Status === 'Approved' || task.Status === 'Done'
                ? task.Comment
                : task.Status === 'Pending'
                  ? '审批中'
                  : '未开始'
            }}
          </div>
          <div class="text-xs">{{ formatDate(task.CreatedAt, 'medium') || '' }}</div>
        </div>
      </div>
      <!-- 连接线到下一个节点 -->
      <div
        class="position-absolute bg-success"
        style="top: 1rem; left: 50%; width: 50%; height: 2px; z-index: 1"
        v-if="index < approvalTasks.length - 1"
      ></div>
    </div>
  </div>
</template>

<script setup>
  import { formatDate } from '@/utils/date';

  const props = defineProps({
    approvalNodes: {
      type: Array,
      default: () => [],
    },
    approvalTasks: {
      type: Array,
      default: () => [],
    },
  });

  const getApprovalTask = nodeCode => {
    const approvalTask = props.approvalTasks.find(task => task.NodeId === nodeCode);
    return approvalTask;
  };
</script>

<style scoped></style>
