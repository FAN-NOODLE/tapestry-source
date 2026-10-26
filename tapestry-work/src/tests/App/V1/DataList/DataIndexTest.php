<?php

class DataIndexTest extends BaseTest
{

    /**
     * 详情
     */
    public function testDetail()
    {
        $this->callApi('GET', BaseTest::BASE_URL.'/data-list/data-index');
        $this->assertEquals(200, $this->response->getStatusCode());
    }
}
