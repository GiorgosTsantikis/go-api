package service

import (
	"api/internal/database"
	"api/model"
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type StoreService interface {
	IsValidURL(url string) bool
	GenerateQRCodes(storeId int) ([][]byte, error)
	GenerateQRCodePDF(storeId int) (*bytes.Buffer, error)
	DeleteMenuItem(storeId int, menuItemId int) error
	CreateMenuItem(storeId int32, menuItem *model.MenuItem) (id int64, err error)
}

type storeService struct {
	DB   *database.Queries
	conn *sql.DB
}

func NewStoreService(database *database.Queries, conn *sql.DB) StoreService {
	return &storeService{
		DB:   database,
		conn: conn,
	}
}

func (s *storeService) IsValidURL(url string) bool {
	return false
}

func (s *storeService) CreateMenuItem(storeId int32, menuItem *model.MenuItem) (id int64, err error) {
	menuItem.StoreId = storeId
	fmt.Printf("StoreService.CreateMenuItem received: %v %v \n", menuItem, storeId)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
				fmt.Printf("rollback error: %v\n", rbErr)
			}
		}
	}()

	qtx := s.DB.WithTx(tx)

	for _, co := range menuItem.ConfigOptions {
		for _, opt := range co.Options {
			_, err = qtx.CreateOption(ctx, database.CreateOptionParams{
				ConfigOptionID: sql.NullInt64{
					Int64: int64(opt.ID),
					Valid: true,
				},
				Name:       opt.Name,
				PriceDelta: float64(opt.PriceDelta),
			})
			if err != nil {
				return 0, err
			}
		}
		_, err = qtx.CreateConfigOption(ctx, database.CreateConfigOptionParams{
			Name:          co.Name,
			MaxSelectable: int32(co.MaxSelectable),
			MenuItemID: sql.NullInt64{
				Int64: int64(storeId),
				Valid: true,
			},
		})
		if err != nil {
			return 0, err
		}
	}
	createdMenuItem, err := qtx.CreateMenuItem(ctx, database.CreateMenuItemParams{
		Name:  menuItem.Name,
		Price: menuItem.Price,
	})
	err = tx.Commit()
	return createdMenuItem.ID, err
}

func (s *storeService) DeleteMenuItem(storeId int, menuItemId int) error {

	return nil
}

func (s *storeService) GenerateQRCodes(storeId int) ([][]byte, error) {
	fmt.Printf("GenerateQRCode %d", storeId)
	domain := os.Getenv("FRONTEND_URL")
	store, dbError := s.DB.GetStoreByID(context.Background(), int32(storeId))
	if dbError != nil {
		fmt.Printf("error fetching store by id ", dbError)
		return nil, dbError
	}
	var err error
	var qrImageBytes [][]byte
	var tmpImage []byte
	for i := 1; int32(i) <= store.Tables; i++ {
		url := fmt.Sprintf("%s/%s/%d", domain, store.Name, i)
		tmpImage, err = qrcode.Encode(url, qrcode.Medium, 256)
		qrImageBytes = append(qrImageBytes, tmpImage)
		if err != nil {
			fmt.Println("error encoding qr urls", err)
			return nil, err
		}
	}

	return qrImageBytes, nil
}

func (s *storeService) GenerateQRCodePDF(storeId int) (*bytes.Buffer, error) {
	qrImageBytes, err := s.GenerateQRCodes(storeId)

	if err != nil {
		fmt.Printf("Error generating qr code bytes")
		return nil, err
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetFont("Arial", "", 12)

	const qrSize = 50.0
	const cellPadding = 10.0
	const numCols = 3
	xStart := 10.0
	y := 10.0
	col := 0

	pdf.AddPage()

	for i, qrBytes := range qrImageBytes {
		imgOpt := gofpdf.ImageOptions{
			ImageType:             "PNG",
			ReadDpi:               false,
			AllowNegativePosition: false,
		}

		// Decode PNG from byte slice
		img, decodeErr := png.Decode(bytes.NewReader(qrBytes))
		if decodeErr != nil {
			return nil, decodeErr
		}

		// Convert image.Image into a reader again
		buf := new(bytes.Buffer)
		encodeErr := png.Encode(buf, img)
		if encodeErr != nil {
			fmt.Printf("Error encoding into image")
		}
		imgName := fmt.Sprintf("qr%d", i)

		// Register the image
		pdf.RegisterImageOptionsReader(imgName, imgOpt, buf)

		// Calculate x position
		x := xStart + float64(col)*(qrSize+cellPadding)

		// Place image
		pdf.ImageOptions(imgName, x, y, qrSize, qrSize, false, imgOpt, 0, "")

		// Add label under the image
		pdf.SetXY(x, y+qrSize+2)
		pdf.CellFormat(qrSize, 10, fmt.Sprintf("Table %d", i+1), "", 0, "C", false, 0, "")

		col++
		if col >= numCols {
			col = 0
			y += qrSize + 20 // move to next row
			if y > 260 {     // new page if too low
				pdf.AddPage()
				y = 10
			}
		}
	}

	var out bytes.Buffer
	outputError := pdf.Output(&out)
	if outputError != nil {
		fmt.Printf("Error generating the pdf file")
	}
	return &out, nil
}
