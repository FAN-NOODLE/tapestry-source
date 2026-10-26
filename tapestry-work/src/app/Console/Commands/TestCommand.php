<?php

namespace App\Console\Commands;

use App\Library\Uc\EnterprisePlugin;
use App\Sdk\ShenCe\Client;
use Exception;
use Illuminate\Console\Command;
use SimpleSoftwareIO\QrCode\BaconQrCodeGenerator;

class TestCommand extends Command
{
    protected $signature = 'test';

    public function handle()
    {
        $this->shenCe();
    }

    private function shenCe()
    {
        # 从神策分析配置页面中获取的数据接收的 URL
        $SA_SERVER_URL = 'https://analytics-sink.coach.com.cn/sa?project=sw_test';

        try {
            $client = new Client(Client::DebugConsumer, $SA_SERVER_URL);

            // 用户信息上报
//            $distinct_id = '654322';
//            $client->instance()->profile_set($distinct_id, true, [
//                "userid"     => $distinct_id,
//                "mobile"     => "18012345678",
//                "position"   => "Store Manager",
//                "type"       => "sm",
//                "store_code" => "OC007",
//                "gender"     => 1,
//                "avatar"     => "https://wework.qpic.cn/wwpic/842167_Ekfb5t-hRTGmxO8_1668058895/0",
//            ]);


            // 维度上报
            $client->instance()->item_set("reserve_service", "3361318419622720", [
                "reserve_name" => "折扣专场",
                "service_type" => 1,
                "start_time"   => "2006-01-02 15:04",
                "end_time"     => "2006-01-02 15:04",
            ]);

            $client->instance()->flush();
            $client->instance()->close();
        } catch (Exception $e) {
            print_r($e->getMessage());
        }
    }

    private function genQrcode()
    {
        $qrcode = new BaconQrCodeGenerator();
        $str    = base64_encode($qrcode->format("png")->backgroundColor(0, 255, 255)->size(220)->generate('111111111111111111'));
        dd($str);
    }

    private function saveStoreDepartmentInfo()
    {


//        $localFilePath = storage_path() . '/123.csv' ;
//        $handle = fopen($localFilePath, 'r');
//
//        while (($row = fgetcsv($handle)) !== false) {
//
//            // 过滤空行
//            if (empty(array_filter($row))) continue;
//
//            //处理字符串前后空格
////            foreach ($row as $k => &$v) {
////                $v = iconv('gbk', 'utf-8', $v); //解决中文字符乱码问题
////            }
//            StoreDepartment::query()->where('third_store_code', $row[0])
//                ->update(['store_type' => trim($row[2]) == '零售' ? 2 : 1, 'store_name' => $row[1] , 'sync_status' => 0]);
//        }
    }

    /**
     * 注册智汇导购
     * @throws Exception
     */
    private function installMiniInstall()
    {
        $result = (new EnterprisePlugin())->appInstall([
            'corpId'            => 'wwf899a3508c60daa5',
            'installType'       => 1,
            'agentId'           => 1000003,
            'plName'            => '智慧导购小程序',
            'plIdentifier'      => 'crs_store_assistant_self',
            'token'             => $this->randCode(6),
            'secret'            => 'REDACTED_FOR_SOURCE_SHARING',
            'plHost'            => 'tapestry-api-uat.coach.com.cn',
            'encodingAeskey'    => $this->randCode(43),
            'imgName'           => md5(time()),
            'enterpriseChannel' => 4
        ]);
        dd($result);
    }

    private function randCode($length = 5, $type = 0)
    {
        // 生成字符串的值
        $arr = [
            1 => "0123456789",
            2 => "abcdefghijklmnopqrstuvwxyz",
            3 => "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
            4 => "~@#$%^&*(){}[]|"
        ];
        // 初始化
        $code = '';
        // 当前搭配
        if ($type == 0) {
            array_pop($arr);
            $string = implode("", $arr);
        } else if ($type == "-1") {
            $string = implode("", $arr);
        } else {
            $string = $arr[$type];
        }
        // 生成字符串
        $count = strlen($string) - 1;
        for ($i = 0; $i < $length; $i++) {
            $str[$i] = $string[rand(0, $count)];
            $code    .= $str[$i];
        }

        return $code;
    }

    /**
     * 注册通讯录
     * @throws Exception
     */
    private function installCommunalContact()
    {
        $result = (new EnterprisePlugin())->appInstall([
            'corpId'            => 'wwf899a3508c60daa5',
            'installType'       => 2, //通讯录应用
            'agentId'           => '1', // 小程序agentId因为没有，所以是0
            'plName'            => '通用-通讯录标识',
            'plIdentifier'      => 'communalContact',
            'token'             => $this->randCode(6),
            'plHost'            => 'tapestry-api-uat.coach.com.cn',
            'encodingAeskey'    => $this->randCode(43),
            'imgName'           => 'imgs/client/images/appicon/news.png',
            'secret'            => 'REDACTED_FOR_SOURCE_SHARING',
            'corpSecret'        => 'REDACTED_FOR_SOURCE_SHARING',
            "wxqyCorpsecret"    => 'REDACTED_FOR_SOURCE_SHARING',
            'enterpriseChannel' => 6
        ]);
        dd($result);
    }
}