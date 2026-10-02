<template>
  <div class="site-config">
    <el-col :span="12">
      <div class="website-info">
        <div class="page-title">{{ t('system.config.site.websiteInfo') }}</div>
        <div class="content">
          <el-form
              :model="websiteInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.site.logo')">
              <el-upload
                  :action="`${path}/image/upload`"
                  drag
                  with-credentials
                  :headers="{'x-access-token':userStore.state.accessToken}"
                  :show-file-list="false"
                  :on-success="handleLogoSuccess"
                  :on-error="handleLogoSuccess"
                  name="image"
              >

                <el-image v-if="websiteInfo.logo" :src="websiteInfo.logo"
                          alt=""/>
                <div v-else class="upload-content">
                  <div class="container">
                    <component is="UploadFilled" class="upload-filled"></component>
                    <div class="el-upload__text">
                      {{ t('system.config.site.dragPrefix') }}<em>{{ t('system.config.site.dragClick') }}</em>
                    </div>
                  </div>
                </div>

                <template #tip>
                  <div class="clear-button">
                    <el-button v-if="websiteInfo.logo" icon="Delete" type="danger" @click="clearLogo"/>
                  </div>
                  <div class="el-upload__tip">
                    {{ t('system.config.site.uploadTip') }}
                  </div>
                </template>
              </el-upload>

              <el-input
                  v-model="websiteInfo.logo"
                  size="large"
                  disabled
              />
            </el-form-item>
            <el-form-item :label="t('system.config.site.fullLogo')">
              <el-upload
                  :action="`${path}/image/upload`"
                  drag
                  with-credentials
                  :headers="{'x-access-token':userStore.state.accessToken}"
                  :show-file-list="false"
                  :on-success="handleFullLogoSuccess"
                  :on-error="handleFullLogoSuccess"
                  name="image"
              >

                <el-image v-if="websiteInfo.full_logo" :src="websiteInfo.full_logo"
                          alt=""/>

                <div v-else class="upload-content">
                  <div class="container">
                    <component is="UploadFilled" class="upload-filled"></component>
                    <div class="el-upload__text">
                      {{ t('system.config.site.dragPrefix') }}<em>{{ t('system.config.site.dragClick') }}</em>
                    </div>
                  </div>
                </div>

                <template #tip>
                  <div class="clear-button">
                    <el-button v-if="websiteInfo.full_logo" icon="Delete" type="danger" @click="clearFullLogo"/>
                  </div>
                  <div class="el-upload__tip">
                    {{ t('system.config.site.uploadTip') }}
                  </div>
                </template>
              </el-upload>

              <el-input
                  v-model="websiteInfo.full_logo"
                  size="large"
                  disabled
              />
            </el-form-item>
            <el-form-item :label="t('system.config.site.siteTitle')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.title"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.slogan')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.slogan"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.sloganEn')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.slogan_en"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.description')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.description" type="textarea" :rows="4"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.version')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.version"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.createdAt')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.created_at"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.icp')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.icp_filing"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.police')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.public_security_filing"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.bilibili')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.bilibili_url"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.gitee')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.gitee_url"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.github')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.github_url"/>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-col>
    <el-col :span="12">
      <div class="personal-info">
        <div class="page-title">{{ t('system.config.site.personalInfo') }}</div>
        <div class="content">
          <el-form
              :model="websiteInfo"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.site.nickname')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.name"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.job')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.job"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.address')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.address"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.email')">
              <el-input @change="updateWebsiteInfo" v-model="websiteInfo.email"/>
            </el-form-item>
            <el-form-item :label="t('system.config.site.qqImage')">
              <el-upload
                  :action="`${path}/image/upload`"
                  drag
                  with-credentials
                  :headers="{'x-access-token':userStore.state.accessToken}"
                  :show-file-list="false"
                  :on-success="handleQQImageSuccess"
                  :on-error="handleQQImageSuccess"
                  name="image"
              >

                <el-image v-if="websiteInfo.qq_image" :src="websiteInfo.qq_image"
                          alt=""/>

                <div v-else class="upload-content">
                  <div class="container">
                    <component is="UploadFilled" class="upload-filled"></component>
                    <div class="el-upload__text">
                      {{ t('system.config.site.dragPrefix') }}<em>{{ t('system.config.site.dragClick') }}</em>
                    </div>
                  </div>
                </div>

                <template #tip>
                  <div class="clear-button">
                    <el-button v-if="websiteInfo.qq_image" icon="Delete" type="danger" @click="clearQQImageLogo"/>
                  </div>
                  <div class="el-upload__tip">
                    {{ t('system.config.site.uploadTip') }}
                  </div>
                </template>
              </el-upload>

              <el-input
                  v-model="websiteInfo.qq_image"
                  size="large"
                  disabled
              />
            </el-form-item>
            <el-form-item :label="t('system.config.site.wechatImage')">
              <el-upload
                  :action="`${path}/image/upload`"
                  drag
                  with-credentials
                  :headers="{'x-access-token':userStore.state.accessToken}"
                  :show-file-list="false"
                  :on-success="handleWechatImageSuccess"
                  :on-error="handleWechatImageSuccess"
                  name="image"
              >

                <el-image v-if="websiteInfo.wechat_image" :src="websiteInfo.wechat_image"
                          alt=""/>

                <div v-else class="upload-content">
                  <div class="container">
                    <component is="UploadFilled" class="upload-filled"></component>
                    <div class="el-upload__text">
                      {{ t('system.config.site.dragPrefix') }}<em>{{ t('system.config.site.dragClick') }}</em>
                    </div>
                  </div>
                </div>

                <template #tip>
                  <div class="clear-button">
                    <el-button v-if="websiteInfo.wechat_image" icon="Delete" type="danger"
                               @click="clearWechatImageLogo"/>
                  </div>
                  <div class="el-upload__tip">
                    {{ t('system.config.site.uploadTip') }}
                  </div>
                </template>
              </el-upload>

              <el-input
                  v-model="websiteInfo.wechat_image"
                  size="large"
                  disabled
              />
            </el-form-item>
          </el-form>
        </div>
      </div>

      <div class="footer-link">
        <div class="page-title">{{ t('system.config.site.footerLink') }}</div>
        <div class="content">
          <el-form
              :model="footerLinkList"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <template v-for="item in footerLinkList">
              <el-form-item :label="item.title">
                <el-input
                    v-model="item.link"
                    size="large"
                    disabled
                />
              </el-form-item>
              <div class="delete-button">
                <el-button icon="Delete" type="danger" @click="handleDeleteFooterLink(item)"/>
              </div>
            </template>
          </el-form>
          <div class="button-group">
            <el-button v-if="!isShow" type="success" @click="isShow=true">{{ t('system.config.site.create') }}</el-button>
            <el-button v-if="isShow" type="primary" @click="isShow=false;handleCreateFooterLink(footerLink)">{{ t('common.confirm') }}</el-button>
            <el-button v-if="isShow" @click="isShow=false">{{ t('common.cancel') }}</el-button>
          </div>
          <el-form
              v-if="isShow"
              :model="footerLink"
              :validate-on-rule-change="false"
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.site.title')">
              <el-input
                  v-model="footerLink.title"
                  size="large"
              />
            </el-form-item>
            <el-form-item :label="t('system.config.site.link')">
              <el-input
                  v-model="footerLink.link"
                  size="large"
              />
            </el-form-item>
          </el-form>
        </div>
      </div>

      <div class="carousel-info">
        <div class="page-title">{{ t('system.config.site.homeImage') }}</div>
        <div class="content">
          <template v-for="item in carouselList">
            <div class="carousel-item">
              <el-image :src="item" alt=""/>
              <el-button icon="Delete" type="danger" @click="cancelCarousel(item)"/>
            </div>
          </template>

          <el-form
              label-width="auto"
              style="max-width: 400px"
          >
            <el-form-item :label="t('system.config.site.imageUpload')">
              <el-upload
                  :action="`${path}/image/upload`"
                  drag
                  with-credentials
                  :headers="{'x-access-token':userStore.state.accessToken}"
                  :show-file-list="false"
                  :on-success="handleCarouselSuccess"
                  :on-error="handleCarouselSuccess"
                  name="image"
              >

                <div class="upload-content">
                  <div class="container">
                    <component is="UploadFilled" class="upload-filled"></component>
                    <div class="el-upload__text">
                      {{ t('system.config.site.dragPrefix') }}<em>{{ t('system.config.site.dragClick') }}</em>
                    </div>
                  </div>
                </div>

                <template #tip>
                  <div class="el-upload__tip">
                    {{ t('system.config.site.uploadTip') }}
                  </div>
                </template>
              </el-upload>
            </el-form-item>
          </el-form>
        </div>
      </div>
    </el-col>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref, watch} from "vue";
import {getWebsite, updateWebsite, type Website} from "@/api/config";
import {useUserStore} from "@/stores/user";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";
import {ElMessage} from "element-plus";
import {useWebsiteStore} from "@/stores/website";
import {
  type FooterLink,
  websiteAddCarousel,
  websiteCancelCarousel,
  websiteCarousel,
  type WebsiteCarouselOperation, websiteCreateFooterLink, websiteDeleteFooterLink, websiteFooterLink
} from "@/api/website";
import {useI18n} from "vue-i18n";

const {t} = useI18n()

const path = ref(import.meta.env.VITE_BASE_API)
const userStore = useUserStore()

const websiteInfo = ref<Website>({
  logo: '',
  full_logo: '',
  title: '',
  slogan: '',
  slogan_en: '',
  description: '',
  version: '',
  created_at: '',
  icp_filing: '',
  public_security_filing: '',
  bilibili_url: '',
  gitee_url: '',
  github_url: '',
  name: '',
  job: '',
  address: '',
  email: '',
  qq_image: '',
  wechat_image: '',
})

const getWebsiteInfo = async () => {
  const res = await getWebsite()
  if (res.code === 0) {
    websiteInfo.value = res.data
    useWebsiteStore().state.websiteInfo = res.data
  }
}

getWebsiteInfo()

const shouldRefreshInfo = ref(false)
watch(() => shouldRefreshInfo.value, (newVal) => {
  if (newVal) {
    getWebsiteInfo()
    shouldRefreshInfo.value = false
  }
})

const updateWebsiteInfo = async () => {
  const res = await updateWebsite(websiteInfo.value)
  console.log(websiteInfo.value)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    shouldRefreshInfo.value = true
  }
}

const handleLogoSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    websiteInfo.value.logo = res.data.url
    ElMessage.success(res.msg)
    updateWebsiteInfo()
  }
}

const clearLogo = () => {
  ElMessageBox.confirm(
      t('system.config.site.clearLogoConfirm'),
      t('system.config.site.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      })
      .then(() => {
        websiteInfo.value.logo = ''
        updateWebsiteInfo()
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('system.config.site.operationCancelled'),
        })
      })
}

const handleFullLogoSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    websiteInfo.value.full_logo = res.data.url
    ElMessage.success(res.msg)
    updateWebsiteInfo()
  }
}

const clearFullLogo = () => {
  ElMessageBox.confirm(
      t('system.config.site.clearFullLogoConfirm'),
      t('system.config.site.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      })
      .then(() => {
        websiteInfo.value.full_logo = ''
        updateWebsiteInfo()
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('system.config.site.operationCancelled'),
        })
      })
}

const handleQQImageSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    websiteInfo.value.qq_image = res.data.url
    ElMessage.success(res.msg)
    updateWebsiteInfo()
  }
}

const clearQQImageLogo = () => {
  ElMessageBox.confirm(
      t('system.config.site.clearQqImageConfirm'),
      t('system.config.site.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      })
      .then(() => {
        websiteInfo.value.qq_image = ''
        updateWebsiteInfo()
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('system.config.site.operationCancelled'),
        })
      })
}

const handleWechatImageSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    websiteInfo.value.wechat_image = res.data.url
    ElMessage.success(res.msg)
    updateWebsiteInfo()
  }
}

const clearWechatImageLogo = () => {
  ElMessageBox.confirm(
      t('system.config.site.clearWechatImageConfirm'),
      t('system.config.site.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      })
      .then(() => {
        websiteInfo.value.wechat_image = ''
        updateWebsiteInfo()
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('system.config.site.operationCancelled'),
        })
      })
}

const footerLinkList = ref<FooterLink[]>([])

const getFooterLinkList = async () => {
  const res = await websiteFooterLink()
  if (res.code === 0) {
    footerLinkList.value = res.data
  }
}

getFooterLinkList()

const shouldRefreshFooterLinkInfo = ref(false)
watch(() => shouldRefreshFooterLinkInfo.value, (newVal) => {
  if (newVal) {
    getFooterLinkList()
    shouldRefreshFooterLinkInfo.value = false
  }
})

const handleDeleteFooterLink = async (item: FooterLink) => {
  const res = await websiteDeleteFooterLink(item)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    shouldRefreshFooterLinkInfo.value = true
  }
}

const isShow = ref(false)

const footerLink = reactive<FooterLink>({
  title: '',
  link: '',
})

const handleCreateFooterLink = async (footerLink: FooterLink) => {
  const res = await websiteCreateFooterLink(footerLink)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    shouldRefreshFooterLinkInfo.value = true
  }
}


const carouselList = ref<string[]>([])

const getCarouseList = async () => {
  const res = await websiteCarousel()
  if (res.code === 0) {
    carouselList.value = res.data
  }
}

getCarouseList()

const addCarouse = async (url: string) => {
  const req: WebsiteCarouselOperation = {
    url: url,
  }
  const res = await websiteAddCarousel(req)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  }
}

const handleCarouselSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    addCarouse(res.data.url)
    ElMessage.success(res.msg)
    getCarouseList()
  }
}

const cancelCarousel = (url: string) => {
  ElMessageBox.confirm(
      t('system.config.site.removeCarouselConfirm'),
      t('system.config.site.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
      })
      .then(async () => {
        const req: WebsiteCarouselOperation = {
          url: url,
        }
        const res = await websiteCancelCarousel(req)
        if (res.code === 0) {
          ElMessage.success(res.msg)
          await getCarouseList()
        }
      })
      .catch(() => {
        ElMessage({
          type: 'info',
          message: t('system.config.site.operationCancelled'),
        })
      })

}
</script>

<style scoped lang="scss">
.site-config {
  .page-title {
    font-size: var(--fs-20);
    font-weight: 600;
    color: var(--text-primary);
    line-height: var(--lh-title);
    margin-bottom: var(--sp-4);
  }

  .content {
    background: var(--bg-elevated);
    border: 1px solid var(--border);
    border-radius: var(--radius-md);
    padding: var(--sp-5);
    margin-bottom: var(--sp-5);

    .el-form {
      .el-form-item {
        .el-image {
          height: 120px;
          border-radius: var(--radius-sm);
        }

        .upload-content {
          display: flex;
          height: 60px;

          .container {
            margin: auto;

            .upload-filled {
              height: 32px;
              width: 32px;
              color: var(--text-muted);
            }
          }
        }
      }
    }

    .carousel-item {
      display: flex;
      align-items: center;
      padding: var(--sp-3);
      max-width: 400px;

      .el-image {
        height: 160px;
        margin-right: var(--sp-4);
        border-radius: var(--radius-sm);
      }
    }
  }
}
</style>
<style lang="scss">
.el-upload {
  --el-upload-dragger-padding-horizontal: 0px;
  --el-upload-dragger-padding-vertical: 0px;
  line-height: 0;
  border: 2px dashed var(--border);
  border-radius: var(--radius-sm);
  transition: border-color 150ms ease-out;
}

.el-upload:hover {
  border-color: var(--accent);
}
</style>
