/**
 * 通用 mixins 方法
 * 提供项目中常用的工具方法和生命周期钩子
 */
export default {
  data() {
    return {
    }
  },

  methods: {
    /**
     * 通用拖拽封装
     * @param {Object} opt
     * @param {string} [opt.containerSelector='.editable-wp'] - 监听拖拽的容器选择器
     * @param {(files:Array,ctx:Object)=>void} [opt.onFiles] - 文件落下后的处理回调
     */
    $setupDragAndDrop(opt = {}) {
      const { containerSelector = '.editable-wp', onFiles } = opt
      const wrap = this.$el && this.$el.querySelector ? this.$el.querySelector(containerSelector) : null
      if (!wrap) return () => {}

      const prevent = (e) => { e.preventDefault(); e.stopPropagation(); wrap.classList.add('is-dropping'); }
      const leave = () => { wrap.classList.remove('is-dropping'); }
      const onDrop = async (e) => {
        prevent(e)
        try {
          const dt = e && e.dataTransfer
          const fileList = (dt && dt.files) ? dt.files : []
          const rawFiles = Array.prototype.slice.call(fileList)
          if (!rawFiles.length) return

          // 安全限制：数量/大小/类型白名单
          const maxFiles = Number(opt.maxFiles || 3)
          const maxSizeMB = Number(opt.maxSizeMB || 50) // 单个文件默认 50MB
          const maxSize = maxSizeMB * 1024 * 1024
          const allowExt = (opt.acceptExt || ['jpg','jpeg','png','gif','webp','bmp','svg','mp3','wav','ogg','txt','pdf','doc','docx','xlsx','xls','pptx','csv','html']).map(function(s){return String(s).toLowerCase()})

          // 过滤非法/过大文件
          const safeFiles = []
          const rejected = []
          for (var i = 0; i < rawFiles.length && safeFiles.length < maxFiles; i++) {
            var f = rawFiles[i]
            var ext = (f.name && f.name.split('.').pop() || '').toLowerCase()
            var okType = allowExt.indexOf(ext) > -1 || (f.type && (f.type.indexOf('image/') === 0 || f.type.indexOf('audio/') === 0))
            if (!okType || (typeof f.size === 'number' && f.size > maxSize)) {
              rejected.push(f)
              continue
            }
            safeFiles.push(f)
          }

          // 提示被拒文件
          if (rejected.length && this && this.$message && this.$message.warning) {
            this.$message.warning('部分文件类型不支持或体积过大，已自动忽略')
          }

          if (!safeFiles.length) return

          // 覆盖前释放旧的 ObjectURL，避免内存泄漏
          try {
            var currentList = this && this.fileList
            if (currentList && currentList.forEach) {
              currentList.forEach(function(f){
                try { if (f && f.fileUrl) URL.revokeObjectURL(f.fileUrl) } catch(e) {}
                try { if (f && f.imgUrl) URL.revokeObjectURL(f.imgUrl) } catch(e) {}
              })
            }
          } catch(err) {}

          if (typeof onFiles === 'function') {
            onFiles(safeFiles, { event: e, wrap })
          }
        } finally {
          leave()
        }
      }

      wrap.addEventListener('dragenter', prevent)
      wrap.addEventListener('dragover', prevent)
      wrap.addEventListener('dragleave', leave)
      wrap.addEventListener('drop', onDrop)

      const cleanup = () => {
        wrap.removeEventListener('dragenter', prevent)
        wrap.removeEventListener('dragover', prevent)
        wrap.removeEventListener('dragleave', leave)
        wrap.removeEventListener('drop', onDrop)
      }

      this.$once('hook:beforeDestroy', () => {
        try { cleanup() } catch (e) {}
      })

      return cleanup
    },
    /**
     * 格式化日期
     * @param {Date|string|number} date - 日期
     * @param {string} format - 格式字符串
     * @returns {string} - 格式化后的日期字符串
     */
    $formatDate(date, format = 'YYYY-MM-DD HH:mm:ss') {
      if (!date) return ''
      const d = new Date(date)
      if (isNaN(d.getTime())) return ''
      
      const year = d.getFullYear()
      const month = String(d.getMonth() + 1).padStart(2, '0')
      const day = String(d.getDate()).padStart(2, '0')
      const hours = String(d.getHours()).padStart(2, '0')
      const minutes = String(d.getMinutes()).padStart(2, '0')
      const seconds = String(d.getSeconds()).padStart(2, '0')
      
      return format
        .replace('YYYY', year)
        .replace('MM', month)
        .replace('DD', day)
        .replace('HH', hours)
        .replace('mm', minutes)
        .replace('ss', seconds)
    },

    /**
     * 深拷贝对象
     * @param {any} obj - 要拷贝的对象
     * @returns {any} - 拷贝后的对象
     */
    $deepClone(obj) {
      if (obj === null || typeof obj !== 'object') return obj
      if (obj instanceof Date) return new Date(obj.getTime())
      if (obj instanceof Array) return obj.map(item => this.$deepClone(item))
      if (typeof obj === 'object') {
        const clonedObj = {}
        for (const key in obj) {
          if (obj.hasOwnProperty(key)) {
            clonedObj[key] = this.$deepClone(obj[key])
          }
        }
        return clonedObj
      }
    },

    /**
     * 防抖函数
     * @param {Function} func - 要防抖的函数
     * @param {number} delay - 延迟时间（毫秒）
     * @returns {Function} - 防抖后的函数
     */
    $debounce(func, delay = 300) {
      let timeoutId
      return function (...args) {
        clearTimeout(timeoutId)
        timeoutId = setTimeout(() => func.apply(this, args), delay)
      }
    },

    /**
     * 节流函数
     * @param {Function} func - 要节流的函数
     * @param {number} delay - 延迟时间（毫秒）
     * @returns {Function} - 节流后的函数
     */
    $throttle(func, delay = 300) {
      let lastCall = 0
      return function (...args) {
        const now = Date.now()
        if (now - lastCall >= delay) {
          lastCall = now
          func.apply(this, args)
        }
      }
    },
    /**
     * 获取文件大小格式化字符串
     * @param {number} bytes - 字节数
     * @returns {string} - 格式化后的文件大小
     */
    $formatFileSize(bytes) {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    },

    /**
     * 验证邮箱格式
     * @param {string} email - 邮箱地址
     * @returns {boolean} - 是否有效
     */
    $isValidEmail(email) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      return emailRegex.test(email)
    },

    /**
     * 验证手机号格式
     * @param {string} phone - 手机号
     * @returns {boolean} - 是否有效
     */
    $isValidPhone(phone) {
      const phoneRegex = /^1[3-9]\d{9}$/
      return phoneRegex.test(phone)
    },

    /**
     * 滚动到页面顶部
     */
    $scrollToTop() {
      window.scrollTo({
        top: 0,
        behavior: 'smooth'
      })
    },

    /**
     * 复制文本到剪贴板
     * @param {string} text - 要复制的文本
     * @returns {Promise} - 复制结果
     */
    async $copyToClipboard(text) {
      try {
        if (navigator.clipboard) {
          await navigator.clipboard.writeText(text)
          this.$success('复制成功')
        } else {
          // 兼容旧浏览器
          const textArea = document.createElement('textarea')
          textArea.value = text
          document.body.appendChild(textArea)
          textArea.select()
          document.execCommand('copy')
          document.body.removeChild(textArea)
          this.$success('复制成功')
        }
      } catch (error) {
        this.$error('复制失败')
        console.error('复制失败:', error)
      }
    },

    /**
     * 处理引用点击事件
     * @param {Event} e - 点击事件
     * @param {Object} options - 配置选项
     * @param {number} options.sessionStatus - 会话状态
     * @param {Object} options.sessionData - 会话数据
     * @param {string} options.citationSelector - 引用元素选择器，默认为 '.citation'
     * @param {string} options.subTagSelector - 子标签选择器，默认为 '.subTag'
     * @param {string} options.scrollElementId - 滚动容器ID，默认为 'timeScroll'
     * @param {Function} options.onToggleCollapse - 切换折叠状态的回调函数
     */
    $handleCitationClick(e, options = {}) {
      const {
        sessionStatus = 0,
        sessionData = null,
        citationSelector = '.citation',
        scrollElementId = 'timeScroll',
        onToggleCollapse = null
      } = options
      // 检查会话状态
      if (sessionStatus === 0) return

      // 查找最近的引用元素
      const citationElement = e.target.closest(citationSelector)
      if (!citationElement) return

      // 获取标签索引
      const tagIndex = parseInt(citationElement.textContent, 10)
      if (isNaN(tagIndex) || tagIndex <= 0) return

      // 获取父级索引和折叠状态
      const parentsIndexAttr = citationElement.getAttribute('data-parents-index')
      const parentsIndex = parentsIndexAttr ? parseInt(parentsIndexAttr, 10) : null
      // 检查 parentsIndex 是否有效
      if (isNaN(parentsIndex)) return
      
      // 检查会话数据结构
      if (!sessionData || 
          !sessionData.history || 
          !sessionData.history[parentsIndex] || 
          !sessionData.history[parentsIndex].searchList || 
          !sessionData.history[parentsIndex].searchList[tagIndex - 1]
        ) {
        return
      }
      // 切换折叠状态 - 严格按照组件中的collapseClick方法逻辑
      const searchItem = sessionData.history[parentsIndex].searchList[tagIndex - 1]
      const currentCollapse = searchItem.collapse
      const newCollapse = !currentCollapse
      if (onToggleCollapse && typeof onToggleCollapse === 'function') {
        onToggleCollapse(searchItem, newCollapse)
      } else {
        const updatedItem = { ...searchItem, collapse: newCollapse }
        if (this.$set) {
          this.$set(sessionData.history[parentsIndex].searchList, tagIndex - 1, updatedItem)
        } else {
          sessionData.history[parentsIndex].searchList[tagIndex - 1] = updatedItem
        }
      }

      // 滚动到底部
      const timeScrollElement = document.getElementById(scrollElementId)
      if (timeScrollElement) {
        timeScrollElement.scrollTop = timeScrollElement.scrollHeight
      }

      // 阻止事件冒泡
      e.stopPropagation()
    }
  },

  computed: {
    /**
     * 是否为空对象
     * @returns {Function} - 判断函数
     */
    $isEmpty() {
      return (obj) => {
        if (obj === null || obj === undefined) return true
        if (typeof obj === 'string') return obj.trim() === ''
        if (Array.isArray(obj)) return obj.length === 0
        if (typeof obj === 'object') return Object.keys(obj).length === 0
        return false
      }
    }
  },

  mounted() {

  },

  beforeDestroy() {

  }
}
