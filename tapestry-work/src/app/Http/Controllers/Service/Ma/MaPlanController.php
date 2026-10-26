<?php

namespace App\Http\Controllers\Service\Ma;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\Ma\MaPlanService;
use Exception;
use GuzzleHttp\Exception\GuzzleException;
use Illuminate\Http\JsonResponse;


class MaPlanController extends BaseController
{
    private $enterpriseId;
    private $originType;

    public function __construct()
    {
        parent::__construct();
        $this->enterpriseId = $this->request->post('enterprise_id');
        $this->originType = $this->request->post('origin_type');
    }

    /**
     * @return JsonResponse
     * @throws GuzzleException
     */
    public function sendResult(): JsonResponse
    {
        $params = $this->request->all();
        $result = [];
        try {
            $result = (new MaPlanService($this->enterpriseId,$this->originType))->sendResult($params);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return successJson($result);
    }

    /**
     * @return JsonResponse
     */
    public function sendFinish(): JsonResponse
    {
        $params = $this->request->all();
        $result = [];
        try {
            $result = (new MaPlanService($this->enterpriseId,$this->originType))->sendFinish($params);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return successJson($result);
    }

    /**
     * @return JsonResponse
     * @throws GuzzleException
     */
    public function momentSendResult(): JsonResponse
    {
        $params = $this->request->all();
        $result = [];
        try {
            $result = (new MaPlanService($this->enterpriseId,$this->originType))->momentSendResult($params);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return successJson($result);
    }

    /**
     * @return JsonResponse
     * @throws GuzzleException
     */
    public function momentComments(): JsonResponse
    {
        $params = $this->request->all();
        try {
            $result = (new MaPlanService($this->enterpriseId,$this->originType))->momentComments($params);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return successJson($result);
    }
}