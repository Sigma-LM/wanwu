<template>
  <div class="tempSquare-management">
    <div class="tempSquare-content-box tempSquare-third">
      <div class="tempSquare-main">
        <div class="tempSquare-content">
          <div class="tempSquare-card-box">
            <div class="card-search card-search-cust">
              <search-input
                style="margin-right: 2px"
                :placeholder="$t('tempSquare.searchText')"
                ref="searchInput"
                @handleSearch="doGetWorkflowTempList"
              />
            </div>

            <div class="card-loading-box" v-if="list.length">
              <div class="card-box" v-loading="loading">
                <div
                  class="card"
                  v-for="(item, index) in list"
                  :key="index"
                  @click.stop="handleClick(item)"
                >
                  <div class="card-title">
                    <img
                      class="card-logo"
                      v-if="item.avatar && item.avatar.path"
                      :src="item.avatar.path"
                    />
                    <div class="mcp_detailBox">
                      <span class="mcp_name">{{ item.name }}</span>
                      <span class="mcp_from">
                        <label>
                          {{ $t('tempSquare.author') }}：{{ item.author }}
                        </label>
                      </span>
                    </div>
                  </div>
                  <div class="card-des">{{ item.desc }}</div>
                  <div class="card-bottom">
                    <div class="card-bottom-left">
                      {{ $t('tempSquare.downloadCount') }}：{{
                        item.downloadCount || 0
                      }}
                    </div>
                    <div class="card-bottom-right">
                      <i
                        class="el-icon-download"
                        :title="$t('tempSquare.download')"
                        @click.stop="downloadTemplate(item)"
                      ></i>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="empty">
              <el-empty :description="$t('common.noData')"></el-empty>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { getWorkflowTempList, downloadWorkflow } from '@/api/templateSquare';
import SearchInput from '@/components/searchInput.vue';
export default {
  components: { SearchInput },
  props: {
    isPublic: true,
    type: '',
  },
  data() {
    return {
      basePath: this.$basePath,
      category: this.$t('square.all'),
      list: [],
      templateUrl: '',
      loading: false,
    };
  },
  mounted() {
    this.doGetWorkflowTempList();
  },
  methods: {
    doGetWorkflowTempList() {
      const searchInput = this.$refs.searchInput;
      const params = {
        name: searchInput.value,
      };

      getWorkflowTempList(params)
        .then(res => {
          const { list } = res.data || {};
          this.list = list || [];
          this.loading = false;
        })
        .catch(() => (this.loading = false));
    },
    downloadTemplate(item) {
      downloadWorkflow({ templateId: item.templateId }).then(response => {
        const blob = new Blob([response], { type: response.type });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = item.name + '.json';
        link.click();
        window.URL.revokeObjectURL(link.href);
        this.doGetWorkflowTempList();
      });
    },
    handleClick(val) {
      const path = `${this.isPublic ? '/public' : ''}/templateSquare/detail`;
      this.$router.push({
        path,
        query: { templateSquareId: val.templateId, type: this.type },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
@import '@/style/tempSquare.scss';
.tempSquare-management {
  .card-search-cust {
    justify-content: flex-end;
    margin-top: -35px;
  }
}
</style>
