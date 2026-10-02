<template>
  <div class="user-list">
    <div class="page-header">
      <div>
        <div class="page-title">{{ t('system.users.title') }}</div>
        <div class="page-desc">{{ t('system.users.desc') }}</div>
      </div>
    </div>

    <div class="user-list-request">
      <el-form :inline="true" :model="userListRequest">
        <el-form-item label="uuid">
          <el-input v-model="userListRequest.uuid" :placeholder="t('system.users.uuidPlaceholder')" clearable/>
        </el-form-item>
        <el-form-item :label="t('system.users.registerSource')">
          <el-select
              v-model="userListRequest.register"
              :placeholder="t('system.users.selectPlaceholder')"
              style="width: 200px"
          >
            <el-option
                v-for="item in registerOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="getUserTableData">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
        :data="userTableData"
    >
      <el-table-column :label="t('system.users.avatar')">
        <template #default="scope:{ row: User, column: any, $index: number }">
          <el-avatar :src="scope.row.avatar" alt=""/>
        </template>
      </el-table-column>
      <el-table-column prop="username" :label="t('system.users.username')"/>
      <el-table-column prop="uuid" label="UUID" width="320px"/>
      <el-table-column prop="address" :label="t('system.users.address')"/>
      <el-table-column :label="t('system.users.registerTime')" width="250px">
        <template #default="scope:{ row: User, column: any, $index: number }">
          {{ getTime(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column :label="t('system.users.role')">
        <template #default="scope:{ row: User, column: any, $index: number }">
          {{ scope.row.role_id === 2 ? t('system.users.admin') : t('system.users.normalUser') }}
        </template>
      </el-table-column>
      <el-table-column prop="register" :label="t('system.users.registerSource')"/>
      <el-table-column :label="t('common.actions')">
        <template #default="scope:{ row: User, column: any, $index: number }">
          <el-button
              v-if="scope.row.role_id===1"
              link
              type="warning"
              @click="userFreezeVisible=true;userInfo=scope.row"
          >
            {{ scope.row.freeze === false ? t('system.users.freeze') : t('system.users.unfreeze') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
        v-model="userFreezeVisible"
        width="500"
        align-center
        destroy-on-close
    >
      <template #header>
        {{ t('system.users.freezeDialogTitle') }}
      </template>
      {{ t('system.users.freezeConfirm', {action: userInfo.freeze === false ? t('system.users.freeze') : t('system.users.unfreeze'), username: userInfo.username}) }}
      <template #footer>
        <el-button type="primary" @click="handleFreeze(userInfo.id,userInfo.freeze)">
          {{ t('common.confirm') }}
        </el-button>
        <el-button @click="userFreezeVisible = false">{{ t('common.cancel') }}</el-button>
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
import {type User, userFreeze, userList, type UserListRequest, type UserOperation, userUnfreeze} from "@/api/user";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const userTableData = ref<User[]>()
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const layoutStore = useLayoutStore()

const registerOptions = computed(() => [
  {
    value: '',
    label: t('common.all'),
  },
  {
    value: '邮箱',
    label: t('system.users.registerEmail'),
  },
])

const userListRequest = reactive<UserListRequest>({
  uuid: null,
  register:null,
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

onMounted(() => {
  userListRequest.uuid = route.query.uuid as string || null
  userListRequest.register = route.query.register as string || null
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

const getTime = (date: Date): string => {
  const time = new Date(date)
  return time.toLocaleString()
}

const getUserTableData = async () => {
  if (userListRequest.uuid === "") {
    userListRequest.uuid = null
  }
  if (userListRequest.register === "") {
    userListRequest.register = null
  }

  userListRequest.page = page.value
  userListRequest.page_size = page_size.value

  const table = await userList(userListRequest)

  if (table.code === 0) {
    userTableData.value = table.data.list
    total.value = table.data.total

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        uuid: userListRequest.uuid,
        register: userListRequest.register,
        page: userListRequest.page,
        page_size: userListRequest.page_size,
      },
    })
  }
}

watch(() => route.query, (newQuery) => {
  userListRequest.uuid = newQuery.uuid as string || null
  userListRequest.register = newQuery.register as string || null
  userListRequest.page = Number(newQuery.page) || 1
  userListRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getUserTableData()
})

let userInfo: User

const userFreezeVisible = ref(false)

const handleFreeze = async (id: number, freeze: boolean) => {
  const userOperationRequest: UserOperation = {
    id: id,
  }
  if (freeze) {
    const res = await userUnfreeze(userOperationRequest)
    if (res.code === 0) {
      ElMessage.success(res.msg)
    }
    userFreezeVisible.value = false
    layoutStore.state.shouldRefreshUserTable = true
  } else {
    const res = await userFreeze(userOperationRequest)
    if (res.code === 0) {
      ElMessage.success(res.msg)
    }
    userFreezeVisible.value = false
    layoutStore.state.shouldRefreshUserTable = true
  }
}

watch(() => layoutStore.state.shouldRefreshUserTable, (newVal) => {
  if (newVal) {
    getUserTableData()
    layoutStore.state.shouldRefreshUserTable = false
  }
})

const handleSizeChange = (val: number) => {
  page_size.value = val
  getUserTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getUserTableData()
}
</script>

<style scoped lang="scss">
.user-list {
  .page-header {
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
  }

  .user-list-request {
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
  }

  .el-pagination {
    display: flex;
    justify-content: flex-end;
    margin-top: var(--sp-5);
  }
}
</style>
