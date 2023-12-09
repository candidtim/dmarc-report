package dmarc

import (
	"encoding/xml"
	"io"
	"time"

	"github.com/candidtim/dmarc-report/util"
)

type Feedback struct {
	ReportMetadata  ReportMetadata  `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Record          Record          `xml:"record"`
}

type ReportMetadata struct {
	OrgName   string    `xml:"org_name"`
	ReportID  string    `xml:"report_id"`
	DateRange DateRange `xml:"date_range"`
}

type DateRange struct {
	Begin time.Time `xml:"begin"`
	End   time.Time `xml:"end"`
}

type PolicyPublished struct {
	Domain string `xml:"domain"`
}

type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

type Row struct {
	SourceIP        string          `xml:"source_ip"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

type PolicyEvaluated struct {
	DKIM string `xml:"dkim"`
	SPF  string `xml:"spf"`
}

type Identifiers struct {
	HeaderFrom string `xml:"header_from"`
}

type AuthResults struct {
	DKIM DKIMResult `xml:"dkim"`
	SPF  SPFResult  `xml:"spf"`
}

type DKIMResult struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
}

type SPFResult struct {
	Domain string `xml:"domain"`
	Result string `xml:"result"`
}

func (d *DateRange) UnmarshalXML(dcr *xml.Decoder, start xml.StartElement) error {
	var data struct {
		BeginUnix int64 `xml:"begin"`
		EndUnix   int64 `xml:"end"`
	}

	if err := dcr.DecodeElement(&data, &start); err != nil {
		return err
	}

	d.Begin = time.Unix(data.BeginUnix, 0)
	d.End = time.Unix(data.EndUnix, 0)

	return nil
}

func ParseDMARCReport(filePath string) (Feedback, error) {
	file, err := util.DecompressOpen(filePath)
	if err != nil {
		return Feedback{}, err
	}

	defer file.Close()
	return parseDMARCReportFile(file)
}

func parseDMARCReportFile(reader io.Reader) (Feedback, error) {
	var feedback Feedback
	decoder := xml.NewDecoder(reader)
	err := decoder.Decode(&feedback)
	if err != nil {
		return feedback, err
	}

	return feedback, nil
}
