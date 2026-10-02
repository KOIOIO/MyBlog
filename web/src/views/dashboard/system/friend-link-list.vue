<template>
  <div class="friend-link-list">
    <div class="page-header">
      <div>
        <div class="page-title">{{ t('system.friendLink.title') }}</div>
        <div class="page-desc">{{ t('system.friendLink.desc') }}</div>
      </div>
      <div class="page-actions">
        <el-button type="primary" icon="Plus" @click="layoutStore.state.friendLinkCreateVisible = true">
          {{ t('system.friendLink.create') }}
        </el-button>

        <el-dialog
            v-model="friendLinkCreateVisible"
            width="500"
            align-center
            destroy-on-close
            :before-close="friendLinkCreateVisibleSynchronization"
        >
          <template #header>
            {{ t('system.friendLink.create') }}
          </template>
          <friend-link-create-form/>
          <template #footer>
          </template>
        </el-dialog>

        <el-button type="danger" plain icon="Delete" @click="friendLinkBulkDeleteVisible = true;handleIdsToDelete()">
          {{ t('common.batchDelete') }}
        </el-button>

        <el-dialog
            v-model="friendLinkBulkDeleteVisible"
            width="500"
            align-center
            destroy-on-close
        >
          <template #header>
            {{ t('system.friendLink.delete') }}
          </template>
          {{ t('system.friendLink.confirmDeleteWithCount', {count: idsToDelete.length}) }}
          <template #footer>
            <el-button type="primary" @click="handleBulkDelete(idsToDelete)">
              {{ t('common.confirm') }}
            </el-button>
            <el-button @click="friendLinkBulkDeleteVisible = false">{{ t('common.cancel') }}</el-button>
          </template>
        </el-dialog>
      </div>
    </div>

    <div class="friend-link-list-request">
      <el-form :inline="true" :model="friendLinkListRequest">
        <el-form-item :label="t('system.friendLink.name')">
          <el-input v-model="friendLinkListRequest.name" :placeholder="t('system.friendLink.namePlaceholder')" clearable/>
        </el-form-item>
        <el-form-item :label="t('system.friendLink.description')">
          <el-input v-model="friendLinkListRequest.description" :placeholder="t('system.friendLink.descriptionPlaceholder')" clearable/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="getFriendLinkTableData">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
        ref="multipleFriendLinkTableRef"
        :data="friendLinkTableData"
    >
      <el-table-column type="selection" width="60"/>
      <el-table-column label="Logo">
        <template #default="scope:{ row: any, column: any, $index: number }">
          <el-image :src="scope.row.logo" alt=""/>
        </template>
      </el-table-column>
      <el-table-column prop="link" :label="t('system.friendLink.link')"/>
      <el-table-column prop="name" :label="t('common.name')"/>
      <el-table-column prop="description" :label="t('system.friendLink.descCol')"/>
      <el-table-column :label="t('common.actions')">
        <template #default="scope:{ row: any, column: any, $index: number }">
          <el-button
              link
              type="primary"
              @click="layoutStore.state.friendLinkUpdateVisible=true;friendLinkInfo=scope.row"
          >
            {{ t('system.friendLink.updateAction') }}
          </el-button>
          <el-button
              link
              type="danger"
              @click="friendLinkDeleteVisible=true;friendLinkInfo=scope.row"
          >
            {{ t('common.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
        v-model="friendLinkUpdateVisible"
        width="500"
        align-center
        destroy-on-close
        :before-close="friendLinkUpdateVisibleSynchronization"
    >
      <template #header>
        {{ t('system.friendLink.update') }}
      </template>
      <friend-link-update-form :friendLink=friendLinkInfo />
      <template #footer>
      </template>
    </el-dialog>

    <el-dialog
        v-model="friendLinkDeleteVisible"
        width="500"
        align-center
        destroy-on-close
    >
      <template #header>
        {{ t('system.friendLink.delete') }}
      </template>
      {{ t('system.friendLink.confirmDeleteWithCount', {count: 1}) }}
      <template #footer>
        <el-button type="primary" @click="handleDelete(friendLinkInfo.id)">
          {{ t('common.confirm') }}
        </el-button>
        <el-button @click="friendLinkDeleteVisible = false">{{ t('common.cancel') }}</el-button>
      </template>
    </el-dialog>

    <el-pagination
        :current-page="page"
        :page-size="page_size"
        :page-sizes="[10, 30, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
    />
  </div>
</template>

<script setup lang="ts">
import {nextTick, onMounted, reactive, ref, watch} from 'vue'
import {
  type FriendLinkListRequest,
  type FriendLink,
  friendLinkDelete,
  type FriendLinkDeleteRequest
} from "@/api/friend-link";
import {friendLinkList} from "@/api/friend-link";
import FriendLinkCreateForm from "@/components/forms/FriendLinkCreateForm.vue";
import {useLayoutStore} from "@/stores/layout";
import FriendLinkUpdateForm from "@/components/forms/FriendLinkUpdateForm.vue";
import {useRoute, useRouter} from "vue-router";
import {ElMessage} from "element-plus";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const layoutStore = useLayoutStore()

/*
* 这里写成 const friendLinkCreateVisible = computed(() => layoutStore.state.friendLinkCreateVisible)
* 就会导致，上面不能直接赋值 friendLinkCreateVisible = false
* 写成 const friendLinkCreateVisible = ref(layoutStore.state.friendLinkCreateVisible)
* 就会导致，不能根据 layoutStore.state.friendLinkCreateVisible 自动变化
* const friendLinkCreateVisible = ref(computed(() => layoutStore.state.friendLinkCreateVisible)) 不被允许的
* */
const friendLinkCreateVisible = ref(layoutStore.state.friendLinkCreateVisible);
watch(
    () => layoutStore.state.friendLinkCreateVisible,
    (newValue) => {
      friendLinkCreateVisible.value = newValue;
    }
);

const friendLinkBulkDeleteVisible = ref(false)
let idsToDelete: number[]

const friendLinkCreateVisibleSynchronization = () => {
  layoutStore.state.friendLinkCreateVisible = false
}

const page = ref(1)
const page_size = ref(10)
const total = ref(0)
const friendLinkTableData = ref<FriendLink[]>()


const friendLinkListRequest = reactive<FriendLinkListRequest>({
  name: null,
  description: null,
  page: 1,
  page_size: 10,
})

const route = useRoute();
const router = useRouter();

onMounted(() => {
  friendLinkListRequest.name = route.query.name as string || null;
  friendLinkListRequest.description = route.query.description as string || null;
  page.value = Number(route.query.page) || 1;
  page_size.value = Number(route.query.page_size) || 10;
});

const getFriendLinkTableData = async () => {
  if (friendLinkListRequest.name === "") {
    friendLinkListRequest.name = null;
  }
  if (friendLinkListRequest.description === "") {
    friendLinkListRequest.description = null;
  }

  friendLinkListRequest.page = page.value;
  friendLinkListRequest.page_size = page_size.value;

  const table = await friendLinkList(friendLinkListRequest);

  if (table.code === 0) {
    friendLinkTableData.value = table.data.list;
    total.value = table.data.total;

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        name: friendLinkListRequest.name,
        description: friendLinkListRequest.description,
        page: friendLinkListRequest.page,
        page_size: friendLinkListRequest.page_size,
      },
    })
  }
}

nextTick(() => {
  getFriendLinkTableData();
});

const handleSizeChange = (val: number) => {
  page_size.value = val
  getFriendLinkTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getFriendLinkTableData()
}

const friendLinkUpdateVisible = ref(layoutStore.state.friendLinkUpdateVisible);
watch(
    () => layoutStore.state.friendLinkUpdateVisible,
    (newValue) => {
      friendLinkUpdateVisible.value = newValue;
    }
);

const friendLinkUpdateVisibleSynchronization = () => {
  layoutStore.state.friendLinkUpdateVisible = false
}

const friendLinkDeleteVisible = ref(false)

const handleDelete = async (id: number) => {
  let ids: number[] = []
  ids.push(id)

  const requestData: FriendLinkDeleteRequest = {
    ids: ids
  };

  const res = await friendLinkDelete(requestData);
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  friendLinkDeleteVisible.value = false
  layoutStore.state.shouldRefreshFriendLinkTable = true
}

const multipleFriendLinkTableRef = ref()

const handleIdsToDelete = () => {
  idsToDelete = []

  const rows: FriendLink[] = multipleFriendLinkTableRef.value.getSelectionRows()
  rows.forEach(row => {
    idsToDelete.push(row.id)
  })
}

const handleBulkDelete = async (ids: number[]) => {
  const requestData: FriendLinkDeleteRequest = {
    ids: ids
  };


  const res = await friendLinkDelete(requestData);
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  friendLinkBulkDeleteVisible.value = false
  layoutStore.state.shouldRefreshFriendLinkTable = true
}

watch(() => layoutStore.state.shouldRefreshFriendLinkTable, (newVal) => {
  if (newVal) {
    getFriendLinkTableData();
    layoutStore.state.shouldRefreshFriendLinkTable = false;
  }
});

let friendLinkInfo: FriendLink


watch(() => route.query, (newQuery) => {
  friendLinkListRequest.name = newQuery.name as string || null;
  friendLinkListRequest.description = newQuery.description as string || null;
  friendLinkListRequest.page = Number(newQuery.page) || 1;
  friendLinkListRequest.page_size = Number(newQuery.page_size) || 10;
}, {immediate: true});
</script>

<style scoped lang="scss">
.friend-link-list {
  .page-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: var(--sp-4);
    margin-bottom: var(--sp-5);

    .page-title {
      font-size: var(--fs-24);
      font-weight: 600;
      color: var(--text-primary);
      line-height: var(--lh-title);
    }

    .page-desc {
      margin-top: var(--sp-1);
      font-size: var(--fs-14);
      color: var(--text-muted);
    }

    .page-actions {
      display: flex;
      gap: var(--sp-2);
      flex-shrink: 0;
    }
  }

  .friend-link-list-request {
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    padding: var(--sp-4) var(--sp-5);
    margin-bottom: var(--sp-5);

    .el-form {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: var(--sp-2);
    }
  }

  .el-table {
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);

    .el-image {
      height: 48px;
      border-radius: var(--radius-sm);
    }
  }

  .el-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--sp-5);
  }
}
</style>
