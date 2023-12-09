package dmarc

import (
	"os"
	"strings"
	"testing"
	"time"
)

const SampleReportXMLContent = `
<?xml version="1.0" encoding="UTF-8" ?>
<feedback>
  <report_metadata>
    <org_name>example.com</org_name>
    <OPTIONAL_ATTRIBUTE>OPTIONAL_VALUE</OPTIONAL_ATTRIBUTE>
    <report_id>123</report_id>
    <date_range>
      <begin>1672527600</begin>
      <end>1672614000</end>
    </date_range>
  </report_metadata>
  <policy_published>
    <domain>example.com</domain>
  </policy_published>
  <record>
    <row>
      <source_ip>192.168.0.1</source_ip>
      <policy_evaluated>
        <dkim>pass</dkim>
        <spf>pass</spf>
      </policy_evaluated>
    </row>
    <identifiers>
      <header_from>example.com</header_from>
    </identifiers>
    <auth_results>
      <dkim>
        <domain>example.com</domain>
        <result>pass</result>
      </dkim>
      <spf>
        <domain>example.com</domain>
        <result>pass</result>
      </spf>
    </auth_results>
  </record>
</feedback>
`

var SampleFeedback = Feedback{
	ReportMetadata: ReportMetadata{
		OrgName:  "example.com",
		ReportID: "123",
		DateRange: DateRange{
			Begin: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.Local),
			End:   time.Date(2023, time.January, 2, 0, 0, 0, 0, time.Local),
		},
	},
	PolicyPublished: PolicyPublished{
		Domain: "example.com",
	},
	Record: Record{
		Row: Row{
			SourceIP: "192.168.0.1",
			PolicyEvaluated: PolicyEvaluated{
				DKIM: "pass",
				SPF:  "pass",
			},
		},
		Identifiers: Identifiers{
			HeaderFrom: "example.com",
		},
		AuthResults: AuthResults{
			DKIM: DKIMResult{
				Domain: "example.com",
				Result: "pass",
			},
			SPF: SPFResult{
				Domain: "example.com",
				Result: "pass",
			},
		},
	},
}

func TestParseDMARCReport(t *testing.T) {
	reader := strings.NewReader(SampleReportXMLContent)
	feedback, err := parseDMARCReportFile(reader)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if feedback != SampleFeedback {
		t.Errorf("Expected %v, got %v", "example.com", feedback.ReportMetadata.OrgName)
	}
}

func TestParseDMARCReportFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "dmarc-report*.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name()) // Clean up the temporary file

	if _, err := tempFile.WriteString(SampleReportXMLContent); err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	feedback, err := ParseDMARCReport(tempFile.Name())

	if feedback != SampleFeedback {
		t.Errorf("Expected %v, got %v", "example.com", feedback.ReportMetadata.OrgName)
	}
}
