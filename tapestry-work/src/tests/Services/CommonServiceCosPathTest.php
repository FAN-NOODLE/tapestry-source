<?php

use App\Services\CommonService;

class CommonServiceCosPathTest extends TestCase
{
    public function testDefaultCosKeyUsesPublicUploadsDirectory()
    {
        $method = new ReflectionMethod(CommonService::class, 'defaultCosKey');
        $method->setAccessible(true);

        $key = $method->invoke(null);

        $this->assertRegExp('#^uploads/\d{4}/\d{2}/\d{2}/[a-f0-9]{32}\.png$#', $key);
    }
}
