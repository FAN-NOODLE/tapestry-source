<?php
/**
 * 引导页浏览记录
 */
namespace App\Http\Controllers\Service\V1\Common;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\BootPageLog;
use CrsSdk\Boss\Auth;
use Illuminate\Http\Request;

class BootPageLogController extends BaseController
{
    /**
     * 引导页浏览记录 - 写入
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $param = $request->all([
            'enterprise_id',
            'version_code',
            'worker_id'
        ]);
        //引导页这样验证是为了不影响程序正常向下执行
        if (!empty($param['enterprise_id'])
            && !empty($param['version_code'])
            && !empty($param['worker_id'])) {
            $data = [
                'enterprise_id' => $param['enterprise_id'],
                'version_code' => $param['version_code'],
                'worker_id' => $param['worker_id']
            ];
            //查询并创建模型对象
            $logObj = BootPageLog::query()
                ->firstOrNew($data);
            //如果用户之前在该版本从没有浏览过记录，则填充数据写入，否则直接返回
            if (!isset($logObj->id)) {
                $logObj->fill($data)->save();
            }
        }

        return success([]);
    }
    /**
     * 引导页浏览记录 - 检查是否需要展示
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function check(Request $request)
    {
        $param = $request->all([
            'enterprise_id',
            'version_code',
            'worker_id'
        ]);
        $show = false;
        if (!empty($param['enterprise_id'])
            && !empty($param['version_code'])
            && !empty($param['worker_id'])) {
            $info = BootPageLog::query()
                ->where('enterprise_id', $param['enterprise_id'])
                ->where('version_code', $param['version_code'])
                ->where('worker_id', $param['worker_id'])
                ->first();
            if (!isset($info->id)) {
                $show = true;
            }
        }

        return success([
            'show' => boolval($show)
        ]);
    }
}
