<template>
  <layout-form>
    <template slot="content">
      <edit-form ref="editForm" :data="formInfo"></edit-form>
    </template>
    <template slot="footer">
      <div>
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </div>
    </template>
  </layout-form>
</template>
<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import EditForm, { IFormInfo, defaultFormInfo } from '../components/edit-form.vue'
import Api from '@/api/wisdom/communication';
import { getActivityDetailApi, saveActivityApi } from '@/api/enterprise/activity';
@Component({
  components: {
    EditForm
  }
})
export default class LaunchEdit extends Vue {
  /** 表单数据 */
  formInfo: IFormInfo = JSON.parse(JSON.stringify(defaultFormInfo))

  mounted() {
    this.getDetail()
  }

  /**
   * 获取详情
   */
  getDetail() {
    getActivityDetailApi({
      id: Number(this.$route.query.id)
    }).then((res: any) => {
      if (res.code == 200) {
        res.data.skip_verify = res.data.skip_verify == 1
        res.data.add_name_tag = res.data.add_name_tag == 1
        res.data.add_category_tag = res.data.add_category_tag == 1
        res.data.scope_area_ids = res.data.scope_area_list.map(item => item.id)
        res.data.scope_member_ids = res.data.scope_member_list.map(item => item.id)
        res.data.scope_member_list = res.data.scope_member_list.map(item => {
          return {
            type: 5,
            ...item
          }
        })
        res.data.tag_ids = res.data.tag_list.map(item => item.id)
        res.data.category_list = res.data.category_list.map(item => {
          return {
            ...item,
            points: item.point_list
          }
        })

        res.data.has_welcome = res.data.welcome_list.length > 0

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
    (this.$refs.editForm as any).validate(async (value: IFormInfo) => {
      const params = JSON.parse(JSON.stringify(value));
      const now = new Date().getTime()/1000
      
      if(params.time_type === 0 || (params.time_type === 1 && params.start_time < now &&  now < params.end_time)) {
        await this.$confirm('当前活动已处于进行中，编辑可能会影响投放活动数据统计，确定保存当前编辑内容?', '编辑提醒', {
          confirmButtonText: '确定',
          cancelButtonText: '我再想想',
          type: 'warning'
        })
      }
      
      saveActivityApi(params)
        .then((res: any) => {
          if (res.code == 200) {
            this.$message.success('保存成功')
            this.$router.back()
          } else {
            this.$message.error(res.message);
          }
        })
    })
  }
}
</script>