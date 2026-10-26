<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddSourceTypeToChannelTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        if (Schema::hasTable('channel')) {
            Schema::table('channel', function (Blueprint $table) {
                $table->unsignedTinyInteger('source_type')->default(0)->after('channel')
                    ->comment('建立方式，0=默认正常建立，1=内购自建,2=门店商城,3=支付既导购');
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
        if (Schema::hasTable('channel')) {
            Schema::table('channel', function (Blueprint $table) {
                $table->dropColumn('source_type');
            });
        }
    }
}
