<?php

use Faker\Factory as Faker;

abstract class TestCase extends Laravel\Lumen\Testing\TestCase
{
    protected $faker;
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
     * Set up the test
     */
    public function setUp()
    {
        parent::setUp();
        $this->faker = Faker::create();
    }
}
