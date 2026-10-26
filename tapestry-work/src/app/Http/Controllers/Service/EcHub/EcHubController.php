<?php

namespace App\Http\Controllers\Service\EcHub;

use App\Http\Controllers\Service\V1\BaseController;
use App\Sdk\EcHub\Request;
use App\Services\CommonService;
use Exception;
use GuzzleHttp\Exception\GuzzleException;
use Illuminate\Http\JsonResponse;

class EcHubController extends BaseController
{

    private $enterpriseId;

    public function __construct()
    {
        parent::__construct();
        $this->enterpriseId = $this->request->input('enterprise_id');

    }


    /**
     * @return JsonResponse
     * @throws GuzzleException
     */
    public function saAssignment(): JsonResponse
    {
        $params = $this->request->all();
        $data=$params["params"];
        $result = [];
        try {
            $result =  (new Request($this->enterpriseId))->setUrl("saAssignment")->sendPost($data);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return successJson($result);
    }

    /**
     * @return JsonResponse
     * @throws GuzzleException
     */
    public function saUnAssignment(): JsonResponse
    {
        $params = $this->request->all();
        $data=$params["params"];
        $result = [];
        try {
            $result =  (new Request($this->enterpriseId))->setUrl("saUnAssignment")->sendPost($data);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return successJson($result);
    }

    /**
     * @return array
     * @throws GuzzleException
     */
    public function getEcUrlScheme(): array
    {
        $params = $this->request->all();
        return success($this->getUrlScheme($params));
    }

    /**
     * @return array
     * @throws GuzzleException
     */
    public function getEcQrCode(): array
    {
        $params = $this->request->all();
        return success($this->getQrCode($params));
    }

    /**
     * @return array
     * @throws GuzzleException
     */
    public function getQrCodeAndLink(): array
    {
        $params = $this->request->all();
        return success(array_merge($this->getUrlScheme($params), $this->getQrCode($params)));
    }


    /**
     * @param $params
     * @return array
     * @throws GuzzleException
     */
    private function getUrlScheme($params): array
    {
        $result['link'] = '';
        try {
            $res = (new Request($this->enterpriseId))->setUrl("url-scheme")->sendPost([
                'path' => $params['path'],
                'parameterString' => 'scene=' . $params['scene'],
//                'envVersion' => $params['env_version'] ?: 'release'
            ]);
            if ($res['msg'] == 'success' && isset($res['data']['url'])) {
                $result['link']  = $res['data']['url'];
            }
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return $result;
    }

    /**
     * @param $params
     * @return array
     * @throws GuzzleException
     */
    private function getQrCode($params): array
    {
        $result['code_url'] = '';
        try {
            $res = (new Request($this->enterpriseId))->setUrl("qr-code")->sendPost([
                'page' => $params['path'],
                'scene' => 'scene=' . $params['scene'],
//                'envVersion' => $params['env_version'] ?: 'release',
                'width' => $params['width'] ?? 1417
            ]);
            if ($res['msg'] == 'success' && isset($res['data']['qrCode'])) {
                $result['code_url']  = CommonService::uploadContent(base64_decode($res['data']['qrCode']));
            }
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        return $result;
    }
}
