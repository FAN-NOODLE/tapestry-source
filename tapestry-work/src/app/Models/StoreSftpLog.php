<?php


namespace App\Models;


class StoreSftpLog extends Base
{

    protected $table = 'store_sftp_log';

    protected $guarded = [];


    public function errorRecord()
    {
        return $this->hasMany(StoreSftpErrorRecord::class, 'file_log_id');
    }

}
