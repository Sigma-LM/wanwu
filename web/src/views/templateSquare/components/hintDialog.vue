<template>
  <el-dialog
    title=""
    :visible.sync="dialogVisible"
    width="500px"
    append-to-body
    :close-on-click-modal="false"
  >
    <div style="text-align: center; margin-bottom: 25px">
      <div>我们提供了丰富的工作流模板，欢迎您前往查看并下载。</div>
      <div class="link-text">跳转</div>
    </div>
  </el-dialog>
</template>

<script>

export default {
  props: {
    linkUrl: ''
  },
  data() {
    return {
      dialogVisible: false
    }
  },
  methods: {
    openDialog() {
      this.dialogVisible = true
    }
  },
};
</script>

<style lang="scss" scoped>
.link-text {
  color: $color;
  text-decoration: underline;
  cursor: pointer;
}
</style>
