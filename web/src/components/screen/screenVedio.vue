<template>
    <div class="cvedio">
        <div class="cvedio-left ve">
            <div class="txt">主驾侧摄像头</div>
            <div class="player"></div>
        </div>
        <div class="cvedio-right ve">
            <div class="txt">侧驾侧摄像头</div>
            <div class="player"></div>
        </div>
    </div>
</template>
<script>
export default {
    name:'screenVedio',
    data(){
        return{

        }
    }
}
</script>
<style scoped lang="scss">
    .cvedio{
        display: flex;
        padding: 20px 32px;
        height: 340px;
        // line-height: 340px;
        .ve{
            border:1px solid rgb(161, 161, 88);
            flex:1;
            color: #fff;
            display: flex;
            font-size: 18px;
            align-items: center;
            padding: 10px;
            .txt{
                width: 18px;
                margin: 0 12px;
            }
            .player{
                flex: 1;
            }
        }   
        &-left{
            
        }
        &-right{
            
        }
    }
</style>