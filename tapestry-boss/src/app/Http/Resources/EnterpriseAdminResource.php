<?php

namespace App\Http\Resources;

use Illuminate\Http\Resources\Json\Resource;
use App\Models\EnterpriseAdmin;

class EnterpriseAdminResource extends Resource
{
    /**
     * Transform the resource collection into an array.
     *
     * @param \Illuminate\Http\Request $request
     *
     * @return array
     */
    public function toArray($request)
    {
        $username = EnterpriseAdmin::query()->where('id', $this->operator_id)->first(['username'])['username'];

        return [
            'id'              => $this->id,
            'enterprise_id'   => $this->enterprise_id,
            'username'        => $this->username,
            'account'         => $this->account,
            'last_login_time' => $this->last_login_time,
            'status'          => $this->status,
            'role'            => $this->role,
            'operator_id'     => $this->operator_id,
            'operator_name'   => $this->operator_id ? $username : '',
            'created_at'      => $this->created_at,
            'updated_at'      => $this->updated_at
        ];
    }
}
