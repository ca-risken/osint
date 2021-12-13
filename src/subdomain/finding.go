package main

import (
	"context"
	"encoding/hex"

	"crypto/sha256"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/osint/pkg/common"
	"github.com/ca-risken/osint/pkg/message"
	"github.com/vikyd/zero"
)

func (s *sqsHandler) putFindings(ctx context.Context, findingMap map[string][]*finding.FindingForUpsert, resourceName string) error {
	for key, findings := range findingMap {
		for _, f := range findings {
			res, err := s.findingClient.PutFinding(ctx, &finding.PutFindingRequest{Finding: f})
			if err != nil {
				appLogger.Errorf("Failed to put finding. finding: %v, error: %v", f, err)
				return err
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagOsint); err != nil {
				appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagOsint, err)
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, common.TagDomain); err != nil {
				appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", common.TagDomain, err)
			}
			var tagFindingType string
			switch key {
			case "Takeover":
				tagFindingType = common.TagTakeover
			case "PrivateExpose":
				tagFindingType = common.TagPrivateExpose
			case "CertificateExpiration":
				tagFindingType = common.TagCertificateExpiration
			}
			if !zero.IsZeroVal(tagFindingType) {
				if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, tagFindingType); err != nil {
					appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", tagFindingType, err)
				}
			}
			if err = s.tagFinding(ctx, res.Finding.ProjectId, res.Finding.FindingId, resourceName); err != nil {
				appLogger.Errorf("Failed to tag finding. tag: %v, error: %v", resourceName, err)
			}

			if err = s.putRecommend(ctx, res.Finding.ProjectId, res.Finding.FindingId, key); err != nil {
				appLogger.Errorf("Failed to put recommend. key: %v, error: %v", key, err)
				return err
			}

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

func (s *sqsHandler) putRecommend(ctx context.Context, projectID uint32, findingID uint64, category string) error {
	r := getRecommend(category)
	if r.Type == "" || (r.Risk == "" && r.Recommendation == "") {
		appLogger.Warnf("Failed to get recommendation, Unknown category=%s", category)
		return nil
	}
	if _, err := s.findingClient.PutRecommend(ctx, &finding.PutRecommendRequest{
		ProjectId:      projectID,
		FindingId:      findingID,
		DataSource:     message.SubdomainDataSource,
		Type:           r.Type,
		Risk:           r.Risk,
		Recommendation: r.Recommendation,
	}); err != nil {
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
		findingTakeover, err := osintResult.Takeover.makeFinding(isDown, message.ProjectID, message.DataSource)
		if err != nil {
			appLogger.Errorf("Error occured when make Takeover finding. error: %v", err)
			// その他のfindingを登録するため、ログだけ吐いて続行する
		}
		if findingTakeover != nil {
			findingsTakeover = append(findingsTakeover, findingTakeover)
		}
		findingPrivateExpose, err := osintResult.PrivateExpose.makeFinding(message.ProjectID, message.DataSource)
		if err != nil {
			appLogger.Errorf("Error occured when make PrivateExpose finding. error: %v", err)
			// その他のfindingを登録するため、ログだけ吐いて続行する
		}
		if findingPrivateExpose != nil {
			findingsPrivateExpose = append(findingsPrivateExpose, findingPrivateExpose)
		}
		findingCertificateExpiration, err := osintResult.CertificateExpiration.makeFinding(message.ProjectID, message.DataSource)
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
