<?php

namespace App\Models;

class MobileWhite extends BaseModel
{
    protected $table = 'mobile_white';

    Const WHITE_ALLOW = 1;  //白名单允许
    Const ALL_ALLOW = 2;    //全部允许

    Const CUSTOMER_SEARCH = 1;  //会员查询按钮
    Const WECHAT_MOMENTS_SOURCE_UPLOAD = 2; //微信朋友圈素材上传按钮
    Const INVENTORY_SEARCH = 3; //库存查询按钮

    Const NEW_CUSTOMER_LIST = 4;// 新版顾客列表

    Const IS_ESA = 5;// 是否ESA
}
