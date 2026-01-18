<template>
    <div class="category-tree card h-100 border-0">
        <div class="card-header bg-transparent border-bottom-0 py-2 d-flex justify-content-between align-items-center">
            <small class="mb-0 text-muted fw-semibold">{{ title }}</small>
            <button class="btn btn-sm btn-link text-decoration-none p-0" @click="fetchData" :title="$t('Refresh')">
                <i class="bi bi-arrow-clockwise"></i>
            </button>
        </div>
        <div class="card-body p-0 overflow-auto custom-scrollbar">
            <div v-if="loading" class="text-center py-3 text-muted">
                <div class="spinner-border spinner-border-sm" role="status"></div>
            </div>
            <div v-else-if="roots.length === 0" class="text-center py-3 text-muted small">
                {{ $t('No Data') }}
            </div>
            <ul v-else class="list-unstyled mb-0 tree-list">
                <tree-node v-for="node in roots" :key="node.id" :node="node" :level="0" :selected-id="selectedId"
                    @select="onSelect" @toggle="onToggle">
                </tree-node>
            </ul>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed, defineComponent, h, watch } from 'vue';
import { getEntityList } from '@/api/entity';
import { getTableList } from '@/api/table';

// Recursive Sub-component for Node
const TreeNode = defineComponent({
    name: 'TreeNode',
    props: {
        node: Object,
        level: Number,
        selectedId: [String, Number]
    },
    emits: ['select', 'toggle'],
    setup(props, { emit }) {
        const toggle = () => {
            if (props.node.children && props.node.children.length > 0) {
                if (!props.node.isExpanded) props.node.isExpanded = false;
                props.node.isExpanded = !props.node.isExpanded;
            }
            emit('toggle', props.node);
        };

        const select = () => {
            emit('select', props.node);
        };

        return { toggle, select };
    },
    render() {
        const { node, level, selectedId } = this;
        const hasChildren = node.children && node.children.length > 0;

        return h('li', { class: 'tree-item' }, [
            h('div', {
                class: {
                    'tree-content': true,
                    'd-flex': true,
                    'align-items-center': true,
                    'py-1': true,
                    'px-2': true,
                    'active': (node.code && String(selectedId) === String(node.code)) || String(selectedId) === String(node.id),
                    'ps-0': true
                },
                style: { paddingLeft: (level * 20) + 'px !important' },
                onClick: this.select
            }, [
                h('span', {
                    class: 'tree-toggle-icon me-1 text-center',
                    style: { width: '20px', cursor: 'pointer' },
                    onClick: (e) => {
                        e.stopPropagation();
                        this.toggle();
                    }
                }, [
                    hasChildren
                        ? h('i', {
                            class: ['bi', node.isExpanded ? 'bi-caret-down-fill' : 'bi-caret-right-fill', 'text-muted']
                        })
                        : h('i', { class: 'bi bi-dot opacity-25' })
                ]),
                h('i', {
                    class: [
                        'bi',
                        hasChildren
                            ? (node.isExpanded ? 'bi-folder2-open' : 'bi-folder')
                            : 'bi-file-text',
                        'me-2'
                    ],
                    style: { color: hasChildren ? '#0d6efd' : '#6c757d' }
                }),
                h('span', {
                    class: 'text-truncate flex-grow-1 user-select-none',
                    title: node.name || node.code
                }, node.name || node.code)
            ]),
            hasChildren && node.isExpanded
                ? h('ul', { class: 'list-unstyled mb-0' },
                    node.children.map(child =>
                        h(TreeNode, {
                            key: child.id,
                            node: child,
                            level: level + 1,
                            selectedId: selectedId,
                            onSelect: (e) => this.$emit('select', e),
                            onToggle: (e) => this.$emit('toggle', e)
                        })
                    )
                )
                : null
        ]);
    }
});

const props = defineProps({
    tableCode: {
        type: String,
        required: true
    },
    selectedId: {
        type: [String, Number],
        default: null
    }
});

const emit = defineEmits(['select']);

const loading = ref(false);
const rawData = ref([]);
const selectedId = ref(null);
const title = ref('');
const forceUpdateKey = ref(0);

const roots = computed(() => {
    // Access forceUpdateKey to trigger recomputation
    forceUpdateKey.value;

    const items = rawData.value;
    const map = {};
    const rootNodes = [];

    // Initialize - work directly with items to maintain reactivity
    items.forEach(item => {
        if (item.isExpanded === undefined) item.isExpanded = true;
        map[item.id] = item;
        // Reset children array
        item.children = [];
    });

    // Build tree
    items.forEach(item => {
        if (item.parent_id && map[item.parent_id]) {
            map[item.parent_id].children.push(item);
        } else {
            rootNodes.push(item);
        }
    });

    return rootNodes;
});

const fetchData = async () => {
    if (!props.tableCode) return;

    loading.value = true;
    try {
        // 1. Get Table Name
        getTableList({ code: props.tableCode }).then(res => {
            if (res && res.data && res.data.length > 0) title.value = res.data[0].Name;
        });

        // 2. Get Data (assume max 10000 for category tree)
        const res = await getEntityList({
            table_code: props.tableCode,
            page: 1,
            pageSize: 10000
        });

        if (res && res.data) {
            rawData.value = res.data.map(item => ({ ...item, isExpanded: true })); // Default expand all? Or false.
            // Default expand all is better for sidebar trees usually.
        }
    } catch (err) {
        console.error(err);
    } finally {
        loading.value = false;
    }
};

const onSelect = (node) => {
    selectedId.value = node.code || node.id;
    emit('select', node);
};

const onToggle = (node) => {
    // Force re-render by incrementing the key
    forceUpdateKey.value++;
};

// Watch for external selectedId changes
watch(() => props.selectedId, (newId) => {
    selectedId.value = newId;
}, { immediate: true });

onMounted(() => {
    fetchData();
});
</script>

<style>
.category-tree .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
}

.category-tree .custom-scrollbar::-webkit-scrollbar-thumb {
    background-color: rgba(0, 0, 0, 0.1);
    border-radius: 4px;
}

.category-tree .tree-content {
    cursor: pointer;
    border-radius: 4px;
    transition: background-color 0.2s;
    font-size: 0.8rem;
}

.category-tree .tree-content:hover {
    background-color: #f8f9fa;
}

.category-tree .tree-content.active {
    background-color: #e7f1ff;
    color: #0d6efd;
}

.category-tree .tree-toggle-icon {
    font-size: 0.75rem;
}
</style>
