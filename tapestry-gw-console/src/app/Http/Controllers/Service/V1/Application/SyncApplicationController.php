<?php
/**
 * 同步boss后台企业应用
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/2/20
 * Time: 15:19
 */
namespace App\Http\Controllers\Service\V1\Application;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\ApplicationPower;
use App\Services\ApplicationPowerService;
use Illuminate\Http\Request;

class SyncApplicationController extends BaseController
{
    /**
     *  同步boss后台企业应用设置
     * @param Request $request
     * @param ApplicationPowerService $applicationPowerS
     * @return array
     */
    public function index(Request $request, ApplicationPowerService $applicationPowerS)
    {
        $applications = $request->post('applications');
        $enterpriseId = $request->post('enterprise_id');

        if (!$enterpriseId) {
            throw new  Error(10000, '企业id不能为空喔');
        }

        if (!$applications) {
            throw new  Error(10000, '应用参数不能为空喔');
        }

        // 组装入库数据格式
        $data = $this->dealDifferApplication($applications, $enterpriseId);

        if (!empty($data)) {
            $res = $applicationPowerS->syncApplicationToConsole($data);
            if ($res) {
                return success(['status' => true]);
            } else {
                throw  new Error(10000, '应用数据同步失败');
            }
        } else {
            return success(['status' => true]);
        }
    }


    /**
     * 处理即将同步应用数据与当前企业开通应用数据差异
     * @param $applications
     * @param $enterpriseId
     * @return array
     */
    public function dealDifferApplication($applications, $enterpriseId)
    {
        $alreadyData = ApplicationPower::query()
                        ->where('enterprise_id', '=', $enterpriseId)
                        ->where('parent', '=', 'myApplication')
                        ->get(['path'])
                        ->toArray();
        // 首次同步 直接格式化插入
        if (empty($alreadyData)) {
            $result = $applications;
        } else {
            $result = $this->getApplicationData($alreadyData, $applications, $enterpriseId);
        }
        // 获取配置数据
        $applicationInfo = config('application');
        if (empty($applicationInfo)) {
            throw new Error(10000, '应用配置数据不存在');
        }
        // 获取应用配置数据
        $result = $this->getFixData($result, $applicationInfo);
        // 组装应用父级数据
        $result = $this->packParentData($result, $applicationInfo, $enterpriseId);
        // 组装企业信息
        $result = $this->packEnterpriseInfo($result, $enterpriseId);
        return $result;
    }

    /**
     * 获取应该数据库操作数据
     * @param $alreadyData
     * @param $applications
     * @param $enterpriseId
     * @return array
     */
    public function getApplicationData($alreadyData, $applications, $enterpriseId)
    {
        $result = [];
        $alApplications = array_column($alreadyData, 'path');
        $intersect = array_intersect($applications, $alApplications);

        // 交集不存在，删除原应用数据，新数据入库
        if (empty($intersect)) {
            $result = $applications;
            ApplicationPower::query()->where('enterprise_id', '=', $enterpriseId)->delete();
        }

        // 取原数据差集, 当前数据差集
        $alDiffer = array_diff($alApplications, $intersect);
        $differ = array_diff($applications, $intersect);

        if (!empty($alDiffer) && !empty($differ)) {
            // 删$alDiffer
            $result = $differ;
            ApplicationPower::query()->where('enterprise_id', '=', $enterpriseId)
                ->whereIn('path', $alDiffer)
                ->delete();
        }
        if (!empty($alDiffer) && empty($differ)) {
            // 删$alDiffer
            ApplicationPower::query()->where('enterprise_id', '=', $enterpriseId)
                ->whereIn('path', $alDiffer)
                ->delete();
        }
        if (empty($alDiffer) && !empty($differ)) {
            $result = $differ;
        }
        return $result;
    }

    /**
     * 组装数据
     * @param $data
     * @param $enterpriseId
     * @return array
     */
    public function packEnterpriseInfo($data, $enterpriseId)
    {
        if (empty($data)) {
            return [];
        }
        // 拼接企业、时间到应用配置信息
        $arr = [
            'enterprise_id' => $enterpriseId,
            'created_at' => time()
        ];

        array_walk($data, function (&$value) use ($arr) {
            $value = array_merge($value, $arr);
        });
        return $data;
    }

    /**
     * 获取配置数据格式
     * @param $applications
     * @param $applicationInfo
     * @return array
     */
    public function getFixData($applications, $applicationInfo)
    {
        // 初始化返回数据
        $data = [];
        foreach ($applications as $application) {
            if (isset($applicationInfo[$application])) {
                $data[] = $applicationInfo[$application];
            }
        }
        return $data;
    }

    /**
     * 拼接应用父级数据
     * @param $data
     * @param $applicationInfo
     * @param $enterpriseId
     * @return mixed
     */
    public function packParentData($data, $applicationInfo, $enterpriseId)
    {
        // 获取应用父级配置数据
        $parent = ApplicationPower::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->where('path', '=', 'myApplication')
            ->get()
            ->toArray();
        if (empty($parent)) {
            // 检索应用服务类型，生成上一层服务状态
            $serveTypeArr = array_unique(array_column($data, 'serve_type'));
            $parentServeType = ApplicationPowerService::SERVE_TYPE_ALL;
            // 若只有一种服务类型， 且服务类型是 商城
            if ((count($serveTypeArr) == 1) && ($serveTypeArr[0] == ApplicationPowerService::SERVE_TYPE_SHOP)) {
                $parentServeType = ApplicationPowerService::SERVE_TYPE_SHOP;
            }

            // 若只有一种服务类型， 且服务类型是 crs
            if ((count($serveTypeArr) == 1) && ($serveTypeArr[0] == ApplicationPowerService::SERVE_TYPE_CRS)) {
                $parentServeType = ApplicationPowerService::SERVE_TYPE_CRS;
            }
            // 若超过一种服务类型， 且服务类型是 商城+crs
            if (count($serveTypeArr) > 1) {
                $parentServeType = ApplicationPowerService::SERVE_TYPE_ALL;
            }

            // 重置应用配置 中 前两级服务类型
            $applicationInfo['application']['serve_type'] = $parentServeType;
            $applicationInfo['myApplication']['serve_type'] = $parentServeType;
            array_unshift($data, $applicationInfo['application'], $applicationInfo['myApplication']);
        }
        return $data;
    }
}
