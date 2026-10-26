<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/17 16:42
 */

namespace App\Http\Controllers\Service\V1\PartnerMember;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerMember;

class DetailAllController extends BaseController
{
    protected  $select = ['partner_sn','partner_name','partner_mobile','open_id','avatar',
        'mobile_short','partner_id','enterprise_id','worker_id','worker_name','worker_uid','partner_status','mobile_show','map.id as partner_map_id'];
    public function index()
    {
        $rules    = [
            'mobile'=>'required',//手机号
            'enterprise_id'=>'required',
        ];
        $messages = [
            'mobile.*'     => '手机号不合法',
            'enterprise_id.*'     => '企业ID必须存在',
        ];
        $this->validate($this->request, $rules, $messages);
        $input = $this->request->all();

        $result = PartnerMember::query()
            ->leftJoin('partner_enterprise_map as map','map.partner_id','partner_member.id')
            ->where('partner_member.partner_mobile',$input['mobile'])
            ->where('map.enterprise_id',$input['enterprise_id'])
            ->select($this->select)
            ->first();

        $result = !empty($result) ? $result->toArray() : [];

        return success($result);
    }
}