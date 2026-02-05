<template>
  <div class="mark__render">
    <div class="markdown-body" v-html="md.render(content || '')"></div>
  </div>
</template>
<script>
import { md } from '@/mixins/markdown-it';

export default {
  props: {
    content: '',
  },
  data() {
    return {
      md: md,
    };
  },
  mounted() {
    this.addCopyClick();
  },
  beforeDestroy() {
    if (this.timer) clearTimeout(this.timer);
  },
  methods: {
    addCopyClick() {
      let copyList = document.getElementsByClassName('copy-btn') || [];
      for (let i = 0; i < copyList.length; i++) {
        copyList[i].addEventListener('click', e => {
          let innerText = e.target.parentNode.nextElementSibling.innerText;
          this.$copy(innerText);
          e.target.innerText = this.$t('common.copy.copySuccess');
          this.timer = setTimeout(() => {
            e.target.innerText = this.$t('common.button.copy');
          }, 1500);
        });
      }
    },
  },
};
</script>

<style lang="scss">
@import '@/style/markdown.scss';
.mark__render .markdown-body {
  font-family: 'Microsoft YaHei', Arial, sans-serif;
  color: #333;
}
</style>
