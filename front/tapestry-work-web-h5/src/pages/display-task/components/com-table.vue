<template>
  <div class="table">
    <div class="scrollbar__view">
      <table class="table__header">
        <colgroup>
          <col 
            v-for="(column, index) in columns" 
            :key="index"
            :name="className(index)"
            :width="column.width"
          />
        </colgroup>
        <thead>
          <tr class="table__row">
            <th 
              :class="[
                'table__col',
                className(index),
                {
                  'table__col--sort': column.sortable,
                  'table__col--sort--desc': column.sortable && column.sortType === 'desc',
                  'table__col--sort--asc': column.sortable && column.sortType === 'asc'
                },
                {
                  'align-center': column.align === 'center',
                  'align-left': column.align === 'left',
                  'align-right': column.align === 'right'
                }
              ]" 
              v-for="(column, index) in columns" 
              :key="index"
              :rowspan="1"
              :colspan="1"
              @click="handleHeaderColumnClick(column)"
              >
              <div class="cell">
                <span>{{column.label}}</span>
                <span class="sort"></span>
              </div>
            </th>
          </tr>
        </thead>
      </table>
    </div>
    <div v-if="data.length" class="scrollbar__view">
      <table class="table__body">
        <colgroup>
          <col 
            v-for="(column, index) in columns" 
            :key="index"
            :name="className(index)"
            :width="column.width"
          />
        </colgroup>
        <tbody>
          <tr class="table__row" 
            v-for="(item, index) in innerData" 
            :key="index"
            @click="handleRowClick(item, index, $event)"
          >
            <td 
              :class="[
                'table__col',
                className(index),
                {
                  'align-center': column.align === 'center',
                  'align-left': column.align === 'left',
                  'align-right': column.align === 'right'
                }
              ]" 
              v-for="(column, columnIndex) in columns" 
              :key="columnIndex"
              :rowspan="1"
              :colspan="1"
              @click="handleColClick(item, column, index, $event)"
            >
              <div class="cell">
                <slot :name="column.prop" :row="item" :index="index">{{item[column.prop]}}</slot>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <slot name="empty" v-else>
      <div class="empty">
        <img class="empty-image" src="@/assets/material-library/no-data.png" alt="">
        <p class="empty-text">暂无数据</p>
      </div>
    </slot>
  </div>
</template>
<script lang="ts">
import { Component, Prop, Vue, Watch } from 'vue-property-decorator'
@Component({})
export default class ComTable extends Vue {
  @Prop({ type: Array, default: () => [] }) columns: any[]
  @Prop({ type: Array, default: () => [] }) data: any[]

  innerData:any[] = []

  @Watch('data')
  onDataChange() {
    this.innerData = JSON.parse(JSON.stringify(this.data))
  }

  mounted() {
    this.innerData = JSON.parse(JSON.stringify(this.data))
  }

  get className() {
    return (index: number) => {
      return `el-table_${1}_column_${index}`
    }
  }

  handleRowClick(row: any, index: number, event: any) {
    this.$emit('row-click', row, index, event)
  }

  handleColClick(row: any, column: any, index: number, event: any) {
    this.$emit('col-click', row, column, index, event)
  }

  handleHeaderColumnClick(column: any) {
    // 排序
    if(column.sortable) {
      column.sortType = {
        '': 'desc',
        'desc': 'asc',
        'asc': ''
      }[column.sortType];

      this.handleSort(column.prop, column.sortType)
    }
  }

  handleSort(prop: string, type: string) {
    const data = JSON.parse(JSON.stringify(this.data))
    if(type) {
      data.sort((a: any, b: any) => {
        if (a[prop] < b[prop]) {
          return type === 'desc' ? 1 : -1;
        } else if (a[prop] > b[prop]) {
          return type === 'desc' ? -1 : 1;
        } else {
          return 0;
        }
      })
    }
    this.innerData = data
  }
}
</script>

<style lang="scss" scoped>
.table {
  width: 100%;
  font-size: 14px;
  .cell {
    font-size: 12px;
  }
  &__col {
    padding: 10px 0;
    box-sizing: border-box;
    text-overflow: ellipsis;
    vertical-align: middle;
    position: relative;
    text-align: left;
    font-weight: normal;
    .cell {
      padding: 0 15px;
    }

    &.align-center {
      text-align: center;
    }

    .sort {
      padding: 7px;
      display: none;
      background: url('~@/assets/display-task/icon_sort.png') center/100% no-repeat;
    }

    &--sort {
      .sort {
        display: inline;
      }
    }

    &--sort--desc {
      .sort {
        background: url('~@/assets/display-task/icon_sort_desc.png') center/100% no-repeat;
      }
    }
    &--sort--asc {
      .sort {
        background: url('~@/assets/display-task/icon_sort_asc.png') center/100% no-repeat;
      }
    }
  }
  &__header {
    width: 100%;
    text-align: left;
  }

  &__body {
    width: 100%;
  }
}
.empty {
  padding: 50px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  &-image {
    width: 100px;
    height: 100px;
  }
  &-text {
    margin-top: 5px;
    text-align: center;
    // color: ;
  }
}
</style>