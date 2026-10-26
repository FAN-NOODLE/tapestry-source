import Vue from "vue";
import {
  Pagination, // 分页
  Dialog, // 模态框
  Autocomplete,
  Dropdown, // 下拉菜单
  DropdownMenu, // 下拉菜单
  DropdownItem, // 下拉菜单项
  Menu, // 菜单
  Submenu, // 菜单
  MenuItem, // 菜单项
  MenuItemGroup, // 菜单项组
  Input, // 输入框
  InputNumber, // number 输入框
  Radio, // 单选框
  RadioGroup, // 单选框组
  RadioButton, // 单选按钮
  Checkbox, // 多选框
  CheckboxButton, //多选框按钮
  CheckboxGroup, // 多选框组
  Switch, // 开关
  Select, // 选择器
  Option, // 选择器 item
  OptionGroup, // 选择器组
  Button, // 按钮
  ButtonGroup, // 按钮组
  Table, // 表格
  TableColumn, // 列
  DatePicker, // 日期选择器
  TimeSelect, // 时间选择器
  TimePicker, // 时间日期选择器
  Popover, // 弹出框
  Tooltip, // 文字提示
  Avatar,
  Form, // 表单
  FormItem, // 表单项
  Tabs, // 标签页
  TabPane, // 标签页
  Tag, // 标签
  Tree, // 树
  Alert, // 警告
  Icon, //图标
  Row, // 行
  Col, // 列
  Upload, // 上传
  Progress, // 进度条
  Card, // 卡片
  Steps, // 步骤
  Step, // 步骤项
  Carousel, // 走马灯
  CarouselItem, // 走马灯 item
  Collapse, // 折叠面板
  CollapseItem, // 折叠面板 item
  Cascader, // 级联选择器
  ColorPicker, // color select
  Container, // 容器
  Header, // 头
  Aside, // left
  Main, // main
  // Footer, // footer
  Loading, // 加载中
  MessageBox, // 提示
  Message, // 提示
  Image,
  Popconfirm,
  Breadcrumb, // 面包屑
  BreadcrumbItem, // 面包屑 item
  InfiniteScroll, // 无限滚动
  Scrollbar, // 区域滚动
  Link,
  Timeline,
  TimelineItem,
  Notification,
  Descriptions,
  DescriptionsItem,
  Divider
} from "element-ui";
import "../element-variables.scss";
import "element-ui/lib/theme-chalk/base.css";
import CollapseTransition from "element-ui/lib/transitions/collapse-transition";

Vue.use(Link);
Vue.use(Pagination);
Vue.use(Dialog);
Vue.use(Autocomplete);
Vue.use(Dropdown);
Vue.use(DropdownMenu);
Vue.use(DropdownItem);
Vue.use(Menu);
Vue.use(Submenu);
Vue.use(MenuItem);
Vue.use(MenuItemGroup);
Vue.use(Input);
Vue.use(InputNumber);
Vue.use(Radio);
Vue.use(RadioGroup);
Vue.use(RadioButton);
Vue.use(Checkbox);
Vue.use(CheckboxButton);
Vue.use(CheckboxGroup);
Vue.use(Switch);
Vue.use(Select);
Vue.use(Option);
Vue.use(OptionGroup);
Vue.use(Button);
Vue.use(ButtonGroup);
Vue.use(Table);
Vue.use(TableColumn);
Vue.use(DatePicker);
Vue.use(TimeSelect);
Vue.use(TimePicker);
Vue.use(Popover);
Vue.use(Tooltip);
Vue.use(Form);
Vue.use(FormItem);
Vue.use(Tabs);
Vue.use(TabPane);
Vue.use(Tag);
Vue.use(Tree);
Vue.use(Alert);
Vue.use(Icon);
Vue.use(Row);
Vue.use(Col);
Vue.use(Upload);
Vue.use(Progress);
Vue.use(Card);
Vue.use(Avatar);
Vue.use(Steps);
Vue.use(Step);
Vue.use(Carousel);
Vue.use(CarouselItem);
Vue.use(Collapse);
Vue.use(CollapseItem);
Vue.use(Cascader);
Vue.use(ColorPicker);
Vue.use(Container);
Vue.use(Header);
Vue.use(Aside);
Vue.use(Main);
Vue.use(Image);
Vue.use(Popconfirm);
Vue.use(Breadcrumb);
Vue.use(BreadcrumbItem);
Vue.use(Timeline);
Vue.use(TimelineItem);
Vue.use(Loading.directive);
Vue.use(InfiniteScroll);
Vue.use(CollapseTransition);
Vue.use(Scrollbar);
// Vue.use(Notification);
Vue.use(Descriptions);
Vue.use(DescriptionsItem);
Vue.use(Divider);
Vue.prototype.$loading = Loading.service;
Vue.prototype.$msgbox = MessageBox;
Vue.prototype.$alert = MessageBox.alert;
Vue.prototype.$confirm = MessageBox.confirm;
Vue.prototype.$prompt = MessageBox.prompt;
Vue.prototype.$notify = Notification;
Vue.prototype.$message = Message;
Vue.component(CollapseTransition.name, CollapseTransition);
