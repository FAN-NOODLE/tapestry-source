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
class ImportManualCoinsController extends BaseController
{
    const TITLE_LINE = 5;//标题行是第6行,数据下标是5

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
        if ($data[self::TITLE_LINE][0] != "*导购ID"
            || $data[self::TITLE_LINE][1] != "*金币变动原因"
            || $data[self::TITLE_LINE][2] != "*金币") {
            error(1002, [1002 => "文件格式错误"], [], 200);
        }

        $paramsData = [];
        foreach ($data as $k => $item) {
            foreach ($item as &$v) {
                $v = trim($v);
            }
                //从第7行开始获取实际数据
            if ($k > self::TITLE_LINE) {
                if (empty($item[0])) {
                    continue;
                }
                if (empty($item[1])) {
                    error(1003, [1003 => "金币获取条件不能为空"], [], 200);
                }
                if (empty(floatval($item[2]))) {
                    error(1003, [1003 => "金币必须为不为0的数字"], [], 200);
                }
                $paramsData[] = [
                    "user_id" => $item[0],
                    "coins_remark" => $item[1],
                    "coins" => floatval($item[2]),
                ];
            }

        }
        if (empty($paramsData)) {
            return success();
        }
        //调用scrm服务
        $result = [];
        try {
            $result = GoMicroScrm::importManualCoins([
                'enterprise_id' => $this->enterpriseId,
                'manual_coins_data' => $paramsData
            ]);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        $this->writeOperatorLog("手动金币导入成功,操作人:" . $this->enterpriseAdminId);
        return success($result);
    }

}