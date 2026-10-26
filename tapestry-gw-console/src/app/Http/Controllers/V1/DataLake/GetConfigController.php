<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/1 0001
 * Time: 10:43
 */

namespace App\Http\Controllers\V1\DataLake;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\DataLake;
use Illuminate\Http\JsonResponse;

class GetConfigController extends BaseController
{
    private $sdk;

    public function __construct()
    {
        parent::__construct();
        $this->sdk = new DataLake();
    }

    /**
     * @return JsonResponse
     */
    public function index():JsonResponse
    {
        $params['enterprise_id'] = $this->enterpriseId;
        $params['page']          = 1;
        $params['size']          = 1;
        $info                    = $this->sdk->list($params)['list'][0] ?? [];
        if(!$info){
            return success([]);
        }
        return success([
            'access_token'=>$this->accessToken($info),
            'brand_config'=>$this->brandConfig($info)
        ]);
    }

    /**
     * @param array $info
     * @return array
     */
    private function accessToken(array $info):array
    {
        return [
            'vendor'        => $info['vendor'],
            'client_id'     => $info['client_id'],
            'client_secret' => $info['client_secret'],
            'seller_name'   => $info['seller_name'],
            'auth_user'     => $info['auth_user'],
            'auth_password' => $info['auth_password']
        ];
    }

    /**
     * @param array $info
     * @return array
     */
    private function brandConfig(array $info):array
    {
        return [
            'brand_code'   => $info['brand_code'],
            'program_code' => $info['program_code'],
            'store_code'   => $info['store_code']
        ];
    }
}