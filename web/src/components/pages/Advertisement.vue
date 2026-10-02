<template>
  <el-card class="advertisement">
    <el-row class="title">{{ t('components.advertisement.title') }}</el-row>
    <el-carousel :interval="5000" type="card" height="320px">
      <el-carousel-item v-for="advertisement in advertisementList" :key="advertisement.ad_image">
        <el-image :src="advertisement.ad_image" alt=""></el-image>
        <el-row class="ad-title">{{ advertisement.title }}</el-row>
        <el-text class="ad-content">{{ advertisement.content }}</el-text>
      </el-carousel-item>
    </el-carousel>
  </el-card>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {type Advertisement, advertisementInfo} from "@/api/advertisement";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const advertisementList = ref<Advertisement[]>()
const getAdvertisementList = async () => {
  const res = await advertisementInfo()
  if (res.code == 0) {
    advertisementList.value = res.data.list
  }
}
getAdvertisementList()
</script>

<style scoped lang="scss">
.advertisement {
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

  .el-carousel__item {
    background-color: var(--bg-elevated);
    border-radius: var(--radius-md);
    overflow: hidden;

    .el-image {
      height: 240px;
      width: 100%;
    }
  }

  .ad-title {
    font-size: var(--fs-14);
    font-weight: 600;
    color: var(--text-primary);
    padding: var(--sp-2) var(--sp-3) 0;
  }

  .ad-content {
    display: block;
    font-size: var(--fs-12);
    color: var(--text-muted);
    padding: 0 var(--sp-3) var(--sp-2);
  }
}

</style>
