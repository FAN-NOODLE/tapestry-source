<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/12 0012
 * Time: 17:26
 */

namespace App\Http\Controllers\V1\ActivityCommon;


use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Activity;
use CrsSdk\MicroActivity\BuyOutPrice;

class DownloadSunCodeController extends ActivityBaseController
{
    /**
     * @throws \Exception
     */
    public function index()
    {
        $input        = $this->getParams();
        $enterpriseId = $this->enterpriseId;

        $input['enterprise_id'] = $enterpriseId;
        $input['size']          = 1000;

        if ($input['type'] == ActivityBaseController::TYPE_BUYOUT_PRICE) {
            $qrCodeInfo = (new BuyOutPrice())->qrCode([
                'rule_id'       => $input['activity_rule_id'],
                'enterprise_id' => $enterpriseId
            ]);

            $zipFileUrl = $this->downLoadBuyoutPricePackage($qrCodeInfo, $input['activity_name']);
        } else {
            // 获取活动信息
            $activityList = (new Activity())->activitySpuList($input);

            /*$activityList = [
                'list' => [
                    [
                        'spu_id'   => 1,
                        'spu_name' => '商品名称1',
                        'qr_code'  => 'https://img-crs-test.vchangyi.com/2020/05/28/61af71bfd2b0811d09ae576ded67e48b.png',
                        'path'     => 1,
                    ],
                    [
                        'spu_id'   => 2,
                        'spu_name' => '商品名称2',
                        'qr_code'  => 'https://img-crs-test.vchangyi.com/2020/05/28/478cbd4b74733be0d047d87536bb9131.png',
                        'path'     => 1,
                    ]
                ]
            ];*/

            $zipFileUrl = $this->DownloadPackage($activityList['list'], $input['activity_name']);
        }

        $this->DownloadZipFile($zipFileUrl, $input['activity_name'] . '.zip');
    }

    public function DownloadZipFile($fileUrl, $fileName)
    {
        ob_start();
        header('Content-Type:text/html;charset=utf-8');
        header('Content-disposition:attachment;filename=' . $fileName);
        $filesize = filesize($fileUrl);
        readfile($fileUrl);
        header('Content-length:' . $filesize);
        ob_end_flush();
        unlink($fileUrl);
    }

    /**
     * 下载一口价二维码
     * @param $qrCodeInfo
     * @param $fileName
     * @return string
     */
    private function downLoadBuyoutPricePackage($qrCodeInfo, $fileName)
    {
        $zip      = new \ZipArchive();
        $path     = base_path('/resources/upload/tmp/');
        $dataFile = [];
        try {
            $filename = $path . $fileName . ".zip"; // 压缩包所在的位置路径

            $zip->open($filename, \ZipArchive::CREATE);   //打开压缩包

            $fileUrl = '';
            if (!empty($qrCodeInfo['qr_code'])) {
                $fileUrl = $this->downloadFile($qrCodeInfo['qr_code'], $fileName);
                $zip->addFile($fileUrl, basename($fileUrl));   //向压缩包中添加文件
            }

            $zip->close();  // 关闭压缩包
        } catch (\Exception $exception) {
            return $exception->getMessage();
        }

        if (file_exists($fileUrl)) {
            unlink($fileUrl);
        }

        return $filename;
    }

    public function DownloadPackage($list, $zipName)
    {
        $zip = new \ZipArchive();

        $path = base_path('/resources/upload/tmp/');

        $dataFile = [];

        try {
            $filename = $path . $zipName . ".zip"; // 压缩包所在的位置路径

            $zip->open($filename, \ZipArchive::CREATE);   //打开压缩包
            foreach ($list as $value) {
                if (!empty($value['qr_code'])) {
                    $fileUrl = $this->downloadFile($value['qr_code'], $value['spu_name']);
                    $zip->addFile($fileUrl, basename($fileUrl));   //向压缩包中添加文件
                    $dataFile[] = $fileUrl;
                }
            }
            $zip->close();  // 关闭压缩包
        } catch (\Exception $exception) {
            return $exception->getMessage();
        }

        foreach ($dataFile as $file) {
            if (file_exists($file)) {
                unlink($file);
            }
        }

        return $filename;
    }

    public function downloadFile($imageUrl, $fileName)
    {
        $path = base_path('/resources/upload/tmp');

        if (!is_dir($path)) {
            mkdir($path, 0777, true);
        }

        $fileUrl = $path . '/' . $fileName . '.jpg';

        $content = file_get_contents($imageUrl);

        if (false === file_put_contents($fileUrl, $content)) {
            throw new Error(20000, '本地图片存储失败!');
        }
        return $fileUrl;
    }

    /**
     * 参数验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams()
    {
        $input = $this->request->input();

        // 参数验证规则
        $rules = [
            'activity_rule_id' => 'required|integer',
            'type'             => 'required|integer',
            'activity_name'    => 'required|string',
        ];

        // 提示语
        $messages = [
            'activity_rule_id.*' => '活动Id有误',
            'type.*'             => '活动类型有误',
            'activity_name.*'    => '活动名称有误'
        ];

        $this->validate($this->request, $rules, $messages);

        return $input;
    }

}
