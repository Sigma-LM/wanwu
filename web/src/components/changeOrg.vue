<template>
  <div class="header__org_container">
    <div class="header__org_wrapper" :title="getCurrentOrgName()">
      <el-select
        v-model="org.orgId"
        :placeholder="$t('header.org.placeholder')"
        filterable
        class="header__org_select"
        v-if="orgList && orgList.length"
        @change="changeOrg"
      >
        <el-option
          v-for="(item, index) in orgList"
          :command="index"
          :key="item.id + index"
          :class="org.orgId === item.id ? 'header__org_active' : ''"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ChangeOrg',
  props: {
    orgList: [],
    org: { orgId: '' },
    changeOrg: { type: Function, required: true },
    getCurrentOrgName: { type: Function, required: true },
  },
};
</script>

<style lang="scss" scoped>
/*.header__org_container {
  padding: 10px 15px 0 15px;
  .header__org_wrapper {
    padding-bottom: 8px;
    !*border-bottom: 1px solid #ebebeb;*!
  }
}*/

.header__org_active {
  color: $color !important;
}
.header__org_select,
.menu__org_select ::v-deep {
  width: 100%;
  .el-input__inner:focus,
  .el-input__inner:hover,
  .el-input.is-focus .el-input__inner {
    border-color: #fff !important; // #dcdfe6
  }
  .el-input__inner {
    background-color: rgba(255, 255, 255, 0);
    border: 1px solid #fff;
    color: $color_title;
    font-weight: bold;
    padding-left: 10px;
  }
  .el-input__inner::placeholder {
    color: rgba(18, 18, 18, 0.7);
  }
  .el-input {
    .el-select__caret {
      color: #aaa;
      font-size: 15px;
    }
  }
}
.menu__org_select ::v-deep {
  width: 190px;
  .el-input__inner {
    background-color: rgba(255, 255, 255, 0);
    border: none !important;
    color: $color_title !important;
    font-weight: normal;
    padding-left: 0 !important;
    margin-left: 0 !important;
  }
}
</style>
