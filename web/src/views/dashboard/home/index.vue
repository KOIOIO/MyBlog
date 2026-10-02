<template xmlns="http://www.w3.org/1999/html">
  <div class="home">
    <div class="header">
      <el-card class="user-card">
        <el-row>
          {{ t('dashboard.home.greeting', { username: userStore.state.userInfo.username }) }}
        </el-row>
        <div class="weather">
          {{ weatherInfo }}
        </div>
      </el-card>
    </div>

    <div class="content">
      <el-col :span="14">
        <el-card class="entrance-card">
          <el-row class="title">
            {{ t('dashboard.home.quickEntry') }}
          </el-row>
          <div class="button-group">
            <div class="button-item" v-for="item in entranceList">
              <el-button :icon="item.icon" :type="item.type" plain @click="handleClick(item)"/>
              {{ item.title }}
            </div>
          </div>
        </el-card>
        <el-card class="chart-card">
          <el-row class="title">
            {{ t('dashboard.home.userData') }}
          </el-row>
          <div class="time-select">
            <el-select
                @change="getChartInfo"
                v-model="userChartRequest.date"
                :placeholder="t('dashboard.home.chartPlaceholder')"
                style="width: 200px"
            >
              <el-option
                  v-for="item in userChartOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
              />
            </el-select>
          </div>
          <user-activity-chart v-if="isShow" :chart="chart"/>
        </el-card>
      </el-col>
      <el-col :span="10">
        <el-card class="aside">
          <el-row class="title">
            {{ t('dashboard.home.declaration.title') }}
          </el-row>
          <div class="text">
            <el-text>
              {{ t('dashboard.home.declaration.intro') }}<br>
              <h3>{{ t('dashboard.home.declaration.copyright.title') }}</h3>
              {{ t('dashboard.home.declaration.copyright.c1') }}<br>
              {{ t('dashboard.home.declaration.copyright.c2') }}<br>
              {{ t('dashboard.home.declaration.copyright.c3') }}<br>
              <h3>{{ t('dashboard.home.declaration.terms.title') }}</h3>
              {{ t('dashboard.home.declaration.terms.t1') }}<br>
              {{ t('dashboard.home.declaration.terms.t2') }}<br>
              <h3>{{ t('dashboard.home.declaration.privacy.title') }}</h3>
              {{ t('dashboard.home.declaration.privacy.p1') }}<br>
              {{ t('dashboard.home.declaration.privacy.p2') }}<br>
              <h3>{{ t('dashboard.home.declaration.contact.title') }}</h3>
              {{ t('dashboard.home.declaration.contact.c1') }}<br>
              {{ t('dashboard.home.declaration.contact.emailLabel') }}：[{{ useWebsiteStore().state.websiteInfo.email }}]<br>
              {{ t('dashboard.home.declaration.contact.c2') }}<br>
            </el-text>
          </div>
        </el-card>
      </el-col>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useUserStore} from "@/stores/user";
import {computed, ref} from "vue";
import {useI18n} from "vue-i18n";
import {userChart, type UserChartRequest, type UserChartResponse, userWeather} from "@/api/user";
import UserActivityChart from "@/components/widgets/UserActivityChart.vue";
import {useWebsiteStore} from "@/stores/website";
import router from "@/router";
import {type Tag, useTagStore} from "@/stores/tag";

const {t} = useI18n()
const userStore = useUserStore()

const weatherInfo = ref('')

const getWeatherInfo = async () => {
  const res = await userWeather()
  if (res.code === 0) {
    weatherInfo.value = res.data
  }
}

getWeatherInfo()

interface Entrance {
  title: string
  name: string;
  icon: string;
  type: string;
}

const entranceList = computed<Entrance[]>(() => [
  {
    title: t('menu.userCenter.info'),
    name: 'user-info',
    icon: 'Postcard',
    type: 'primary',
  },
  {
    title: t('menu.userCenter.star'),
    name: 'user-star',
    icon: 'Star',
    type: 'warning',
  },
  {
    title: t('menu.userCenter.comment'),
    name: 'user-comment',
    icon: 'ChatDotRound',
    type: 'info',
  },
  {
    title: t('menu.userCenter.feedback'),
    name: 'user-feedback',
    icon: 'Message',
    type: 'success',
  }
])

const tagStore = useTagStore()

function handleClick(item: Entrance) {
  const newTag: Tag = {
    title: item.title,
    name: item.name
  }
  const exists = tagStore.state.tags.some(tag => tag.name === newTag.name);
  if (exists) {
    return;
  }
  tagStore.state.tags.push(newTag);
  router.push({name:item.name})
}

const userChartRequest = ref<UserChartRequest>({
  date: 7,
})

const userChartOptions = computed(() => [
  {
    value: 7,
    label: t('dashboard.home.range.7'),
  },
  {
    value: 30,
    label: t('dashboard.home.range.30'),
  },
  {
    value: 90,
    label: t('dashboard.home.range.90'),
  },
  {
    value: 180,
    label: t('dashboard.home.range.180'),
  },
  {
    value: 365,
    label: t('dashboard.home.range.365'),
  },
])

const chart = ref<UserChartResponse>(
    {
      date_list: [],
      login_data: [],
      register_data: [],
    }
)

const isShow = ref(false)

const getChartInfo = async () => {
  isShow.value = false
  const res = await userChart(userChartRequest.value)
  if (res.code === 0) {
    chart.value = res.data
    isShow.value = true
  }
}

getChartInfo()

</script>

<style scoped lang="scss">
.home {
  .header {
    margin-bottom: var(--sp-5);

    .user-card {
      background-color: var(--bg-elevated);
      border: 1px solid var(--border);
      border-radius: var(--radius-md);
      color: var(--text-body);

      .el-row {
        font-size: var(--fs-24);
        font-weight: 600;
        color: var(--text-primary);
      }

      .weather {
        margin-top: var(--sp-3);
        margin-bottom: var(--sp-2);
        font-size: var(--fs-14);
        color: var(--text-muted);
      }
    }
  }

  .content {
    display: flex;
    gap: var(--sp-5);

    .entrance-card {
      background-color: var(--bg-elevated);
      border: 1px solid var(--border);
      border-radius: var(--radius-md);
      margin-bottom: var(--sp-5);

      .title {
        font-size: var(--fs-18);
        font-weight: 600;
        color: var(--text-primary);
      }

      .button-group {
        display: flex;
        gap: var(--sp-4);
        margin-top: var(--sp-4);
        margin-bottom: var(--sp-4);

        .button-item {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: var(--sp-2);
          width: 80px;
          font-size: var(--fs-12);
          color: var(--text-muted);

          .el-button {
            border: none;
            --el-font-size-base: 22px;
            height: 48px;
            width: 48px;
            padding: 0;
            background-color: var(--accent-weak);
            color: var(--accent);
            border-radius: var(--radius-md);
            transition: background-color 180ms ease-out, color 180ms ease-out;
          }

          .el-button:hover {
            background-color: var(--accent);
            color: var(--bg);
          }
        }
      }
    }

    .chart-card {
      background-color: var(--bg-elevated);
      border: 1px solid var(--border);
      border-radius: var(--radius-md);

      .title {
        font-size: var(--fs-18);
        font-weight: 600;
        color: var(--text-primary);
      }

      .time-select {
        display: flex;
        margin-top: var(--sp-3);
        margin-bottom: var(--sp-4);

        .el-select {
          margin-left: auto;
          width: 200px;
        }
      }

      .user-activity-chart {
        position: relative;
        width: 100%;
      }
    }

    .aside {
      background-color: var(--bg-elevated);
      border: 1px solid var(--border);
      border-radius: var(--radius-md);
      margin-left: var(--sp-5);
      height: 100%;

      .title {
        font-size: var(--fs-18);
        font-weight: 600;
        color: var(--text-primary);
      }

      .text {
        margin-top: var(--sp-4);
        font-size: var(--fs-14);
        line-height: var(--lh-body);
        color: var(--text-body);

        h3 {
          font-size: var(--fs-16);
          color: var(--text-primary);
        }
      }
    }
  }
}
</style>
