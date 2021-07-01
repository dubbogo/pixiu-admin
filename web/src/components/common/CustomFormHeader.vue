<!--
 * @Description: 表单标题
 * @Author: chenyaru
 * @Date: 2020-02-24 14:16:56
 -->
<template>
<div class="title">
    <div class="line--vertical"></div>
    <div class="name">{{ title }}</div>
    <div class="search">
        <slot name="searchContent"></slot>
    </div>
</div>
</template>

<script>
export default {
    name: 'CustomFormHeader',
    props: ['title']
}
</script>

<style scoped>
.title {
    /* height: 46px; */
    /* line-height: 46px; */
    color: rgba(48, 49, 51, 1);
    margin-bottom: 20px;
    height: 40px;
    background: rgba(248, 248, 248, 1);
}

.name {
    font-size: 15px;
    font-family: PingFangSC-Regular, PingFang SC;
    font-weight: 400;
    line-height: 40px;
}

.title div {
    float: left;
}

.line--vertical {
    width: 3px;
    height: 14px;
    background: rgba(64, 158, 255, 1);
    margin: 13px 7px 13px 12px;
}

.search {
    float: right !important;
    margin-top: 4px;
    height: 32px;
    padding-right: 16px;
}
</style>
