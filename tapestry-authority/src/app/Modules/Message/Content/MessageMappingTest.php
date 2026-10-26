<?php

namespace App\Modules\Message\Content;


use Doctrine\Inflector\Rules\French\Inflectible;
use PHPUnit\Framework\TestCase;

class MessageMappingTest extends TestCase
{
    public $messageMapping;

    public function __construct($name = null, array $data = [], $dataName = '')
    {
        parent::__construct($name, $data, $dataName);
        $this->messageMapping = new MessageMapping();

    }

    protected function getMeta($id)
    {

        $path = __DIR__.'/mini_message.json';
        $json    = file_get_contents($path);
        $metaArr = json_decode($json, 1);

        return collect($metaArr)->where('id', $id)->first()['key_values'] ?? [];
    }


    public function testLive()
    {
        $meta = $this->getMeta(10);

        $d   = "2020-08-04 15:32:00";
        $msg = [
            "live_title" => 111,
            "start_time" => strtotime($d),
        ];

        $r = $this->messageMapping->translate($meta, $msg);


        $this->assertEquals('您预约的直播“111”开始了，点击...', $r["thing2"]['value']);
        $this->assertEquals('2020年08月04日 15:32', $r["date3"]['value']);
    }

    public function testCoupon786()
    {
        $meta = $this->getMeta(12);

        $d   = "2020-08-04 15:32:00";
        $msg = [
            "price"      => 1111,
            "use_status" => '待使用',
            "end_at"     => strtotime($d),
        ];

        $r = $this->messageMapping->translate($meta, $msg);

        $this->assertEquals('¥11.11', $r["amount3"]['value']);
        $this->assertEquals('待使用', $r["phrase1"]['value']);
        $this->assertEquals('2020年08月04日 15:32', $r["date2"]['value']);
        $this->assertEquals('优惠不等人，请在有效期内尽快使用', $r["thing4"]['value']);
    }

    public function testCoupon307()
    {
        $meta = $this->getMeta(13);
        $d    = "2020-08-04 15:32:00";
        $msg  = [
            'miniProgramName' => '天航信息',
            "coupon_name"     => '111111',
            "end_at"          => strtotime($d),
        ];

        $r = $this->messageMapping->translate($meta, $msg);

        $this->assertEquals('天航信息', $r["thing1"]['value']);
        $this->assertEquals('111111', $r["thing2"]['value']);
        $this->assertEquals('2020年08月04日 15:32', $r["time3"]['value']);
    }

    public function testOrderWaitPay818()
    {
        $meta = $this->getMeta(16);

        $d   = "2020-08-04 15:32:00";
        $msg = [
            'wait_pay_money'   => '111111',
            "create_at"        => strtotime($d),
            "order_no"         => '111111',
            "order_goods_info" => '111111',
            "remain_at"        => strtotime($d),
        ];

        $r = $this->messageMapping->translate($meta, $msg);

        $this->assertEquals('¥1111.11', $r["amount1"]['value']);
        $this->assertEquals('2020年08月04日 15:32', $r["date2"]['value']);
        $this->assertEquals('111111', $r["character_string3"]['value']);
        $this->assertEquals('111111', $r["thing4"]['value']);
        $this->assertEquals('商品为您预留到2020年08月04...', $r["thing5"]['value']);
    }


    public function testGroup786()
    {
        $meta = $this->getMeta(59);

        $d   = "2020-08-04 15:32:00";
        $msg = [
            'spu_name'       => '111111',
            "group_price"    => '111111',
            "group_progress" => 1,
            "need_num"       => 10,
            "miss_num"       => 9,
            "end_at"         => strtotime($d),
        ];

        $r = $this->messageMapping->translate($meta, $msg);


        $this->assertEquals('111111', $r["thing1"]['value']);
        $this->assertEquals('¥1111.11', $r["amount2"]['value']);
        $this->assertEquals(10, $r["number3"]['value']);
        $this->assertEquals(9, $r["number4"]['value']);
        $this->assertEquals('00:00', $r["time5"]['value']);
    }

    public function testGroup612()
    {
        $meta = $this->getMeta(186);

        $d   = "2020-08-04 15:32:00";
        $msg = [
            "spu_name"       => "南极人四件套",
            "need_num"       => "4",
            "group_progress" => 1,
            "start_at"       => 1597720637,
            "group_price"    => "1",
            "end_at"         => strtotime($d),
        ];

        $r = $this->messageMapping->translate($meta, $msg);


        $this->assertEquals('111111', $r["thing1"]['value']);
        $this->assertEquals('¥1111.11', $r["amount2"]['value']);
        $this->assertEquals(10, $r["number3"]['value']);
        $this->assertEquals(9, $r["number4"]['value']);
        $this->assertEquals('00:00', $r["time5"]['value']);
    }

    public function test1103()
    {
        $meta = $this->getMeta(198);

        $msg = [
            "sku_name"       => "助力2",
            "sell_price"       => 1,
            "stock_number" => 100,
        ];

        $r = $this->messageMapping->translate($meta, $msg);



        $this->assertEquals('111111', $r["thing1"]['value']);
        $this->assertEquals('¥1111.11', $r["amount2"]['value']);
        $this->assertEquals(10, $r["number3"]['value']);
    }
}
