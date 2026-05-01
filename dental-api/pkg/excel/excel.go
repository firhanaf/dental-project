package excel

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/yourusername/dental-api/internal/model"
)

const headerBgColor = "0F6E56"

func headerStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill: excelize.Fill{Type: "pattern", Color: []string{headerBgColor}, Pattern: 1},
	})
	return style
}

func freezeFirstRow(f *excelize.File, sheet string) {
	f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})
}

func ExportPatients(patients []model.PatientExportRow) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Pasien"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"No. RM", "Nama", "Tgl Lahir", "Usia", "JK", "Telepon", "Alamat", "Alergi", "Status", "Cabang", "Tgl Input"}
	style := headerStyle(f)
	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}

	for i, h := range headers {
		cell := cols[i] + "1"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, style)
	}
	freezeFirstRow(f, sheet)

	for i, p := range patients {
		row := i + 2

		gender := "Laki-laki"
		if p.Gender == model.GenderFemale {
			gender = "Perempuan"
		}
		address := ""
		if p.Address != nil {
			address = *p.Address
		}
		allergy := ""
		if p.AllergyNotes != nil {
			allergy = *p.AllergyNotes
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.NoRM)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.DateOfBirth.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), p.Age)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), gender)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), address)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), allergy)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), string(p.Status))
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), p.BranchName)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), p.CreatedAt.Format("2006-01-02"))
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return &buf, nil
}

func ExportVisits(visits []model.VisitExportRow) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	sheet := "Kunjungan"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"No. RM", "Nama Pasien", "Tgl Kunjungan", "Dokter", "Keluhan", "Diagnosis", "Tindakan", "Gigi", "Biaya", "Cabang"}
	style := headerStyle(f)
	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	for i, h := range headers {
		cell := cols[i] + "1"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, style)
	}
	freezeFirstRow(f, sheet)

	rupiahFmt := `"Rp "#,##0`
	currencyStyle, _ := f.NewStyle(&excelize.Style{CustomNumFmt: &rupiahFmt})

	for i, v := range visits {
		row := i + 2

		diagnosis := ""
		if v.Diagnosis != nil {
			diagnosis = *v.Diagnosis
		}
		treatment := ""
		if v.Treatment != nil {
			treatment = *v.Treatment
		}
		teeth := ""
		if v.TeethInvolved != nil {
			teeth = *v.TeethInvolved
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), v.PatientNoRM)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.PatientName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.VisitDate.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), v.DoctorName)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), v.ChiefComplaint)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), diagnosis)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), treatment)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), teeth)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), v.Cost)
		f.SetCellStyle(sheet, fmt.Sprintf("I%d", row), fmt.Sprintf("I%d", row), currencyStyle)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), v.BranchName)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return &buf, nil
}
