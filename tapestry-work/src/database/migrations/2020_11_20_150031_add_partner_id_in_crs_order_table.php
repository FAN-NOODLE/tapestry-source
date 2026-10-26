<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

/**
 * 添加partner_id到crs_order表
 * Created by PhpStorm.
 * User: Hui Yang
 * Email: 2361313833@qq.com
 * DateTime: 2020/11/20 15:04
 * Class AddPartnerIdInCrsOrderTable
 */
class AddPartnerIdInCrsOrderTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('crs_order')) {
            Schema::table('crs_order', function (Blueprint $table) {
                if (!Schema::hasColumn('crs_order', 'partner_id')) {
                    $table->unsignedBigInteger('partner_id')
                        ->default(0)
                        ->after('store_staff_id')
                        ->comment('分销合伙人ID');
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
