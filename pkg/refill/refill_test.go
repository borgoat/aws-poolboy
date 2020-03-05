package refill_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"

	"github.com/giorgioazzinnaro/aws-poolboy/pkg/refill"
)

func TestRefill(t *testing.T) {

	iterations := 10

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	org := NewMockOrganizationsAPI(ctrl)

	org.EXPECT().
		CreateAccount(gomock.Any()).
		Return(&organizations.CreateAccountOutput{
			CreateAccountStatus: &organizations.CreateAccountStatus{
				Id:          aws.String("car-abcd1234"),
				AccountId:   aws.String("123456123456"),
				AccountName: aws.String("a-random-account-name"),
			},
		}, nil).
		Times(iterations)

	org.EXPECT().
		DescribeCreateAccountStatus(gomock.Any()).
		DoAndReturn(
			func(in *organizations.DescribeCreateAccountStatusInput) (*organizations.DescribeCreateAccountStatusOutput, error) {
				var out organizations.DescribeCreateAccountStatusOutput
				var s organizations.CreateAccountStatus

				s.SetId(*in.CreateAccountRequestId)
				s.SetAccountId("123456123456")
				s.SetAccountName("a-random-account-name")
				s.SetCompletedTimestamp(time.Now())

				out.SetCreateAccountStatus(&s)

				return &out, nil
			},
		).
		AnyTimes()

	r := refill.New(org, &refill.RefillerOpts{
		Parallelism:         5,
		AccountPrefix:       "pool",
		AccountRoleName:     "PoolboyAccountAccessRole",
		AccountRootDomain:   "example.com",
		AccountRootUsername: "poolboy",
		CleanupOU:           "",
		TargetOU:            "",
	})

	r.Create(iterations)
}
