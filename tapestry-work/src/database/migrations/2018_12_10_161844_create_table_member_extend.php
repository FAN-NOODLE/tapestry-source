<?php

use \Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateTableMemberExtend extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (!Schema::hasTable('member_extend')) {
            Schema::table('member_avatar', function (Blueprint $table) {
                $table->string('guide_id', 32)->default('')->after('id')->comment('服务人员在支付即服务系统中的唯一标识')->index();
            });

            DB::statement("ALTER TABLE `member_avatar` RENAME as `member_extend`");
        }
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('member_avatar');
    }
}
