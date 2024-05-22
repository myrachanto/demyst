package product

// import (
// 	"context"
// 	"reflect"
// 	"testing"

// 	httperrors "github.com/myrachanto/erroring"
// 	"github.com/myrachanto/estate/src/support"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// func Test_productrepository_GetProductsbyMajorcategory(t *testing.T) {
// 	type fields struct {
// 		Mongodb *mongo.Database
// 		Cancel  context.CancelFunc
// 	}
// 	type args struct {
// 		p *support.Paginator
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *Results
// 		want1  httperrors.HttpErr
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &productrepository{
// 				Mongodb: tt.fields.Mongodb,
// 				Cancel:  tt.fields.Cancel,
// 			}
// 			got, got1 := r.GetProductsbyMajorcategory(tt.args.p)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("productrepository.GetProductsbyMajorcategory() got = %v, want %v", got, tt.want)
// 			}
// 			if !reflect.DeepEqual(got1, tt.want1) {
// 				t.Errorf("productrepository.GetProductsbyMajorcategory() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }
