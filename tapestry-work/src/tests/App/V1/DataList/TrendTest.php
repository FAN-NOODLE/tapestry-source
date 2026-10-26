<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/3/22
 * Time: 18:09
 */

class TrendTest extends BaseTest
{
    public function testData()
    {
        $data['type'] = 1;
        $this->assertEquals(1, $data['type']);
        return $data;
    }
    /**
     * 详情
     *
     * @depends testData
     */
    public function testDetail($data)
    {
        $this->callApi('POST', BaseTest::BASE_URL.'/data-list/trend', $data);
        $this->assertEquals(200, $this->response->getStatusCode());
    }
}