<template>
  <div class="mark__content">
    <div
      class="markdown-body"
      v-html="marked(content || $t('common.noData'))"
    ></div>
  </div>
</template>
<script>
const highlight = require('highlight.js');
import 'highlight.js/styles/github.css';
import { marked } from 'marked';

marked.setOptions({
  renderer: new marked.Renderer(),
  gfm: true,
  tables: true,
  breaks: false,
  pedantic: false,
  sanitize: false,
  smartLists: true,
  smartypants: false,
  highlight: function (code) {
    return highlight.highlightAuto(code).value;
  },
});

export default {
  props: {
    content: '',
  },
  methods: {
    marked,
  },
};
</script>

<style lang="scss" scoped>
@import '@/assets/showDocs/showdoc.scss';
.mark__content .markdown-body {
  background: rgba(255, 255, 255, 0);
}
</style>
