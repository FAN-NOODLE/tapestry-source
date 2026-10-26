<?php
/**
 * 营销计划可见范围表
 */

namespace App\Models;

class MarketingPlanVisibleRange extends Base
{

    /**
     * 表名称
     * @var string
     */
    protected $table = 'marketing_plan_visible_range';

    // 主键
    protected $primaryKey = 'id';

    /**
     * 1、按门店所属区域(默认),2、按导购所属部门
     */
    const CHOICE_AREA = 1;
    /**
     * 1、按门店所属区域(默认),2、按导购所属部门
     */
    const CHOICE_DEPARTMENT = 2;

    const range_select = [
        //可见范围表
        "worker_id",
        "object_type"
    ];
}
