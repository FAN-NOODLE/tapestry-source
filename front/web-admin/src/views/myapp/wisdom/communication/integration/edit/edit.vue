<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" :data="formInfo"></edit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
        <!-- <el-button type="primary" @click="handlePreview()">预览周传讯PDF</el-button> -->
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import EditForm, { IFormInfo, defaultFormInfo } from '../components/edit-form.vue'
import Api from '@/api/wisdom/communication';
@Component({
  name: 'IntegrationEdit',
  components: {
    EditForm
  }
})
export default class IntegrationEdit extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))

  created() {
    this.getDetail()
  }

  /**
   * 获取详情
   */
  getDetail() {
    Api.getMergeDetailApi({
      id: Number(this.$route.query.id)
    }).then((res: any) => {
      if (res.code == 200) {
        if(res.data.send_type === 0) {
          res.data.send_time = parseInt(res.data.send_time * 1000+'')
        }
        res.data.content = res.data.content_list
        res.data.crowd_list = (res.data.crowd_list||[]).map(item => {
          if(res.data.crowd_type === 1) {
            return {
              id: item.id,
              name: item.title
            }
          } else {
            return {
              member_id: item.id,
              name: item.title,
              dp_parent_id: item.department_id
            }
          }
        })

        res.data.week_num = res.data.week_num||''

        this.formInfo = res.data
      } else {
        this.$message.error(res.message);
      }
    })
  }

  handleCancel() {
    this.$router.back()
  }
  handleSave() {
    (this.$refs.editForm as any).validate(async(value: IFormInfo) => {
      const { crowd_list, ...params}: IFormInfo = JSON.parse(JSON.stringify(value));
      
      // 处理关联内容ID
      params.content = params.content.map(item => {
        const {info_list, ...other} = item
        return {
          ...other,
          info_ids: info_list.map(item => item.id)
        }
      })
      // 获取通知人群ID
      const crowd_ids = crowd_list.map(item => {
        if(params.crowd_type === 1) {
          return item.id
        } else {
          return item.member_id
        }
      })
      // 发布时间
      if(params.send_type === 0) {
        params.send_time = parseInt(params.send_time / 1000+'')
      }

      params.week_num = params.week_num || 0
      
      const res: any = await Api.saveMergeApi({...params, crowd_ids})
      if (res.code == 200) {
        this.$message.success('保存成功')
        this.$router.back()
      } else {
        this.$message.error(res.message);
      }
    })
  }
  /**
   * 预览
   */
  handlePreview() {
    (this.$refs.editForm as any).handlePreview()
  }
}
</script>