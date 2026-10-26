<?php

namespace App\Models;

class ContactQrCodeParams extends Base
{
    /**
     * 数据来源
     */
    const SOURCE_TYPE_DEFAULT = 0; //默认正常创建
    const SOURCE_TYPE_AUTO = 1; //内购员工自建

    /**
     * type
     */
    const MEMBER = 1;     // 员工
    const DEPARTMENT = 2; // 部门

    /**
     * 表名
     * @var string
     */
    protected $table = 'contact_qr_code_params';

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
        'enterprise_id',
        'contact_qr_code_id',
        'object',
        'object_extra',
        'type',
    ];


    /**
     * 内购自建的联系人二维码url
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     * @author: DLong
     */
    public function autoContactCode()
    {
        return $this->belongsTo(ContactQrCode::class, 'contact_qr_code_id')
            ->where([
                ['type', ContactQrCode::CODE_TYPE_ONE],
                ['source_type', ContactQrCode::SOURCE_TYPE_AUTO]
            ])
            ->select(['id', 'transform_url']);
    }
}
