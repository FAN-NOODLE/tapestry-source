<?php

namespace App\Http\Controllers\V1\Notification;

use App\Http\Controllers\V1\BaseController;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Traits\ExcelTrait;
use Exception;
use Illuminate\Http\JsonResponse;

/**
 * 导入门店编码获取员工
 */
class ImportCrowdController extends BaseController
{
    const TITLE_LINE = 0;

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
        if ($data[self::TITLE_LINE][0] != "*门店code") {
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
                $paramsData[] = strval($item[0]);
            }
        }
        if (empty($paramsData)) {
            error(1004, [1004 => "数据为空"], [], 200);
        }
        $uniqueData = array_values(array_unique($paramsData));
        //调用scrm服务
        $result = [];
        try {
            $result = GoMicroScrm::getMemberListByStoreCodes([
                'enterprise_id' => $this->enterpriseId,
                'store_codes' => $uniqueData
            ]);
        } catch (Exception $e) {
            error($e->getCode(), $e->getMessage());
        }
        $result['store_codes'] = $result['store_codes'] ?? [];
        //处理重复数量
        $failStoreCodes = array_diff($uniqueData, $result['store_codes']);

        $res = [
            'result' => [
                'import' => count($paramsData),
                'success' => count($result['store_codes']),
                'repeat' => count(fetchRepeatMemberInArray($paramsData)),
                'fail' => count($failStoreCodes),
                'fail_store_code' => array_values($failStoreCodes)
            ],
            'member_list' => $result['member_list'] ?? [],
        ];
        return success($res);
    }


}