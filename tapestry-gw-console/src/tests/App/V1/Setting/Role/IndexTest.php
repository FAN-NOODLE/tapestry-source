<?php

namespace Tests\App\V1\Setting\Role;

use Tests\App\V1\BaseTest;

class IndexTest extends BaseTest
{

	/**
	 * 新增角色
	 */
	public function testAdd()
	{
		$this->login();
        $response = $this->call('POST', parent::BASE_URL.'/setting/role/add', [
        	'name' => '测试角色名称',
        	'powers' => '',
        	'type' => 1,
        	'store_ids' => [],
        	'describe' => '测试角色描述'
        ], $this->buildHeader());
        
        $this->assertEquals(200, $response->status());
	}

	/**
	 * 角色详情
	 */
	public function testDetail()
	{
		$this->login();
        $response = $this->call('GET', parent::BASE_URL.'/setting/role/detail', [
        	'role_id' => 1
        ], $this->buildHeader());

        $this->assertEquals(200, $response->status());
	}

	/**
	 * 角色删除
	 */
	public function testDelete()
	{
		$this->login();
        $response = $this->call('POST', parent::BASE_URL.'/setting/role/delete', [
        	'role_id' => 1
        ], $this->buildHeader());

        $this->assertEquals(200, $response->status());
	}

	/**
	 * 角色编辑
	 */
	public function testUpdate()
	{
		$this->login();
        $response = $this->call('POST', parent::BASE_URL.'/setting/role/update', [
        	'role_id' => 1,
        	'name' => '测试角色名称修改',
        	'powers' => '',
        	'type' => 1,
        	'store_ids' => [],
        	'describe' => '测试角色描述'
        ], $this->buildHeader());

        $this->assertEquals(200, $response->status());
	}

}
