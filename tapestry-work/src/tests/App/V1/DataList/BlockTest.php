<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2019/3/22
 * Time: 18:09
 */

class BlockTest extends BaseTest
{
    /**
     * 详情
     */
    public function testDetail()
    {
        $this->callApi('GET', BaseTest::BASE_URL.'/data-list/block');
        $this->assertEquals(200, $this->response->getStatusCode());
    }
}