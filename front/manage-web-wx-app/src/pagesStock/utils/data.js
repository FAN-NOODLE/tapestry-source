export const skuItemHeader = [{
  prop: 'image_url',
  label: '图片',
  type: 'image',
  width: 110,
  fixed: 'left'
}, {
  prop: 'upc_code',
  label: '条形码',
  type: 'link',
  width: 110
}, {
  prop: 'color',
  label: '颜色',
  width: 140
}, {
  prop: 'style',
  label: '尺寸',
  width: 80
},
{
  prop: 'sku_tag_price',
  label: '吊牌价',
  width: 140
},
{
  prop: 'store_price',
  label: '零售价',
  width: 140
},
{
  prop: 'on_hand',
  label: '净库存',
  sort: true,
  sort_type: 'asc',
  width: 140
},
{
  prop: 'in_transit',
  label: '在途库存',
  sort_type: '',
  sort: true,
  width: 140
}, {
  prop: 'dcl_on_hand',
  label: '大仓库存',
  width: 140,
  sort_type: '',
  sort: true
}, {
  prop: 'search',
  label: '跨店库存',
  type: 'button',
  width: 120,
  fixed: 'right'
}]
export const upcCodeHeader = [{
  prop: 'image_url',
  label: '图片',
  type: 'image',
  width: 120,
  fixed: 'left'
}, {
  prop: 'upc_code',
  label: '条形码',
  type: 'link',
  width: 220
},
{
  prop: 'sku_tag_price',
  label: '吊牌价',
  width: 140
},
{
  prop: 'store_price',
  label: '零售价',
  width: 140
}, {
  prop: 'on_hand',
  label: '净库存',
  sort: true,
  sort_type: 'asc',
  width: 140
}, {
  prop: 'in_transit',
  label: '在途库存',
  sort: true,
  sort_type: '',
  width: 140
}, {
  prop: 'dcl_on_hand',
  label: '大仓库存',
  width: 140,
  sort_type: '',
  sort: true
}, {
  prop: 'search',
  label: '跨店库存',
  type: 'button',
  width: 140,
  fixed: 'right'
}]
export const otherHeader = [
  {
    prop: 'city',
    label: '城市',
    width: 150
  },
  {
    prop: 'loc_code',
    label: '店铺代码',
    width: 120
  },
  {
    prop: 'inventory_name',
    label: '店铺名称',
    width: 120,
    showTooltip: true
  },
  {
    prop: 'on_hand',
    label: '净库存',
    width: 140,
    sort: true,
    sort_type: 'asc'
  },
  {
    prop: 'in_transit',
    label: '在途库存',
    sort: true,
    width: 140,
    sort_type: ''
  },
  // {
  //   prop: 'dcl_on_hand',
  //   label: '大仓库存',
  //   width: 140
  // },
  {
    prop: 'phone',
    label: '电话',
    width: 140,
    type: 'icon'
  }
]
