# crs-admin

### 项目目录结构
├─ account ------------------- 本地运行脚本     
│   
├─ create ------------------- 自动创建模块脚本
│  
├─ deploy ------------------- Docker部署      
│   
├─ public ------------------- 创建页面执行文件   
│  ├─ favicon.ico ----------- 标题icon   
│  └─ index.html ------------ 静态页面  
│  
├─ dist --------------------- 打包资源   
│  
├─ src ---------------------- 目录结构             
│  ├─ assets ---------------- 静态图片                
│  ├─ api ------------------- 接口文件               
│  ├─ components ------------ 公共组件          
│  ├─ directives ------------ 公共指令        
│  ├─ filters --------------- 公共过滤器  
│  ├─ layout ---------------- 全局布局            
│  ├─ icon ------------------ svg图标                
│  ├─ mixins ---------------- 全局mixins                      
│  ├─ router ---------------- 路由配置       
│  ├─ store ----------------- vuex目录                  
│  │  └─modules ------------- store                     
│  ├─ style ----------------- 公共样式                  
│  │  ├─base ---------------- 基本样式                    
│  │  ├─color --------------- 颜色相关  
│  │  └─cover --------------- 覆盖         
│  ├─ utils ----------------- 公共方法    
│  ├─ views ----------------- 页面            
│  │  ├─  base -------------- 登录/402...               
│  │  └─  myapp ------------- 业务页面                                       
│  ├─ App.vue --------------- 页面入口              
│  ├─ element-variables.scss- 自定义 element style 入口  
│  ├─ main.ts --------------- 入口文件  
│  ├─ shims-vue.d.ts -------- 帮助理解.vue文件     
│  └─ registerServiceWorker.ts简化缓存机制     
│                                                   
├─ test --------------------- 单元测试                                
│  
├─ .browserslistrc ---------- 目标浏览器配置    
│  
├─ .env --------------------- 公共环境变量    
│  
├─ .env.development --------- 本地开发环境变量配置  
│    
├─ .env.prod ------------- 本地生产环境变量配置   
│                                            
├─ .env.test ------------- 本地测试环境变量配置   
│                                            
├─ .eslintignore ------------ eslint 忽略目录配置   
│      
├─ .eslintrc.js ------------- eslint配置  
│      
├─ .gitignore --------------- git忽略目录配置  
│                                                     
├─ tsconfig.json ------------ typescript配置  
│  
└─ vue.config.js ------------ vue项目全局配置(webpack)

### 安装依赖
```
npm install
```

### 本地运行dev环境
```
npm run serve
```

### 本地运行测试环境
```
npm run test
```

### 本地运行生产环境
```
npm run prod
```

### 代码打包生成dist
```
npm run build
```

### 格式化你的代码
```
npm run lint
```

### 创建新模块

```
npm run new-module

```
>之后按提示步骤执行即可，输入路径时，需要输入views地下的目录
>例如：myapp/app/myapplication/activity

node v14.17.5
41