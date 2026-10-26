<?php


namespace App\Models;


class StoreDepartment extends Base
{

    CONST STORE_PARENT_DEPT_ID_KEY = 'store_root_dept_id';

    CONST SYNC_SUCCESS = 1; // 同步成功

    CONST POPUP_STORE_CLOSE = 0; //关闭 popue store

    CONST POPUP_STORE_OPEN = 1; //开启 popue store

    CONST SA_CHECK_ROLE_BY_SM = 1; //只允许SM查看全部订单

    CONST SA_CHECK_ROLE_BY_OTHERS = 2; //配置了其他人也可以查看全部订单

    protected $table = 'store_department';

    protected $guarded = [];

}
