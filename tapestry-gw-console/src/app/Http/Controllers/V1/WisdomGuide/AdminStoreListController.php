<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/29 0029
 * Time: 14:37
 */

namespace App\Http\Controllers\V1\WisdomGuide;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\Shop\StoreArea;

class AdminStoreListController extends BaseController
{
    private $accountSdk;

    public function __construct()
    {
        parent::__construct();
        $this->accountSdk = new Account();
    }

    public function index()
    {
//        $this->checkParam();
//        $adminStores = [];
//        $params = [
//            'admin_id' => $this->getAdminId()
//        ];
//        $res = $this->accountSdk->dataRange($params);
//
//        foreach ($res['data_range_list'] as $object) {
//            if ($object['type'] == 2) {
//                $adminStores[] = $object['object_id'];
//            }
//        }
//        if ($adminStores) {
//            $stores = (new StoreArea())->getStore(['store_ids' => $adminStores]);
//        }

        return success([
            'range_type' => 2,
//            'list' => $stores['list'] ?? []
            'list' => []
        ]);
    }

    /**
     * @throws \Illuminate\Validation\ValidationException
     */
    private function checkParam()
    {
        $this->validate($this->request, [
            'guide_mobile' => 'required'
        ]);
    }

    private function getAdminId()
    {
        $mobile = $this->request->input('guide_mobile');
        $input = [
            'mobile' => str_replace("****","%",$mobile),
            'enterprise_id' => $this->enterpriseId,
            'only' => 2
        ];
        return $this->accountSdk->detailAccount($input)['id'];
    }
}



