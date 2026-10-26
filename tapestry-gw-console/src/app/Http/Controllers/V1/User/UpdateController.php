<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use App\Common\Constant;
use Illuminate\Http\Request;
use CrsSdk\Work\MemberExtend;
use CrsSdk\MicroPay\Cmchcircle;
use Illuminate\Http\JsonResponse;
use App\Http\Requests\UserRequest;
use Illuminate\Support\Facades\Log;
use App\Http\Controllers\V1\BaseController;

class UpdateController extends BaseController
{
    /**
     * 编辑员工
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request)
    {
        // 参数列表
        $params = $request->post();

        // 参数验证
        $validator = (new UserRequest())->update($params);
        if ($validator->fails()) {
            error($validator->messages()->first());
        }

        // 部门是否被设置为门店（0:否，1:是）
        $isStore = isset($params['is_store']) ? intval($params['is_store']) : 0;
        // 员工角色
        $memRole = $params['mem_role'] ?? '';

        // 不是门店
        if (Constant::DEPARTMENT_IS_STORE_FALSE === $isStore) {
            $memRole = '';
        }

        // 头像
        $memFace = $params['mem_face'] ?? '';

        // 调用sdk
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'corpid'         => $this->enterprise['corp_id'],
            'app_identifier' => $this->appIdentifier,
            'enterprise_id'  => $this->enterpriseId,
            'dp_id'          => $params['dp_id'],
            'params'         => [
                'memUid'      => $params['mem_uid'],
                'memUsername' => $params['mem_username'],
                'memFace'     => $memFace,
                'dpIdList'    => [$params['dp_id']],
                'memGender'   => $params['mem_gender'],
                'memMobile'   => $params['mem_mobile'],
                'memRole'     => $memRole
            ]
        ];

        $detail = (new Member())->edit($condition);

        if ($detail) {

            $condition = [
                'qy_domain'      => $this->qyDomain,
                'app_identifier' => $this->appIdentifier,
                'params'         => ['memUid' => $detail['memUid']]
            ];

            $memDetail = (new Member())->detail($condition);

            $guideId = '';

            // 二维码获取
            $qr       = '';
            $configId = '';
            try {
                $where = [
                    'qy_domain'   => $this->qyDomain,
                    'mem_user_id' => $memDetail['memUserid'],
                    'ecw_state'   => Constant::GUIDE_ADDITION // 企业自定义的state参数，用于区分不同的添加渠道
                ];

                $result = (new Member())->contactQrCode($where);

                $qr       = $result['data']['qr_code'] ?? '';
                $configId = $result['data']['config_id'] ?? '';
            } catch (\Exception $e) {
                Log::error('update-qr-select-error', ['code' => $e->getCode(), 'msg' => $e->getMessage(), 'line' => $e->getLine()]);
            }

            if (self::GUIDE_MEMBER_ROLE == trim($memRole)) {
                try {
                    $where = [
                        'enterprise_id' => $this->enterpriseId,
                        'mem_uids'      => [$detail['memUid']]
                    ];

                    // // 先查询是否是已注册导购
                    // $where = [
                    //     'shop_id'       => 0,
                    //     'userid'        => $memDetail['memUserid'],
                    //     'mobile'        => $params['mem_mobile'],
                    //     'enterprise_id' => $this->enterpriseId,
                    //     'hash'          => $this->enterprise['hash']
                    // ];

                    // $result = (new Cmchcircle())->queryGuide($where);

                    $result = (new MemberExtend())->extendList($where);

                    $guideId = isset($result['data'][0]) ? $result['data'][0]['guide_id'] : '';
                } catch (\Exception $e) {
                    Log::error('guide-select-error', ['code' => $e->getCode(), 'msg' => $e->getMessage(), 'line' => $e->getLine()]);
                }

                // 获取到导购id就更新导购基本信息
                if ($guideId) {
                    try {
                        $update = [
                            'guide_id'      => $guideId,
                            'mobile'        => $params['mem_mobile'],
                            'name'          => $params['mem_username'],
                            'avatar'        => $memDetail['memFace'],
                            'qr_code'       => $qr,
                            'enterprise_id' => $this->enterpriseId,
                            'hash'          => $this->enterprise['hash']
                        ];

                        $result = (new Cmchcircle())->updateGuiderInfo($update);

                    } catch (\Exception $e) {
                        Log::error('guide-update-error', ['code' => $e->getCode(), 'msg' => $e->getMessage(), 'line' => $e->getLine()]);
                    }

                } else {
                    // 没有导购并且获取到二维码就注册导购
                    if ($qr) {
                        try {
                            $create = [
                                'userid'        => $memDetail['memUserid'],
                                'name'          => $memDetail['memUsername'],
                                'corpid'        => $this->enterprise['corp_id'],
                                'shop_id'       => 0,
                                'mobile'        => $params['mem_mobile'],
                                'avatar'        => $memDetail['memFace'],
                                'qr_code'       => $qr,
                                'enterprise_id' => $this->enterpriseId,
                                'hash'          => $this->enterprise['hash']
                            ];

                            $result = (new Cmchcircle())->createGuide($create);

                            $guideId = $result['guide_id'];
                        } catch (\Exception $e) {
                            Log::error('guide-create-error', ['code' => $e->getCode(), 'msg' => $e->getMessage(), 'line' => $e->getLine()]);
                        }
                    }
                }
            }

            $save = [
                'enterprise_id'   => $this->enterpriseId,
                'guide_id'        => $guideId,
                'mem_uid'         => $detail['memUid'],
                'avatar'          => $memFace,
                'contact_qr_code' => $qr,
                'config_id'       => $configId,
                'role'            => trim($memRole),
            ];

            // 保存或更新员工信息
            (new MemberExtend())->saveExtend($save);
        }
        // 写入记录日志
        $this->writeOperatorLog('编辑员工：' . $params['mem_username'] . '，角色：' . $memRole);

        return success();
    }
}
