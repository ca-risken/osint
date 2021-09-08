package main

import (
	"context"
	"encoding/hex"

	"crypto/sha256"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/osint/pkg/common"
	"github.com/ca-risken/osint/pkg/message"
)

func (s *sqsHandler) putFindings(ctx context.Context, findingMap map[string][]*finding.FindingForUpsert) error {
	for key, findings := range findingMap {
		for _, f := range findings {
			res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
			if err != nil {
				return err
			}
			s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagOsint)
			s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDomain)
			switch key {
			case "Takeover":
				s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagTakeover)
			case "PrivateExpose":
				s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagPrivateExpose)
			case "CertificateExpiration":
				s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagCertificateExpiration)
			}

			appLogger.Infof("Success to PutFinding. finding: %v", res)
		}
	}
	return nil
}

func (s *sqsHandler) tagFinding(ctx context.Context, projectID uint32, findingID uint64, tag string) error {

	_, err := s.findingClient.TagFinding(ctx, &finding.TagFindingRequest{
		ProjectId: projectID,
		Tag: &finding.FindingTagForUpsert{
			FindingId: findingID,
			ProjectId: projectID,
			Tag:       tag,
		}})
	if err != nil {
		appLogger.Errorf("Failed to TagFinding. error: %v", err)
		return err
	}
	return nil
}

func makeFindings(osintResults *[]osintResult, message *message.OsintQueueMessage) (map[string][]*finding.FindingForUpsert, error) {
	findings := map[string][]*finding.FindingForUpsert{}
	findingsTakeover := []*finding.FindingForUpsert{}
	findingsPrivateExpose := []*finding.FindingForUpsert{}
	findingsCertificateExpiration := []*finding.FindingForUpsert{}
	for _, osintResult := range *osintResults {
		isDown := osintResult.Host.isDown()
		findingTakeover, err := osintResult.Takeover.makeFinding(isDown, message.ProjectID, message.DataSource, message.ResourceName)
		if err != nil {
			appLogger.Errorf("Error occured when make Takeover finding. error: %v", err)
			// その他のfindingを登録するため、ログだけ吐いて続行する
		}
		if findingTakeover != nil {
			findingsTakeover = append(findingsTakeover, findingTakeover)
		}
		findingPrivateExpose, err := osintResult.PrivateExpose.makeFinding(osintResult.Host.HostName, message.ProjectID, message.DataSource, message.ResourceName)
		if err != nil {
			appLogger.Errorf("Error occured when make PrivateExpose finding. error: %v", err)
			// その他のfindingを登録するため、ログだけ吐いて続行する
		}
		if findingPrivateExpose != nil {
			findingsPrivateExpose = append(findingsPrivateExpose, findingPrivateExpose)
		}
		findingCertificateExpiration, err := osintResult.CertificateExpiration.makeFinding(message.ProjectID, message.DataSource, message.ResourceName)
		if err != nil {
			appLogger.Errorf("Error occured when make Certificate Expiration finding. error: %v", err)
			// その他のfindingを登録するため、ログだけ吐いて続行する
		}
		if findingCertificateExpiration != nil {
			findingsCertificateExpiration = append(findingsCertificateExpiration, findingCertificateExpiration)
		}
	}
	findings["Takeover"] = findingsTakeover
	findings["PrivateExpose"] = findingsPrivateExpose
	findings["CertificateExpiration"] = findingsCertificateExpiration
	return findings, nil
}

func generateDataSourceID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
