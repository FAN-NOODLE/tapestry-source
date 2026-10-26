<?php

/**
 * 微服务SDK接口
 */
namespace App\Micros\Activity;

use CrsSdk\MicroActivity\Luck;

class LuckMicro
{

    public function exportActivity($params)
    {
        return (new Luck())->exportActivity($params);
    }

    public function joinUserList($params)
    {
        return (new Luck())->joinUserList($params);
    }

    public function initRankList($params)
    {
        return (new Luck())->initRankList($params);
    }

    public function winList($params)
    {
        return (new Luck())->winList($params);
    }

}
