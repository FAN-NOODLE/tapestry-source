<?php


namespace App\Micros\DownloadCenter;

use CrsSdk\GoDownload\Download;
use CrsSdk\Shop\DownloadCenter;
use Illuminate\Support\Facades\Redis;

class DownloadCenterMicro
{
    public function lists($params)
    {
        if ($this->getPhpVersion($params["enterprise_id"])) {
            return (new DownloadCenter())->lists($params);
        } else {
            return (new Download())->downloadList($params);
        }
    }

    public function delete($params)
    {
        if ($this->getPhpVersion($params["enterprise_id"])) {
            return (new DownloadCenter())->delete($params);
        } else {
            return (new Download())->downloadDelete($params);
        }
    }

    public function detail($params)
    {
        if ($this->getPhpVersion($params["enterprise_id"]) || $this->getPhpVersion($params["hash"])) {
            return (new DownloadCenter())->download($params);
        } else {
            $params['id'] = intval($params['id']);
            return (new Download())->downloadDetail($params);
        }
    }

    public function getPhpVersion($enterpriseId)
    {

        $key = 'micro:php-download:release:enterprise:list';

        $redis = Redis::connection();

        $enterpriseList = $redis->lrange($key, 0, -1);

        if (in_array($enterpriseId, $enterpriseList)) {
            return true;
        } else {
            return false;
        }
    }
}
