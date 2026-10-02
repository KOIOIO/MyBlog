<template>
  <div class="image-list">
    <div class="title">
      <el-row>{{ t('dashboard.images.title') }}</el-row>
      <el-button-group>
        <el-button type="danger" icon="Delete" @click="imageBulkDeleteVisible = true;handleIdsToDelete()">
          {{ t('common.batchDelete') }}
        </el-button>

        <el-dialog
            v-model="imageBulkDeleteVisible"
            width="500"
            align-center
            destroy-on-close
        >
          <template #header>
            {{ t('dashboard.images.deleteTitle') }}
          </template>
          {{ t('dashboard.common.deleteConfirm', { count: idsToDelete?.length ?? 0 }) }}
          <template #footer>
            <el-button type="primary" @click="handleBulkDelete(idsToDelete)">
              {{ t('common.confirm') }}
            </el-button>
            <el-button @click="imageBulkDeleteVisible = false">{{ t('common.cancel') }}</el-button>
          </template>
        </el-dialog>
      </el-button-group>
    </div>

    <div class="image-list-request">
      <el-form :inline="true" :model="imageListRequest">
        <el-form-item :label="t('dashboard.images.filter.name')">
          <el-input v-model="imageListRequest.name" :placeholder="t('dashboard.images.filter.namePh')" clearable/>
        </el-form-item>
        <el-form-item :label="t('dashboard.images.filter.category')">
          <el-select
              v-model="imageListRequest.category"
              :placeholder="t('dashboard.common.select')"
              style="width: 200px"
          >
            <el-option
                v-for="item in categoryOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('dashboard.images.filter.storage')">
          <el-select
              v-model="imageListRequest.storage"
              :placeholder="t('dashboard.common.select')"
              style="width: 200px"
          >
            <el-option
                v-for="item in storageOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="getImageTableData">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
        ref="multipleImageTableRef"
        :data="imageTableData"
    >
      <el-table-column type="selection" :selectable="selectable" width="60"/>
      <el-table-column :label="t('dashboard.images.column.image')" width="100">
        <template #default="scope:{ row: Image, column: any, $index: number }">
          <el-image :src="scope.row.url" alt=""/>
        </template>
      </el-table-column>
      <el-table-column prop="name" :label="t('dashboard.images.column.name')" width="320"/>
      <el-table-column prop="url" label="URL" width="340"/>
      <el-table-column prop="category" :label="t('dashboard.images.column.category')"/>
      <el-table-column prop="storage" :label="t('dashboard.images.column.storage')"/>
      <el-table-column :label="t('dashboard.images.column.actions')">
        <template #default="scope:{ row: Image, column: any, $index: number }">
          <el-button
              v-if="scope.row.category==='未使用'"
              type="danger"
              @click="imageDeleteVisible=true;imageInfo=scope.row"
          >
            {{ t('common.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
        v-model="imageDeleteVisible"
        width="500"
        align-center
        destroy-on-close
    >
      <template #header>
        {{ t('dashboard.images.deleteTitle') }}
      </template>
      {{ t('dashboard.common.deleteConfirm', { count: 1 }) }}
      <template #footer>
        <el-button type="primary" @click="handleDelete(imageInfo.id)">
          {{ t('common.confirm') }}
        </el-button>
        <el-button @click="imageDeleteVisible = false">{{ t('common.cancel') }}</el-button>
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
import {computed, nextTick, onMounted, reactive, ref, watch} from "vue";
import {useLayoutStore} from "@/stores/layout";
import {ElMessage} from "element-plus";
import {useRoute, useRouter} from "vue-router";
import {useI18n} from "vue-i18n";
import {type Image, imageDelete, type ImageDeleteRequest, imageList, type ImageListRequest} from "@/api/image";

const {t} = useI18n()

const multipleImageTableRef = ref()
const imageTableData = ref<Image[]>()
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const layoutStore = useLayoutStore()

const selectable = (row: Image) => row.category==='未使用'

const imageBulkDeleteVisible = ref(false)
let idsToDelete: number[]

const handleIdsToDelete = () => {
  idsToDelete = []

  const rows: Image[] = multipleImageTableRef.value.getSelectionRows()
  rows.forEach(row => {
    idsToDelete.push(row.id)
  })
}

const handleBulkDelete = async (ids: number[]) => {
  const requestData: ImageDeleteRequest = {
    ids: ids
  }
  const res = await imageDelete(requestData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  imageBulkDeleteVisible.value = false
  layoutStore.state.shouldRefreshImageTable = true
}

const categoryOptions = computed(() => [
  {
    value: '',
    label: t('dashboard.images.category.all'),
  },
  {
    value: '未使用',
    label: t('dashboard.images.category.unused'),
  },
  {
    value: '系统',
    label: t('dashboard.images.category.system'),
  },
  {
    value: '背景',
    label: t('dashboard.images.category.background'),
  },
  {
    value: '封面',
    label: t('dashboard.images.category.cover'),
  },
  {
    value: '插图',
    label: t('dashboard.images.category.illustration'),
  },
  {
    value: '广告',
    label: t('dashboard.images.category.ad'),
  },
  {
    value: '友链',
    label: t('dashboard.images.category.friendLink'),
  },
])

const storageOptions = computed(() => [
  {
    value: '',
    label: t('dashboard.images.storage.all'),
  },
  {
    value: '本地',
    label: t('dashboard.images.storage.local'),
  },
  {
    value: '七牛云',
    label: t('dashboard.images.storage.qiniu'),
  },
])

const imageListRequest = reactive<ImageListRequest>({
  name: null,
  category: null,
  storage: null,
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

onMounted(() => {
  imageListRequest.name = route.query.name as string || null
  imageListRequest.category = route.query.category as string || null
  imageListRequest.storage = route.query.storage as string || null
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

const getImageTableData = async () => {
  if (imageListRequest.name === "") {
    imageListRequest.name = null
  }
  if (imageListRequest.category === "") {
    imageListRequest.category = null
  }
  if (imageListRequest.storage === "") {
    imageListRequest.storage = null
  }

  imageListRequest.page = page.value
  imageListRequest.page_size = page_size.value

  const table = await imageList(imageListRequest)

  if (table.code === 0) {
    imageTableData.value = table.data.list
    total.value = table.data.total

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        title: imageListRequest.name,
        content: imageListRequest.category,
        storage: imageListRequest.storage,
        page: imageListRequest.page,
        page_size: imageListRequest.page_size,
      },
    })
  }
}

watch(() => route.query, (newQuery) => {
  imageListRequest.name = newQuery.name as string || null
  imageListRequest.category = newQuery.category as string || null
  imageListRequest.storage = newQuery.storage as string || null
  imageListRequest.page = Number(newQuery.page) || 1
  imageListRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getImageTableData()
})

let imageInfo: Image
const imageDeleteVisible = ref(false)

const handleDelete = async (id: number) => {
  let ids: number[] = []
  ids.push(id)

  const requestData: ImageDeleteRequest = {
    ids: ids
  }

  const res = await imageDelete(requestData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
  imageDeleteVisible.value = false
  layoutStore.state.shouldRefreshImageTable = true
}

watch(() => layoutStore.state.shouldRefreshImageTable, (newVal) => {
  if (newVal) {
    getImageTableData()
    layoutStore.state.shouldRefreshImageTable = false
  }
})

const handleSizeChange = (val: number) => {
  page_size.value = val
  getImageTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getImageTableData()
}
</script>

<style scoped lang="scss">
.image-list {
  .title {
    display: flex;
    align-items: center;
    margin-bottom: var(--sp-5);

    .el-row {
      font-size: var(--fs-24);
      font-weight: 600;
      color: var(--text-primary);
    }

    .el-button-group {
      margin-left: auto;

      .el-button {
        margin-left: var(--sp-3);
      }
    }
  }

  .image-list-request {
    background-color: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    padding: var(--sp-4) var(--sp-5);
    margin-bottom: var(--sp-5);

    .el-form {
      display: flex;
      flex-wrap: wrap;
      align-items: center;

      .el-form-item {
        margin-right: var(--sp-4);
        margin-bottom: 0;
      }

      .el-input__wrapper {
        height: 40px;
        border-radius: var(--radius-sm);
      }
    }
  }

  .el-table {
    --el-table-border-color: var(--border);

    .el-image {
      height: 48px;
      width: 48px;
      border-radius: var(--radius-sm);
    }
  }

  /* 操作列：文字删除按钮 */
  .el-table .cell .el-button--danger {
    --el-button-bg-color: transparent;
    --el-button-border-color: transparent;
    --el-button-text-color: var(--el-color-danger);
    --el-button-hover-bg-color: var(--el-color-danger-light-9);
    --el-button-hover-text-color: var(--el-color-danger);
    --el-button-hover-border-color: transparent;
    --el-button-active-bg-color: var(--el-color-danger-light-9);
    --el-button-active-text-color: var(--el-color-danger);
    padding: 4px 0;
  }

  .el-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--sp-5);
  }
}
</style>
