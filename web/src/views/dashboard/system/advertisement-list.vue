<template>
  <div class="advertisement-list">
    <div class="page-header">
      <div>
        <div class="page-title">{{ t('system.advertisement.title') }}</div>
        <div class="page-desc">{{ t('system.advertisement.desc') }}</div>
      </div>
      <div class="page-actions">
        <el-button type="primary" icon="Plus" @click="layoutStore.state.advertisementCreateVisible = true">
          {{ t('system.advertisement.create') }}
        </el-button>

        <el-dialog
            v-model="advertisementCreateVisible"
            width="500"
            align-center
            destroy-on-close
            :before-close="advertisementCreateVisibleSynchronization"
        >
          <template #header>
            {{ t('system.advertisement.create') }}
          </template>
          <advertisement-create-form/>
          <template #footer>
          </template>
        </el-dialog>

        <el-button type="danger" plain icon="Delete" @click="advertisementBulkDeleteVisible = true;handleIdsToDelete()">
          {{ t('common.batchDelete') }}
        </el-button>

        <el-dialog
            v-model="advertisementBulkDeleteVisible"
            width="500"
            align-center
            destroy-on-close
        >
          <template #header>
            {{ t('system.advertisement.delete') }}
          </template>
          {{ t('system.advertisement.confirmDeleteWithCount', {count: idsToDelete.length}) }}
          <template #footer>
            <el-button type="primary" @click="handleBulkDelete(idsToDelete)">
              {{ t('common.confirm') }}
            </el-button>
            <el-button @click="advertisementBulkDeleteVisible = false">{{ t('common.cancel') }}</el-button>
          </template>
        </el-dialog>
      </div>
    </div>

    <div class="advertisement-list-request">
      <el-form :inline="true" :model="advertisementListRequest">
        <el-form-item :label="t('system.advertisement.adTitle')">
          <el-input v-model="advertisementListRequest.title" :placeholder="t('system.advertisement.adTitlePlaceholder')" clearable/>
        </el-form-item>
        <el-form-item :label="t('system.advertisement.adContent')">
          <el-input v-model="advertisementListRequest.content" :placeholder="t('system.advertisement.adContentPlaceholder')" clearable/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="getAdvertisementTableData">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
        ref="multipleAdvertisementTableRef"
        :data="advertisementTableData"
    >
      <el-table-column type="selection" width="60"/>
      <el-table-column :label="t('system.advertisement.image')">
        <template #default="scope:{ row: any, column: any, $index: number }">
          <el-image :src="scope.row.ad_image" alt=""/>
        </template>
      </el-table-column>
      <el-table-column prop="link" :label="t('system.advertisement.link')"/>
      <el-table-column prop="title" :label="t('common.title')"/>
      <el-table-column prop="content" :label="t('common.content')"/>
      <el-table-column :label="t('common.actions')">
        <template #default="scope:{ row: any, column: any, $index: number }">
          <el-button
              link
              type="primary"
              @click="layoutStore.state.advertisementUpdateVisible=true;advertisementInfo=scope.row"
          >
            {{ t('system.advertisement.updateAction') }}
          </el-button>
          <el-button
              link
              type="danger"
              @click="advertisementDeleteVisible=true;advertisementInfo=scope.row"
          >
            {{ t('common.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
        v-model="advertisementUpdateVisible"
        width="500"
        align-center
        destroy-on-close
        :before-close="advertisementUpdateVisibleSynchronization"
    >
      <template #header>
        {{ t('system.advertisement.update') }}
      </template>
      <advertisement-update-form :advertisement=advertisementInfo />
      <template #footer>
      </template>
    </el-dialog>

    <el-dialog
        v-model="advertisementDeleteVisible"
        width="500"
        align-center
        destroy-on-close
    >
      <template #header>
        {{ t('system.advertisement.delete') }}
      </template>
      {{ t('system.advertisement.confirmDeleteWithCount', {count: 1}) }}
      <template #footer>
        <el-button type="primary" @click="handleDelete(advertisementInfo.id)">
          {{ t('common.confirm') }}
        </el-button>
        <el-button @click="advertisementDeleteVisible = false">{{ t('common.cancel') }}</el-button>
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
import {nextTick, onMounted, reactive, ref, watch} from "vue";
import {
  type Advertisement,
  advertisementDelete, type AdvertisementDeleteRequest,
  advertisementList,
  type AdvertisementListRequest
} from "@/api/advertisement";
import {useLayoutStore} from "@/stores/layout";
import AdvertisementCreateForm from "@/components/forms/AdvertisementCreateForm.vue";
import {ElMessage} from "element-plus";
import {useRoute, useRouter} from "vue-router";
import AdvertisementUpdateForm from "@/components/forms/AdvertisementUpdateForm.vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const multipleAdvertisementTableRef = ref()
const advertisementTableData = ref<Advertisement[]>()
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const layoutStore = useLayoutStore()

const advertisementCreateVisible = ref(layoutStore.state.advertisementCreateVisible)
watch(
    () => layoutStore.state.advertisementCreateVisible,
    (newValue) => {
      advertisementCreateVisible.value = newValue
    }
)

const advertisementCreateVisibleSynchronization = () => {
  layoutStore.state.advertisementCreateVisible = false
}

const advertisementBulkDeleteVisible = ref(false)
let idsToDelete: number[]

const handleIdsToDelete = () => {
  idsToDelete = []

  const rows: Advertisement[] = multipleAdvertisementTableRef.value.getSelectionRows()
  rows.forEach(row => {
    idsToDelete.push(row.id)
  })
}

const handleBulkDelete = async (ids: number[]) => {
  const requestData: AdvertisementDeleteRequest = {
    ids: ids
  }
  const res = await advertisementDelete(requestData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  advertisementBulkDeleteVisible.value = false
  layoutStore.state.shouldRefreshAdvertisementTable = true
}

const advertisementListRequest = reactive<AdvertisementListRequest>({
  title: null,
  content: null,
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

onMounted(() => {
  advertisementListRequest.title = route.query.title as string || null
  advertisementListRequest.content = route.query.content as string || null
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

const getAdvertisementTableData = async () => {
  if (advertisementListRequest.title === "") {
    advertisementListRequest.title = null
  }
  if (advertisementListRequest.content === "") {
    advertisementListRequest.content = null
  }

  advertisementListRequest.page = page.value
  advertisementListRequest.page_size = page_size.value

  const table = await advertisementList(advertisementListRequest)

  if (table.code === 0) {
    advertisementTableData.value = table.data.list
    total.value = table.data.total

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        title: advertisementListRequest.title,
        content: advertisementListRequest.content,
        page: advertisementListRequest.page,
        page_size: advertisementListRequest.page_size,
      },
    })
  }
}

watch(() => route.query, (newQuery) => {
  advertisementListRequest.title = newQuery.title as string || null
  advertisementListRequest.content = newQuery.content as string || null
  advertisementListRequest.page = Number(newQuery.page) || 1
  advertisementListRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getAdvertisementTableData()
})

let advertisementInfo: Advertisement
const advertisementDeleteVisible = ref(false)

const advertisementUpdateVisible = ref(layoutStore.state.advertisementUpdateVisible)
watch(
    () => layoutStore.state.advertisementUpdateVisible,
    (newValue) => {
      advertisementUpdateVisible.value = newValue
    }
)

const advertisementUpdateVisibleSynchronization = () => {
  layoutStore.state.advertisementUpdateVisible = false
}

const handleDelete = async (id: number) => {
  let ids: number[] = []
  ids.push(id)

  const requestData: AdvertisementDeleteRequest = {
    ids: ids
  }

  const res = await advertisementDelete(requestData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  advertisementDeleteVisible.value = false
  layoutStore.state.shouldRefreshAdvertisementTable = true
}

watch(() => layoutStore.state.shouldRefreshAdvertisementTable, (newVal) => {
  if (newVal) {
    getAdvertisementTableData()
    layoutStore.state.shouldRefreshAdvertisementTable = false
  }
})

const handleSizeChange = (val: number) => {
  page_size.value = val
  getAdvertisementTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getAdvertisementTableData()
}
</script>

<style scoped lang="scss">
.advertisement-list {
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

  .advertisement-list-request {
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
