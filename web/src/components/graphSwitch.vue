<template>
  <div class="graph-switch">
    <div class="graph-switch-item">
      <div class="graph-switch-label">
        <span class="label_name">{{ label }}</span>
        <el-tooltip
          class="item"
          effect="dark"
          placement="top-start"
          popper-class="knowledge-graph-tooltip"
        >
          <span class="el-icon-question question-icon"></span>
          <template #content>
            <p
              v-for="(item, i) in tips"
              :key="i"
              class="tooltip-item"
            >
              <span class="tooltip-title">{{ item.title }}</span>
              <span class="tooltip-content">{{ item.content }}</span>
            </p>
          </template>
        </el-tooltip>
      </div>
      <div class="graph-switch-content">
        <el-switch 
          v-model="switchValue" 
          :disabled="disabled"
          @change="handleChange"
        ></el-switch>
      </div>
    </div>
  </div>
</template>

<script>
import { KNOWLEDGE_GRAPH_TIPS } from '@/views/knowledge/config'

export default {
  name: 'GraphSwitch',
  props: {
    value: {
      type: Boolean,
      default: false
    },
    label: {
      type: String,
      default: '知识图谱'
    },
    disabled: {
      type: Boolean,
      default: false
    },
    tips: {
      type: Array,
      default: function() {
        return KNOWLEDGE_GRAPH_TIPS
      }
    }
  },
  computed: {
    switchValue: {
      get() {
        return this.value
      },
      set(val) {
        this.$emit('input', val)
      }
    }
  },
  methods: {
    handleChange(val) {
      this.$emit('change', val)
    }
  }
}
</script>

<style lang="scss" scoped>
.graph-switch {
  .graph-switch-item {
    display: flex;
    align-items: center;
    justify-content:space-between;
    .graph-switch-label {
      display: flex;
      align-items: center;
      min-width: 140px;
      padding-right: 12px;
      font-size: 14px;
      line-height: 40px;
      box-sizing: border-box;
      span{
        font-size:15px;
      }
      .label_name{
        font-weight:600;
      }
      .question-icon {
        cursor: pointer;
        margin-left: 8px;
      }
    }
  }
}
</style>

<style lang="scss">
.knowledge-graph-tooltip {
  max-width: 400px !important;

  .tooltip-item {
    margin: 0;
    padding: 4px 0;

    .tooltip-title {
      font-weight: bold;
      margin-right: 8px;
    }

    .tooltip-content {
      display: inline-block;
    }
  }
}
</style>

