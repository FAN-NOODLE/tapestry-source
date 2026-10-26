<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/16 16:26
 */

namespace App\Http\Controllers\Service\V1\PartnerMember;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerEnterpriseMap;
use App\Models\PartnerMember;
use App\Models\PartnerMemberAuth;

class CheckController extends BaseController
{
    const   REGISTER  = 1;//去注册
    const   SECOND_REG = 2;//二次注册
    const   REGISTERED = 3;//已注册

    public function index()
    {
        $open_id = $this->request->post('open_id');
        $enterprise_id =$this->request->post('enterprise_id');
        if (empty($open_id)) {
            throw new Error(110002, 'open_id不能为空');
        }
        $partnerAuthInfo = PartnerMemberAuth::query()
            ->where('open_id',$open_id)
            ->first();

        $result = [];
        if(empty($partnerAuthInfo)){//未注册
            $result['status'] = self::REGISTER;//去注册
            $result['partner_info'] = [];

            return success($result);
        }

        $partnerInfo = PartnerMember::query()
            ->where('partner_mobile', $partnerAuthInfo->partner_mobile)
            ->first();
        if(empty($partnerInfo)){//未注册
            $result['status'] = self::REGISTER;//去注册
            $result['partner_info'] = [];

            return success($result);
        }
        $partnerInfo = $partnerInfo->toArray();

        $map = PartnerEnterpriseMap::query()
            ->where('partner_id',$partnerInfo['id'])
            ->select(['enterprise_id'])
            ->get()
            ->toArray();
        $enterprise_ids = array_column($map,'enterprise_id');
        //判断企业id是否在$map中
        if(!in_array($enterprise_id,$enterprise_ids)){
            $result['status'] = self::SECOND_REG;//二次注册
            //获得合伙人信息
            $result['partner_info'] = $partnerInfo;
        }else{
            //存在
            $result['status'] = self::REGISTERED;//已注册过
            $result['partner_info'] = [];
        }


        return success($result);
    }
}