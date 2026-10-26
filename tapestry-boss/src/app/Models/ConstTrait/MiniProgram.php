<?php
/*
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Models\ConstTrait;

trait MiniProgram
{
    /**
     * 审核状态
     *
     * @var array
     */
    public static $AUDIT_STATUS = [
        'PENDING'  => 'PENDING',
        'AUDITING' => 'AUDITING',
        'AUDITED'  => 'AUDITED',
        'REFUSED'  => 'REFUSED',
        'WITHDRAW' => 'WITHDRAW'
    ];

    /**
     * 上线状态
     *
     * @var array
     */
    public static $RELEASE_STATUS = [
        'PENDING' => 'PENDING',
        'ONLINE'  => 'ONLINE',
        'OFFLINE' => 'OFFLINE',
    ];
}
