<template>
  <div id="app">
    <router-view/>
  </div>
</template>


<script>

export default {
  name: "app",
  components: {},
  mounted() {
    this.init()
  },
  methods: {
    async init() {
      await this.getData()
      this.getNowData()
    },
    async getData() {
      setTimeout(function(){
        console.log(111)
      }, 5000)
      clearTimeout()
    },
    getNowData() {
      console.log(222)
    }
  }
};
</script>

<style>
#app {
  position:relative;
  width: 100%;
  height: 100%;
  /* margin: 0 auto;
  max-width: 1560px; */
}
</style>
