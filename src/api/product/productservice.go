package product

import (
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
)

var (
	ProductService ProductServiceInterface = &productService{}
)

type ProductServiceInterface interface {
	Create(product *Product) httperrors.HttpErr
	GetOneE(code string) (*Producto, httperrors.HttpErr)
	GetOneE1(code support.Paginator2) (*Response, httperrors.HttpErr)
	GetOneE2(code support.Paginator2) (*Response2, httperrors.HttpErr)
	GetLand(code support.Paginator2) (*Response, httperrors.HttpErr)
	GetRental(code support.Paginator2) (*Response, httperrors.HttpErr)
	GetProperty(code support.Paginator2) (*Response, httperrors.HttpErr)
	GetPropertyType(code support.Paginator2) (*Response, httperrors.HttpErr)
	GetProductsByLocation(loaction string) ([]Product, httperrors.HttpErr)
	GetAll(search support.Paginator) (*Results, httperrors.HttpErr)
	Search(search support.Paginator2) (*Response, httperrors.HttpErr)
	Search2(search support.Paginator2) ([]Product, httperrors.HttpErr)

	GetNavs() (navs *Navs, errors httperrors.HttpErr)

	GetProductsbyMajorcategory(paginators *support.Paginator) (*Results, httperrors.HttpErr)
	GetProductsbycategory(search support.Paginator2) (*Response, httperrors.HttpErr)
	GetProperties() (*Property, httperrors.HttpErr)
	GetProductsFlavours(code string) ([]*Product, httperrors.HttpErr)
	GetProductsbyarrival(code string) (*Newarrivals, httperrors.HttpErr)
	GetProductshotdeals() (*Newarrivals, httperrors.HttpErr)
	GetOne(code string) (*Product, httperrors.HttpErr)
	Results(searcher support.Paginator) (*Results, httperrors.HttpErr)
	GetThree() ([]*Product, httperrors.HttpErr)
	GetFeatured() ([]Product, httperrors.HttpErr)
	Update(code string, product *Product) httperrors.HttpErr
	UpdateFeatured(code string, status bool) httperrors.HttpErr
	UpdateSold(code string, status bool) httperrors.HttpErr
	UpdateCompleted(code string, status bool) httperrors.HttpErr
	UpdateHotdeals(code string, status bool) httperrors.HttpErr
	UpdatePromotion(code string, status bool) httperrors.HttpErr
	Likes(code string, likes int64) httperrors.HttpErr
	AUpdate(code string, b, old, new, buy float64) httperrors.HttpErr
	Delete(id string) (string, httperrors.HttpErr)
	GetSize(size string) ([]Product, httperrors.HttpErr)
}
type productService struct {
	repo ProductRepoInterface
}

func NewProductService(repository ProductRepoInterface) ProductServiceInterface {
	return &productService{
		repository,
	}
}

func (service *productService) Create(product *Product) httperrors.HttpErr {
	err1 := service.repo.Create(product)
	return err1

}

func (service *productService) GetOneE(code string) (*Producto, httperrors.HttpErr) {
	product, err1 := service.repo.GetOneE(code)
	return product, err1
}
func (service *productService) GetOneE1(searcher support.Paginator2) (*Response, httperrors.HttpErr) {
	response, err1 := service.repo.GetOneE1(searcher)
	return response, err1
}
func (service *productService) GetNavs() (navs *Navs, errors httperrors.HttpErr) {
	response, err1 := service.repo.GetNavs()
	return response, err1
}
func (service *productService) GetOneE2(searcher support.Paginator2) (*Response2, httperrors.HttpErr) {
	response, err1 := service.repo.GetOneE2(searcher)
	return response, err1
}
func (service *productService) GetLand(searcher support.Paginator2) (*Response, httperrors.HttpErr) {
	response, err1 := service.repo.GetLand(searcher)
	return response, err1
}
func (service *productService) GetRental(searcher support.Paginator2) (*Response, httperrors.HttpErr) {
	response, err1 := service.repo.GetRental(searcher)
	return response, err1
}
func (service *productService) GetProperty(searcher support.Paginator2) (*Response, httperrors.HttpErr) {
	response, err1 := service.repo.GetProperty(searcher)
	return response, err1
}
func (service *productService) GetPropertyType(searcher support.Paginator2) (*Response, httperrors.HttpErr) {
	response, err1 := service.repo.GetPropertyType(searcher)
	return response, err1
}

func (service *productService) GetProductsbyMajorcategory(paginators *support.Paginator) (*Results, httperrors.HttpErr) {
	product, err1 := service.repo.GetProductsbyMajorcategory(paginators)
	return product, err1
}
func (service *productService) GetProductsbycategory(search support.Paginator2) (*Response, httperrors.HttpErr) {
	product, err1 := service.repo.GetProductsbycategory(search)
	return product, err1
}
func (service *productService) GetProperties() (*Property, httperrors.HttpErr) {
	product, err1 := service.repo.GetProperties()
	return product, err1
}
func (service *productService) GetProductsFlavours(code string) ([]*Product, httperrors.HttpErr) {
	product, err1 := service.repo.GetProductsFlavours(code)
	return product, err1
}
func (service *productService) GetProductsbyarrival(code string) (*Newarrivals, httperrors.HttpErr) {
	product, err1 := service.repo.GetProductsbyarrival(code)
	return product, err1
}
func (service *productService) GetProductshotdeals() (*Newarrivals, httperrors.HttpErr) {
	product, err1 := service.repo.GetProductshotdeals()
	return product, err1
}
func (service *productService) GetOne(code string) (*Product, httperrors.HttpErr) {
	product, err1 := service.repo.GetOne(code)
	return product, err1
}

func (service *productService) Search(search support.Paginator2) (*Response, httperrors.HttpErr) {
	products, err := service.repo.Search(search)
	return products, err
}
func (service *productService) Search2(search support.Paginator2) ([]Product, httperrors.HttpErr) {
	products, err := service.repo.Search2(search)
	return products, err
}
func (service *productService) GetAll(search support.Paginator) (*Results, httperrors.HttpErr) {
	products, err := service.repo.GetAll(search)
	return products, err
}

func (service *productService) Results(search support.Paginator) (*Results, httperrors.HttpErr) {
	products, err := service.repo.Results(search)
	return products, err
}
func (service *productService) GetThree() ([]*Product, httperrors.HttpErr) {
	products, err := service.repo.GetThree()
	return products, err
}

//	func (service *productService) Flavours() ([]*Flavours, httperrors.HttpErr) {
//		products, err := service.repo.Flavours()
//		return products, err
//	}
//
//	func (service *productService) Themed() ([]*Flavours, httperrors.HttpErr) {
//		products, err := service.repo.Themed()
//		return products, err
//	}
func (service *productService) GetFeatured() ([]Product, httperrors.HttpErr) {
	products, err := service.repo.GetFeatured()
	return products, err
}
func (service *productService) GetProductsByLocation(location string) ([]Product, httperrors.HttpErr) {
	products, err := service.repo.GetProductsByLocation(location)
	return products, err
}

func (service *productService) GetSize(size string) ([]Product, httperrors.HttpErr) {
	products, err := service.repo.GetSize(size)
	return products, err
}

func (service *productService) Update(code string, product *Product) httperrors.HttpErr {
	err1 := service.repo.Update(code, product)
	return err1
}

func (service *productService) UpdateFeatured(code string, status bool) httperrors.HttpErr {
	err1 := service.repo.UpdateFeatured(code, status)
	return err1
}
func (service *productService) UpdateSold(code string, status bool) httperrors.HttpErr {
	err1 := service.repo.UpdateSold(code, status)
	return err1
}
func (service *productService) UpdateCompleted(code string, status bool) httperrors.HttpErr {
	err1 := service.repo.UpdateCompleted(code, status)
	return err1
}
func (service *productService) UpdateHotdeals(code string, status bool) httperrors.HttpErr {
	err1 := service.repo.UpdateHotdeals(code, status)
	return err1
}
func (service *productService) UpdatePromotion(code string, status bool) httperrors.HttpErr {
	err1 := service.repo.UpdatePromotion(code, status)
	return err1
}
func (service *productService) Likes(code string, likes int64) httperrors.HttpErr {
	err1 := service.repo.Likes(code, likes)
	return err1
}
func (service *productService) AUpdate(code string, b, old, new, buy float64) httperrors.HttpErr {
	err1 := service.repo.AUpdate(code, b, old, new, buy)
	return err1
}
func (service *productService) Delete(id string) (string, httperrors.HttpErr) {
	success, failure := service.repo.Delete(id)
	return success, failure
}
