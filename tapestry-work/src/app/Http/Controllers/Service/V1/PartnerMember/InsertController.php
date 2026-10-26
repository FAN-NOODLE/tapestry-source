<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/17 19:01
 */

namespace App\Http\Controllers\Service\V1\PartnerMember;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerMember;
use App\Models\PartnerMemberAuth;

class InsertController extends BaseController
{
    public function index()
    {
        $rules    = [
            'avatar'      => 'required',//头像
            'partner_name'=>'required',//合伙人昵称
            'partner_mobile'=>'required',//手机号
            'open_id'=>'required',
            'mobile_short'=>'required',//手机短号
            'partner_sn'=>'required',//合伙人编号

        ];
        $messages = [
            'avatar.*'       => '头像信息必须存在',
            'partner_mobile.*'     => '手机号不合法',
            'partner_name.*'       => '合伙人昵称必须存在',
            'open_id.*'       => 'openID必须存在',
            'mobile_short.*'       => '手机短号必须存在',
            'partner_sn.*'       => '合伙人编号必须存在',

        ];
        $this->validate($this->request, $rules, $messages);
        $input = $this->request->all();

        $input['created_at'] = time();
        $input['updated_at'] = time();
        $input['id'] = idGenerate();

        PartnerMemberAuth::query()
            ->insert($input);

        unset($input['app_id']);

        $result = PartnerMember::query()
            ->where('partner_mobile',$input['partner_mobile'])
            ->value('id');
        if(!$result){
            $input['id'] = idGenerate();
            $result = PartnerMember::query()
                ->insertGetId($input);
        }

        return success(['data'=>$result]);

    }

}