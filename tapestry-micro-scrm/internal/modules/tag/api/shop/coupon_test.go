package shop

//func TestGetWxCouponListByIds(t *testing.T) {
//	type args struct {
//		c   ctx.Context
//		ids []int64
//	}
//	config.CfgData.Service.Shop.Domain = "https://shop-api-test-crs.vchangyi.com/gw-shop"
//	tests := []struct {
//		name    string
//		args    args
//		want    []entity.TagCouponInfo
//		wantErr bool
//	}{
//		//
//		{
//			name: "接口连通性测试",
//			args: args{
//				c:   ctx.New(),
//				ids: []int64{91077255979264},
//			},
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := GetWxCouponListByIds(tt.args.c, tt.args.ids)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetWxCouponListByIds() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetWxCouponListByIds() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
