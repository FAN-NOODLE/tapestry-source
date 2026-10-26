<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" @change="handleChange"></edit-form>
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
import EditForm, { IFormInfo } from '../components/edit-form.vue'
import Api from '@/api/wisdom/communication';
@Component({
  name: 'IntegrationAdd',
  components: {
    EditForm
  }
})
export default class IntegrationAdd extends Vue {
  formInfo: IFormInfo = null
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
        this.$message.success('新建成功')
        this.$router.back()
      } else {
        this.$message.error(res.message);
      }
    })
  }

  handleChange(value: IFormInfo) {
    console.log(11111, value);
    
    this.formInfo = value
  }
  /**
   * 预览
   */
  handlePreview() {
    (this.$refs.editForm as any).handlePreview()
  }
}
</script>