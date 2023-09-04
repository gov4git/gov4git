package github

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"

	"github.com/google/go-github/v54/github"
	"github.com/gov4git/gov4git/gov4git/api"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/vendor4git"
	vendor "github.com/gov4git/vendor4git/github"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
)

func Deploy(
	ctx context.Context,
	token string, // permissions: read project issues, create/write govPrefix
	project GithubRepo,
	govPrefix GithubRepo,
) api.Config {

	// create authenticated GitHub client
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	ghClient := github.NewClient(tc)

	// create governance public and private repos
	v := vendor.NewGithubVendorWithClient(ctx, ghClient)

	govPublic := GithubRepo{Owner: govPrefix.Owner, Name: govPrefix.Name + ":gov.public"}
	base.Infof("creating GitHub repository %v", govPublic)
	govPublicURLs, err := v.CreateRepo(ctx, govPublic.Name, govPublic.Owner, false)
	must.NoError(ctx, err)

	govPrivate := GithubRepo{Owner: govPrefix.Owner, Name: govPrefix.Name + ":gov.private"}
	base.Infof("creating GitHub repository %v", govPrivate)
	govPrivateURLs, err := v.CreateRepo(ctx, govPrivate.Name, govPrivate.Owner, false)
	must.NoError(ctx, err)

	govOwnerAddr := id.OwnerAddress{
		Public: id.PublicAddress{
			Repo:   git.URL(govPublicURLs.HTTPSURL),
			Branch: git.MainBranch,
		},
		Private: id.PrivateAddress{
			Repo:   git.URL(govPrivateURLs.HTTPSURL),
			Branch: git.MainBranch,
		},
	}

	// attach access token authentication to context for git use
	git.SetAuth(ctx, govOwnerAddr.Public.Repo, git.MakeTokenAuth(ctx, token))
	git.SetAuth(ctx, govOwnerAddr.Private.Repo, git.MakeTokenAuth(ctx, token))

	// initialize governance identity
	base.Infof("initializing governance for %v", project)
	id.Init(ctx, govOwnerAddr)

	// create GitHub environment for governance
	base.Infof("creating GitHub environment for governance in %v", govPublic)
	createDeployEnvironment(ctx, ghClient, token, project, govPublic, govPublicURLs, govPrivateURLs)

	// install github automation in the public governance repo
	base.Infof("deploying GitHub actions for governance in %v, targetting %v", govPublic, project)
	XXX

	// return config for gov4git admin
	return api.Config{
		Auth: XXX,
		//
		GovPublicURL:     XXX,
		GovPublicBranch:  XXX,
		GovPrivateURL:    XXX,
		GovPrivateBranch: XXX,
		//
		MemberPublicURL:     XXX,
		MemberPublicBranch:  XXX,
		MemberPrivateURL:    XXX,
		MemberPrivateBranch: XXX,
		//
		CacheDir:        XXX,
		CacheTTLSeconds: XXX,
	}
}

func createDeployEnvironment(
	ctx context.Context,
	ghClient *github.Client,
	token string,
	project GithubRepo,
	govPublic GithubRepo,
	govPublicURLs *vendor4git.Repository,
	govPrivateURLs *vendor4git.Repository,
) {

	// fetch repo id
	ghRepo, _, err := ghClient.Repositories.Get(ctx, project.Owner, project.Name)
	must.NoError(ctx, err)

	// create deploy environment
	createEnv := &github.CreateUpdateEnvironment{}
	env, _, err := ghClient.Repositories.CreateUpdateEnvironment(ctx, govPublic.Owner, govPublic.Name, GithubDeployEnvName, createEnv)
	must.NoError(ctx, err)

	// create environment secrets
	envSecrets := map[string]string{
		"GOV_AUTH_USER":  "",
		"GOV_AUTH_TOKEN": token,
	}

	govPubPubKey, _, err := ghClient.Actions.GetRepoPublicKey(ctx, govPublic.Owner, govPublic.Name)
	must.NoError(ctx, err)

	for k, v := range envSecrets {
		encryptedValue := encryptValue(ctx, govPubPubKey, v)
		_, err := ghClient.Actions.CreateOrUpdateEnvSecret(ctx, int(*ghRepo.ID), env.GetName(),
			&github.EncryptedSecret{
				Name:           k,
				KeyID:          govPubPubKey.GetKeyID(),
				EncryptedValue: encryptedValue,
			})
		must.NoError(ctx, err)
	}

	// create environment variables
	envVars := map[string]string{
		"GITHUB_PROJECT_OWNER": project.Owner,
		"GITHUB_PROJECT_REPO":  project.Name,
		"GOV_PUB_REPO":         govPublicURLs.HTTPSURL,
		"GOV_PRIV_REPO":        govPrivateURLs.HTTPSURL,
	}
	for k, v := range envVars {
		_, err := ghClient.Actions.CreateEnvVariable(ctx, int(*ghRepo.ID), env.GetName(), &github.ActionsVariable{Name: k, Value: v})
		must.NoError(ctx, err)
	}

	XXX
}

func encryptValue(ctx context.Context, pubKey *github.PublicKey, secretValue string) string {

	decodedPubKey, err := base64.StdEncoding.DecodeString(pubKey.GetKey())
	must.NoError(ctx, err)

	var boxKey [32]byte
	copy(boxKey[:], decodedPubKey)
	secretBytes := []byte(secretValue)
	encryptedBytes, err := box.SealAnonymous([]byte{}, secretBytes, &boxKey, crypto_rand.Reader)
	must.NoError(ctx, err)

	return base64.StdEncoding.EncodeToString(encryptedBytes)
}
