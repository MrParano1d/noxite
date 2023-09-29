package adapters

import (
	"context"
	"io"

	json "github.com/bytedance/sonic"

	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
)

type PackageAdapter struct {
	usersAdapter ports.UserPort
}

var _ ports.PackagePort = (*PackageAdapter)(nil)

func NewPackageAdapter(usersAdapter ports.UserPort) *PackageAdapter {
	return &PackageAdapter{
		usersAdapter: usersAdapter,
	}
}

func (a *PackageAdapter) ParseManifest(ctx context.Context, r io.Reader) (*entities.PackageVersion, error) {
	var m manifest
	if err := json.ConfigDefault.NewDecoder(r).Decode(&m); err != nil {
		return nil, &ports.PackageAdapterManifestParseError{Err: err}
	}
	manifest, contributorsToCheck, err := ManifestFromPackageJSON(m)
	if err != nil {
		return nil, &ports.PackageAdapterManifestConvertError{Err: err}
	}

	// check if contributors exist
	contributors := make(fields.MixedAuthors, len(contributorsToCheck))
	if len(contributorsToCheck) > 0 {
		knownContributors, err := a.usersAdapter.FindUsersByEmailAddress(ctx, contributorsToCheck)
		if err != nil {
			return nil, &ports.PackageAdapterManifestConvertError{Err: err}
		}

		for i := range contributors {
			packageContributor := m.Versions[manifest.Version.String()].Contributors[i]

			foreignAuthorBuilder := fields.ForeignAuthorBuilder()
			foreignAuthorBuilder.Name(packageContributor.Name)
			if packageContributor.Email != nil {
				foreignAuthorBuilder.Email(*packageContributor.Email)
			}
			if packageContributor.URL != nil {
				foreignAuthorBuilder.Website(*packageContributor.URL)
			}
			foreignAuthor, err := foreignAuthorBuilder.Build()
			if err != nil {
				return nil, &ports.PackageAdapterManifestConvertError{Err: err}
			}
			contributors[i] = foreignAuthor.ToMixedAuthor()
		}

		// override with known contributors
		for i, contributor := range knownContributors {
			contributors[i] = fields.AuthorFromEntityID(contributor.ID).ToMixedAuthor()
		}
	}

	manifest.Contributors = contributors

	// check if author exists
	if m.Versions[manifest.Version.String()].Author != nil {
		email, err := fields.EmailFromString(m.Versions[manifest.Version.String()].Author.Email)
		if err != nil {
			return nil, &ports.PackageAdapterManifestConvertError{Err: err}
		}
		knownAuthor, err := a.usersAdapter.FindUsersByEmailAddress(ctx, []fields.Email{email})
		if err != nil {
			return nil, &ports.PackageAdapterManifestConvertError{Err: err}
		}
		if len(knownAuthor) > 0 {
			mixedAuthor := fields.AuthorFromEntityID(knownAuthor[0].ID).ToMixedAuthor()
			manifest.Author = &mixedAuthor
		} else {
			manifestAuthor, err := fields.ForeignAuthorBuilder().
				Name(m.Versions[manifest.Version.String()].Author.Name).
				Email(m.Versions[manifest.Version.String()].Author.Email).
				Website(m.Versions[manifest.Version.String()].Author.URL).
				Build()
			if err != nil {
				return nil, &ports.PackageAdapterManifestConvertError{Err: err}
			}
			mixedAuthor := manifestAuthor.ToMixedAuthor()
			manifest.Author = &mixedAuthor
		}
	}

	return manifest, nil
}
