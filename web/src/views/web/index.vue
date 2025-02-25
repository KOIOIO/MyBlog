<template>
{{ url }}
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {qqLoginURL} from "@/api/base.ts";
const usr=ref
const getUrl = async ()=>{
  const res =await qqLoginURL()
  if(res.code===0){
    url.value=res.data
  }
}
</script>
getUrl()
<style scoped lang="scss">
</style>