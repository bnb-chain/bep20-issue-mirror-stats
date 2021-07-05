package fix

import (
	"encoding/csv"
	"github.com/shopspring/decimal"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func ParseCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(record) > 0 {
		return record, nil
	}
	return nil, nil
}

func WriteCSV(fname string, records [][]string) error {
	// write the file
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	if err = w.WriteAll(records); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func main() {
	records, err := ParseCSV("bep20-issue-mirror-stats.csv")
	if err != nil {
		panic(err)
	}
	if len(records) <= 1 {
		panic("no records found")
	}
	for row := range records {
		if row == 0 {
			continue
		}
		totalApplyStr := records[row][2]
		if strings.Contains(totalApplyStr, ",") {
			records[row][2] = strings.ReplaceAll(totalApplyStr, ",", "")
		}
	}

	fix(records)

	if err = WriteCSV("../bep20-issue-mirror-stats.csv", records); err != nil {
		panic(err)
	}

}

func fix(records [][]string) {
	for row := range records {
		if row == 0 {
			continue
		}
		totalApplyStr := records[row][2]
		if noNeedChange(records[row][1]) {
			//fmt.Printf("origin: %s \n, after: %s \n \n", totalApplyStr, totalApplyStr)
			continue
		}
		decimals, err := strconv.Atoi(records[row][3])
		totalSupply, err := decimal.NewFromString(totalApplyStr)
		if err == nil {
			decimals := decimal.NewFromBigInt(big.NewInt(1), int32(decimals))
			totalSupplyNew := totalSupply.Mul(decimals)
			//fmt.Printf("origin: %s \n, after: %s \n \n", totalApplyStr, totalSupplyNew.String())
			records[row][2] = totalSupplyNew.String()
		} else {
			panic(err)
		}
	}
}

func noNeedChange(asset string) bool {
	noChange := []string{"WRX", "SHIB", "JST", "OLDSUN", "SUN"}
	for _, item := range noChange {
		if item == asset {
			return true
		}
	}
	return false

}
