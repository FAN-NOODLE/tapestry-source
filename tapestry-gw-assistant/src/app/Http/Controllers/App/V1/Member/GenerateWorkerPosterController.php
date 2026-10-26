<?php
/**
 * Created by PHPStorm.
 * author: 刘朋波
 * Date: 2020/4/10
 * Time: 16:38
 */

namespace App\Http\Controllers\App\V1\Member;

use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Models\Scrm\Members;
use CrsSdk\Work\StoreWorker;

/**
 * 生成海报接口
 * Class EditStatusController
 * @package App\Http\Controllers\App\V1\Member
 */
class GenerateWorkerPosterController extends Controller
{
    /**
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function __invoke()
    {

        try {
            $input['enterprise_id'] = $this->enterpriseId;
            $input['mobile']        = $this->mobile;
            $input['store_id']      = request()->input('store_id', 0);
            $storeWorker            = new StoreWorker();
//            $detail                 = $storeWorker->detail($input);
            $detail = GoMicroScrm::getMemberInfo([
                'enterprise_id' => $this->enterpriseId,
                'mobile'        => $this->mobile,
            ]);

            if (!empty($detail)) {
                $storeWorker->generateWorkerPoster([
                    'enterprise_id' => $this->enterpriseId,
                    'id'            => $detail['id'],
                    'domain'        => $this->loginAdminEnterprise['domain']
                ]);
            }
            return success();
        } catch (\Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
    }
}
