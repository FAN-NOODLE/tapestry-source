<?php


namespace App\Common;


class MessageConstant
{

    const EXTERNAL_CONTACT_CREATE = 1;
    const EXTERNAL_CONTACT_DELETE = 2;


    const SHARE_LINK_CREATE = 1;
    const SHARE_LINK_DELETE = 2;
    const SHARE_LINK_UPDATE = 3;

    // 业绩集成
    const SHARE_LINK_BATCH_INHERIT = 1;
    // 员工离职
    const SHARE_LINK_BATCH_LEAVE = 2;
    // 超过保护期自动解绑
    const SHARE_LINK_BATCH_AUTO_UNBINDING = 3;
    // 员工变导购
    const SHARE_LINK_BATCH_CREATE_GUIDE = 4;
}