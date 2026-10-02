<template>
<div class="chart-card">
  <div id="chart" :style="{height: '300px'}">
  </div>
</div>
</template>

<script setup lang="ts">
import * as echarts from 'echarts/core';
import {
  TooltipComponent,
  GridComponent,
  LegendComponent
} from 'echarts/components';
import { LineChart } from 'echarts/charts';
import { UniversalTransition } from 'echarts/features';
import { CanvasRenderer } from 'echarts/renderers';
import {defineProps, onMounted} from "vue";
import type {UserChartResponse} from "@/api/user";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const props = defineProps<{
  chart: UserChartResponse;
}>();

echarts.use([
  TooltipComponent,
  GridComponent,
  LegendComponent,
  LineChart,
  CanvasRenderer,
  UniversalTransition
]);

onMounted(() => {
  const chartDom = document.getElementById('chart');
  const myChart = echarts.init(chartDom);
  let option;

  const cssVar = (name: string) =>
    getComputedStyle(document.documentElement).getPropertyValue(name).trim();
  const accentColor = cssVar('--accent');
  const mutedColor = cssVar('--text-muted');
  const borderColor = cssVar('--border');
  const elevatedBg = cssVar('--bg-elevated');
  const bodyColor = cssVar('--text-body');

  option = {
    backgroundColor: 'transparent',
    tooltip: {
      trigger: 'axis',
      backgroundColor: elevatedBg,
      borderColor: borderColor,
      textStyle: {
        color: bodyColor,
      },
    },
    legend: {
      data: [t('components.userActivityChart.login'), t('components.userActivityChart.register')],
      textStyle: { color: mutedColor },
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: props.chart.date_list,
      axisLine: { lineStyle: { color: borderColor } },
      axisLabel: { color: mutedColor },
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: mutedColor },
      splitLine: { lineStyle: { color: borderColor } },
    },
    series: [
      {
        name: t('components.userActivityChart.login'),
        type: 'line',
        data: props.chart.login_data,
        itemStyle: { color: accentColor },
        lineStyle: { color: accentColor },
      },
      {
        name: t('components.userActivityChart.register'),
        type: 'line',
        data: props.chart.register_data,
        itemStyle: { color: mutedColor },
        lineStyle: { color: mutedColor },
      }
    ]
  };

  option && myChart.setOption(option);
})

</script>

<style scoped lang="scss">
.chart-card {
  background-color: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  padding: var(--sp-4);
}
</style>
