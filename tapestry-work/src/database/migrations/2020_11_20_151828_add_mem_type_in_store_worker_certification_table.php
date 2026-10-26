<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

/**
 * 添加mem_type字段到store_worker_certification表
 * Created by PhpStorm.
 * User: Hui Yang
 * Email: 2361313833@qq.com
 * DateTime: 2020/11/20 15:29
 * Class AddMemTypeInStoreWorkerCertificationTable
 */
class AddMemTypeInStoreWorkerCertificationTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('store_worker_certification')) {
            Schema::table('store_worker_certification', function (Blueprint $table) {
                if (!Schema::hasColumn('store_worker_certification', 'mem_type')) {
                    $table->unsignedBigInteger('mem_type')
                        ->default(1)
                        ->after('enterprise_id')
                        ->comment('身份类型：1:导购(默认) 2:合伙人');
                }
            });
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
    }
}
