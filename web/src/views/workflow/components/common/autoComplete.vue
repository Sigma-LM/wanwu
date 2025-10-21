<template>
  <div class="auto-complete" ref="container">
    <textarea
      class="auto-complete-textarea"
      rows="2"
      ref="input"
      :placeholder="placeholder"
      v-model="inputValue"
      @input="handleInput"
      @keydown.down="moveDown"
      @keydown.up="moveUp"
      @keydown.enter="e => selectItem(e)"
      @keydown.esc="closeList"
      @blur="handleBlur"
    />
    <ul v-show="showSuggestions" class="suggestions">
      <li
        v-for="(item, index) in filteredSuggestions"
        :key="index"
        :class="{ active: index === activeIndex }"
        @mousedown="e => selectItem(e, item)"
      >
        {{ item }}
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  props: {
    suggestions: {
      type: Array,
      default: () => []
    },
    trigger: {
      type: String,
      default: '{{'
    },
    placeholder: ''
  },

  data() {
    return {
      inputValue: '',
      showSuggestions: false,
      activeIndex: -1,
      lastTriggerPosition: -1
    }
  },

  computed: {
    filteredSuggestions() {
      const searchTerm = this.getSearchTerm();
      return this.suggestions.filter(item => item.toLowerCase().includes(searchTerm.toLowerCase()))
    }
  },

  methods: {
    handleInput() {
      const cursorPos = this.$refs.input.selectionStart;
      const value = this.inputValue;

      // 检查是否输入了触发字符
      if (value.slice(cursorPos - this.trigger.length, cursorPos) === this.trigger) {
        this.lastTriggerPosition = cursorPos - this.trigger.length;
        this.showSuggestions = true;
        this.activeIndex = 0;
        return;
      }

      // 检查是否在触发字符后输入
      if (this.showSuggestions && cursorPos > this.lastTriggerPosition) {
        return;
      }

      this.showSuggestions = false;
      this.$emit('changeValue', this.inputValue);
    },
    getSearchTerm() {
      if (this.lastTriggerPosition === -1) return '';
      const cursorPos = this.$refs.input.selectionStart;
      return this.inputValue.slice(this.lastTriggerPosition + this.trigger.length, cursorPos);
    },
    moveDown() {
      if (!this.showSuggestions) return;
      this.activeIndex = Math.min(this.activeIndex + 1, this.filteredSuggestions.length - 1);
    },
    moveUp() {
      if (!this.showSuggestions) return;
      this.activeIndex = Math.max(this.activeIndex - 1, 0);
    },
    selectItem(e, item) {
      e.stopPropagation()
      e.preventDefault()

      if (!item && this.activeIndex >= 0) {
        item = this.filteredSuggestions[this.activeIndex];
      }

      if (item) {
        // 替换触发字符后的内容为选中的建议项
        const newValue =
          this.inputValue.slice(0, this.lastTriggerPosition) +
          this.trigger + item + '}}' +
          this.inputValue.slice(this.$refs.input.selectionStart) + ' ';

        this.inputValue = newValue;
        this.$emit('select', item);
        this.$emit('changeValue', this.inputValue);
      }

      this.closeList();
    },
    closeList() {
      this.showSuggestions = false;
      this.activeIndex = -1;
      this.lastTriggerPosition = -1;
    },
    handleBlur() {
      // 延迟关闭以允许选择列表项
      setTimeout(() => this.closeList(), 200);
    }
  }
}
</script>

<style lang="scss" scoped>
.auto-complete {
  position: relative;
  display: inline-block;
  width: 100%;
  .auto-complete-textarea {
    padding: 8px;
    width: 100%;
    border: 1px solid #dcdfe6;
    border-radius: 4px;
    outline: none;
    font-family: inherit;
  }
  .auto-complete-textarea::placeholder {
    color: #ccd0d6;
  }

  .auto-complete-textarea:focus, .auto-complete-textarea:hover {
    border-color: $color;
  }

  .suggestions {
    position: absolute;
    bottom: 100%;
    left: 5px;
    right: 5px;
    margin-top: 4px;
    padding: 5px 10px;
    list-style: none;
    border-radius: 4px;
    background: #fff;
    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, 0.15);
    max-height: 200px;
    overflow-y: auto;
    z-index: 1000;
  }

  .suggestions li {
    padding: 8px 12px;
    border-radius: 8px;
    cursor: pointer;
  }

  .suggestions li.active,
  .suggestions li:hover {
    background-color: #F5F7FA;
  }
}
</style>