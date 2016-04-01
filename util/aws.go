package util

import (
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/route53"
)

func NewRoute53Service(keyId string, accessKey string) *route53.Route53 {
	provider := credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     keyId,
			SecretAccessKey: accessKey,
			SessionToken:    "",
		},
	}

	return route53.New(session.New(aws.NewConfig().WithCredentials(credentials.NewCredentials(&provider))))
}
