package refill

import (
	"errors"
	"fmt"
	"sync"
	"time"

	retry "github.com/avast/retry-go"
	petname "github.com/dustinkirkland/golang-petname"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/organizations/organizationsiface"
)

type Refiller interface {
	Create(count int)
}

type generic struct {
	org    organizationsiface.OrganizationsAPI
	config *RefillerOpts
}

type RefillerOpts struct {
	Parallelism         int
	TargetOU            string
	CleanupOU           string
	AccountPrefix       string
	AccountRootUsername string
	AccountRootDomain   string
	AccountRoleName     string
}

func New(organizationsClient organizationsiface.OrganizationsAPI, opts *RefillerOpts) Refiller {

	if opts.Parallelism > 5 {
		opts.Parallelism = 5
	}

	return &generic{
		org:    organizationsClient,
		config: opts,
	}
}

func (r *generic) Create(count int) {

	if count < 1 {
		count = 1
	}

	var waitGroup sync.WaitGroup

	waitGroup.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer waitGroup.Done()
			r.createOneAccount()
		}()
	}

	waitGroup.Wait()

}

func (r *generic) createOneAccount() {
	names := generateAccountNames(r.config.AccountPrefix, r.config.AccountRootUsername, r.config.AccountRootDomain)

	createAccount, _ := r.org.CreateAccount(&organizations.CreateAccountInput{
		AccountName:            &names.name,
		Email:                  &names.email,
		IamUserAccessToBilling: aws.String("ALLOWED"),
		RoleName:               &r.config.AccountRoleName,
	})

	retry.Do(
		func() error {
			status, err := r.org.DescribeCreateAccountStatus(&organizations.DescribeCreateAccountStatusInput{
				CreateAccountRequestId: createAccount.CreateAccountStatus.Id,
			})
			if err != nil {
				return err
			}

			if status.CreateAccountStatus.CompletedTimestamp == nil {
				return errors.New("creation not complete")
			}

			return nil
		},
		retry.Attempts(50),
		retry.Delay(10*time.Second),
	)
}

type accountNames struct {
	email string
	name  string
}

func generateAccountNames(prefix, username, domain string) *accountNames {

	pet := petname.Generate(4, "-")

	return &accountNames{
		email: fmt.Sprintf("%s+%s@%s", username, pet, domain),
		name:  fmt.Sprintf("%s-%s", prefix, pet),
	}
}
