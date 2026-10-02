<template>
  <el-card class="calendar">
    <el-row class="title">{{ t('components.calendar.title') }}</el-row>
    <div class="calendar-rows">
      <el-row>{{ t('components.calendar.timeLabel') }}：{{ calendarInfo.date }} {{ currentTime }}</el-row>
      <el-row>{{ t('components.calendar.lunarLabel') }}：{{ calendarInfo.lunar_date }}</el-row>
      <el-row>{{ t('components.calendar.ganzhiLabel') }}：{{ calendarInfo.ganzhi }}</el-row>
      <el-row>{{ t('components.calendar.zodiacLabel') }}：{{ calendarInfo.zodiac }}</el-row>
      <el-row>{{ t('components.calendar.dayOfYearLabel') }}：{{ calendarInfo.day_of_year }}</el-row>
      <el-row>{{ t('components.calendar.solarTermLabel') }}：{{ calendarInfo.solar_term }}</el-row>
      <el-row>{{ t('components.calendar.auspiciousLabel') }}：{{ calendarInfo.auspicious }}</el-row>
      <el-row>{{ t('components.calendar.inauspiciousLabel') }}：{{ calendarInfo.inauspicious }}</el-row>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import {onUnmounted, ref} from "vue";
import {websiteCalendar, type WebsiteCalendarResponse} from "@/api/website";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const calendarInfo = ref<WebsiteCalendarResponse>({
  date: '',
  lunar_date: '',
  ganzhi: '',
  zodiac: '',
  day_of_year: '',
  solar_term: '',
  auspicious: '',
  inauspicious: '',
})

let timerId: number | null = null
const currentTime = ref('')

function updateCurrentTime() {
  currentTime.value = new Date().toLocaleTimeString()
}

function initializeTimer() {
  updateCurrentTime()
  timerId = setInterval(updateCurrentTime, 1000)
}

onUnmounted(() => {
  clearInterval(timerId as number)
})

initializeTimer()

const getCalendarInfo = async () => {
  const res = await websiteCalendar()
  if (res.code == 0) {
    calendarInfo.value = res.data
  }
}
getCalendarInfo()
</script>

<style scoped lang="scss">
.calendar {
  margin-bottom: var(--sp-4);
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  animation: kf-fade-up var(--dur-mid) var(--ease-out) backwards;

  .title {
    font-size: var(--fs-16);
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: var(--sp-3);
  }

  .calendar-rows {
    font-size: var(--fs-14);
    color: var(--text-body);
    line-height: 2;
  }
}

</style>
