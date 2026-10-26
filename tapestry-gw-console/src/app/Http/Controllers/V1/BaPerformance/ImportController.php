<?php


namespace App\Http\Controllers\V1\BaPerformance;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\GoMicroScrm\BaPerformance;
use CrsSdk\Work\BaseClient;
use CrsSdk\Work\Member;
use CrsSdk\Work\StoreWorker;

class ImportController extends BaseController
{
    use ExcelTrait;

    protected $requireLogin = false;
    /**
     * 模板表头
     */
    const EXCEL_TITLE = [
        '区域（大区）',
        '区域',
        '城市',
        '柜台code',
        '柜台名称',
        '导购BA姓名',
        '手机号',
        '月度销售目标（元）',
        '新二购且一码三关注月度目标（人）',
        '老客购且一码三关注月度目标（人）',
        '新客购且一码三关注月度目标（人）',
    ];

    /**
     * 批量导入
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        if (!$this->request->hasFile('excel')) {
            notice(100000, ['未上传excel文件']);
        }

        $file = $this->request->file('excel');

        // 保存临时文件
        $saveName = uniqid('excel-') . '.' . $file->getClientOriginalExtension();
        $file->move(storage_path('logs/gw-console'), $saveName);
        $savePath = storage_path('logs/gw-console/' . $saveName);

        // 解析excel
        $excelData = $this->excelRead($savePath);
        // 删除临时文件
        unlink($savePath);

        if (count($excelData) > 500) {
            notice(100000, ['批量处理最多支持500条']);
        }

        // 判断excel模板是否被修改
        if (array_slice($excelData[0], 0, count(self::EXCEL_TITLE)) !== self::EXCEL_TITLE) {
            notice(100000, ['模板表头内容错误,请重新下载模板']);
        }

        //所有智慧导购手机号
        $allMobileArr = $this->allStormMem();
        // 开始处理数据
        $result = [
            'success_count' => 0,
            'error_count' => 0,
            'error_list' => [],
            'title' => array_merge(self::EXCEL_TITLE, ['失败原因'])
        ];
        $params['enterprise_id'] = $this->enterpriseId;
        $mobileArr = [];
        $row = 0;
        try {
            foreach ($excelData as $key => $value) {
                // 过滤表头
                if ($key === 0) {
                    continue;
                }
                $row = $key + 1;
                //字段验证
                if (!$this->is_mobile($value[6])) {
                    notice(100000, ['手机格式有误']);
                }
                if (!is_numeric($value[7]) || floor($value[8]) != $value[8] || floor($value[9]) != $value[9] || floor($value[10]) != $value[10]) {
                    notice(100000, ['目标业绩值应为数字格式，人数应为正整数']);
                }
                if ($value[7] < 0 || $value[8] < 0 || $value[9] < 0 || $value[10] < 0) {
                    notice(100000, ['目标值应大于0']);
                }
                if ($value[7] > 1000000 || $value[8] > 100000 || $value[9] > 100000 || $value[10] > 100000) {
                    notice(100000, ['业绩值范围仅限百万内，人数范围仅限十万内的整数']);
                }
                if (empty($value[0]) || empty($value[1]) || empty($value[2]) || empty($value[3]) || empty($value[4])) {
                    notice(100000, ['不能有空项']);
                }
                if (empty($allMobileArr[$value[6]]) || $allMobileArr[$value[6]] != $value[5]) {
                    notice(100000, ['此手机号和导购名称不匹配或不在导购列表中']);
                }
                //忽略重复手机号
                if (in_array($value[6], $mobileArr)) {
                    continue;
                }
                $mobileArr[] = $value[6];
                $params['model_array'][] = [
                    //'userid'                              => $userIdList[$value[6]] ?? '',
                    'regional' => $value[0],
                    'area' => $value[1],
                    'city' => $value[2],
                    'counter_code' => $value[3],
                    'counter_name' => $value[4],
                    'ba_name' => $value[5],
                    'mobile' => $value[6],
                    'monthly_sales' => (int)$value[7],
                    'monthly_second_buy_and_three_scan' => (int)$value[8],
                    'monthly_old_customer_and_three_scan' => (int)$value[9],
                    'monthly_new_customer_and_three_scan' => (int)$value[10],
                    'operation_user_id' => $this->enterpriseAdminId,
                    'month' => substr(date("Ymd"), 0, 6)
                ];
                $result['success_count']++;
            }
        } catch (\Exception $e) {
            $value = array_slice($value, 0, count(self::EXCEL_TITLE));
            $value[] = $e->getMessage();
            $result['error_list'][] = $value;
            $result['error_count']++;
            $result['success_count'] = 0;
            error($e->getCode(), $row . '行错误:' . $e->getMessage(), $result);
        }
        if (!empty($params['model_array']) && $result['error_count'] == 0) {
            try {
                (new BaPerformance())->Create($params);
            } catch (\Exception $e) {
                error($e->getCode(), $e->getMessage());
            }
        }

        return success($result);
    }

    /**
     * 下载模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function download()
    {

        $path = resource_path('static/计划任务标签人群包模板.csv');
        return response()->download($path);
    }

    /**
     * 下载模板(客户分组导入名单模板)
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function groupDownload()
    {

        $path = resource_path('static/客户分组导入模板.xlsx');
        return response()->download($path);
    }

    /**
     * 下载模板(计划任务导入名单模板)
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function planTaskDownload()
    {
        $task_type = $this->request->get("task_type");
        $path = resource_path('static/计划任务客户人群包模板.xlsx');
        if (!empty($task_type) && $task_type == 4) {
            $path = resource_path('static/计划任务客户&&SaCode人群包模板.xlsx');
        }

        return response()->download($path);
    }

    /**
     * 验证输入的手机号码
     *
     * @return bool
     */
    function is_mobile($user_mobile)
    {
        $chars = "/^((\(\d{2,3}\))|(\d{3}\-))?1(3|5|6|7|8|9)\d{9}$/";
        if (preg_match($chars, $user_mobile)) {
            return true;
        } else {
            return false;
        }
    }

    /**
     *获取所有智慧导购,手机和姓名
     */
    function allStormMem()
    {
        // 导购搜索条件
        $search = [
            'page' => $params['page'] ?? Constant::DEFAULT_PAGE,
            'size' => $params['size'] ?? Constant::MAX_DEFAULT_SIZE,
            'enterprise_id' => $this->enterpriseId,
            'worker_status' => Constant::IN_SERVICE
        ];
        $list = (new StoreWorker())->getList($search);
        return array_column($list['list'] ?? [], 'name', 'mobile');
    }
}