<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

/**
 * 添加mem_type字段到withdrawal_automatic_review表
 * Created by PhpStorm.
 * User: Hui Yang
 * Email: 2361313833@qq.com
 * DateTime: 2020年12月2日 18:12:28
 * Class AddMemTypeInWithdrawalAutomaticReviewTable
 */
class AddMemTypeToWithdrawalAutomaticReviewTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        //
        if (Schema::hasTable('withdrawal_automatic_review')) {
            Schema::table('withdrawal_automatic_review', function (Blueprint $table) {
                if (!Schema::hasColumn('withdrawal_automatic_review', 'mem_type')) {
                    $table->tinyInteger('mem_type')
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
