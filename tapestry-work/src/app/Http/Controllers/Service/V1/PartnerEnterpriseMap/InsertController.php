<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/17 20:16
 */

namespace App\Http\Controllers\Service\V1\PartnerEnterpriseMap;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerEnterpriseMap;
use App\Models\PartnerMember;
use Illuminate\Support\Facades\Redis;

class InsertController extends BaseController
{
    public function index()
    {
        $rules    = [
            'partner_id'      => 'required|integer',//合伙人ID
            'enterprise_id'=>'required',//企业ID
            'worker_id'=>'required|integer',//导购ID
            'worker_uid'=>'required',//导购ID
            'worker_name'=>'required',//导购名称
            'mobile_sensitive'=>'required|in:0,1',//是否脱敏
            'mobile_show'=>'required',//是否脱敏展示的手机号
        ];
        $messages = [
            'partner_id.*'       => '合伙人不合法',
            'enterprise_id.*'     => '企业ID必须存在',
            'worker_id.*'       => '导购ID必须存在',
            'worker_uid.*'       => '导购mem_uid必须存在',
            'worker_name.*'       => '导购名称必须存在',
            'mobile_sensitive.*'       => '是否脱敏不合法',
            'mobile_show.*'       => '是否脱敏展示的手机号',

        ];
        $this->validate($this->request, $rules, $messages);
        $input = $this->request->all();

        //加分布式锁防止重复提交
        $lockKey = 'add_partner:'.$input['partner_id'].'-'.$input['enterprise_id'];
        $getCache = Redis::get($lockKey);
        //锁是否会过期 释放(防止死锁)
        if(!empty($getCache) && $getCache <time()){
            Redis::del($lockKey);
        }
        //加锁
        $setCache = Redis::setnx($lockKey,time()+60);
        //如果有锁 直接跳出
        if(!$setCache){
            notice(1100032);
        }
        //企业唯一性效验
        $partnerId = PartnerEnterpriseMap::query()
            ->where('partner_id',$input['partner_id'])
            ->where('enterprise_id',$input['enterprise_id'])
            ->first();
        if(!empty($partnerId)){
            notice(1100030);
        }
        $input['created_at'] = time();
        $input['updated_at'] = time();
        $input['id'] = idGenerate();
        $result = PartnerEnterpriseMap::query()
            ->insertGetId($input);
        //删除锁
        Redis::del([$lockKey]);
        return success(['data'=>$result]);
    }

}