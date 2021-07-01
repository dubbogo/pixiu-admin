<template>
  <div class="common">
    <div class="common-left">{{title}}</div>
    <div class="common-right">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: 'CommonTitle',
  props: ['title']
}
</script>

<style lang="less" scoped>
.common {
  height: 40px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-left: 12px;
  padding-right: 16px;
  background: #fff;
  box-sizing: border-box;
}
.common-left {
  display: flex;
  align-items: center;
  font-size: 18px;
  color: #303133;
  font-weight: 600;
}
.common-left::before {
  content: '';
  display: inline-block;
  width: 3px;
  height: 14px;
  background: #409EFF;
  margin-right: 5px;
}
</style>