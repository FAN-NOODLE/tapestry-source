<?php
/**
 * Created by PhpStorm.
 * User: Yang Hui
 * Date: 2020年10月17日11:16:11
 */

namespace App\Models;

class StoreWorkerCertification extends Base
{
    /**
     * 证件类型(0：默认无 1:居民身份证, 3:港澳通行证 4：港澳居住证 5：台湾通行证 6：台湾居住证)
     */
    const ID_TYPE_DEFAULT = 0;
    /**
     * 证件类型(0：默认无 1:居民身份证, 3:港澳通行证 4：港澳居住证 5：台湾通行证 6：台湾居住证)
     */
    const ID_TYPE_MAINLAND_RESIDENT = 1;
    /**
     * 证件类型(0：默认无 1:居民身份证, 3:港澳通行证 4：港澳居住证 5：台湾通行证 6：台湾居住证)
     */
    const ID_TYPE_HK_MA_PASS = 3;
    /**
     * 证件类型(0：默认无 1:居民身份证, 3:港澳通行证 4：港澳居住证 5：台湾通行证 6：台湾居住证)
     */
    const ID_TYPE_HK_MA_RESIDENT = 4;
    /**
     * 证件类型(0：默认无 1:居民身份证, 3:港澳通行证 4：港澳居住证 5：台湾通行证 6：台湾居住证)
     */
    const ID_TYPE_TW_PASS = 5;
    /**
     * 证件类型(0：默认无 1:居民身份证, 3:港澳通行证 4：港澳居住证 5：台湾通行证 6：台湾居住证)
     */
    const ID_TYPE_TW_RESIDENT = 6;

    /**
     * 状态(0-默认 1:认证通过 2:认证失败)
     */
    const AUTH_STATUS_DEFAULT = 0;
    /**
     * 状态(0-默认 1:认证通过 2:认证失败)
     */
    const AUTH_STATUS_PASS = 1;
    /**
     * 状态(0-默认 1:认证通过 2:认证失败)
     */
    const AUTH_STATUS_FAILED = 2;

    /**
     * 状态(1:启用,2:禁用)
     */
    const STATUS_ENABLED = 1;
    /**
     * 状态(1:启用,2:禁用)
     */
    const STATUS_DISABLED = 2;

    /**
     * 用户身份标识：1、导购   2、合伙人
     */
    const MEM_TYPE_WORKER = 1;

    const MEM_TYPE_PARTNER = 2;


    protected $table = 'store_worker_certification';

    protected $primaryKey = 'id';

    protected $guarded = [];
}
