<?php

namespace App\Models;

class WelcomeSpeech extends Base
{
    /**
     * 是否默认 0 否
     * @var int
     */
    const IS_DEFAULT_NO = 0;

    /**
     * 是否默认 1 是
     * @var int
     */
    const IS_DEFAULT_YES = 1;

    // 默认欢迎语数据
    const DEFAULT_CONTENT = [
        'id'         => 0,
        'content'    => '您好，欢迎加我',
        'created_at' => 0
    ];

    const NONE_LINK = 1; // 无链接
    const IMAGE = 2; // 图片
    const LINK = 3; // 网页
    const MINI_PROGRAM = 4; // 小程序

    /**
     * 表名
     * @var string
     */
    protected $table = 'welcome_speech';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'contact_qr_code_id',
        'enterprise_id',
        'channel_id',
        'type',
        'content',
        'object',
        'is_default',
        'original_data',
    ];
}
