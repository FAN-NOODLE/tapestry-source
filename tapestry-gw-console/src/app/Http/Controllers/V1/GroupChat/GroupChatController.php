<?php

namespace App\Http\Controllers\V1\GroupChat;

use App\Http\Controllers\V1\BaseController;

class GroupChatController extends BaseController
{
    protected $requireLogin = false;
    /**
     * 下载模板
     *
     * @return \Symfony\Component\HttpFoundation\BinaryFileResponse
     */
    public function download()
    {
        $fileType = $this->request->get("file_type");
        switch ($fileType) {
            case "chat-setting":
                $path = resource_path('static/企微社群批量设置模板.xlsx');
                break;
            case "mars_code":
                $path = resource_path('static/计划子任务导入mars_code模板.xlsx');
                break;
            case "userids":
                $path = resource_path('static/计划子任务导入userid模板.xlsx');
                break;
            default:
                notice(100000, ['参数错误']);
        }
        return response()->download($path);
    }
}
