<?php

class HandoverCountTestTest extends TestCase
{
    protected $url;

    const BASE_URL = 'http://work-api-dev.crs.dev-test.vchangyi.com/service/v1/';

    /**
     * 设定请求url
     * 父级设定请求头签名
     */
    public function setUp()
    {
        parent::setUp();

        $this->url = 'customer/handover-count';
    }

    /**
     * 获取会员集合的 - 待分配和待确认的总数
     *
     */
    public function testHandoverCount()
    {
        $arg = [
            'enterprise_id' => 2,
            'mem_uids' => [
                ''
            ]
        ];

        $assertArg = [
            'data' => [
                'mem_uid' => [
                    '' => [
                        'no_distributed' => 0,
                        'wait_confirm' => 0,
                    ]
                ]
            ]
        ];
        $response = $this->post(self::BASE_URL . $this->url, $arg, $this->head);
        $response->seeJsonContains($assertArg);
        $response->assertResponseStatus(200);
    }
}
