<template>
  <div class="web-footer">
    <div class="container">
      <div class="footer-col">
        <div class="col-title">{{ t('footer.aboutSite') }}</div>
        <div class="footer-logo">
          <el-image
              :src="websiteStore.state.websiteInfo.full_logo===''?'/image/full_logo.png':websiteStore.state.websiteInfo.full_logo"
              alt=""/>
        </div>
        <p class="desc">{{ websiteStore.state.websiteInfo.description }}</p>
      </div>

      <div class="footer-col">
        <div class="col-title">{{ t('footer.friendLinks') }}</div>
        <ul class="link-list">
          <li v-for="item in footerLinkList" :key="item.title">
            <el-link :href=item.link :underline="false">{{ item.title }}</el-link>
          </li>
        </ul>
      </div>

      <div class="footer-col">
        <div class="col-title">{{ t('footer.contact') }}</div>
        <div class="social-link">
          <el-link v-for="socialLink in socialLinks" :href=socialLink.url :underline="false" :aria-label="socialLink.alt">
            <el-image :src=socialLink.src :alt=socialLink.alt></el-image>
          </el-link>
        </div>
        <div class="version">
          <el-tag size="small">{{ t('footer.version') }}</el-tag>
          <el-tag size="small" type="info">{{ websiteStore.state.websiteInfo.version }}</el-tag>
        </div>
      </div>
    </div>

    <div class="copyright">
      <div class="runtime">{{ t('footer.buildDate') }}：{{ websiteStore.state.websiteInfo.created_at }} · {{ t('footer.uptime') }} {{ elapsedTime }}</div>
      <div class="filing">
        <el-image src="/image/filing.png" alt=""/>
        <el-link href="https://beian.miit.gov.cn/#/Integrated/index" :underline="false">
          {{ websiteStore.state.websiteInfo.icp_filing }}
        </el-link>
        <el-link :href=publicSecurityFilingLink :underline="false">
          {{ websiteStore.state.websiteInfo.public_security_filing }}
        </el-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useWebsiteStore} from "@/stores/website";
import {computed} from "vue";
import {ref} from "vue";
import {onUnmounted} from "vue";
import {type FooterLink, websiteFooterLink} from "@/api/website";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const footerLinkList=ref<FooterLink[]>([])

const getFooterLinkList =async ()=>{
  const res = await websiteFooterLink()
  if (res.code===0){
    footerLinkList.value=res.data
  }
}

getFooterLinkList()

const websiteStore = useWebsiteStore()

let timerId: number | null = null;
const elapsedTime = ref("");

function updateElapsedTime() {
  let creationDate = websiteStore.state.websiteInfo.created_at;
  if (creationDate) {
    let creationTimestamp = new Date(creationDate).getTime();
    let currentTimestamp = new Date().getTime();

    let totalDays = (currentTimestamp - creationTimestamp) / 1000 / (60 * 60 * 24);
    let daysPassed = Math.floor(totalDays);
    let hoursRemaining = Math.floor((totalDays - daysPassed) * 24);
    let minutesRemaining = Math.floor((totalDays - daysPassed - (hoursRemaining / 24)) * 24 * 60);
    let secondsRemaining = Math.floor((totalDays - daysPassed - (hoursRemaining / 24) - (minutesRemaining / 24 / 60)) * 24 * 60 * 60);

    elapsedTime.value = `${daysPassed}${t('footer.days')}${hoursRemaining}${t('footer.hours')}${minutesRemaining}${t('footer.minutes')}${secondsRemaining}${t('footer.seconds')}`;
  }
}

function initializeTimer() {
  updateElapsedTime();
  timerId = setInterval(updateElapsedTime, 1000);
}

onUnmounted(() => {
  clearInterval(timerId as number);
});

initializeTimer();

const publicSecurityFilingLink = computed(() => "http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=" + websiteStore.state.websiteInfo.public_security_filing.match(/\d+/))
const socialLinks = computed(() => [
  {
    src: "/image/bilibili.png",
    alt: "",
    url: websiteStore.state.websiteInfo.bilibili_url
  },
  {
    src: "/image/gitee.png",
    alt: "",
    url: websiteStore.state.websiteInfo.gitee_url
  },
  {
    src: "/image/github.png",
    alt: "",
    url: websiteStore.state.websiteInfo.github_url
  },
])

</script>


<style scoped lang="scss">
.web-footer {
  border-top: 1px solid var(--border);
  background-color: var(--bg-elevated);

  .container {
    display: grid;
    grid-template-columns: 1.5fr 1fr 1fr;
    gap: var(--sp-6);
    max-width: var(--content-width);
    width: 100%;
    margin: 0 auto;
    padding: var(--sp-7) var(--sp-5) var(--sp-6);

    .footer-col {
      .col-title {
        font-size: var(--fs-12);
        font-weight: 600;
        color: var(--text-muted);
        letter-spacing: 0.05em;
        margin-bottom: var(--sp-3);
      }

      .footer-logo {
        .el-image {
          height: 40px;
          width: auto;
        }
      }

      .desc {
        margin-top: var(--sp-3);
        font-size: var(--fs-14);
        color: var(--text-muted);
        line-height: var(--lh-body);
      }

      .link-list {
        list-style: none;
        padding: 0;
        margin: 0;

        li {
          margin-bottom: var(--sp-2);

          .el-link {
            font-size: var(--fs-14);
            color: var(--text-body);

            &:hover {
              color: var(--accent);
            }
          }
        }
      }

      .social-link {
        display: flex;
        gap: var(--sp-3);
        margin-bottom: var(--sp-4);

        .el-link {
          .el-image {
            height: 24px;
            width: 24px;
          }
        }
      }

      .version {
        display: flex;
        gap: var(--sp-2);
      }
    }
  }

  .copyright {
    max-width: var(--content-width);
    width: 100%;
    margin: 0 auto;
    padding: var(--sp-4) var(--sp-5);
    border-top: 1px solid var(--border);
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    gap: var(--sp-2);
    font-size: var(--fs-12);
    color: var(--text-muted);

    .filing {
      display: flex;
      align-items: center;
      gap: var(--sp-2);

      .el-image {
        height: 14px;
        width: auto;
      }

      .el-link {
        font-size: var(--fs-12);
        color: var(--text-muted);
      }
    }
  }
}
</style>
