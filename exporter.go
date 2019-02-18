package main

import (
    "bufio"
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/Danceiny/go.fastjson"
    . "github.com/Danceiny/go.utils"
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "reflect"
    "strings"
    "time"
)

func openRecordFile(siteName string) (f *os.File) {
    files, err := ioutil.ReadDir("./")
    for _, fn := range files {
        fname := fn.Name()
        if strings.HasPrefix(fname, siteName) && strings.HasSuffix(fname, "txt") {
            fmt.Printf("[%s]'s record file [%s] found\n", siteName, fname)
            f, err = os.Open(fn.Name())
            break
        }
    }
    PanicError(err)
    return f
}

func readJson(scanner *bufio.Scanner) *fastjson.JSONObject {
    scanner.Scan()
    jsonStr := scanner.Text()
    err := scanner.Err()
    if err != nil {
        logrus.Error(err)
    }
    if jsonStr == "" {
        return nil
    }
    return fastjson.ParseObject(jsonStr)
}

func Export2Excel() {
    var jo *fastjson.JSONObject
    const siteName = "dianping"
    f := openRecordFile(siteName)
    xlsx := CreateExcel(siteName)
    writeHeader(xlsx, siteName)
    scanner := bufio.NewScanner(f)
    scanner.Split(bufio.ScanLines)
    var i = 1
    jo = readJson(scanner)
    idRowMap := make(map[string]int)
    for ; jo != nil; jo = readJson(scanner) {
        id, _ := jo.GetString("id")
        index, written := idRowMap[id]
        if written {
            WriteJson2Excel(xlsx, siteName, jo, index)
        } else {
            WriteJson2Excel(xlsx, siteName, jo, i)
            idRowMap[id] = i
            i++
        }
    }
    logrus.Infof("lines count: %d", i)
    // Save xlsx file by the given path.
    err := xlsx.Save()
    PanicError(err)

}
func WriteJson2Excel(xlsx *excelize.File, sheetName string, jo *fastjson.JSONObject, row int) {
    for i, field := range fields {
        xlsx.SetCellValue(sheetName,
            fmt.Sprintf("%s%d", string(A_Z[i]), row+1),
            jo.Get(field))
    }
}

func excelFn(siteName string) string {
    y, m, d := time.Now().Date()
    return fmt.Sprintf("%s_%d_%d_%d.xlsx", siteName, y, m, d)
}

func writeHeader(xlsx *excelize.File, sheetName string) {
    for i, field := range fields {
        xlsx.SetCellValue(sheetName,
            fmt.Sprintf("%s%d", string(A_Z[i]), 1),
            field)
    }
    err := xlsx.Save()
    PanicError(err)
}

func CreateExcel(siteName string) (xlsx *excelize.File) {
    var fn = excelFn(siteName)
    var err error
    if _, err = os.Stat(fn); os.IsNotExist(err) {
        xlsx = excelize.NewFile()
        xlsx.Path = fn
        // Create a new sheet.
        xlsx.SetSheetName("Sheet1", siteName)
        // Set active sheet of the workbook.
        // xlsx.SetActiveSheet(index)
    } else {
        xlsx, err = excelize.OpenFile(fn)
        PanicError(err)
    }
    return xlsx

}

const A_Z = "ABCDEFGHIJKLMNOPQRST"

var fields []string

func main() {
    fields = getShopFields()
    logrus.Infof("shop fields: %v", fields)
    Export2Excel()
}

func getShopFields() []string {
    t := reflect.TypeOf(Shop{})
    c := t.NumField()
    var ret = make([]string, c)
    for i := 0; i < c; i++ {
        ret[i] = t.Field(i).Tag.Get("json")
    }
    return ret
}
