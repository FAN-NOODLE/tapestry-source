/*
 * @Author: your name
 * @Date: 2020-10-13 17:46:25
 * @LastEditTime: 2020-11-06 14:30:10
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /crs-admin-base/src/store/modules/user.ts
 */
import {
  VuexModule,
  Module,
  Action,
  Mutation,
  getModule
} from "vuex-module-decorators";
import store from "@/store";

interface IQueueItem {
  /** 任务ID */
  id: number;
  /** 任务标题 */
  title: string;
  /** 任务描述 */
  description: string;
  /** 任务进度 */
  percent: number;
  /** 任务状态 */
  state: 'waiting' | 'running' | 'success' | 'fail' | 'aborted';
  /** 任务上下文 */
  run: Function;
  /** 时间戳 */
  timestamp: number;
}

interface IStatesInterface {
  queue: IQueueItem[];
}
@Module({ name: "downloadTask", dynamic: true, store })
class DownloadTask extends VuexModule implements IStatesInterface {
  public queue: IQueueItem[] = [];

  get allTask() {
    return this.queue;
  }

  get total() {
    return this.queue.length;
  }

  get getCurrentTask() {
    return this.queue.find((item) => item.state === 'waiting');
  }

  get hasRunningTask () {
    return this.queue.some((item) => item.state === 'running');
  }

  get progress() {
    const current = this.queue.find((item) => item.state === 'running')
    
    const success = this.queue.filter((item) => item.state === 'success').length;
    const fail = this.queue.filter((item) => item.state === 'fail').length;
    if (this.total === 0) {
      return {
        success: 0,
        total: 0,
        percent: 0,
      };
    }

    let percent = success / this.total

    if(current) {
      percent += current.percent/this.total;
    }

    return {
      fail: fail,
      success: success,
      total: this.total,
      percent: parseFloat((percent * 100).toFixed(1)),
    };
  }


  @Mutation
  private ADD_TASK(data: IQueueItem) {
    this.queue.push(data);
  }
  @Mutation
  private DELETE_TASK(index: number) {
    this.queue.splice(index, 1);
  }
  @Mutation
  private CLEAR_TASK() {
    this.queue = []
  }

  @Action({ rawError: true })
  public addTask(params: IQueueItem) {
    this.context.commit("ADD_TASK", params);
  }

  @Action({ rawError: true })
  public deleteTask(index: number) {
    this.context.commit("DELETE_TASK", index);
  }
  @Action({ rawError: true })
  public clearAll() {
    this.context.commit("CLEAR_TASK");
  }
}

export default getModule(DownloadTask);

export class Task implements IQueueItem {
  id: number;
  title: string;
  description: string;
  percent: number;
  state: "waiting" | "running" | "success" | "fail" | "aborted";
  run: Function;
  timestamp: number;

  constructor(options: Partial<IQueueItem>) {
    this.id = Date.now();
    this.title = options.title || '';
    this.description = options.description || '';
    this.state = options.state || 'waiting';
    this.percent = 0;
    this.timestamp = parseInt(Date.now()/1000+'');
  }

  onRun(callback: Function) {
    this.run = callback.bind(this);
    return this
  }
}