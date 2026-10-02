<template>
  <div class="login-logs">
    <div class="page-header">
      <div>
        <div class="page-title">{{ t('system.loginLogs.title') }}</div>
        <div class="page-desc">{{ t('system.loginLogs.desc') }}</div>
      </div>
    </div>

    <div class="user-login-list-request">
      <el-form :inline="true" :model="userLoginListRequest">
        <el-form-item label="uuid">
          <el-input v-model="userLoginListRequest.uuid" :placeholder="t('system.loginLogs.uuidPlaceholder')" clearable/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="getUserLoginTableData">{{ t('common.query') }}</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
        :data="userLoginTableData"
    >
      <el-table-column :label="t('system.loginLogs.user')" width="80">
        <template #default="scope:{ row: Login, column: any, $index: number }">
          <el-popover width="280">
            <template #reference>
              <el-avatar :src="scope.row.user.avatar"/>
            </template>
            <template #default>
              <user-card :uuid="''"
                         :user-card-info="{uuid:scope.row.user.uuid,username:scope.row.user.username,avatar:scope.row.user.avatar,address:scope.row.user.address,signature:scope.row.user.signature}"/>
            </template>
          </el-popover>
        </template>
      </el-table-column>
      <el-table-column :label="t('system.loginLogs.username')" width="80">
        <template #default="scope:{ row: Login, column: any, $index: number }">
          {{ scope.row.user.username }}
        </template>
      </el-table-column>
      <el-table-column :label="t('system.loginLogs.loginTime')">
        <template #default="scope:{ row: Login, column: any, $index: number }">
          {{ getTime(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="login_method" :label="t('system.loginLogs.loginMethod')"/>
      <el-table-column prop="ip" label="IP"/>
      <el-table-column prop="address" :label="t('system.loginLogs.loginAddress')"/>
      <el-table-column prop="os" :label="t('system.loginLogs.os')"/>
      <el-table-column prop="device_info" :label="t('system.loginLogs.deviceInfo')"/>
      <el-table-column prop="browser_info" :label="t('system.loginLogs.browserInfo')"/>
      <el-table-column prop="status" :label="t('system.loginLogs.loginStatus')"/>
    </el-table>

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
import {useRoute, useRouter} from "vue-router";
import {
  type Login,
  userLoginList,
  type UserLoginListRequest,
} from "@/api/user";
import UserCard from "@/components/widgets/UserCard.vue";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const userLoginTableData = ref<Login[]>()
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const userLoginListRequest = reactive<UserLoginListRequest>({
  uuid: null,
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

onMounted(() => {
  userLoginListRequest.uuid = route.query.uuid as string || null
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

const getUserLoginTableData = async () => {
  if (userLoginListRequest.uuid === "") {
    userLoginListRequest.uuid = null
  }

  userLoginListRequest.page = page.value
  userLoginListRequest.page_size = page_size.value

  const table = await userLoginList(userLoginListRequest)

  if (table.code === 0) {
    userLoginTableData.value = table.data.list
    total.value = table.data.total

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        uuid: userLoginListRequest.uuid,
        page: userLoginListRequest.page,
        page_size: userLoginListRequest.page_size,
      },
    })
  }
}

watch(() => route.query, (newQuery) => {
  userLoginListRequest.uuid = newQuery.uuid as string || null
  userLoginListRequest.page = Number(newQuery.page) || 1
  userLoginListRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getUserLoginTableData()
})

const getTime = (date: Date): string => {
  const time = new Date(date)
  return time.toLocaleString()
}


const handleSizeChange = (val: number) => {
  page_size.value = val
  getUserLoginTableData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getUserLoginTableData()
}
</script>

<style scoped lang="scss">
.login-logs {
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

  .user-login-list-request {
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
