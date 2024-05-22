package product

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	// "mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	httperrors "github.com/myrachanto/erroring"
	"github.com/myrachanto/estate/src/support"
	"github.com/myrachanto/imagery"
	///images manipulation
)

// ProductController ..
var (
	ProductController ProductControllerInterface = &productController{}
	Bizname           string
)

type ProductControllerInterface interface {
	Create(c echo.Context) error
	GetAll(c echo.Context) error
	Results(c echo.Context) error
	GetThree(c echo.Context) error
	GetFeatured(c echo.Context) error
	GetProductsByLocation(c echo.Context) error
	GetOne(c echo.Context) error
	GetOneE(c echo.Context) error
	GetOneE1(c echo.Context) error
	GetOneE2(c echo.Context) error
	GetLand(c echo.Context) error
	GetRental(c echo.Context) error
	GetProperty(c echo.Context) error
	GetPropertyType(c echo.Context) error
	Search(c echo.Context) error
	Search2(c echo.Context) error

	GetNavs(c echo.Context) error

	GetProductsbyMajorcategory(c echo.Context) error
	GetProductsbycategory(c echo.Context) error
	GetProductsFlavours(c echo.Context) error
	GetProductsbyarrival(c echo.Context) error
	GetProductshotdeals(c echo.Context) error
	Update(c echo.Context) error
	UpdateFeatured(c echo.Context) error
	UpdateHotdeals(c echo.Context) error
	UpdatePromotion(c echo.Context) error
	UpdateCompleted(c echo.Context) error
	UpdateSold(c echo.Context) error
	AUpdate(c echo.Context) error
	Likes(c echo.Context) error
	Delete(c echo.Context) error
	GetProperties(c echo.Context) error
}

type productController struct {
	service ProductServiceInterface
}

func NewProductController(service ProductServiceInterface) ProductControllerInterface {
	return &productController{
		service,
	}
}

// ///////controllers/////////////////

// Create godoc
// @Summary Create a product
// @Description Create a new product item
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /api/products [post]
func (controller productController) Create(c echo.Context) error {

	fmt.Println("---------------create product")
	product := &Product{}
	product.Name = c.FormValue("name")
	product.Description = c.FormValue("description")
	product.Footer = c.FormValue("footer")
	product.Title = c.FormValue("title")
	product.Meta = c.FormValue("meta")
	product.Altertag = c.FormValue("altertag")
	product.Category = c.FormValue("category")
	product.Location = c.FormValue("location")
	product.Kind = c.FormValue("kind")
	product.Majorcategory = c.FormValue("majorcategory")
	product.SubLocation = c.FormValue("sublocation")
	product.Video = c.FormValue("video")
	product.Url = strings.ToLower(c.FormValue("url"))
	serv := c.FormValue("services")

	price, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid buying price")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Price = price
	bedrooms, err := strconv.ParseInt(c.FormValue("bedrooms"), 10, 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid bedrooms")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Bedrooms = bedrooms
	bathrooms, err := strconv.ParseInt(c.FormValue("bathrooms"), 10, 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid bathrooms")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Bathrooms = bathrooms
	sqft, err := strconv.ParseFloat(c.FormValue("sqrFt"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid sqft")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Sqft = sqft
	product.Bathrooms = bathrooms
	leng, err := strconv.ParseFloat(c.FormValue("length"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid len")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Length = leng
	wid, err := strconv.ParseFloat(c.FormValue("width"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid width")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Width = wid

	// cols := c.FormValue("colors")
	// sizs := c.FormValue("sizes")
	tags := c.FormValue("tags")
	features := c.FormValue("features")
	// color := c.FormValue("colors")
	// fmt.Println("./////////////////", serv)
	// c.Set("services", serv)
	feat, err := strconv.ParseBool(c.FormValue("featured"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	hot, err1 := strconv.ParseBool(c.FormValue("hotdeals"))
	if err1 != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	prom, err2 := strconv.ParseBool(c.FormValue("promotion"))
	if err2 != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	product.Featured = feat
	product.Hotdeals = hot
	product.Promotion = prom
	var ser []map[string]interface{}
	err3 := json.Unmarshal([]byte(serv), &ser)
	if err3 != nil {
		httperror := httperrors.NewBadRequestError("something went wrong unmarshalling")
		return c.JSON(httperror.Code(), err3.Error())
	}
	// fmt.Println("sssssssssssssssssssssssssssssssssss", ser)
	s := Service{}
	ss := []Service{}
	for _, v := range ser {
		s.Name = fmt.Sprintf("%s", v["name"])
		s.Code = fmt.Sprintf("%s", v["code"])
		ss = append(ss, s)
	}
	product.Services = ss
	// var sproduct []map[string]interface{}
	// err3s := json.Unmarshal([]byte(products), &sproduct)
	// if err3s != nil {
	// 	httperror := httperrors.NewBadRequestError("something went wrong unmarshalling")
	// 	return c.JSON(httperror.Code(), err3s)
	// }
	// t := product{}
	// tt := []product{}
	// for _, v := range sproduct {
	// 	t.Name = fmt.Sprintf("%s", v["name"])
	// 	tt = append(tt, t)
	// }
	// product.product = tt

	t := Tag{}
	ts := []Tag{}
	if string(tags) != "0" {

		var producti []map[string]interface{}
		err4 := json.Unmarshal([]byte(tags), &producti)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("something went wrong unmarshalling products")
			return c.JSON(httperror.Code(), err4.Error())
		}
		fmt.Println("./////////////////step4")

		for _, v := range producti {
			t.Name = fmt.Sprintf("%s", v["name"])
			// t.Code = fmt.Sprintf("%s", v["code"])
			ts = append(ts, t)
		}
		// fmt.Println("./////////////////", ts)
	}
	// fmt.Println("./////////////////step5")
	product.Tag = ts

	fet := Feature{}
	fets := []Feature{}
	if string(features) != "0" {

		var producti []map[string]interface{}
		err4 := json.Unmarshal([]byte(features), &producti)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("something went wrong unmarshalling features")
			return c.JSON(httperror.Code(), err4.Error())
		}

		for _, v := range producti {
			fet.Name = fmt.Sprintf("%s", v["name"])
			// t.Code = fmt.Sprintf("%s", v["code"])
			fets = append(fets, fet)
		}
		// fmt.Println("./////////////////", ts)
	}
	// fmt.Println("./////////////////step5")
	product.Features = fets

	form, errs := c.MultipartForm()
	if errs != nil {
		return errs
	}
	files := form.File["pictures"]
	fmt.Println("files", files)
	img := Picture{}
	imgs := []Picture{}
	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			httperror := httperrors.NewBadRequestError("Invalid picture")
			return c.JSON(httperror.Code(), err.Error())
		}
		defer src.Close()

		// Destination
		Original_Path := strings.Split(file.Filename, ".")
		name1 := Original_Path[len(Original_Path)-1]
		nameSplit := strings.Join(strings.Split(product.Name, " "), "-")
		updated := fmt.Sprintf("-updated-%v", strconv.FormatInt(time.Now().UTC().Unix(), 10))
		imagename := Original_Path[0] + "-" + nameSplit + updated + "." + name1

		filePath := "./src/public/imgs/products/" + file.Filename
		// filePath1 := "/imgs/products/" +"_" + file.Filename
		filepath3 := "./src/public/imgs/products/" + imagename
		dst, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		if len(files) > 0 {
			imagename = fmt.Sprintf("%s-%d.%s", nameSplit, len(imgs), name1)
			imagery.Imageryrepository.Imagetype(filePath, filePath, 500, 800)
			filepath3 = "./src/public/imgs/products/" + imagename
			filePath4 := support.RenameImage(filePath, filepath3)
			filepath5 := "/" + strings.Join(strings.Split(filePath4, "/")[3:], "/")
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>", filepath5)
			img.Name = filepath5
			img.Productcode = product.Code
			imgs = append(imgs, img)
		} else {
			imagery.Imageryrepository.Imagetype(filePath, filePath, 500, 800)
			filePath4 := support.RenameImage(filePath, filepath3)
			filepath5 := "/" + strings.Join(strings.Split(filePath4, "/")[3:], "/")
			fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>", filepath5)
			img.Name = filepath5
			img.Productcode = product.Code
			imgs = append(imgs, img)
		}
	}
	product.Images = imgs
	errs1 := controller.service.Create(product)
	if errs1 != nil {
		return c.JSON(errs1.Code(), errs1.Message())
	}
	return c.JSON(http.StatusCreated, "created successifuly")
	// }
	// err1 := service.controller.service.Create(product)
	// if err1 != nil {
	// 	return c.JSON(err1.Code(), err1)
	// }
	// return c.JSON(http.StatusCreated, "created successifuly")
}

// GetAll godoc
// @Summary GetAll a product
// @Description Getall products
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /front/products [get]
func (controller productController) GetAll(c echo.Context) error {

	search := c.QueryParam("search")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	// fmt.Println("----------------------sdfgghh")
	if err != nil {
		fmt.Println("Invalid pagesize")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator{Page: page, Pagesize: pagesize, Search: search}
	// fmt.Println("---------------", searcher)
	products, err3 := controller.service.GetAll(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, products)
}

// Search godoc
// @Summary Search a product
// @Description Search products
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /front/Search [get]
func (controller productController) Search(c echo.Context) error {

	location := c.QueryParam("location")
	sublocation := c.QueryParam("sublocation")
	kind := c.QueryParam("kind")
	majorcat := c.QueryParam("majorcat")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	// fmt.Println("----------------------sdfgghh")
	if err != nil {
		fmt.Println("Invalid pagesize")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Sublocation: sublocation, Kind: kind, Majorcat: majorcat}
	fmt.Println("---------------", searcher)
	products, err3 := controller.service.Search(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, products)
}

// Search godoc
// @Summary Search a product
// @Description Search products
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /front/Search [get]
func (controller productController) Search2(c echo.Context) error {

	location := c.QueryParam("location")
	majorcat := c.QueryParam("majorcat")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	// fmt.Println("----------------------sdfgghh")
	if err != nil {
		fmt.Println("Invalid pagesize")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Majorcat: majorcat}
	// fmt.Println("---------------", searcher)
	products, err3 := controller.service.Search2(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, products)
}

// @Summary Results a product
// @Description Results products
// @Tags products
// @Accept json
// @Produce json
// @Param        search   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /api/products [get]
func (controller productController) Results(c echo.Context) error {

	search := c.QueryParam("search")
	style := c.Param("style")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	// fmt.Println("----------------------sdfgghh")
	if err != nil {
		fmt.Println("Invalid pagesize")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator{Page: page, Pagesize: pagesize, Search: search, Style: style}
	// fmt.Println("-----------------lasks", style)
	products, err3 := controller.service.Results(searcher)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, products)
}

// @Summary Get Navs a product
// @Description Get Navs
// @Tags products
// @Accept json
// @Produce json
// @Param        search   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/productlimit/four [get]
func (controller productController) GetNavs(c echo.Context) error {

	navs, err3 := controller.service.GetNavs()
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, navs)
}

// @Summary Get Four a product
// @Description Get four products
// @Tags products
// @Accept json
// @Produce json
// @Param        search   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/productlimit/four [get]
func (controller productController) GetThree(c echo.Context) error {

	products, err3 := controller.service.GetThree()
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, products)
}

// @Summary Get Four a product
// @Description Get four products
// @Tags products
// @Accept json
// @Produce json
// @Param        search   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/productlimit/four [get]
func (controller productController) GetSize(c echo.Context) error {

	size := c.Param("size")
	fmt.Println("----------------step 1", size)
	products, err3 := controller.service.GetSize(size)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, products)
}

// // @Summary Get Flavours a product
// // @Description Get Flavours products
// // @Tags products
// // @Accept json
// // @Produce json
// // @Param        search   query     string  false  "code"
// // @Success 201 {object} []Product
// // @Failure 400 {object} support.HttpError
// // @Router /front/Flavours [get]
// func (controller productController) Flavours(c echo.Context) error {
// 	// fmt.Println("----------------step 1")
//
// 	if Bizname == "" {
// 		Bizname = c.QueryParam("bizname")
// 	}
// 	products, err3 := controller.service.Flavours()
// 	if err3 != nil {
// 		return c.JSON(err3.Code(), err3.Message())
// 	}
// 	return c.JSON(http.StatusOK, products)
// }

// // @Summary Get Themed a product
// // @Description Get Themed products
// // @Tags products
// // @Accept json
// // @Produce json
// // @Param        search   query     string  false  "code"
// // @Success 201 {object} []Product
// // @Failure 400 {object} support.HttpError
// // @Router /front/Themed [get]
// func (controller productController) Themed(c echo.Context) error {
// 	// fmt.Println("----------------step 1")
//
// 	if Bizname == "" {
// 		Bizname = c.QueryParam("bizname")
// 	}
// 	products, err3 := controller.service.Themed()
// 	if err3 != nil {
// 		return c.JSON(err3.Code(), err3.Message())
// 	}
// 	return c.JSON(http.StatusOK, products)
// }

// @Summary Get four featured a product
// @Description Get four featured products
// @Tags products
// @Accept json
// @Produce json
// @Param        search   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/featuredproduct/four [get]
func (controller productController) GetFeatured(c echo.Context) error {

	products, err3 := controller.service.GetFeatured()
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, products)
}

// @Summary Geteight a product
// @Description Get eight products
// @Tags products
// @Accept json
// @Produce json
// @Param        search   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/productlimit8 [get]
func (controller productController) GetProductsByLocation(c echo.Context) error {

	location := c.Param("location")
	products, err3 := controller.service.GetProductsByLocation(location)
	if err3 != nil {
		return c.JSON(err3.Code(), err3.Message())
	}
	return c.JSON(http.StatusOK, products)
}

// @Summary Get a product
// @Description Get item
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/products [get]
func (controller productController) GetOne(c echo.Context) error {

	code := c.Param("code")
	product, problem := controller.service.GetOne(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}

func (controller productController) GetOneE(c echo.Context) error {

	code := c.Param("code")
	product, problem := controller.service.GetOneE(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}
func (controller productController) GetOneE2(c echo.Context) error {

	code := c.Param("sublocation")
	location := c.QueryParam("location")
	// sublocation := c.QueryParam("sublocation")
	kind := c.QueryParam("kind")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	if err != nil {
		fmt.Println("Invalid page number")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Sublocation: code, Kind: kind}
	fmt.Printf("+++++++++++++++++++++++++++ %#v \n", code)
	product, problem := controller.service.GetOneE2(searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}
func (controller productController) GetOneE1(c echo.Context) error {

	code := c.Param("town")
	location := c.QueryParam("location")
	sublocation := c.QueryParam("sublocation")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	if err != nil {
		fmt.Println("Invalid page number")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Sublocation: sublocation, Kind: code}
	product, problem := controller.service.GetOneE1(searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}
func (controller productController) GetLand(c echo.Context) error {

	code := c.Param("town")
	location := c.QueryParam("location")
	sublocation := c.QueryParam("sublocation")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	if err != nil {
		fmt.Println("Invalid page number")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Sublocation: sublocation, Kind: code}
	product, problem := controller.service.GetLand(searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}
func (controller productController) GetRental(c echo.Context) error {

	code := c.Param("town")
	location := c.QueryParam("location")
	sublocation := c.QueryParam("sublocation")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	if err != nil {
		fmt.Println("Invalid page number")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Sublocation: sublocation, Kind: code}
	product, problem := controller.service.GetRental(searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}
func (controller productController) GetProperty(c echo.Context) error {

	code := c.Param("town")
	location := c.QueryParam("location")
	sublocation := c.QueryParam("sublocation")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	if err != nil {
		fmt.Println("Invalid page number")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Sublocation: sublocation, Kind: code}
	product, problem := controller.service.GetProperty(searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}
func (controller productController) GetPropertyType(c echo.Context) error {

	code := c.Param("town")
	location := c.QueryParam("location")
	sublocation := c.QueryParam("sublocation")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	if err != nil {
		fmt.Println("Invalid page number")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Sublocation: sublocation, Kind: code}
	product, problem := controller.service.GetPropertyType(searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}

// @Summary GetProductsbyMajorcategory product
// @Description Get Products by Majorcategory
// @Tags products
// @Accept json
// @Produce json
// @Param        type   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/products/majorcategory [get]
func (controller productController) GetProductsbyMajorcategory(c echo.Context) error {

	code := c.Param("type")
	pag := c.QueryParam("page")
	page, err := strconv.Atoi(pag)
	// fmt.Println("----------------------sdfgghh")
	if err != nil {
		fmt.Println("Invalid pagesize")
		page = 1
	}
	ps := c.QueryParam("pageSize")
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	paginators := &support.Paginator{
		Search:   code,
		Page:     page,
		Pagesize: pagesize,
	}
	// fmt.Println("----------------------sdfgghh", paginators)
	products, problem := controller.service.GetProductsbyMajorcategory(paginators)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	fmt.Println()
	// fmt.Println("----------------------sdfgghh", products)
	fmt.Println()
	return c.JSON(http.StatusOK, products)
}

// @Summary GetProductsbycategory product
// @Description Get Products by category
// @Tags products
// @Accept json
// @Produce json
// @Param        type   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/products/category [get]
func (controller productController) GetProductsbycategory(c echo.Context) error {

	code := c.Param("type")
	location := c.QueryParam("location")
	sublocation := c.QueryParam("sublocation")
	ps := c.QueryParam("pagesize")
	pn := c.QueryParam("pagenumber")
	page, err := strconv.Atoi(pn)
	if err != nil {
		fmt.Println("Invalid page number")
		page = 1
	}
	pagesize, err := strconv.Atoi(ps)
	if err != nil {
		fmt.Println("Invalid pagesize")
		pagesize = 10
	}
	searcher := support.Paginator2{Page: page, Pagesize: pagesize, Location: location, Sublocation: sublocation, Kind: code}
	product, problem := controller.service.GetProductsbycategory(searcher)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}

// @Summary GetProductsbycategory product
// @Description Get Products by category
// @Tags products
// @Accept json
// @Produce json
// @Param        type   query     string  false  "code"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/products/category [get]
func (controller productController) GetProperties(c echo.Context) error {

	product, problem := controller.service.GetProperties()
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}

// @Summary GetProductsFlavours product
// @Description Get Products by flavours
// @Tags products
// @Accept json
// @Produce json
// @Param        type   query     string  false  "type"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/products/flavours [get]
func (controller productController) GetProductsFlavours(c echo.Context) error {

	code := c.Param("type")
	product, problem := controller.service.GetProductsFlavours(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}

// @Summary GetProductsbyarrival product
// @Description Get Products by Newest arrivals
// @Tags products
// @Accept json
// @Produce json
// @Param        type   query     string  false  "type"
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/products/newarrivals [get]
func (controller productController) GetProductsbyarrival(c echo.Context) error {

	code := c.Param("type")
	product, problem := controller.service.GetProductsbyarrival(code)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}

// @Summary GetProductshotdeals product
// @Description Get Products by Hot deals
// @Tags products
// @Accept json
// @Produce json
// @Success 201 {object} []Product
// @Failure 400 {object} support.HttpError
// @Router /front/products/hotdeals [get]

func (controller productController) GetProductshotdeals(c echo.Context) error {

	product, problem := controller.service.GetProductshotdeals()
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, product)
}

// GetOne godoc
// @Summary Update a product
// @Description Update a product item
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /api/products [put]
func (controller productController) Update(c echo.Context) error {
	fmt.Println("----------------------product update step 1")

	product := &Product{}
	product.Name = c.FormValue("name")
	product.Description = c.FormValue("description")
	product.Footer = c.FormValue("footer")
	product.Meta = c.FormValue("meta")
	product.Altertag = c.FormValue("altertag")
	product.Title = c.FormValue("title")
	product.Category = c.FormValue("category")
	product.Majorcategory = c.FormValue("majorcategory")
	product.Url = c.FormValue("url")
	// image := c.FormValue("image")
	tamo := c.FormValue("tamo")
	code := c.Param("code")
	tags := c.FormValue("tags")
	features := c.FormValue("features")
	product.Kind = c.FormValue("kind")
	product.SubLocation = c.FormValue("sublocation")
	product.Location = c.FormValue("location")
	product.Video = c.FormValue("video")
	price, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid buying price")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Price = price
	bedrooms, err := strconv.ParseInt(c.FormValue("bedrooms"), 10, 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid bedrooms")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Bedrooms = bedrooms
	bathrooms, err := strconv.ParseInt(c.FormValue("bathrooms"), 10, 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid bathrooms")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Bathrooms = bathrooms
	sqft, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid sqft")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Sqft = sqft
	leng, err := strconv.ParseFloat(c.FormValue("length"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid len")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Length = leng
	wid, err := strconv.ParseFloat(c.FormValue("width"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid width")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	product.Width = wid
	fmt.Println("----------------------product update step 2")
	//products
	t := Tag{}
	ts := []Tag{}
	if tamo == "true" {

		var producti []map[string]interface{}
		err4 := json.Unmarshal([]byte(tags), &producti)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("something went wrong unmarshalling products")
			return c.JSON(httperror.Code(), err4.Error())
		}
		// fmt.Println("./////////////////step4 tags")

		for _, v := range producti {
			t.Name = strings.ToLower(fmt.Sprintf("%s", v["name"]))
			// t.Code = fmt.Sprintf("%s", v["code"])
			ts = append(ts, t)
		}
		// fmt.Println("./////////////////", ts)
	}
	product.Tag = ts
	fmt.Println("./////////////////", len(features))

	fet := Feature{}
	fets := []Feature{}

	if len(features) > 0 {

		var producti []map[string]interface{}
		err4 := json.Unmarshal([]byte(features), &producti)
		if err4 != nil {
			httperror := httperrors.NewBadRequestError("something went wrong unmarshalling features")
			return c.JSON(httperror.Code(), err4.Error())
		}

		for _, v := range producti {
			fet.Name = fmt.Sprintf("%s", v["name"])
			// t.Code = fmt.Sprintf("%s", v["code"])
			fets = append(fets, fet)
		}
		// fmt.Println("./////////////////", ts)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["pictures"]
	fmt.Println("=++++++++++++++++++++++++++===================files", len(files))
	if len(files) > 0 {
		img := Picture{}
		imgs := []Picture{}
		for _, file := range files {
			// Source
			src, err := file.Open()
			if err != nil {
				httperror := httperrors.NewBadRequestError("Invalid picture")
				return c.JSON(httperror.Code(), err.Error())
			}
			defer src.Close()

			// Destination
			Original_Path := strings.Split(file.Filename, ".")
			name1 := Original_Path[len(Original_Path)-1]
			nameSplit := strings.Join(strings.Split(product.Name, " "), "-")
			updated := fmt.Sprintf("-updated-%v", strconv.FormatInt(time.Now().UTC().Unix(), 10))
			imagename := Original_Path[0] + "-" + nameSplit + updated + "." + name1

			filePath := "./src/public/imgs/products/" + file.Filename
			// filePath1 := "/imgs/products/" +"_" + file.Filename
			filepath3 := "./src/public/imgs/products/" + imagename
			dst, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer dst.Close()

			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}
			if len(files) > 0 {
				imagename = fmt.Sprintf("%s-%d.%s", nameSplit, len(imgs), name1)
				imagery.Imageryrepository.Imagetype(filePath, filePath, 500, 800)
				// filepath3 = "./src/public/imgs/products/" + imagename
				filePath4 := support.RenameImage(filePath, filepath3)
				filepath5 := "/" + strings.Join(strings.Split(filePath4, "/")[3:], "/")
				// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>aaa", filepath5)
				img.Name = filepath5
				img.Productcode = product.Code
				imgs = append(imgs, img)
			} else {
				imagery.Imageryrepository.Imagetype(filePath, filePath, 500, 800)
				filePath4 := support.RenameImage(filePath, filepath3)
				filepath5 := "/" + strings.Join(strings.Split(filePath4, "/")[3:], "/")
				// fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>bbbb", filepath5)
				img.Name = filepath5
				img.Productcode = product.Code
				imgs = append(imgs, img)
			}
		}
		// fmt.Println("+++++++++++++++++++++++++++++step huh", imgs[0])
		product.Images = imgs
		err1 := controller.service.Update(code, product)
		if err1 != nil {
			return c.JSON(err1.Code(), err1.Message())
		}
		return c.JSON(http.StatusOK, "update successifuly")
	}
	problem := controller.service.Update(code, product)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusCreated, "Updated successifuly")
}

// GetOne godoc
// @Summary Update Featured a product
// @Description Update a Featured product
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /api/products [put]
func (controller productController) UpdateCompleted(c echo.Context) error {

	code := c.Param("code")
	status := c.FormValue("status")
	fmt.Println("ccccccccccccccccccccc", status)
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	fmt.Println("jjjjjjjjjjjjjjj--------------", feat, code)
	problem := controller.service.UpdateCompleted(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// GetOne godoc
// @Summary Update Sold status a product
// @Description Update a Sold status product
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /api/updateSold [put]
func (controller productController) UpdateSold(c echo.Context) error {

	code := c.Param("code")
	status := c.FormValue("status")
	fmt.Println("ccccccccccccccccccccc", status)
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	fmt.Println("jjjjjjjjjjjjjjj--------------", feat, code)
	problem := controller.service.UpdateSold(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// GetOne godoc
// @Summary Update Featured a product
// @Description Update a Featured product
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /api/products [put]
func (controller productController) UpdateFeatured(c echo.Context) error {

	code := c.Param("code")
	status := c.FormValue("status")
	fmt.Println("ccccccccccccccccccccc", status)
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	fmt.Println("jjjjjjjjjjjjjjj--------------", feat, code)
	problem := controller.service.UpdateFeatured(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// GetOne godoc
// @Summary Update Hotdeals a product
// @Description Update a Hotdeals product
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /api/products [put]
func (controller productController) UpdateHotdeals(c echo.Context) error {

	code := c.Param("code")
	status := c.FormValue("status")
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	problem := controller.service.UpdateHotdeals(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// GetOne godoc
// @Summary Update Promotion a product
// @Description Update a Promotion product
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /api/products [put]
func (controller productController) UpdatePromotion(c echo.Context) error {

	code := c.Param("code")
	status := c.FormValue("status")
	feat, err := strconv.ParseBool(status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the status!")
	}
	problem := controller.service.UpdatePromotion(code, feat)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}

// GetOne godoc
// @Summary Update Likes a product
// @Description Update a Likes product
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} Product
// @Failure 400 {object} support.HttpError
// @Router /api/products [put]
func (controller productController) Likes(c echo.Context) error {

	code := c.Param("code")
	likes := c.FormValue("likes")
	i, err := strconv.ParseInt(likes, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Unable to parse the likes!")
	}
	problem := controller.service.Likes(code, i)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusOK, "updated succesifully")
}
func (controller productController) AUpdate(c echo.Context) error {

	b, err := strconv.ParseFloat(c.FormValue("quantity"), 64)
	if err != nil {
		httperror := httperrors.NewBadRequestError("Invalid buying price")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	old, err1 := strconv.ParseFloat(c.FormValue("oldprice"), 64)
	if err1 != nil {
		httperror := httperrors.NewBadRequestError("Invalid buying price")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	new, err2 := strconv.ParseFloat(c.FormValue("newprice"), 64)
	if err2 != nil {
		httperror := httperrors.NewBadRequestError("Invalid buying price")
		return c.JSON(httperror.Code(), httperror.Message())
	}
	buy, err3 := strconv.ParseFloat(c.FormValue("buyprice"), 64)
	if err3 != nil {
		httperror := httperrors.NewBadRequestError("Invalid buying price")
		return c.JSON(httperror.Code(), httperror.Message())
	}

	code := c.Param("code")
	fmt.Println(code, b, old, new, buy)
	problem := controller.service.AUpdate(code, b, old, new, buy)
	if problem != nil {
		return c.JSON(problem.Code(), problem.Message())
	}
	return c.JSON(http.StatusCreated, "Updated successifuly")
}

// Delete godoc
// @Summary Delete a product
// @Description Create a new product item
// @Tags products
// @Accept json
// @Produce json
// @Param        code   query     string  false  "code"
// @Success 200 {object} string
// @Failure 400 {object} support.HttpError
// @Router /api/products [delete]
func (controller productController) Delete(c echo.Context) error {

	code := c.Param("code")
	// fmt.Println(">>>>>>>>>>>>>>>>>>>>>> ", "there you go!", code)
	success, failure := controller.service.Delete(code)
	if failure != nil {
		return c.JSON(failure.Code(), failure.Message())
	}
	return c.JSON(http.StatusOK, success)

}
