<template>
  <div class="counter">
    <el-dialog
      :close-on-click-modal="false"
      title="社群设置"
      width="800px"
      :visible.sync="visibleDialog"
      @close="closeEditappDialog()"
    >
      <cy-form
        ref="cform"
        :rules="rules"
        label-width="120px"
      >
        <cy-form-item label="所属渠道：">
          <cy-select
            v-model="staff"
            placeholder="全部"
            search-input
            :optionsList="channelList"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
          <el-button
            type="text"
            class="m-l-20"
            @click="newChannel"
          >新建渠道</el-button>
        </cy-form-item>
        <cy-form-item label="群类型：">
          <cy-select
            v-model="groupType"
            placeholder="全部"
            search-input
            :optionsList="groupTypeList"
            :options="{ label: 'name', value: 'id' }"
          >
          </cy-select>
        </cy-form-item>
        <cy-form-item label="群运营者：">
          <el-button
            type="primary"
            @click="selectStaff = true"
          >选择人员</el-button>
          <div
            class="table-box m-t-15"
            v-if="tableData[0].stat_txt_obj.length > 0"
          >
            <!-- stripe  是否为斑马纹 table -->
            <el-table
              :data="tableData"
              style="width: 100%"
            >
              <el-table-column
                label="员工"
                prop="stat_txt_obj"
              >
                <template slot-scope="{ row }">
                  <el-tag
                    v-for="(item, index) in row.stat_txt_obj"
                    :key="index"
                    class="el-staff-tag"
                  >
                    <img
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/staff-icon.png"
                      alt=""
                      class="table-text-img"
                    />
                    <span>{{ item.name }}</span>
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column
                label="操作"
                width="110"
                fixed="right"
              >
                <template>
                  <el-button
                    type="text"
                    size="small"
                    @click.native.prevent="selectStaff = true"
                  >
                    修改
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </cy-form-item>
        <cy-form-item label="群目标：">
          <el-input
            type="textarea"
            :rows="6"
            :maxlength="50"
            resize="none"
            show-word-limit
            placeholder="请输入50字以内群目标"
            v-model="marks"
          >
          </el-input>
        </cy-form-item>
      </cy-form>
      <span
        slot="footer"
        class="dialog-footer"
      >
        <el-button
          size="small"
          @click="cancel()"
        >取 消</el-button>
        <el-button
          size="small"
          type="primary"
          @click="save()"
        >保 存</el-button>
      </span>
    </el-dialog>
    <selectMember
      v-if="selectStaff"
      :limiteNum="1"
      :hasSelectArr="hasSelectArr"
      @close="selectStaff = false"
      @submit="getMember"
    />
  </div>
</template>

<script lang="ts">
import { Vue, Component, Prop } from "vue-property-decorator";
import {
  chatSettingApi,
  chatInfoApi,
  chatAllChannelApi,
  chatTypeList
} from "@/api/enterprise/association";
import selectMember from "@/components/business/select-member/index.vue";

@Component({
  components: { selectMember }
})
export default class extends Vue {
  @Prop({}) chatId!: string;
  private tableData = [
    {
      stat_txt_obj: []
    }
  ];
  private marks = "";
  // 控制弹框的关闭隐藏
  private selectStaff = false;
  private visibleDialog = true;
  private rules = {};
  private staff = 0;
  private channelList = [
    {
      name: "全部",
      id: 0
    }
  ];
  private groupTypeList: any = [
    {
      name: "全部",
      id: 0
    }
  ];
  private groupType = 0;
  // 发送给选人弹框
  private hasSelectArr: any = [];

  // ----------------------------------------------methods
  // 新建渠道
  private newChannel() {
    const { href } = this.$router.resolve({
      name: "myapp_enterprise_drainage_channel_list",
      query: {
        showAddModal: "1"
      }
    });
    window.open(href, "_blank");
  }
  // 获取渠道列表
  private getChannel() {
    chatAllChannelApi().then(res => {
      console.log(res);
      this.channelList = [{ name: "全部", id: 0 }, ...res.data];
    });
  }
  // 获取社群类型列表
  private getChatType() {
    chatTypeList().then(res => {
      console.log(res);
      this.groupTypeList = [{ name: "全部", id: 0 }, ...res.data];
    });
  }
  // 获取群信息
  private groupDetail(chatId: string) {
    chatInfoApi({
      chat_id: chatId
    }).then(res => {
      console.log(res);
      this.staff = res.data.channel_id;
      this.groupType = res.data.chat_type_id;
      this.marks = res.data.chat_target;
      if (res.data.operator_member_id && res.data.operator_member_id != "0") {
        this.hasSelectArr = [
          { user_id: res.data.operator_member_id, name: res.data.operator }
        ];
        this.tableData[0].stat_txt_obj = this.hasSelectArr;
      }
    });
  }
  // 取消，关闭弹框
  cancel() {
    this.$emit("close");
  }
  // 关闭
  closeEditappDialog() {
    this.$emit("close");
  }
  // 保存
  private save() {
    chatSettingApi({
      chat_id: this.chatId,
      channel_id: this.staff,
      chat_type_id: this.groupType,
      operator_id: (this.hasSelectArr[0] && this.hasSelectArr[0].user_id) || "",
      chat_target: this.marks
    }).then(res => {
      console.log(res);
      this.$emit("close");
    });
  }
  private getMember(list: any) {
    console.log(list);
    this.hasSelectArr = list;
    this.tableData[0].stat_txt_obj = list;
    this.selectStaff = false;
  }
  // 生命周期
  private mounted() {
    console.log("mounted");
    this.getChannel();
    this.getChatType();
    this.groupDetail(this.chatId);
  }
}
</script>
<style lang="scss" scoped>
::v-deep {
  .in-item:last-child {
    margin: 0;
  }

  .el-dialog__body {
    padding-bottom: 15px !important;
  }
}

::v-deep.el-form-item__label {
  width: 100px;
  height: 10px;
}

::v-deep.el-dialog__body {
  padding-bottom: 0;
}

::v-deep {
  .el-upload {
    width: 300px;
    height: 180px;
  }

  .el-upload-dragger {
    @extend .el-upload;
    border: 1px solid $color-second-line !important;
    border-radius: 4px !important;
  }

  .el-upload__tip {
    color: $color-text-second !important;
  }

  .el-upload-dragger {
    .el-icon-upload {
      font-size: 80px !important;
      margin-top: 30px !important;
      margin-bottom: 20px !important;
    }
  }

  .el-dialog {
    max-height: none !important;
  }
}

.tips {
  font-size: 13px;
  color: #909399;
  margin-bottom: 10px;
}

.table-box {
  /deep/.el-table__row {
    cursor: pointer;
  }

  margin-top: 15px;

  .el-staff-tag {
    margin-right: 10px;
    margin-bottom: 10px;
  }

  .table-text-img {
    display: inline-block;
    width: 14px;
    height: 14px;
    vertical-align: middle;
    margin-right: 5px;
  }
}</style>
