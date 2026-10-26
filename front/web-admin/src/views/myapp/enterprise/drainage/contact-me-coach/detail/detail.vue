<template>
  <div class="supernatant">
    <layout-list>
      <template slot="search">
        <el-row :gutter="20">
          <el-col :span="6">
            <div class="item">活动名称：{{ detail.name }}</div>
          </el-col>
          <el-col :span="6">
            <div class="item">渠道：{{detail.channel_name}}</div>
          </el-col>
          <el-col :span="6">
            <div class="item">渠道标签：{{ detail.tags ? detail.tags.join('，') : '-' }}</div>
          </el-col>
          <el-col :span="6">
            <div class="item">二维码类型：{{detail.type===1 ? '员工码+门店码' : '多人共享码'}}</div>
          </el-col>
        </el-row>
      </template>
      <template slot="table">
        <el-table
          :data="contactList"
          stripe
          style="width: 100%"
          tooltip-effect="light"
        >
          <el-table-column
            label="二维码名称"
            align="center"
            show-overflow-tooltip
            prop="name"
            width="150"
          ></el-table-column>
          <el-table-column
            label="二维码"
            width="150"
            align="center"
          >
            <template slot-scope="scope">
              <el-image
                class="code-img"
                :src="scope.row.qr_code_url"
                :preview-src-list="[scope.row.qr_code_url]"
              />
            </template>
          </el-table-column>
          <el-table-column
            prop="members"
            label="部门/员工"
            min-width="320"
            align="center"
          >
              <template slot-scope="scope">
              <!-- <el-tag v-for="(item, index) in scope.row.category_names" :key="index">{{item}}</el-tag> -->
              <el-tooltip effect="light" v-if="memberList(scope.row).length > 3 || scope.row.departments.length > 0">
                <div class="tags" style="cursor: pointer;">
                  <!-- 显示前2个 -->
                  <el-tag class="tag" v-for="item,index in memberList(scope.row).slice(0,3)" :key="index">
                    <!-- icon -> 部门 -->
                    <img
                      class="staff-icon"
                      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                      alt=""
                      v-if="item.is_dp"
                    />
                    <i v-else class="el-icon-user-solid"></i>
                    {{item.name}}
                  </el-tag>
                  <span>...</span>
                </div>
                <div class="" slot="content">
                  <div class="dp-item" v-for="item,index in memberList(scope.row, true)" :key="index">
                    <div class="dp-name">
                      <img
                        class="staff-icon"
                        src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                        alt=""
                        v-if="item.name"
                      />
                        {{item.name || '其他成员'}}</div>
                    <div class="tags tooltip-content">
                      <el-tag v-for="item,index in item.members" :key="index">
                        <img
                          class="staff-icon"
                          src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                          alt=""
                          v-if="item.is_dp"
                        />
                        <i v-else class="el-icon-user-solid"></i>
                        {{item.member_name}}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </el-tooltip>
              <div v-else-if="memberList(scope.row).length <= 3 && memberList(scope.row).length > 0" class="tags">
                <el-tag v-for="item,index in memberList(scope.row)" :key="index">
                  <!-- icon -> 部门 -->
                  <img
                    class="staff-icon"
                    src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
                    alt=""
                    v-if="item.is_dp"
                  />
                  <i v-else class="el-icon-user-solid"></i>
                  {{item.name}}</el-tag>
              </div>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            prop="add_customer_number"
            label="新增好友数"
            min-width="150"
            align="center"
          ></el-table-column>
          <el-table-column
            label="创建时间"
            min-width="200"
            prop="created_at"
            align="center"
          >
            <template slot-scope="scope">
              <span v-if="scope.row.created_at > 0">{{
                scope.row.created_at | formatDate("yyyy-mm-dd hh:mm:ss")
              }}</span>
              <span v-else>--</span>
            </template>
          </el-table-column>
          <el-table-column
            label="操作"
            width="150"
            fixed="right"
          >
            <template slot-scope="scope">
              <el-button
                type="text"
                @click="downLoad(scope.row)"
              >
                下载
              </el-button>
              <el-button
                type="text"
                @click="deleteContact(scope.row.id)"
              >删除</el-button>
            </template>
          </el-table-column>
          <div
            slot="empty"
            class="table-empty-label"
          >
            <img
              class="empty-img"
              :src="searchLoading
                  ? 'https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/empty-search.png'
                  : require('@/assets/images/empty-box.png')
                "
              alt=""
            />
            <div>{{ searchLoading ? "没有找到相关联系我" : "暂无联系我" }}</div>
          </div>
        </el-table>
      </template>
      <template
        slot="footer"
        v-if="meta.total > 0"
      >
        <el-pagination
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          :current-page="meta.page"
          :page-sizes="[10, 20, 30, 40]"
          :page-size="meta.size"
          layout="total, sizes, prev, pager, next, jumper"
          :total="meta.total"
        >
        </el-pagination>
      </template>
    </layout-list>
    <!-- 下载二维码弹框-->
    <download-code
      v-if="showDownLoad"
      :currentDownLoadItem="currentDownLoadItem"
      @closeDialog="showDownLoad = false"
    ></download-code>
  </div>
</template>
<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Api from "@/api/enterprise/drainage-coach";
import downloadCode from "../../components/download-code.vue";
@Component({
  components: {
    downloadCode
  }
})
export default class extends Vue {
  /*************************************** data ****************************************/
  private contactList: Array<any> = []; // 联系我列表数组
  detail: any = {};
  private meta = {
    total: 0, // 总数
    page: 1, // 当前页
    size: 10 // 每页数量
  };
  private searchLoading = false; // 是否正在搜索
  private showDownLoad = false; // 控制下载二维码的弹框显隐
  private currentDownLoadItem: any; // 当前下载二维码的url
  
  get memberList() {
    return (item: any, showLong = false) => {
      const memberList = item.members.map((v: any) => ({name: v, is_dp: false}))
      const departmentList = item.departments.map((v: any) => {

        if(showLong) {
          const relMembers = item.real_members.filter(m => m.third_id === v.third_id).map(m => m.member_name).join('，')
          
          return { 
            name: v.department_name,
            // name: relMembers ? v.department_name + '(' + relMembers + ')' : v.department_name, 
            is_dp: true,
            members: item.real_members.filter(m => m.third_id === v.third_id)
          }
        }

        return { name: v.department_name, is_dp: true, members: [] }
      })
      const dpIds = item.departments.map((v: any) => v.third_id)
      const otherMembers = item.real_members.filter(m => !dpIds.includes(m.third_id))

      if(showLong && otherMembers.length > 0) {
        departmentList.push({
          name: '',
          is_dp: false,
          members: otherMembers
        })
        return [...departmentList];
      }
      
      return [...departmentList, ...memberList];
    };
  }
  /*************************************** method ****************************************/
  // 联系我列表接口
  private getContactList() {
    const params = {
      id:Number(this.$route.query.id),
      page: this.meta.page,
      size: this.meta.size
    };
    Api.contactDetailListApi(params).then((res: any) => {
      if (res.code == 200) {
        const { contact_list, total, ...detail } = res.data
        this.contactList = contact_list;
        this.detail = detail;
        this.meta.total = total;
      } else {
        this.$message.error(res.message);
      }
    });
  }
  
  // 内部绑定函数（获取每页数据量）
  private handleSizeChange(val: number) {
    this.meta.page = 1;
    this.meta.size = val;
    this.getContactList();
  }
  // 内部绑定函数（获取当前第几页）
  private handleCurrentChange(val: number) {
    this.meta.page = val;
    this.getContactList();
  }
  // 点击下载按钮
  private downLoad(item: any) {
    this.currentDownLoadItem = item;
    this.showDownLoad = true;
  }
  // 删除联系我接口
  private deleteContact(id: number) {
    this.$confirm("确认删除该二维码，删除后该二维码将失效？", "删除提醒", {
      confirmButtonText: "确定删除",
      cancelButtonText: "我再想想",
      type: "warning"
    })
      .then(() => {
        Api.contactDeleteContactApi({
          ids: [id]
        }).then((res: any) => {
          if (res.code === 200) {
            this.$message({
              type: "success",
              message: "删除成功!"
            });
            this.getContactList()
          } else {
            this.$message.error(res.message);
          }
        });
      })
      .catch(() => {
        this.$message({
          type: "info",
          message: "已取消删除"
        });
      });
  }
  private mounted() {
    this.getContactList();
  }
}
</script>
<style lang="scss" scoped>

.item {
  line-height: 20px;
  padding: 20px 0;
}

.code-img {
  width: 50px;
  height: 50px;
}

/deep/.el-dropdown-menu__item.color-red {
  color: $color-danger;
}

/deep/.el-dropdown-menu__item.color-edit {
  color: $color-text-button;
}

.table-empty-label {
  display: flex;
  height: 340px;
  flex-direction: column;

  img {
    width: 124px;
    margin: 80px auto 0;
  }

  div {
    line-height: 20px;
    margin-top: 14px;
  }
}

.dp-item {
  margin-bottom: 10px;
  .dp-name {
    line-height: 34px;
    font-weight: bold;
    display: flex;
    align-items: center;

    .staff-icon {
      display: block;
    }
  }
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  justify-content: center;

  &.tooltip-content {
    max-width: 600px;
    justify-content: flex-start;
  }
}

.tag {
  max-width: 100px;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.staff-icon {
  width: 12px;
  height: 12px;
  margin-right: 5px;
  vertical-align: middle;
}
</style>

<style lang="scss"> 
</style>