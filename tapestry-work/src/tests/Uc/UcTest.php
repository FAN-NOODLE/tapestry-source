<?php


class UcTest extends \PHPUnit\Framework\TestCase
{
    /**
     * 命令  ./vendor/bin/phpunit tests/Uc/UcTest.php
     */
    private static $qyDomain = "B1A503057F0000016FD160F1F22E246C";

    public function testMember()
    {
        $this->markTestSkipped();
        $appIdentifier = null;
        $params        = [
            'mobiles'  => [18735931251],
            'pageNum'  => 1,
            'pageSize' => 1,
        ];

        $result = (new \App\Library\Uc\Member())->ignorePermissionsList(
            static::$qyDomain, $appIdentifier, $params);
        \App\Library\Util\Logger::debug("uc结果", $result);
        $this->assertTrue(count($result['list']) == 1);
    }

    /**
     * @throws Exception
     */
    public function testPage()
    {
        $this->markTestSkipped();
        $appIdentifier = null;
        $pageNum       = 1;

        $list = [];
        do {
            $params = [
                'mobiles'  => [18735931251, 15991264641, 13484926808], //两个有，1个没有
                'pageNum'  => $pageNum,
                'pageSize' => 1,
            ];
            $result = (new \App\Library\Uc\Member())->ignorePermissionsList(
                static::$qyDomain, $appIdentifier, $params);
            if ($result['list']) {
                $list = array_merge($list, $result['list']);
            } else {
                break;
            }

            $pageNum++;
        } while (true);

        \App\Library\Util\Logger::debug("uc翻页结果", $list);

        $this->assertTrue(count($list) == 2, "数量不是2");
    }

    /**
     * @throws Exception
     */
    public function testNoDp()
    {
        $this->markTestSkipped();
        $appIdentifier = null;
        $params        = [
//            'mobiles'=>[18735931251,15991264641],
            'pageNum'  => 1,
            'pageSize' => 1,
        ];

        $result = (new \App\Library\Uc\Member())->ignorePermissionsList(
            static::$qyDomain, $appIdentifier, $params);
        \App\Library\Util\Logger::debug("没有部门的结果", $result);
        $this->assertTrue(count($result['list']) > 0, "没有部门没有人");
    }

    public function testMemSubscribeStatus()
    {
        $appIdentifier = null;
        $params        = [
            'memSubscribeStatus' => 4, //没有订阅
            'pageNum'            => 1,
            'pageSize'           => 1,
        ];

        $result = (new \App\Library\Uc\Member())->ignorePermissionsList(
            static::$qyDomain, $appIdentifier, $params);
        \App\Library\Util\Logger::debug("没有订阅的结果", $result);
        $this->assertTrue($result['list'][0]['memSubscribeStatus'] == 4, "订阅状态断言失败");
    }

    public function testById()
    {
        $appIdentifier = null;
        $params        = [
            'mobiles'  => [18710839146],
            'pageNum'  => 1,
            'pageSize' => 1,
        ];

        $result = (new \App\Library\Uc\Member())->ignorePermissionsList(
            static::$qyDomain, $appIdentifier, $params);
        \App\Library\Util\Logger::debug("不是导购的结果", $result);
        $this->assertTrue($result['list'][0]['memSubscribeStatus'] == 1, "不是导购断言失败");
    }


}