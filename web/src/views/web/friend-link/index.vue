<template>
  <div class="friend-link">
    <web-navbar :noScroll="true"/>
    <el-container class="main-content">
      <div class="container">
        <el-main>
          <el-row class="title">友链列表</el-row>
          <div class="list">
            <el-card v-for="item in friendLinkList" @click="handleFriendLinkJumps(item.link)">
              <div class="logo">
                <el-image style="width: 64px; height: 64px" :src="item.logo" alt=""></el-image>
                <el-row class="name">{{ item.name }}</el-row>
              </div>
              <div class="description">
                <el-text>{{ item.description }}</el-text>
              </div>
            </el-card>
          </div>
        </el-main>
      </div>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import WebNavbar from "@/components/layout/WebNavbar.vue";
import {type FriendLink, friendLinkInfo} from "@/api/friend-link";
import {ref} from "vue";

const friendLinkList = ref<FriendLink[]>([])

const getFriendLinkInfo = async () => {
  const res = await friendLinkInfo()
  if (res.code === 0) {
    friendLinkList.value = res.data.list
  }
}

getFriendLinkInfo()

const handleFriendLinkJumps = (link:string)=>{
  window.open(link)
}
</script>

<style scoped lang="scss">
.friend-link {
  background-color: #1e1e2f;
  color: #ffffff;

  .main-content {
    margin-top: 70px;
    display: flex;
    justify-content: center;

    .container {
      display: flex;
      max-width: 1400px;
      width: 100%;

      .el-main {
        .title {
          font-size: 24px;
          margin-bottom: 20px;
          color: #8a8aff;
        }

        .list {
          display: flex;
          flex-wrap: wrap;
          gap: 20px;

          .el-card {
            --el-card-padding: 0;
            width: 25%;
            height: 130px;
            background-color: #2a2a3d;
            border: 1px solid #3a3a5c;
            transition: transform 0.3s, box-shadow 0.3s;

            &:hover {
              transform: translateY(-10px);
              box-shadow: 0 10px 20px rgba(0, 0, 0, 0.5);
            }

            .logo {
              display: flex;

              .name {
                font-size: 24px;
                margin-top: auto;
                margin-bottom: auto;
                color: #8a8aff;
              }
            }

            .description {
              margin-left: 10px;
              margin-right: 10px;
              color: #b0b0ff;
            }
          }
        }
      }
    }
  }
}
</style>
