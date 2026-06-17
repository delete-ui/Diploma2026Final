package seed

import (
	"GolangBackendDiploma26/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func SeedBatteries(ctx context.Context, db *sql.DB) error {
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM batteries").Scan(&count)
	if err != nil {
		return fmt.Errorf("count batteries: %w", err)
	}
	if count > 0 {
		return nil // данные уже есть
	}

	batteries := []models.Battery{
		{
			Title: "EUA901", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/6978623f35bf0.jpg",
			Brand: "EUROSTART", Voltage: 12, Polarity: "Прямая",
			Capacity: 90, Standart: "ASIA", Technology: "SLI", SizeType: "D31",
		},
		{
			Title: "QX060051B12A", Price: 3290, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/54_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 60, Standart: "EURO", Technology: "SLI", SizeType: "L2",
		},
		{
			Title: "QV072064B13A", Price: 3390, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/31_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 72, Standart: "EURO", Technology: "SLI", SizeType: "L3",
		},
		{
			Title: "QX090072B15A", Price: 3990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/56_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 90, Standart: "EURO", Technology: "SLI", SizeType: "L5",
		},
		{
			Title: "A2305410001", Price: 4290, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/5e6f69bc5b863.jpg",
			Brand: "MERCEDES-BENZ", Voltage: 12, Polarity: "Прямая",
			Capacity: 35, Standart: "EURO", Technology: "AGM", SizeType: "POB4",
		},
		{
			Title: "QE060064A12A", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/6_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 60, Standart: "EURO", Technology: "EFB", SizeType: "L2",
		},
		{
			Title: "QE225120A06C", Price: 4390, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/15_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 225, Standart: "ASIA", Technology: "EFB", SizeType: "D6",
		},
		{
			Title: "QE240125A06C", Price: 3990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/16_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 240, Standart: "EURO", Technology: "EFB", SizeType: "D6",
		},
		{
			Title: "QE086072A31B", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/12_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 86, Standart: "ASIA", Technology: "EFB", SizeType: "D31",
		},
		{
			Title: "595402080", Price: 2390, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20250116174857_32782/134_0.jpg",
			Brand: "VARTA", Voltage: 12, Polarity: "Обратная",
			Capacity: 95, Standart: "ASIA", Technology: "SLI", SizeType: "L5",
		},
		{
			Title: "DT 6012", Price: 490, Stock: 10,
			Img:   "https://ir.ozone.ru/s3/multimedia-w/wc1000/6169084760.jpg",
			Brand: "DELTA", Voltage: 6, Polarity: "Универсальная",
			Capacity: 1.2, Standart: "Universal", Technology: "AGM", SizeType: "SM",
		},
		{
			Title: "QE060064A12A", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/6_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 60, Standart: "ASIA", Technology: "EFB", SizeType: "L2",
		},
		{
			Title: "QV075072B13A", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/36_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 75, Standart: "EURO", Technology: "SLI", SizeType: "L3",
		},
		{
			Title: "ZSA751", Price: 2290, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/69786077b201c.jpg",
			Brand: "ZUBR", Voltage: 12, Polarity: "Прямая",
			Capacity: 75, Standart: "ASIA", Technology: "SLI", SizeType: "D26",
		},
		{
			Title: "ZU750", Price: 1990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20260219125125_39939/69_0.jpg",
			Brand: "ZUBR", Voltage: 12, Polarity: "Обратная",
			Capacity: 75, Standart: "EURO", Technology: "SLI", SizeType: "L3",
		},
		{
			Title: "ZU920", Price: 1990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20260219125125_39939/72_0.jpg",
			Brand: "ZUBR", Voltage: 12, Polarity: "Обратная",
			Capacity: 90, Standart: "EURO", Technology: "SLI", SizeType: "LB5",
		},
		{
			Title: "QV070063A26B", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/28_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 70, Standart: "ASIA", Technology: "SLI", SizeType: "D26",
		},
		{
			Title: "QE100093A15A", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/13_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 100, Standart: "EURO", Technology: "SLI", SizeType: "L5",
		},
		{
			Title: "QV060052A23B", Price: 2290, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/17_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 60, Standart: "ASIA", Technology: "SLI", SizeType: "D23",
		},
		{
			Title: "QV060052B23B", Price: 2290, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/18_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 60, Standart: "ASIA", Technology: "SLI", SizeType: "D23",
		},
		{
			Title: "QE072068A26B", Price: 2290, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/9_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 72, Standart: "ASIA", Technology: "EFB", SizeType: "D26",
		},
		{
			Title: "QV063064A82A", Price: 2290, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/25_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 63, Standart: "EURO", Technology: "SLI", SizeType: "LB2",
		},
		{
			Title: "QV063064B12A", Price: 2390, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/26_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 63, Standart: "EURO", Technology: "SLI", SizeType: "L2",
		},
		{
			Title: "QV068063A83A", Price: 2390, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/27_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 68, Standart: "EURO", Technology: "SLI", SizeType: "LB3",
		},
		{
			Title: "QV070063B26B", Price: 2390, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/29_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 70, Standart: "ASIA", Technology: "SLI", SizeType: "D26",
		},
		{
			Title: "QV072064A13A", Price: 2390, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/30_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 72, Standart: "EURO", Technology: "SLI", SizeType: "L3",
		},
		{
			Title: "QV074070A83A", Price: 2390, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/32_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 74, Standart: "EURO", Technology: "SLI", SizeType: "LB3",
		},
		{
			Title: "QV075068A26B", Price: 2590, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/33_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 75, Standart: "ASIA", Technology: "SLI", SizeType: "D26",
		},
		{
			Title: "QV100085B31B", Price: 2790, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/45_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 75, Standart: "ASIA", Technology: "SLI", SizeType: "D31",
		},
		{
			Title: "QV100090B15A", Price: 2790, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/47_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 100, Standart: "EURO", Technology: "SLI", SizeType: "L5",
		},
		{
			Title: "QV110092A16A", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/48_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 110, Standart: "EURO", Technology: "SLI", SizeType: "L6",
		},
		{
			Title: "QV190125A05C", Price: 3490, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/49_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 190, Standart: "EURO", Technology: "SLI", SizeType: "D5",
		},
		{
			Title: "QV190125B05C", Price: 3490, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/50_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Прямая",
			Capacity: 190, Standart: "EURO", Technology: "SLI", SizeType: "D5",
		},
		{
			Title: "QA105095A16A", Price: 3090, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/4_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 105, Standart: "EURO", Technology: "AGM", SizeType: "L6",
		},
		{
			Title: "QA060066A12A", Price: 2690, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/0_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 60, Standart: "EURO", Technology: "AGM", SizeType: "L2",
		},
		{
			Title: "QA080080A14A", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/2_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 80, Standart: "EURO", Technology: "AGM", SizeType: "L4",
		},
		{
			Title: "QE064062A23B", Price: 2990, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/7_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 64, Standart: "ASIA", Technology: "EFB", SizeType: "D23",
		},
		{
			Title: "QV075072A13A", Price: 2690, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/35_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 75, Standart: "EURO", Technology: "SLI", SizeType: "L3",
		},
		{
			Title: "QV100090A15A", Price: 2790, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/46_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 100, Standart: "EURO", Technology: "SLI", SizeType: "L5",
		},
		{
			Title: "EB441 JSU", Price: 3690, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/6978623f35bf0.jpg",
			Brand: "EUROSTART", Voltage: 12, Polarity: "Прямая",
			Capacity: 90, Standart: "ASIA", Technology: "SLI", SizeType: "D31",
		},
		{
			Title: "QE080075A14A", Price: 2590, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/11_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 80, Standart: "EURO", Technology: "EFB", SizeType: "L4",
		},
		{
			Title: "QE070070A13A", Price: 2590, Stock: 10,
			Img:   "https://abstd.ru/static/img/data/20251003173657_35727/8_0.jpg",
			Brand: "ABSEL", Voltage: 12, Polarity: "Обратная",
			Capacity: 70, Standart: "EURO", Technology: "EFB", SizeType: "L3",
		},
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO batteries 
		(id, title, price, stock, img, brand, voltage, polarity, capacity, standart, technology, size_type)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	for _, b := range batteries {
		b.ID = uuid.New()
		_, err = stmt.ExecContext(ctx, b.ID, b.Title, b.Price, b.Stock,
			b.Img, b.Brand, b.Voltage, b.Polarity, b.Capacity,
			b.Standart, b.Technology, b.SizeType)
		if err != nil {
			return fmt.Errorf("insert battery %s: %w", b.Title, err)
		}
	}

	return tx.Commit()
}
