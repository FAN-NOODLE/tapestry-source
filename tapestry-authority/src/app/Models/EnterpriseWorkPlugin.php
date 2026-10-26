<?php

namespace App\Models;

class EnterpriseWorkPlugin extends Base
{
    protected $table = 'enterprise_work_plugin';

    protected $guarded = [];

    public static function firstIdentity($eId)
    {
        $plugin = self::query()->where('enterprise_id', '=', $eId)->first();

        if ($plugin) {
            return $plugin->identifier;
        }

        return '';
    }

    public static function identitySearch($eId, $mark)
    {
        return self::query()
            ->where('enterprise_id', '=', $eId)
            ->where('identifier', '=', $mark)
            ->first();
    }
}