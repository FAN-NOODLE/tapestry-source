<?php

namespace App\Http\Controllers\V1\Excitation;

use App\Http\Controllers\V1\BaseController;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Traits\ExcelTrait;
use Exception;
use Illuminate\Http\JsonResponse;

/**
 * 手动积分导入
 */
class ImportManualPointsController extends BaseController
{
    const TITLE_LINE = 7;//标题行是第8行,数据下标是7

    use ExcelTrait;


    /**
     * 读取模板信息
     * @return JsonResponse
     */
    public function index(): JsonResponse
    {
        // 参数验证
        $fileUrl = $this->request->post('file_url');
        $savePath = storage_path('/') . ltrim(strrchr($fileUrl, '/'), '/');
        $file = "";
        try {
            $file = file_get_contents($fileUrl);
        } catch (Exception $exception) {
            error(1005, [1005 => '文件获取失败'], [], 200);
        }
        file_put_contents($savePath, $file);
        $data = [];
        try {
            $data = $this->excelRead($savePath);
        } catch (Exception $exception) {
            error(1003, [1003 => $exception->getMessage()], [], 200);
        }

        // 判断标题数据
        if (count($data) < self::TITLE_LINE) {
            error(1000, [1000 => "文件头错误"], [], 200);
        }
        if (count($data) == self::TITLE_LINE) {
            error(1001, [1001 => "文件数据为空"], [], 200);
        }
        if ($data[self::TITLE_LINE][0] != "*导购ID" || $data[self::TITLE_LINE][1] != "关联规则ID"
            || $data[self::TITLE_LINE][2] != "规则类目" || $data[self::TITLE_LINE][3] != "*积分获取条件"
            || $data[self::TITLE_LINE][4] != "*积分" || $data[self::TITLE_LINE][5] != "是否补发金币") {
            error(1002, [1002 => "文件格式错误"], [], 200);
        }

        $paramsData = [];
        foreach ($data as $k => $item) {
            foreach ($item as &$v) {
                $v = trim($v);
            }
                //从第8行开始获取实际数据
            if ($k > self::TITLE_LINE) {
                if (empty($item[0])) {
                    continue;
                }
                if (empty($item[3])) {
                    error(1003, [1003 => "积分获取条件不能为空"], [], 200);
                }
                if (empty(intval($item[4]))) {
                    error(1003, [1003 => "积分必须为不为0的数字"], [], 200);
                }
                $item[5] = $item[5] ?? '';
                $paramsData[] = [
                    "user_id" => $item[0],
                    "rule_id" => (int)$item[1],
                    "rule_detail_name" => $item[2],
                    "points_remark" => $item[3],
                    "points" => (int)$item[4],
                    "has_coins" => trim($item[5]) == 'Y',
                ];
            }

        }
        if (empty($paramsData)) {
            return success();
        }
        //调用scrm服务
        $result = [];
        try {
            $result = GoMicroScrm::importManualPoints([
                'enterprise_id' => $this->enterpriseId,
                'medal_points_data' => $paramsData
            ]);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        $this->writeOperatorLog("手动积分导入成功,操作人:" . $this->enterpriseAdminId);
        return success($result);
    }

}