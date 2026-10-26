<template>
  <div v-if="total>0" class="download-task">
    <el-progress type="circle" :percentage="progress.percent" :stroke-width="2" :width="20"></el-progress>
    <div class="download-task__text" @click="show=true">{{progress.total}}</div>
    <float-window v-if="show">
      <div class="task">
        <div class="task-header">
          <span>任务进行中</span>
          <span class="tips">（请勿刷新或关闭页面）</span>
          <el-tooltip content="最小化">
            <div class="clear" @click="show=false"></div>
          </el-tooltip>
        </div>
        <div class="task-body">
          <div v-if="list.length" class="task-list">
            <div class="task-item" v-for="(item, index) in list" :key="index">
              <div class="task-item__info">
                <div class="task-item__name">
                  <span class="">{{ item.title }}</span>
                  <span v-if="item.percent < 1" class="m-r-5">（{{(item.percent*100).toFixed(1)}}%）</span>
                  <span v-else class="m-r-5">（已完成）</span>
                  <el-button class="delete" type="text" icon="el-icon-delete" @click="handleDelete(index)"></el-button>
                </div>
                <div class="tips m-t-5">
                  <span>{{ (item.timestamp) | formatDate("yyyy/mm/dd hh:mm:ss") }}</span>
                  <span v-if="item.state === 'success'" class="m-l-5">(请在浏览器下载中心查看文件)</span>
                </div>
              </div>
              <el-progress
                class="m-t-5"
                :percentage="item.percent * 100"
                :show-text="false"
                :stroke-width="10"
              />
            </div>
          </div>
          <el-empty v-else description="暂无任务" />
        </div>
      </div>
    </float-window>
  </div>
</template>
<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator';
import DownloadModule from "@/store/modules/download";
import { AbortAblePromise } from '@/utils/promise';
import FloatWindow from './FloatWindow.vue'

@Component({
  components: {
    FloatWindow
  }
})
export default class DownloadTask extends Vue {
  [x: string]: any;
  get total() {
    return DownloadModule.total;
  }
  get list() {
    return DownloadModule.allTask;
  }

  get progress() {
    return DownloadModule.progress;
  }

  show = false

  @Watch('total')
  onTotalChange() {
    this.run()
    if(this.total > 0) {
      this.show = true;
    } else {
      this.show = false;
    }
  }

  async run() {
    const hasRunningTask = DownloadModule.hasRunningTask;
    if (hasRunningTask) {
      return;
    }
    const task = DownloadModule.getCurrentTask;
    
    if (task) {
      task.state = 'running';
      console.log('执行中', task?.id);
      const currentTask = new AbortAblePromise(task.run())
      currentTask.then(() => {
          task.state = 'success';
          console.log('执行完成', task?.id);
          // 继续执行下一个任务
          this.run();
        })
        .catch((error) => {
          if (error.message === 'Promise aborted') {
            task.state = 'aborted';
            console.log('任务取消', task?.id);
          } else {
            task.state = 'fail';
            console.log('执行失败', task?.id);
          }
        });
    } else {
      console.log('没有任务了');
    }
  }

  handleDelete(index: number) {
    DownloadModule.deleteTask(index)
  }
}
</script>
<style scoped lang="scss">
::v-deep .el-progress__text {
  display: none;
}

.download-task {
  position: relative;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;


  .download-task__text {
    width: 100%;
    height: 100%;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1;
  }
}

.task {
  width: 400px;
  border-radius: 5px;
  background-color: #fff;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);

  .task-header {
    position: relative;
    height: 40px;
    padding: 0 10px;
    border-bottom: 1px solid #ebebeb;
    line-height: 40px;

    .clear {
      position: absolute;
      top: 50%;
      right: 10px;
      transform: translateY(-50%);
      cursor: pointer;
      width: 40px;
      height: 40px;

      &::after {
        content: "";
        display: block;
        width: 15px;
        height: 3px;
        background-color: #000;
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        border-radius: 3px;
      }
    }
  }

  .task-body {
    padding: 0 15px;
    line-height: normal;
  }
}

.task-item {
  padding: 15px 0;
  &:not(:last-child) {
    border-bottom: 1px solid #ebebeb;
  }
  
  &__name {
    display: flex;
    align-items: center;
    position: relative;
  }

}
.tips {
  font-size:12px;
  color: #999;
}
</style>