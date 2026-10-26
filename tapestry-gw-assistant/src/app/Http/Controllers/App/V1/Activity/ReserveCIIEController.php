<?php
/**
 * Coach的CIIE吊牌预约活动 (用户端)
 */

namespace App\Http\Controllers\App\V1\Activity;

use App\Http\Controllers\Controller;
use App\Micros\Work\Activity\ReserveCIIE;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Psr\SimpleCache\InvalidArgumentException;


class ReserveCIIEController extends Controller
{

    protected $ciie;

    public function __construct()
    {
        parent::__construct();
        $this->ciie = new ReserveCIIE();
    }

    /**
     * 创建预约
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function CreateReserve(Request $request): JsonResponse
    {
        $tokenData               = $request->get('tokenData');
        $params                  = $request->all();
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];
        if (!preg_match('/^1[0-9]{10}$/', $params['mobile'])) {
            error(651003);
        }
        try {
            $result = $this->ciie->create($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage(), 200);
        }
        return success($result);
    }

    /**
     * 获取活动信息
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function Info(Request $request): JsonResponse
    {
        $tokenData               = $request->get('tokenData');
        $params                  = $request->all();
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];
        try {
            $result = $this->ciie->Info($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage(), 200);
        }
        $this->ciie->UserView($params);
        return success($result);
    }

    /**
     * 取消预约
     * @param Request $request
     * @return JsonResponse
     */
    public function CancelReserve(Request $request): JsonResponse
    {
        $tokenData               = $request->get('tokenData');
        $params                  = $request->all();
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];
        try {
            $result = $this->ciie->cancel($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage(), 200);
        }
        return success($result);

    }

    /**
     * 编辑预约
     * @param Request $request
     * @return JsonResponse
     */
    public function EditReserve(Request $request): JsonResponse
    {

        $params                  = $request->all();
        $tokenData               = $request->get('tokenData');
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $params['openid']        = $tokenData['open_id'];
        if (!preg_match('/^1[0-9]{10}$/', $params['mobile'])) {
            error(651003);
        }
        try {
            $result = $this->ciie->edit($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage(), 200);
        }
        return success($result);
    }

    /**
     * 预约列表
     * @param Request $request
     * @return JsonResponse
     */
    public function ReserveLists(Request $request): JsonResponse
    {
        $params         = $request->all();
        $params['type'] = 1;
        try {
            $result = $this->ciie->TimeLists($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage(), 200);
        }
        return success($result);
    }

    /**
     * 获取微信授权的手机号
     * @param Request $request
     * @return JsonResponse
     * @throws InvalidArgumentException
     */
    public function GetPhoneNumber(Request $request): JsonResponse
    {
        $params = $request->all();

        if (empty($params['code'])) {//code的不能为空
            error(652012);
        }
        $tokenData               = $request->get('tokenData');
        $params                  = $request->all();
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        if (empty($params['enterprise_id'])) {
            error(652008);
        }
        try {
            $result = $this->ciie->GetPhoneNumber($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage(), 200);
        }
        return success($result);


    }

}