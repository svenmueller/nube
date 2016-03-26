package util

import (
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws/session"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/route53"
	"github.com/svenmueller/nube/Godeps/_workspace/src/github.com/spf13/viper"
)

func NewRoute53Service() *route53.Route53 {
	provider := credentials.StaticProvider{
		Value: credentials.Value{
			AccessKeyID:     viper.GetString("aws-access-key-id"),
			SecretAccessKey: viper.GetString("aws-secret-access-key"),
			SessionToken:    "",
		},
	}

	return route53.New(session.New(aws.NewConfig().WithCredentials(credentials.NewCredentials(&provider))))
}
