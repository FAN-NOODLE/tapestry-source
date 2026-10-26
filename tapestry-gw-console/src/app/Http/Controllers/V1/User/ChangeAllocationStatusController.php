<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 11:19
 */

namespace App\Http\Controllers\V1\User;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroPay\Cmchcircle;
use CrsSdk\Shop\WxArea;
use CrsSdk\Shop\WxStore;
use CrsSdk\Work\Member;
use CrsSdk\Work\MemberExtend;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class ChangeAllocationStatusController extends BaseController
{
    /**
     * 更新员工支付即导购分配状态
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {

        // 更新支付即导购分配状态：1、开启，2、关闭
        $tradeGuideStatus = trim($request->post('trade_guide_status'));

        // 员工ID
        $memUid = trim($request->post('mem_uid'));

        if (empty($memUid)) {
            error(320151);
        }

        // 开启支付即导购
        if ($tradeGuideStatus == 1) {
            $where = [
                'enterprise_id' => $this->enterpriseId,
                'mem_uids'      => [$memUid]
            ];

            $result = (new MemberExtend())->extendList($where);

            $guideId = isset($result['data'][0]) ? $result['data'][0]['guide_id'] : '';
            if (empty($guideId)) {
                // 获取微信门店编号 v3要用
                $storeNo = (new WxStore())->storeNo(['enterprise_id' => $this->enterpriseId]);

                if (empty($storeNo['store_no'])) {
                    notice(1900001);
                }

                $condition = [
                    'qy_domain'      => $this->qyDomain,
                    'app_identifier' => $this->appIdentifier,
                    'params'         => ['memUid' => $memUid]
                ];

                $memDetail = (new Member())->detail($condition);

                $qrCodeParams = [
                    'enterprise_id' => $this->enterpriseId,
                    'domain'        => $this->qyDomain,
                    'uc_user_id'    => $memDetail['memUserid'],
                    'uc_mem_id'     => $memUid,
                    'uc_mem_name'   => $memDetail['memUsername'],
                    'source_type'   => 3 // 支付即导购
                ];
                $qrCodeData       = (new Member())->addQrCode($qrCodeParams);
                $contactQrCodeData       = $qrCodeData['data'];

                // 没有导购并且获取到二维码就注册导购
                if ($contactQrCodeData) {

                    $create = [
                        'userid'        => $memDetail['memUserid'],
                        'name'          => $memDetail['memUsername'],
                        'corpid'        => $this->enterprise['corp_id'],
                        'store_id'      => (int)$storeNo['store_no'],
                        'mobile'        => $memDetail['memMobile'],
                        'avatar'        => $memDetail['memFace'],
                        'qr_code'       => $contactQrCodeData['official_url'],
                        'enterprise_id' => $this->enterpriseId,
                        'hash'          => $this->enterprise['hash']
                    ];

                    $result = (new Cmchcircle())->createGuide($create);
                    Log::debug('注册导购响应结果', ['result' => $result, 'params' => $create]);

                    // 白名单问题
                    if (isset($result['message']) && $result['message'] !== 'OK') {
                        notice(1900002, $result['message']);
                    }

                    $guideId = $result['guide_id'];

                    $save = [
                        'enterprise_id' => $this->enterpriseId,
                        'guide_id'      => $guideId,
                        'mem_uid'       => $memUid,
                        'contact_qr_code_id' => $contactQrCodeData['id']
                    ];

                    // 保存或更新员工信息
                    (new MemberExtend())->saveExtend($save);
                }
            }
        }

        $returnData = (new MemberExtend)->changeAllocationStatus([
            'trade_guide_status' => $tradeGuideStatus,
            'mem_uid'            => $memUid,
            'enterprise_id'      => $this->enterpriseId
        ]);

      if ($tradeGuideStatus) {
            $tradeGuideStatusName = $tradeGuideStatus == 1 ? '开启' : '关闭';
            $this->writeOperatorLog('更新支付即导购分配状态:' . $tradeGuideStatusName);
        }

        return success($returnData);
    }
}
