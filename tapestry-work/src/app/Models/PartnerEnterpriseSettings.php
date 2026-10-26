<?php
/**
 * @author  HuiYang
 * @date    2020/11/16 16:32
 */

namespace App\Models;


class PartnerEnterpriseSettings extends Base
{
    /**
     * 二级分销针对企业全局设置表
     * @var string
     */
    protected $table = 'partner_enterprise_settings';

    /**
     * 1:开启二级分销 0:关闭二级分销
     */
    const PARTNER_STATUS_ENABLE = 1;
    /**
     * 1:开启二级分销 0:关闭二级分销
     */
    const PARTNER_STATUS_DISABLE = 0;

    /**
     * 二级佣金设置
     * 1: 基于商品分销
     */
    const COMMISSION_TYPE_GOODS = 1;

    /**
     * 二级佣金设置
     * 2: 基于推广佣金总额的分账
     */
    const COMMISSION_TYPE_PROMOTION = 2;

    /**
     * 是否可提现 0 不可提现 1可提现
     */
    const WITHDRAW_STATUS_TRUE = 1;
    const WITHDRAW_STATUS_FALSE = 0;

    protected $fillable = [
        'id',
        'enterprise_id',
        'opt_admin_id',
        'opt_admin_name',
        'partner_status',
        'commission_type',
        'commission_rate',
        'withdraw_status',
        'created_at',
        'updated_at',
        'deleted_at'
    ];

    /**
     * 更新或提交
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/20 11:11
     */
    public function  insertOrUpdate($input)
    {
        //统一提现信息
        $this->distributionConfig($input);
        $exist =  PartnerEnterpriseSettings::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->first();
        if(empty($exist)){
            //一个企业只能设置一次
               $input['created_at'] = time();
               PartnerEnterpriseSettings::query()
                   ->insert($input);
            return [];
        }
        $input['updated_at'] = time();
         PartnerEnterpriseSettings::query()
             ->where('enterprise_id',$input['enterprise_id'])
             ->update($input);
            return [];
    }

    /**
     * 设置提现 提现设置和二级分销提现统一
     * @date    2020/11/24 16:46
     */
    public function distributionConfig($input)
    {

       $setting =  WithdrawalSetting::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->first();

        if(!empty($setting)){
            $array_type =  explode(',',$setting->type);
            //合伙人提现关闭
            if($input['withdraw_status'] == self::WITHDRAW_STATUS_FALSE
             && in_array(WithdrawalSetting::PARTNER_STATUS,$array_type)
            ){
                foreach($array_type as $k=>$v) {
                    if(WithdrawalSetting::PARTNER_STATUS == $v)
                        unset($array_type[$k]);
                }
                $array_type = array_values($array_type);
                $type = implode(',',$array_type);
                $setting->update(['type'=>$type]);
            }
            //合伙人提现关闭
            if($input['withdraw_status'] == self::WITHDRAW_STATUS_TRUE
                && !in_array(WithdrawalSetting::PARTNER_STATUS,$array_type))
            {
                $setting->update(['type'=>$setting->type.','.WithdrawalSetting::PARTNER_STATUS]);
            }
        }
        return true;
    }
}