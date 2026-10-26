<?php

abstract class TestCase extends Laravel\Lumen\Testing\TestCase
{
    protected $key;

    protected $time;

    protected $head;

    protected $client;

    /**
     * Creates the application.
     *
     * @return \Laravel\Lumen\Application
     */
    public function createApplication()
    {
        return require __DIR__.'/../bootstrap/app.php';
    }

    /**
     * 设定基境
     */
    public function setUp()
    {
        parent::setUp();

        $this->key    = 'REDACTED_FOR_SOURCE_SHARING';
        $this->time   = time();
        $this->client = 'work';
        $this->head   = $this->buildHeader();
    }

    /**
     * 设定请求头
     * @return array
     */
    public function buildHeader()
    {
        return [
            'Auth-Sign'      => md5($this->client . $this->key . $this->time),
            'Auth-Client'    => 'work',
            'Auth-Timestamp' => $this->time,
        ];
    }
}
