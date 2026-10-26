<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use CrsSdk\Sdk;
use Exception;
use Illuminate\Http\JsonResponse;

class SowTaskController extends BaseController
{
    /**
     * @var Sdk
     */
    protected $microScrm;

    public function __construct()
    {
        parent::__construct();
        $this->microScrm = Common::microScrm();
    }


    public function list(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_mobile'] = $this->mobile;
        $result = [];
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/sow-task/list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        $codeDescription = ['O' => '非会员',
            'PROSPECT' => '品牌粉丝',
            'C' => '新晋会员',
            'D' => '商场会员',
            'GOLD' => '金卡会员',
            'SILVER' => '银卡会员',
            'DIAMOND' => '钻卡会员'];
        if (!empty($result['list'])) {
            foreach ($result['list'] as $key => $val) {
                $result['list'][$key]['grade'] = $codeDescription[$val['grade']] ?? '非会员';
            }
        }
        return success($result);
    }

    public function create(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_mobile'] = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/sow-task/create')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success(['id' => $result]);
    }

    public function update(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_mobile'] = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/sow-task/update')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    public function finished(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_mobile'] = $this->mobile;
        $params['mobile'] = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/sow-task/finished')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    public function delete(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_mobile'] = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/sow-task/delete')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    public function templateList(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/sow-task/template-list')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    public function newCreate(): JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_mobile'] = $this->mobile;
        try {
            $result = $this->microScrm->setApi('micro-scrm/v1/customer/sow-task/new-create')->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }
}