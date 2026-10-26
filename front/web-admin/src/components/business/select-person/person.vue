<template functional>
  <div :class="[
    'person',
    !props.data.mem_subscribe_status && !props.data.is_dp && 'disabled-person'
  ]">
    <!-- 用户icon -->
    <i
      class="el-icon-user-solid icon-color"
      v-if="!props.data.is_dp"
    ></i>

    <!-- 部门icon -->
    <img
      class="group-img"
      v-else
      src="https://cdn-cos.coach.com.cn/uploads/2023/web-admin/images/group-icon.png"
      alt=""
    />

    <!-- 员工名称/部门名称 -->
    <span class="name text-color">{{ props.data.name }}</span>

    <!-- 员工id -->
    <span
      v-if="!props.data.is_dp"
      class="user-id"
    >({{ props.data.user_id }})</span>

    <!-- 选中icon, props.node.checked 判断是否选中 -->
    <i
      class="el-icon-check icon-color check-margin"
      v-show="props.node ? props.node.checked : false"
    ></i>

    <!-- 右侧删除图标 -->
    <i
      class="el-icon-circle-close del-icon"
      v-if="props.del || false"
      @click="listeners.handleClick"
    ></i>
  </div>
</template>
<style scoped lang="scss">
.person {
  line-height: 30px;
  position: relative;
  padding-right: 24px;

  .name {
    display: inline;
    padding-left: 5px;
  }

  .user-id {
    color: $color-text-second;
    display: inline;
    padding-left: 3px;
  }

  .icon-color {
    color: $color-main;
  }

  .text-color {
    color: $color-text-normal;
  }

  .group-img {
    width: 14px;
    height: 14px;
    position: relative;
    top: 3px;
  }

  .check-margin {
    display: inline-block;
    margin-left: 10px;
  }

  .del-icon {
    color: $color-main;
    cursor: pointer;
    position: absolute;
    top: 50%;
    right: 6px;
    transform: translateY(-50%);
  }
}

.disabled-person {
  .icon-color {
    color: $color-disable;
  }

  .text-color {
    color: $color-text-normal;
  }
}
</style>
