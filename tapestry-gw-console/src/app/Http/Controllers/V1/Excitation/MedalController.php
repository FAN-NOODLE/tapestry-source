<?php

namespace App\Http\Controllers\V1\Excitation;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use Exception;
use Illuminate\Http\JsonResponse;
use Symfony\Component\HttpFoundation\BinaryFileResponse;

/**
 * 奖章
 */
class MedalController extends BaseController
{
    protected $requireLogin = false;

    use ExcelTrait;

    /**
     * 获奖人群包
     * @return BinaryFileResponse
     */
    public function DownloadCrowdPack(): BinaryFileResponse
    {
        $path = resource_path('static/导购勋章-手动上传员工名单.xlsx');
        return response()->download($path);
    }

    /**
     * 读取人群包
     * @return JsonResponse
     */
    public function ReadCrowdPackData(): JsonResponse
    {
        // 参数验证
        $fileUrl  = $this->request->post('file_url');
        $savePath = storage_path('/') . ltrim(strrchr($fileUrl, '/'), '/');
        file_put_contents($savePath, file_get_contents($fileUrl));
        try {
            $data = $this->excelRead($savePath);

            if (count($data) == 1) {
                throw new Exception("文件数据为空", 1000);
            }
            if ($data[0][0] != "员工ID") {
                throw new Exception("文件格式错误", 1001);
            }
        } catch (Exception $exception) {
            error($exception->getCode(), [$exception->getCode() => $exception->getMessage()], [], 200);
        }

        $responseData = [];
        foreach ($data as $key => $value) {
            if ($key >= 1 && $value != "" && $value[0] != "") {
                if (!in_array($value[0], $responseData)) {
                    $responseData[] = $value[0];
                }
            }
        }

        return success($responseData);
    }

    /**
     * 下载补发模板
     * @return BinaryFileResponse
     */
    public function DownloadReissue(): BinaryFileResponse
    {
        $path = resource_path('static/导购勋章-手动补发奖章.xlsx');
        return response()->download($path);
    }

    /**
     * 读取补发模板
     * @return JsonResponse
     */
    public function ReadReissueData(): JsonResponse
    {
        // 参数验证
        $fileUrl  = $this->request->post('file_url');
        $savePath = storage_path('/') . ltrim(strrchr($fileUrl, '/'), '/');
        file_put_contents($savePath, file_get_contents($fileUrl));
        try {
            $data = $this->excelRead($savePath);
            // 头部判断
            if (count($data) < 5) {
                throw new Exception("文件头错误", 1000);
            }
            if (count($data) == 5) {
                throw new Exception("文件数据为空", 1001);
            }
            if ($data[5][0] != "*导购ID" || $data[5][1] != "*关联规则ID" || $data[5][2] != "*奖章ID") {
                throw new Exception("文件格式错误", 1002);
            }
        } catch (Exception $exception) {
            error($exception->getCode(), [$exception->getCode() => $exception->getMessage()], [], 200);
        }

        $uniqueData = [];
        $result     = [];
        foreach ($data as $k => $item) {
            if ($k < 6) continue;
            // 去重
            if (!in_array(trim($item[0]) . trim($item[2]), $uniqueData)) {
                $result[]     = [
                    "member_user_id" => (string)$item[0],
                    "medal_rule_id"  => (int)$item[1],
                    "medal_id"       => (int)$item[2]
                ];
                $uniqueData[] = trim($item[0]) . trim($item[2]);
            }
        }
        return success($result);
    }
}