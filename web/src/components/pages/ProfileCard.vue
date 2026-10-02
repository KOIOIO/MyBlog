<template>
  <el-card class="profile-card">
    <el-row class="title">{{ t('components.profileCard.title') }}</el-row>
    <div class="content">
      <el-row>{{ t('components.profileCard.nicknameLabel') }}：{{ websiteStore.state.websiteInfo.name }}</el-row>
      <el-row>{{ t('components.profileCard.jobLabel') }}：{{ websiteStore.state.websiteInfo.job }}</el-row>
      <el-row>{{ t('components.profileCard.addressLabel') }}：{{ websiteStore.state.websiteInfo.address }}</el-row>
      <el-row>{{ t('components.profileCard.emailLabel') }}：{{ websiteStore.state.websiteInfo.email }}</el-row>
      <div class="contact">
        <div v-if="websiteStore.state.websiteInfo['qq_image']!==''" class="contact-item">
          <el-image :src="websiteStore.state.websiteInfo['qq_image']" alt=""/>
          <el-row>QQ</el-row>
        </div>
        <div v-if="websiteStore.state.websiteInfo['wechat_image']!==''" class="contact-item">
          <el-image :src="websiteStore.state.websiteInfo['wechat_image']" alt=""/>
          <el-row>{{ t('components.profileCard.wechat') }}</el-row>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import {useWebsiteStore} from "@/stores/website";
import {useI18n} from "vue-i18n";

const {t} = useI18n();

const websiteStore = useWebsiteStore()
</script>

<style scoped lang="scss">
.profile-card {
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

  .content {
    font-size: var(--fs-14);
    color: var(--text-body);
    line-height: 2;

    .contact {
      display: flex;
      justify-content: center;
      margin-top: var(--sp-3);

      .contact-item {
        text-align: center;
        margin: 0 var(--sp-3);
        cursor: pointer;
        transition: transform 150ms cubic-bezier(.16,1,.3,1);

        &:hover {
          transform: translateY(-2px);
        }

        &:active {
          transform: translateY(1px);
          transition-duration: 80ms;
        }

        .el-image {
          height: 100px;
        }

        .el-row {
          display: flex;
          justify-content: center;
          font-size: var(--fs-12);
          color: var(--text-muted);
        }
      }
    }
  }
}

</style>
