<?php

namespace App\Traits;

trait IdGenerateTrait
{
    public function initializeIdGenerateTrait()
    {
        $this->{$this->primaryKey} = idGenerate();
    }
}
