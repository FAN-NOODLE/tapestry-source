<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/3/6 10:03
 */
namespace Tests\App\V1\Admin;

use Tests\App\V1\BaseTest;

class IndexTest extends BaseTest
{

    /**
     * 添加
     */
    public function testAdd()
    {
        $this->login();
        $response = $this->call('POST', BaseTest::BASE_URL.'/setting/enterprise-admin/add', [
            'username' => 'testing',
            'account' => '18000000001',
            'password' => 'REDACTED_FOR_SOURCE_SHARING',
            'role_id' => 32
        ], $this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /**
     * 列表
     */
    public function testDetail()
    {
        $this->login();
        $response = $this->call('GET', BaseTest::BASE_URL.'/setting/enterprise-admin/list', [
            'page' => 1,
            'size' => 10
        ], $this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /**
     * 修改密码
     */
    public function testPassword()
    {
        $this->login();
        $response = $this->call('POST', BaseTest::BASE_URL.'/admin/password', [
            'old_password' => 'REDACTED_FOR_SOURCE_SHARING',
            'password' => 'REDACTED_FOR_SOURCE_SHARING',
            're_password' => 'REDACTED_FOR_SOURCE_SHARING'
        ], $this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

    /**
     * 退出登陆
     */
    public function testLogout()
    {
        $this->login();
        $response = $this->call('POST', '/admin/logout', [], $this->buildHeader());
        $this->assertEquals(200, $response->status());
    }

}
