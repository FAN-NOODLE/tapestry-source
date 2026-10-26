<?php

namespace App\Http\Controllers;

use Illuminate\Support\Facades\Log;
use App\Services\EnterpriseLogService;
use App\Services\OpenPlatform\MiniProgramService;

class MonitoringTemplate extends Controller
{
    /**
     * 返回模板数量
     * @param MiniProgramService $programService
     * @return array
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function index(MiniProgramService $programService)
    {
        try {
            $templateList = $programService->getTemplates();
        } catch (\Exception $e) {
            Log::info('自动删除开放平台模版:模版列表获取失败: ' . $e->getMessage());
            EnterpriseLogService::error('自动删除开放平台模版:模版列表获取失败', [
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine()
            ]);
            return ['num' => 0];
        }

        return ['num' => count($templateList)];
    }
}