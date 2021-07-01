<template>
    <div class = "view">
        <div class="view-header">
            <screenHeader></screenHeader>
        </div>
        <div class="view-vedio">
            <screenVedio></screenVedio>
        </div>
        <div class="view-footer">
            <screenBtn></screenBtn>
        </div>
    </div>
</template>
<script>
import screenHeader from '@/components/screen/screenHeader'
import screenVedio from '@/components/screen/screenVedio'
import screenBtn from '@/components/screen/screenBtn'
export default {
    name:'dataView',
    components:{screenHeader, screenVedio, screenBtn},
    data(){
        return{

        }
    }
}
</script>
<style scoped>
    .view{
        background-color: rgb(6, 5, 5);
        width: 100%;
        height: 100%;
    }
</style>