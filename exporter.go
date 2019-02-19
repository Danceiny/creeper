package main

import (
    "bufio"
    "fmt"
    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/Danceiny/go.fastjson"
    . "github.com/Danceiny/go.utils"
    "github.com/sirupsen/logrus"
    "io"
    "io/ioutil"
    "os"
    "strings"
    "time"
)

type FORMAT int
type T int

const (
    _FORMAT = iota
    EXCEL_FORMAT
)
const (
    _T = iota
    T_SHOP
)

type Exporter interface {
    openRecordFile() *os.File
    archiveRecordFile(*os.File)
    exportId(id string, index int, )
    idExportedIndex(id string) (int, bool)
    readJson(scanner *bufio.Scanner) *fastjson.JSONObject
    writeJson(t FORMAT, writer io.WriterTo, sheetName string, row int, object *fastjson.JSONObject)
}

type ExcelExporter struct {
}

func (exporter *ExcelExporter) writeJson2Excel(xlsx *excelize.File, sheetName string, row int,
    object *fastjson.JSONObject, keys []string) {
    if object == nil {
        for i, key := range keys {
            xlsx.SetCellValue(sheetName,
                fmt.Sprintf("%s%d", string(CAPITALS[i]), row+1),
                key)
        }
    } else {
        for i, key := range keys {
            xlsx.SetCellValue(sheetName,
                fmt.Sprintf("%s%d", string(CAPITALS[i]), row+1),
                object.Get(key))
        }
    }

}

type ShopExporter struct {
    ExcelExporter
    siteName      string
    item          TStruct
    keys          []string
    recordFile    *os.File
    exportedIdMap map[string]int
}

func (exporter *ShopExporter) exportId(id string, index int) {
    exporter.exportedIdMap[id] = index
}

func (exporter *ShopExporter) idExportedIndex(id string) (int, bool) {
    v, ok := exporter.exportedIdMap[id]
    return v, ok
}

func (exporter *ShopExporter) writeJson(t FORMAT, writer io.WriterTo, sheetName string, row int, object *fastjson.JSONObject) {
    switch t {
    case EXCEL_FORMAT:
        exporter.writeJson2Excel(writer.(*excelize.File), sheetName, row, object, exporter.keys)
    }
}

func (exporter *ShopExporter) archiveRecordFile(f *os.File) {
    // oldName := f.Name()
    // _ = os.Rename(oldName, oldName+".bak")
    _ = f.Close()
}

func (exporter *ShopExporter) openRecordFile() (f *os.File) {
    var siteName = exporter.siteName
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
    exporter.recordFile = f
    return f
}

func (exporter *ShopExporter) readJson(scanner *bufio.Scanner) *fastjson.JSONObject {
    scanner.Scan()
    jsonStr := scanner.Text()
    var err error
    err = scanner.Err()
    if err != nil {
        logrus.Error(err)
        return nil
    }
    if jsonStr == "" {
        return nil
    }
    return fastjson.ParseObject(jsonStr)
}

func Export2Excel(siteName string, t T) {
    var exporter Exporter
    switch t {
    case T_SHOP:
        exporter = &ShopExporter{siteName: siteName,
            keys:          Shop{}.getShopFields(),
            exportedIdMap: make(map[string]int)}
    }
    f := exporter.openRecordFile()
    xlsx := CreateExcel(siteName)

    // write header
    exporter.writeJson(EXCEL_FORMAT, xlsx, siteName, 0, nil)
    scanner := bufio.NewScanner(f)
    scanner.Split(bufio.ScanLines)
    var i = 1
    st := exporter.readJson(scanner)
    for ; st != nil; st = exporter.readJson(scanner) {
        id, _ := st.GetString("id")
        index, written := exporter.idExportedIndex(id)
        if written {
            exporter.writeJson(EXCEL_FORMAT, xlsx, siteName, index, st)
        } else {
            exporter.writeJson(EXCEL_FORMAT, xlsx, siteName, i, st)
            exporter.exportId(id, i)
            i++
        }
    }
    logrus.Infof("lines count: %d", i)
    // Save xlsx file by the given path.
    err := xlsx.Save()
    PanicError(err)
    exporter.archiveRecordFile(f)
}

func excelFn(sheetName string) string {
    y, m, d := time.Now().Date()
    return fmt.Sprintf("%s_%d_%d_%d.xlsx", sheetName, y, m, d)
}

func CreateExcel(sheetname string) (xlsx *excelize.File) {
    var fn = excelFn(sheetname)
    var err error
    if _, err = os.Stat(fn); os.IsNotExist(err) {
        xlsx = excelize.NewFile()
        xlsx.Path = fn
        // Create a new sheet.
        xlsx.SetSheetName("Sheet1", sheetname)
        // Set active sheet of the workbook.
        // xlsx.SetActiveSheet(index)
    } else {
        xlsx, err = excelize.OpenFile(fn)
        PanicError(err)
    }
    return xlsx

}

func main() {
    Export2Excel("dianping", T_SHOP)
}
